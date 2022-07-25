package server

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"gitlab.com/gitlab-org/labkit/correlation"
	"gitlab.com/gitlab-org/labkit/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
)

const (
	defaultMaxRequestDuration = 15 * time.Second
	shutdownTimeout           = defaultMaxRequestDuration
	readTimeout               = 10 * time.Second
	writeTimeout              = defaultMaxRequestDuration
	idleTimeout               = 1 * time.Minute

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

type kubernetesApiProxy struct {
	log                       *zap.Logger
	api                       modserver.Api
	kubernetesApiClient       rpc.KubernetesApiClient
	gitLabClient              gitlab.ClientInterface
	allowedAgentsCache        *cache.CacheWithErr
	requestCount              usage_metrics.Counter
	metricsHttpHandlerFactory metrics.HandlerFactory
	responseSerializer        runtime.NegotiatedSerializer
	serverName                string
	serverVia                 string
	// urlPathPrefix is guaranteed to end with / by defaulting.
	urlPathPrefix string
}

func (p *kubernetesApiProxy) Run(ctx context.Context, listener net.Listener) error {
	var handler http.Handler
	handler = http.HandlerFunc(p.proxy)
	handler = correlation.InjectCorrelationID(handler, correlation.WithSetResponseHeader())
	handler = p.metricsHttpHandlerFactory(handler)
	srv := &http.Server{ // nolint: gosec
		Handler:      handler,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}
	return httpz.RunServer(ctx, srv, listener, shutdownTimeout)
}

func (p *kubernetesApiProxy) proxy(w http.ResponseWriter, r *http.Request) {
	log, agentId, eResp := p.proxyInternal(w, r)
	if eResp != nil {
		p.writeErrorResponse(log, agentId)(w, r, eResp)
	}
}

func (p *kubernetesApiProxy) proxyInternal(w http.ResponseWriter, r *http.Request) (*zap.Logger, int64 /* agentId */, *grpctool.ErrResp) {
	ctx := r.Context()
	correlationId := correlation.ExtractFromContext(ctx)
	log := p.log.With(logz.CorrelationId(correlationId))
	w.Header()[httpz.ServerHeader] = []string{p.serverName} // It will be removed just before responding with actual headers from upstream

	agentId, jobToken, err := getAgentIdAndJobTokenFromRequest(r)
	if err != nil {
		msg := "Unauthorized"
		log.Debug(msg, logz.Error(err))
		return log, modshared.NoAgentId, &grpctool.ErrResp{
			StatusCode: http.StatusUnauthorized,
			Msg:        msg,
			Err:        err,
		}
	}
	log = log.With(logz.AgentId(agentId))

	allowedForJob, eResp := p.getAllowedAgentsForJob(ctx, log, agentId, jobToken)
	if eResp != nil {
		return log, agentId, eResp
	}

	aa := findAllowedAgent(agentId, allowedForJob)
	if aa == nil {
		msg := "Forbidden: agentId is not allowed"
		log.Debug(msg)
		return log, agentId, &grpctool.ErrResp{
			StatusCode: http.StatusForbidden,
			Msg:        msg,
		}
	}

	if !strings.HasPrefix(r.URL.Path, p.urlPathPrefix) {
		msg := "Bad request: URL does not start with expected prefix"
		log.Debug(msg, logz.UrlPath(r.URL.Path), logz.UrlPathPrefix(p.urlPathPrefix))
		return log, agentId, &grpctool.ErrResp{
			StatusCode: http.StatusBadRequest,
			Msg:        msg,
		}
	}

	p.requestCount.Inc() // Count only authenticated and authorized requests

	impConfig, err := constructImpersonationConfig(allowedForJob, aa)
	if err != nil {
		msg := "Failed to construct impersonation config"
		p.api.HandleProcessingError(ctx, log, agentId, msg, err)
		return log, agentId, &grpctool.ErrResp{
			StatusCode: http.StatusInternalServerError,
			Msg:        msg,
			Err:        err,
		}
	}

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

	p.pipeStreams(log, agentId, w, r, mkClient, impConfig)
	return log, agentId, nil
}

func (p *kubernetesApiProxy) getAllowedAgentsForJob(ctx context.Context, log *zap.Logger, agentId int64, jobToken string) (*gapi.AllowedAgentsForJob, *grpctool.ErrResp) {
	allowedForJob, err := p.allowedAgentsCache.GetItem(ctx, jobToken, func() (interface{}, error) {
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
	return allowedForJob.(*gapi.AllowedAgentsForJob), nil
}

func (p *kubernetesApiProxy) pipeStreams(log *zap.Logger, agentId int64, w http.ResponseWriter, r *http.Request,
	client rpc.KubernetesApi_MakeRequestClient, impConfig *rpc.ImpersonationConfig) {

	// urlPathPrefix is guaranteed to end with / by defaulting. That means / will be removed here.
	// Put it back by -1 on length.
	r.URL.Path = r.URL.Path[len(p.urlPathPrefix)-1:]
	delete(r.Header, httpz.AuthorizationHeader) // Remove Authorization header - we got the CI job token in it
	r.Header[httpz.ViaHeader] = append(r.Header[httpz.ViaHeader], p.serverVia)

	http2grpc := grpctool.InboundHttpToOutboundGrpc{
		Log: log,
		HandleProcessingError: func(msg string, err error) {
			p.api.HandleProcessingError(r.Context(), log, agentId, msg, err)
		},
		WriteErrorResponse: p.writeErrorResponse(log, agentId),
		MergeHeaders: func(outboundResponse, inboundResponse http.Header) {
			delete(inboundResponse, httpz.ServerHeader) // remove the header we've added above. We use Via instead.
			for k, vals := range outboundResponse {
				inboundResponse[k] = vals
			}
			inboundResponse[httpz.ViaHeader] = append(inboundResponse[httpz.ViaHeader], p.serverVia)
		},
	}
	http2grpc.Pipe(client, w, r, &rpc.HeaderExtra{
		ImpConfig: impConfig,
	})
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
		_, _ = fmt.Fprintf(&b, ": %v", err)
	}
	correlationId := correlation.ExtractFromContext(ctx) // can be empty
	if correlationId != "" {
		b.WriteString(". Correlation ID: ")
		b.WriteString(correlationId)
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

func getAgentIdAndJobTokenFromRequest(r *http.Request) (int64, string, error) {
	h := r.Header[httpz.AuthorizationHeader]
	if len(h) == 0 {
		return 0, "", fmt.Errorf("%s header: expecting token", httpz.AuthorizationHeader)
	}
	if len(h) > 1 {
		return 0, "", fmt.Errorf("%s header: expecting a single header, got %d", httpz.AuthorizationHeader, len(h))
	}
	return getAgentIdAndJobTokenFromHeader(h[0])
}

func getAgentIdAndJobTokenFromHeader(header string) (int64, string, error) {
	if !strings.HasPrefix(header, authorizationHeaderBearerPrefix) {
		// "missing" space in message - it's in the authorizationHeaderBearerPrefix constant already
		return 0, "", fmt.Errorf("%s header: expecting %stoken", httpz.AuthorizationHeader, authorizationHeaderBearerPrefix)
	}
	tokenValue := header[len(authorizationHeaderBearerPrefix):]
	tokenValueParts := strings.SplitN(tokenValue, tokenSeparator, 2)
	if len(tokenValueParts) != 2 {
		return 0, "", fmt.Errorf("%s header: invalid value", httpz.AuthorizationHeader)
	}
	switch tokenValueParts[0] {
	case tokenTypeCi:
	default:
		return 0, "", fmt.Errorf("%s header: unknown token type", httpz.AuthorizationHeader)
	}
	agentIdAndToken := tokenValueParts[1]
	agentIdAndTokenParts := strings.SplitN(agentIdAndToken, tokenSeparator, 2)
	if len(agentIdAndTokenParts) != 2 {
		return 0, "", fmt.Errorf("%s header: invalid value", httpz.AuthorizationHeader)
	}
	agentId, err := strconv.ParseInt(agentIdAndTokenParts[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("%s header: failed to parse: %w", httpz.AuthorizationHeader, err)
	}
	token := agentIdAndTokenParts[1]
	if token == "" {
		return 0, "", fmt.Errorf("%s header: empty token", httpz.AuthorizationHeader)
	}
	return agentId, token, nil
}

func constructImpersonationConfig(allowedForJob *gapi.AllowedAgentsForJob, aa *gapi.AllowedAgent) (*rpc.ImpersonationConfig, error) {
	as := aa.GetConfiguration().GetAccessAs().GetAs() // all these fields are optional, so handle nils.
	if as == nil {
		as = &agentcfg.CiAccessAsCF_Agent{} // default value
	}
	switch imp := as.(type) {
	case *agentcfg.CiAccessAsCF_Agent:
		return &rpc.ImpersonationConfig{}, nil
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
		return nil, fmt.Errorf("unexpected impersonation mode: %T", imp)
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
	}
	// 3. The project id.
	groups = append(groups, fmt.Sprintf("gitlab:project:%d", allowedForJob.Project.Id))
	// 4. The slug of the environment this job belongs to, if set.
	if allowedForJob.Environment != nil {
		groups = append(groups, fmt.Sprintf("gitlab:project_env:%d:%s", allowedForJob.Project.Id, allowedForJob.Environment.Slug))
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
		extra = append(extra, &rpc.ExtraKeyVal{
			Key: "agent.gitlab.com/environment_slug",
			Val: []string{allowedForJob.Environment.Slug}, // slug of the environment, if set
		})
	}
	return extra
}
