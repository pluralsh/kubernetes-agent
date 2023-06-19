package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
)

const (
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 1 * time.Minute

	gitLabKasCookieName             = "_gitlab_kas"
	authorizationHeaderBearerPrefix = "Bearer " // must end with a space
	tokenSeparator                  = ":"
	tokenTypeCi                     = "ci"
)

var (
	code2reason = map[int32]metav1.StatusReason{
		// 4xx
		http.StatusBadRequest:            metav1.StatusReasonBadRequest,
		http.StatusUnauthorized:          metav1.StatusReasonUnauthorized,
		http.StatusForbidden:             metav1.StatusReasonForbidden,
		http.StatusNotFound:              metav1.StatusReasonNotFound,
		http.StatusMethodNotAllowed:      metav1.StatusReasonMethodNotAllowed,
		http.StatusNotAcceptable:         metav1.StatusReasonNotAcceptable,
		http.StatusConflict:              metav1.StatusReasonConflict,
		http.StatusGone:                  metav1.StatusReasonGone,
		http.StatusRequestEntityTooLarge: metav1.StatusReasonRequestEntityTooLarge,
		http.StatusUnsupportedMediaType:  metav1.StatusReasonUnsupportedMediaType,
		http.StatusUnprocessableEntity:   metav1.StatusReasonInvalid,
		http.StatusTooManyRequests:       metav1.StatusReasonTooManyRequests,

		// 5xx
		http.StatusInternalServerError: metav1.StatusReasonInternalError,
		http.StatusServiceUnavailable:  metav1.StatusReasonServiceUnavailable,
		http.StatusGatewayTimeout:      metav1.StatusReasonTimeout,
	}
)

type proxyUserCacheKey struct {
	agentId    int64
	accessType string
	accessKey  string
	csrfToken  string
}

type kubernetesApiProxy struct {
	log                     *zap.Logger
	api                     modserver.Api
	kubernetesApiClient     rpc.KubernetesApiClient
	gitLabClient            gitlab.ClientInterface
	allowedOriginUrls       []string
	allowedAgentsCache      *cache.CacheWithErr[string, *gapi.AllowedAgentsForJob]
	authorizeProxyUserCache *cache.CacheWithErr[proxyUserCacheKey, *gapi.AuthorizeProxyUserResponse]
	requestCounter          usage_metrics.Counter
	ciTunnelUsersCounter    usage_metrics.UniqueCounter
	responseSerializer      runtime.NegotiatedSerializer
	traceProvider           trace.TracerProvider
	tracePropagator         propagation.TextMapPropagator
	meterProvider           metric.MeterProvider
	serverName              string
	serverVia               string
	// urlPathPrefix is guaranteed to end with / by defaulting.
	urlPathPrefix       string
	listenerGracePeriod time.Duration
	shutdownGracePeriod time.Duration
}

func (p *kubernetesApiProxy) Run(ctx context.Context, listener net.Listener) error {
	var handler http.Handler
	handler = http.HandlerFunc(p.proxy)
	handler = otelhttp.NewHandler(handler, "k8s-proxy",
		otelhttp.WithTracerProvider(p.traceProvider),
		otelhttp.WithPropagators(p.tracePropagator),
		otelhttp.WithMeterProvider(p.meterProvider),
		otelhttp.WithPublicEndpoint(),
	)
	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}
	return httpz.RunServer(ctx, srv, listener, p.listenerGracePeriod, p.shutdownGracePeriod)
}

// proxy Kubernetes API calls via agentk to the cluster Kube API.
//
// This method also takes care of CORS preflight requests as documented [here](https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request).
func (p *kubernetesApiProxy) proxy(w http.ResponseWriter, r *http.Request) {
	// for preflight and normal requests we want to allow some configured allowed origins and
	// support exposing the response to the client when credentials (e.g. cookies) are included in the request
	header := w.Header()

	requestedOrigin := r.Header.Get(httpz.OriginHeader)
	if requestedOrigin != "" {
		// If the Origin header is set, it needs to match the configured allowed origin urls.
		if !p.isOriginAllowed(requestedOrigin) {
			// Reject the request because origin is not allowed
			p.log.Sugar().Debugf("Received Origin %q is not in configured allowed origins", requestedOrigin)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		header[httpz.AccessControlAllowOriginHeader] = []string{requestedOrigin}
		header[httpz.AccessControlAllowCredentialsHeader] = []string{"true"}
		header[httpz.VaryHeader] = []string{httpz.OriginHeader}
	}
	header[httpz.ServerHeader] = []string{p.serverName} // It will be removed just before responding with actual headers from upstream

	if r.Method == http.MethodOptions {
		// we have a preflight request
		header[httpz.AccessControlAllowHeadersHeader] = r.Header[httpz.AccessControlRequestHeadersHeader]
		// all allowed HTTP methods:
		// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
		header[httpz.AccessControlAllowMethodsHeader] = []string{"GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH"}
		header[httpz.AccessControlMaxAgeHeader] = []string{"86400"}
		w.WriteHeader(http.StatusOK)
	} else {
		log, agentId, eResp := p.proxyInternal(w, r)
		if eResp != nil {
			p.writeErrorResponse(log, agentId)(w, r, eResp)
		}
	}
}

func (p *kubernetesApiProxy) isOriginAllowed(origin string) bool {
	for _, v := range p.allowedOriginUrls {
		if v == origin {
			return true
		}
	}
	return false
}

func (p *kubernetesApiProxy) proxyInternal(w http.ResponseWriter, r *http.Request) (*zap.Logger, int64 /* agentId */, *grpctool.ErrResp) {
	ctx := r.Context()
	log := p.log.With(logz.TraceIdFromContext(ctx))

	if !strings.HasPrefix(r.URL.Path, p.urlPathPrefix) {
		msg := "Bad request: URL does not start with expected prefix"
		log.Debug(msg, logz.UrlPath(r.URL.Path), logz.UrlPathPrefix(p.urlPathPrefix))
		return log, modshared.NoAgentId, &grpctool.ErrResp{
			StatusCode: http.StatusBadRequest,
			Msg:        msg,
		}
	}

	log, agentId, userId, impConfig, eResp := p.authenticateAndImpersonateRequest(ctx, log, r)
	if eResp != nil {
		// If GitLab doesn't authorize the proxy user to make the call,
		// we send an extra header to indicate that, so that the client
		// can differentiate from an *unauthorized* response from GitLab
		// and from an *authorized* response from the proxied K8s cluster.
		if eResp.StatusCode == http.StatusUnauthorized {
			w.Header()[httpz.GitlabUnauthorizedHeader] = []string{"true"}
		}
		return log, agentId, eResp
	}

	p.requestCounter.Inc() // Count only authenticated and authorized requests
	p.ciTunnelUsersCounter.Add(userId)

	md := metadata.Pairs(modserver.RoutingAgentIdMetadataKey, strconv.FormatInt(agentId, 10))
	mkClient, err := p.kubernetesApiClient.MakeRequest(metadata.NewOutgoingContext(ctx, md))
	if err != nil {
		msg := "Proxy failed to make outbound request"
		p.api.HandleProcessingError(ctx, log, agentId, msg, err)
		return log, agentId, &grpctool.ErrResp{
			StatusCode: http.StatusInternalServerError,
			Msg:        msg,
			Err:        err,
		}
	}

	p.pipeStreams(log, agentId, w, r, mkClient, impConfig) // nolint: contextcheck
	return log, agentId, nil
}

func (p *kubernetesApiProxy) authenticateAndImpersonateRequest(ctx context.Context, log *zap.Logger, r *http.Request) (*zap.Logger, int64 /* agentId */, int64 /* userId */, *rpc.ImpersonationConfig, *grpctool.ErrResp) {
	agentId, creds, err := getAuthorizationInfoFromRequest(r)
	if err != nil {
		msg := "Unauthorized"
		log.Debug(msg, logz.Error(err))
		return log, modshared.NoAgentId, 0, nil, &grpctool.ErrResp{
			StatusCode: http.StatusUnauthorized,
			Msg:        msg,
			Err:        err,
		}
	}
	log = log.With(logz.AgentId(agentId))
	trace.SpanFromContext(ctx).SetAttributes(api.TraceAgentIdAttr.Int64(agentId))

	var (
		userId    int64
		impConfig *rpc.ImpersonationConfig // can be nil
	)

	switch c := creds.(type) {
	case ciJobTokenAuthn:
		allowedForJob, eResp := p.getAllowedAgentsForJob(ctx, log, agentId, c.jobToken)
		if eResp != nil {
			return log, agentId, 0, nil, eResp
		}
		userId = allowedForJob.User.Id

		aa := findAllowedAgent(agentId, allowedForJob)
		if aa == nil {
			msg := "Forbidden: agentId is not allowed"
			log.Debug(msg)
			return log, agentId, userId, nil, &grpctool.ErrResp{
				StatusCode: http.StatusForbidden,
				Msg:        msg,
			}
		}

		impConfig, err = constructJobImpersonationConfig(allowedForJob, aa)
		if err != nil {
			msg := "Failed to construct impersonation config"
			p.api.HandleProcessingError(ctx, log, agentId, msg, err)
			return log, agentId, userId, nil, &grpctool.ErrResp{
				StatusCode: http.StatusInternalServerError,
				Msg:        msg,
				Err:        err,
			}
		}
	case sessionCookieAuthn:
		auth, eResp := p.authorizeProxyUser(ctx, log, agentId, "session_cookie", c.encryptedPublicSessionId, c.csrfToken)
		if eResp != nil {
			return log, agentId, 0, nil, eResp
		}
		userId = auth.User.Id
		impConfig, err = constructUserImpersonationConfig(auth, "session_cookie")
		if err != nil {
			msg := "Failed to construct user impersonation config"
			p.api.HandleProcessingError(ctx, log, agentId, msg, err)
			return log, agentId, userId, nil, &grpctool.ErrResp{
				StatusCode: http.StatusInternalServerError,
				Msg:        msg,
				Err:        err,
			}
		}
	default: // This should never happen
		msg := "Invalid authorization type"
		p.api.HandleProcessingError(ctx, log, agentId, msg, err)
		return log, agentId, userId, nil, &grpctool.ErrResp{
			StatusCode: http.StatusInternalServerError,
			Msg:        msg,
		}
	}
	return log, agentId, userId, impConfig, nil
}

func (p *kubernetesApiProxy) getAllowedAgentsForJob(ctx context.Context, log *zap.Logger, agentId int64, jobToken string) (*gapi.AllowedAgentsForJob, *grpctool.ErrResp) {
	allowedForJob, err := p.allowedAgentsCache.GetItem(ctx, jobToken, func() (*gapi.AllowedAgentsForJob, error) {
		return gapi.GetAllowedAgentsForJob(ctx, p.gitLabClient, jobToken)
	})
	if err != nil {
		var status int32
		var msg string
		switch {
		case gitlab.IsUnauthorized(err):
			status = http.StatusUnauthorized
			msg = "Unauthorized: CI job token"
			log.Debug(msg, logz.Error(err))
		case gitlab.IsForbidden(err):
			status = http.StatusForbidden
			msg = "Forbidden: CI job token"
			log.Debug(msg, logz.Error(err))
		case gitlab.IsNotFound(err):
			status = http.StatusNotFound
			msg = "Not found: agents for CI job token"
			log.Debug(msg, logz.Error(err))
		default:
			status = http.StatusInternalServerError
			msg = "Failed to get allowed agents for CI job token"
			p.api.HandleProcessingError(ctx, log, agentId, msg, err)
		}
		return nil, &grpctool.ErrResp{
			StatusCode: status,
			Msg:        msg,
			Err:        err,
		}
	}
	return allowedForJob, nil
}

func (p *kubernetesApiProxy) authorizeProxyUser(ctx context.Context, log *zap.Logger, agentId int64, accessType, accessKey, csrfToken string) (*gapi.AuthorizeProxyUserResponse, *grpctool.ErrResp) {
	key := proxyUserCacheKey{
		agentId:    agentId,
		accessType: accessType,
		accessKey:  accessKey,
		csrfToken:  csrfToken,
	}
	auth, err := p.authorizeProxyUserCache.GetItem(ctx, key, func() (*gapi.AuthorizeProxyUserResponse, error) {
		return gapi.AuthorizeProxyUser(ctx, p.gitLabClient, agentId, accessType, accessKey, csrfToken)
	})
	if err != nil {
		switch {
		case gitlab.IsUnauthorized(err), gitlab.IsForbidden(err), gitlab.IsNotFound(err):
			log.Debug("Authorize proxy user error", logz.Error(err))
			return nil, &grpctool.ErrResp{
				StatusCode: http.StatusUnauthorized,
				Msg:        "Unauthorized",
			}
		default:
			msg := "Failed to authorize user session"
			p.api.HandleProcessingError(ctx, log, agentId, msg, err)
			return nil, &grpctool.ErrResp{
				StatusCode: http.StatusInternalServerError,
				Msg:        msg,
			}
		}

	}
	return auth, nil
}

func (p *kubernetesApiProxy) pipeStreams(log *zap.Logger, agentId int64, w http.ResponseWriter, r *http.Request,
	client rpc.KubernetesApi_MakeRequestClient, impConfig *rpc.ImpersonationConfig) {

	// urlPathPrefix is guaranteed to end with / by defaulting. That means / will be removed here.
	// Put it back by -1 on length.
	r.URL.Path = r.URL.Path[len(p.urlPathPrefix)-1:]

	// remove GitLab authorization headers (job token, session cookie etc)
	delete(r.Header, httpz.AuthorizationHeader)
	delete(r.Header, httpz.CookieHeader)
	delete(r.Header, httpz.GitlabAgentIdHeader)
	delete(r.Header, httpz.CsrfTokenHeader)

	r.Header[httpz.ViaHeader] = append(r.Header[httpz.ViaHeader], p.serverVia)

	http2grpc := grpctool.InboundHttpToOutboundGrpc{
		Log: log,
		HandleProcessingError: func(msg string, err error) {
			p.api.HandleProcessingError(r.Context(), log, agentId, msg, err)
		},
		WriteErrorResponse: p.writeErrorResponse(log, agentId),
		MergeHeaders:       p.mergeProxiedResponseHeaders,
	}
	var extra proto.Message // don't use a concrete type here or extra will be passed as a typed nil.
	if impConfig != nil {
		extra = &rpc.HeaderExtra{
			ImpConfig: impConfig,
		}
	}
	http2grpc.Pipe(client, w, r, extra)
}

func (p *kubernetesApiProxy) mergeProxiedResponseHeaders(outbound, inbound http.Header) {
	delete(inbound, httpz.ServerHeader) // remove the header we've added above. We use Via instead.
	// remove all potential CORS headers from the proxied response
	delete(outbound, httpz.AccessControlAllowOriginHeader)
	delete(outbound, httpz.AccessControlAllowHeadersHeader)
	delete(outbound, httpz.AccessControlAllowCredentialsHeader)
	delete(outbound, httpz.AccessControlAllowMethodsHeader)
	delete(outbound, httpz.AccessControlMaxAgeHeader)

	// set headers from proxied response without overwriting the ones already set (e.g. CORS headers)
	for k, vals := range outbound {
		if len(inbound[k]) == 0 {
			inbound[k] = vals
		}
	}
	// explicitly merge Vary header with the headers from proxies requests.
	// We always set the Vary header to `Origin` for CORS
	if v := append(inbound[httpz.VaryHeader], outbound[httpz.VaryHeader]...); len(v) > 0 {
		inbound[httpz.VaryHeader] = v
	}
	inbound[httpz.ViaHeader] = append(inbound[httpz.ViaHeader], p.serverVia)
}

func (p *kubernetesApiProxy) writeErrorResponse(log *zap.Logger, agentId int64) grpctool.WriteErrorResponse {
	return func(w http.ResponseWriter, r *http.Request, errResp *grpctool.ErrResp) {
		_, info, err := negotiation.NegotiateOutputMediaType(r, p.responseSerializer, negotiation.DefaultEndpointRestrictions)
		ctx := r.Context()
		if err != nil {
			msg := "Failed to negotiate output media type"
			log.Debug(msg, logz.Error(err))
			http.Error(w, formatStatusMessage(ctx, msg, err), http.StatusNotAcceptable)
			return
		}
		message := formatStatusMessage(ctx, errResp.Msg, errResp.Err)
		s := &metav1.Status{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Status",
				APIVersion: "v1",
			},
			Status:  metav1.StatusFailure,
			Message: message,
			Reason:  code2reason[errResp.StatusCode], // if mapping is not present, then "" means metav1.StatusReasonUnknown
			Code:    errResp.StatusCode,
		}
		buf := memz.Get32k() // use a temporary buffer to segregate I/O errors and encoding errors
		defer memz.Put32k(buf)
		buf = buf[:0] // don't care what's in the buf, start writing from the start
		b := bytes.NewBuffer(buf)
		err = info.Serializer.Encode(s, b) // encoding errors
		if err != nil {
			p.api.HandleProcessingError(ctx, log, agentId, "Failed to encode status response", err)
			http.Error(w, message, int(errResp.StatusCode))
			return
		}
		w.Header()[httpz.ContentTypeHeader] = []string{info.MediaType}
		w.WriteHeader(int(errResp.StatusCode))
		_, _ = w.Write(b.Bytes()) // I/O errors
	}
}

// err can be nil.
func formatStatusMessage(ctx context.Context, msg string, err error) string {
	var b strings.Builder
	b.WriteString("GitLab Agent Server: ")
	b.WriteString(msg)
	if err != nil {
		b.WriteString(": ")
		b.WriteString(err.Error())
	}
	traceId := trace.SpanContextFromContext(ctx).TraceID()
	if traceId.IsValid() {
		b.WriteString(". Trace ID: ")
		b.WriteString(traceId.String())
	}
	return b.String()
}

func findAllowedAgent(agentId int64, agentsForJob *gapi.AllowedAgentsForJob) *gapi.AllowedAgent {
	for _, aa := range agentsForJob.AllowedAgents {
		if aa.Id == agentId {
			return aa
		}
	}
	return nil
}

type ciJobTokenAuthn struct {
	jobToken string
}

type sessionCookieAuthn struct {
	encryptedPublicSessionId string
	csrfToken                string
}

func getAuthorizationInfoFromRequest(r *http.Request) (int64 /* agentId */, any, error) {
	if authzHeader := r.Header[httpz.AuthorizationHeader]; len(authzHeader) >= 1 {
		if len(authzHeader) > 1 {
			return 0, nil, fmt.Errorf("%s header: expecting a single header, got %d", httpz.AuthorizationHeader, len(authzHeader))
		}
		agentId, jobToken, err := getAgentIdAndJobTokenFromHeader(authzHeader[0])
		if err != nil {
			return 0, nil, err
		}
		return agentId, ciJobTokenAuthn{
			jobToken: jobToken,
		}, nil
	}
	if cookie, err := r.Cookie(gitLabKasCookieName); err == nil {
		agentId, encryptedPublicSessionId, csrfToken, err := getSessionCookieParams(cookie, r.Header)
		if err != nil {
			return 0, nil, err
		}
		return agentId, sessionCookieAuthn{
			encryptedPublicSessionId: encryptedPublicSessionId,
			csrfToken:                csrfToken,
		}, nil
	}
	return 0, nil, errors.New("no valid credentials provided")
}

func getSessionCookieParams(cookie *http.Cookie, headers http.Header) (int64, string, string, error) {
	if len(cookie.Value) == 0 {
		return 0, "", "", fmt.Errorf("%s cookie value must not be empty", gitLabKasCookieName)
	}
	// NOTE: GitLab Rails uses `rack` as the generic web server framework, which escapes the cookie values.
	// See https://github.com/rack/rack/blob/0b26518acc4c946ca96dfe3d9e68a05ca84439f7/lib/rack/utils.rb#L300
	encryptedPublicSessionId, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return 0, "", "", fmt.Errorf("%s invalid cookie value", gitLabKasCookieName)
	}
	agentIdHeader := headers[httpz.GitlabAgentIdHeader]
	if len(agentIdHeader) != 1 {
		return 0, "", "", fmt.Errorf("%s header must have exactly one value", httpz.GitlabAgentIdHeader)
	}
	agentIdStr := agentIdHeader[0]
	agentId, err := strconv.ParseInt(agentIdStr, 10, 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("%s header: invalid value: %q", httpz.GitlabAgentIdHeader, agentIdStr)
	}
	csrfTokenHeader := headers[httpz.CsrfTokenHeader]
	if len(csrfTokenHeader) != 1 {
		return 0, "", "", fmt.Errorf("%s header must have exactly one value", httpz.CsrfTokenHeader)
	}
	csrfToken := csrfTokenHeader[0]
	return agentId, encryptedPublicSessionId, csrfToken, nil
}

func getAgentIdAndJobTokenFromHeader(header string) (int64, string, error) {
	if !strings.HasPrefix(header, authorizationHeaderBearerPrefix) {
		// "missing" space in message - it's in the authorizationHeaderBearerPrefix constant already
		return 0, "", fmt.Errorf("%s header: expecting %stoken", httpz.AuthorizationHeader, authorizationHeaderBearerPrefix)
	}
	tokenValue := header[len(authorizationHeaderBearerPrefix):]
	tokenType, tokenContents, found := strings.Cut(tokenValue, tokenSeparator)
	if !found {
		return 0, "", fmt.Errorf("%s header: invalid value", httpz.AuthorizationHeader)
	}
	switch tokenType {
	case tokenTypeCi:
	default:
		return 0, "", fmt.Errorf("%s header: unknown token type", httpz.AuthorizationHeader)
	}
	agentIdAndToken := tokenContents
	agentIdStr, token, found := strings.Cut(agentIdAndToken, tokenSeparator)
	if !found {
		return 0, "", fmt.Errorf("%s header: invalid value", httpz.AuthorizationHeader)
	}
	agentId, err := strconv.ParseInt(agentIdStr, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("%s header: failed to parse: %w", httpz.AuthorizationHeader, err)
	}
	if token == "" {
		return 0, "", fmt.Errorf("%s header: empty token", httpz.AuthorizationHeader)
	}
	return agentId, token, nil
}

func constructJobImpersonationConfig(allowedForJob *gapi.AllowedAgentsForJob, aa *gapi.AllowedAgent) (*rpc.ImpersonationConfig, error) {
	as := aa.GetConfiguration().GetAccessAs().GetAs() // all these fields are optional, so handle nils.
	switch imp := as.(type) {
	case nil, *agentcfg.CiAccessAsCF_Agent: // nil means default value, which is Agent.
		return nil, nil
	case *agentcfg.CiAccessAsCF_Impersonate:
		i := imp.Impersonate
		return &rpc.ImpersonationConfig{
			Username: i.Username,
			Groups:   i.Groups,
			Uid:      i.Uid,
			Extra:    impImpersonationExtra(i.Extra),
		}, nil
	case *agentcfg.CiAccessAsCF_CiJob:
		return &rpc.ImpersonationConfig{
			Username: fmt.Sprintf("gitlab:ci_job:%d", allowedForJob.Job.Id),
			Groups:   impCiJobGroups(allowedForJob),
			Extra:    impCiJobExtra(allowedForJob, aa),
		}, nil
	default:
		// Normally this should never happen
		return nil, fmt.Errorf("unexpected job impersonation mode: %T", imp)
	}
}

func constructUserImpersonationConfig(auth *gapi.AuthorizeProxyUserResponse, accessType string) (*rpc.ImpersonationConfig, error) {
	switch imp := auth.GetAccessAs().AccessAs.(type) {
	case *gapi.AccessAsProxyAuthorization_Agent:
		return nil, nil
	case *gapi.AccessAsProxyAuthorization_User:
		return &rpc.ImpersonationConfig{
			Username: fmt.Sprintf("gitlab:user:%s", auth.User.Username),
			Groups:   impUserGroups(auth),
			Extra:    impUserExtra(auth, accessType),
		}, nil
	default:
		// Normally this should never happen
		return nil, fmt.Errorf("unexpected user impersonation mode: %T", imp)
	}
}

func impImpersonationExtra(in []*agentcfg.ExtraKeyValCF) []*rpc.ExtraKeyVal {
	out := make([]*rpc.ExtraKeyVal, 0, len(in))
	for _, kv := range in {
		out = append(out, &rpc.ExtraKeyVal{
			Key: kv.Key,
			Val: kv.Val,
		})
	}
	return out
}

func impCiJobGroups(allowedForJob *gapi.AllowedAgentsForJob) []string {
	// 1. gitlab:ci_job to identify all requests coming from CI jobs.
	groups := make([]string, 0, 3+len(allowedForJob.Project.Groups))
	groups = append(groups, "gitlab:ci_job")
	// 2. The list of ids of groups the project is in.
	for _, projectGroup := range allowedForJob.Project.Groups {
		groups = append(groups, fmt.Sprintf("gitlab:group:%d", projectGroup.Id))

		// 3. The tier of the environment this job belongs to, if set.
		if allowedForJob.Environment != nil {
			groups = append(groups, fmt.Sprintf("gitlab:group_env_tier:%d:%s", projectGroup.Id, allowedForJob.Environment.Tier))
		}
	}
	// 4. The project id.
	groups = append(groups, fmt.Sprintf("gitlab:project:%d", allowedForJob.Project.Id))
	// 5. The slug and tier of the environment this job belongs to, if set.
	if allowedForJob.Environment != nil {
		groups = append(groups,
			fmt.Sprintf("gitlab:project_env:%d:%s", allowedForJob.Project.Id, allowedForJob.Environment.Slug),
			fmt.Sprintf("gitlab:project_env_tier:%d:%s", allowedForJob.Project.Id, allowedForJob.Environment.Tier),
		)
	}
	return groups
}

func impCiJobExtra(allowedForJob *gapi.AllowedAgentsForJob, aa *gapi.AllowedAgent) []*rpc.ExtraKeyVal {
	extra := []*rpc.ExtraKeyVal{
		{
			Key: "agent.gitlab.com/id",
			Val: []string{strconv.FormatInt(aa.Id, 10)}, // agent id
		},
		{
			Key: "agent.gitlab.com/config_project_id",
			Val: []string{strconv.FormatInt(aa.ConfigProject.Id, 10)}, // agent's configuration project id
		},
		{
			Key: "agent.gitlab.com/project_id",
			Val: []string{strconv.FormatInt(allowedForJob.Project.Id, 10)}, // CI project id
		},
		{
			Key: "agent.gitlab.com/ci_pipeline_id",
			Val: []string{strconv.FormatInt(allowedForJob.Pipeline.Id, 10)}, // CI pipeline id
		},
		{
			Key: "agent.gitlab.com/ci_job_id",
			Val: []string{strconv.FormatInt(allowedForJob.Job.Id, 10)}, // CI job id
		},
		{
			Key: "agent.gitlab.com/username",
			Val: []string{allowedForJob.User.Username}, // username of the user the CI job is running as
		},
	}
	if allowedForJob.Environment != nil {
		extra = append(extra,
			&rpc.ExtraKeyVal{
				Key: "agent.gitlab.com/environment_slug",
				Val: []string{allowedForJob.Environment.Slug}, // slug of the environment, if set
			},
			&rpc.ExtraKeyVal{
				Key: "agent.gitlab.com/environment_tier",
				Val: []string{allowedForJob.Environment.Tier}, // tier of the environment, if set
			},
		)
	}
	return extra
}

func impUserGroups(auth *gapi.AuthorizeProxyUserResponse) []string {
	groups := []string{"gitlab:user"}
	for _, accessCF := range auth.AccessAs.GetUser().Projects {
		for _, role := range accessCF.Roles {
			groups = append(groups, fmt.Sprintf("gitlab:project_role:%d:%s", accessCF.Id, role))
		}
	}
	for _, accessCF := range auth.AccessAs.GetUser().Groups {
		for _, role := range accessCF.Roles {
			groups = append(groups, fmt.Sprintf("gitlab:group_role:%d:%s", accessCF.Id, role))
		}
	}
	return groups
}

func impUserExtra(auth *gapi.AuthorizeProxyUserResponse, accessType string) []*rpc.ExtraKeyVal {
	extra := []*rpc.ExtraKeyVal{
		{
			Key: "agent.gitlab.com/id",
			Val: []string{strconv.FormatInt(auth.Agent.Id, 10)},
		},
		{
			Key: "agent.gitlab.com/username",
			Val: []string{auth.User.Username},
		},
		{
			Key: "agent.gitlab.com/access_type",
			Val: []string{accessType},
		},
		{
			Key: "agent.gitlab.com/config_project_id",
			Val: []string{strconv.FormatInt(auth.Agent.ConfigProject.Id, 10)},
		},
	}
	return extra
}
