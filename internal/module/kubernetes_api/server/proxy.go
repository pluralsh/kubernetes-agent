package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
	"gitlab.com/gitlab-org/labkit/correlation"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	defaultMaxRequestDuration = 15 * time.Second
	shutdownTimeout           = defaultMaxRequestDuration
	readTimeout               = 1 * time.Second
	writeTimeout              = defaultMaxRequestDuration
	idleTimeout               = 1 * time.Minute
	maxDataChunkSize          = 32 * 1024

	authorizationHeader             = "Authorization"
	serverHeader                    = "Server"
	viaHeader                       = "Via"
	hostHeader                      = "Host"
	authorizationHeaderBearerPrefix = "Bearer " // must end with a space
	tokenSeparator                  = ":"
	tokenTypeCi                     = "ci"

	headerFieldNumber  protoreflect.FieldNumber = 1
	dataFieldNumber    protoreflect.FieldNumber = 2
	trailerFieldNumber protoreflect.FieldNumber = 3
)

var (
	// See https://httpwg.org/http-core/draft-ietf-httpbis-semantics-latest.html#field.connection
	// See https://tools.ietf.org/html/rfc2616#section-13.5.1
	// See https://github.com/golang/go/blob/81ea89adf38b90c3c3a8c4eed9e6c093a8634d59/src/net/http/httputil/reverseproxy.go#L169-L184
	hopHeaders = []string{
		"Connection",
		"Proxy-Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",      // canonicalized version of "TE"
		"Trailer", // not Trailers as per rfc2616; See errata https://www.rfc-editor.org/errata_search.php?eid=4522
		"Transfer-Encoding",
		"Upgrade",
	}
)

type kubernetesApiProxy struct {
	log                 *zap.Logger
	api                 modserver.Api
	kubernetesApiClient rpc.KubernetesApiClient
	gitLabClient        gitlab.ClientInterface
	streamVisitor       *grpctool.StreamVisitor
	cache               *cache.CacheWithErr
	requestCount        usage_metrics.Counter
	serverName          string
	// urlPathPrefix is guaranteed to end with / by defaulting.
	urlPathPrefix string
}

func (p *kubernetesApiProxy) Run(ctx context.Context, listener net.Listener) error {
	srv := &http.Server{
		Handler: correlation.InjectCorrelationID(
			http.HandlerFunc(p.proxy),
			correlation.WithSetResponseHeader(),
		),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}
	return httpz.RunServer(ctx, srv, listener, shutdownTimeout)
}

func (p *kubernetesApiProxy) proxy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationId := correlation.ExtractFromContext(ctx)
	log := p.log.With(logz.CorrelationId(correlationId))
	w.Header().Set(serverHeader, p.serverName) // It will be removed just before responding with actual headers from upstream

	agentId, jobToken, err := getAgentIdAndJobTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Debug("Unauthorized: header", logz.Error(err))
		return
	}
	log = log.With(logz.AgentId(agentId))

	allowedForJob, err := p.getAllowedAgentsForJob(ctx, jobToken)
	if err != nil {
		switch {
		case gitlab.IsUnauthorized(err):
			w.WriteHeader(http.StatusUnauthorized)
			log.Debug("Unauthorized: CI job token")
		case gitlab.IsForbidden(err):
			w.WriteHeader(http.StatusForbidden)
			log.Debug("Forbidden: CI job token")
		case gitlab.IsNotFound(err):
			w.WriteHeader(http.StatusNotFound)
			log.Debug("Not found: agents for CI job token")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			p.api.HandleProcessingError(ctx, log, agentId, "Failed to get allowed agents for CI job token", err)
		}
		return
	}

	aa := findAllowedAgent(agentId, allowedForJob)
	if aa == nil {
		w.WriteHeader(http.StatusForbidden)
		log.Debug("Forbidden: agentId is not allowed")
		return
	}

	if !strings.HasPrefix(r.URL.Path, p.urlPathPrefix) {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug("Bad request: URL does not start with expected prefix", logz.UrlPath(r.URL.Path), logz.UrlPathPrefix(p.urlPathPrefix))
		return
	}

	p.requestCount.Inc() // Count only authenticated and authorized requests

	impConfig, err := constructImpersonationConfig(aa)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.api.HandleProcessingError(ctx, log, agentId, "Failed to construct impersonation config", err)
		return
	}

	// urlPathPrefix is guaranteed to end with / by defaulting. That means / will be removed here.
	// Put it back by -1 on length.
	r.URL.Path = r.URL.Path[len(p.urlPathPrefix)-1:]
	r.Header.Add(viaHeader, "gRPC/1.0 "+p.serverName)

	headerWritten, errF := p.pipeStreams(ctx, log, agentId, w, r, impConfig)
	if errF != nil {
		if headerWritten {
			// HTTP status has been written already as part of the normal response flow.
			// But then something went wrong and an error happened. To let the client know that something isn't right
			// we have only one thing we can do - abruptly close the connection. To do that we panic with a special
			// error value that the "http" package provides. See its description.
			// If we try to write the status again here, http package would log a warning, which is not nice.
			panic(http.ErrAbortHandler)
		} else {
			errF(w)
		}
	}
}

func (p *kubernetesApiProxy) getAllowedAgentsForJob(ctx context.Context, jobToken string) (*gapi.AllowedAgentsForJob, error) {
	allowedForJob, err := p.cache.GetItem(ctx, jobToken, func() (interface{}, error) {
		return gapi.GetAllowedAgentsForJob(ctx, p.gitLabClient, jobToken)
	})
	if err != nil {
		return nil, err
	}
	return allowedForJob.(*gapi.AllowedAgentsForJob), nil
}

func (p *kubernetesApiProxy) pipeStreams(ctx context.Context, log *zap.Logger, agentId int64, w http.ResponseWriter, r *http.Request, impConfig *rpc.ImpersonationConfig) (bool, errFunc) {
	g, ctx := errgroup.WithContext(ctx)
	md := metadata.Pairs(modserver.RoutingAgentIdMetadataKey, strconv.FormatInt(agentId, 10))
	mkClient, err := p.kubernetesApiClient.MakeRequest(metadata.NewOutgoingContext(ctx, md)) // must use context from errgroup
	if err != nil {
		return false, p.handleProcessingError(ctx, log, agentId, "Proxy failed to make outbound request", err)
	}
	// Pipe client -> remote
	g.Go(func() error {
		errFuncRet := p.pipeClientToRemote(ctx, log, agentId, mkClient, r, impConfig)
		if errFuncRet != nil {
			return errFuncRet
		}
		return nil
	})
	// Pipe remote -> client
	headerWritten := false
	g.Go(func() error {
		var errFuncRet errFunc
		headerWritten, errFuncRet = p.pipeRemoteToClient(ctx, log, agentId, mkClient, w)
		if errFuncRet != nil {
			return errFuncRet
		}
		return nil
	})
	err = g.Wait() // don't inline as headerWritten must be read after Wait() returned
	if err != nil {
		return headerWritten, err.(errFunc) // nolint: errorlint
	}
	return false, nil
}

func (p *kubernetesApiProxy) pipeRemoteToClient(ctx context.Context, log *zap.Logger, agentId int64, mkClient rpc.KubernetesApi_MakeRequestClient, w http.ResponseWriter) (bool, errFunc) {
	writeFailed := false
	headerWritten := false
	err := p.streamVisitor.Visit(mkClient,
		grpctool.WithCallback(headerFieldNumber, func(header *grpctool.HttpResponse_Header) error {
			httpH := header.Response.HttpHeader()
			httpz.RemoveConnectionHeaders(httpH)
			h := w.Header()
			h.Del(serverHeader) // remove the header we've added above. We use Via instead.
			for k, vals := range httpH {
				h[k] = vals
			}
			h.Add(viaHeader, "gRPC/1.0 "+p.serverName)
			w.WriteHeader(int(header.Response.StatusCode))
			headerWritten = true
			return nil
		}),
		grpctool.WithCallback(dataFieldNumber, func(data *grpctool.HttpResponse_Data) error {
			_, err := w.Write(data.Data)
			if err != nil {
				writeFailed = true
			}
			return err
		}),
		grpctool.WithCallback(trailerFieldNumber, func(trailer *grpctool.HttpResponse_Trailer) error {
			return nil
		}),
	)
	if err != nil {
		if writeFailed {
			// there is likely a connection problem so the client will likely not receive this
			return headerWritten, p.handleProcessingError(ctx, log, agentId, "Proxy failed to write response to client", err)
		}
		return headerWritten, p.handleProcessingError(ctx, log, agentId, "Proxy failed to read response from agent", err)
	}
	return headerWritten, nil
}

func (p *kubernetesApiProxy) pipeClientToRemote(ctx context.Context, log *zap.Logger, agentId int64,
	mkClient rpc.KubernetesApi_MakeRequestClient, r *http.Request, impConfig *rpc.ImpersonationConfig) errFunc {
	extra, err := anypb.New(&rpc.HeaderExtra{
		ImpConfig: impConfig,
	})
	if err != nil {
		return p.handleProcessingError(ctx, log, agentId, "Proxy failed to marshal HttpRequestExtra proto", err)
	}
	err = mkClient.Send(&grpctool.HttpRequest{
		Message: &grpctool.HttpRequest_Header_{
			Header: &grpctool.HttpRequest_Header{
				Request: &prototool.HttpRequest{
					Method:  r.Method,
					Header:  headerFromHttpRequestHeader(r.Header),
					UrlPath: r.URL.Path,
					Query:   prototool.UrlValuesToValuesMap(r.URL.Query()),
				},
				Extra: extra,
			},
		},
	})
	if err != nil {
		return p.handleSendError(log, "Proxy failed to send request header to agent", err)
	}

	buffer := make([]byte, maxDataChunkSize)
	for {
		var n int
		n, err = r.Body.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			// There is likely a connection problem so the client will likely not receive this
			return p.handleProcessingError(ctx, log, agentId, "Proxy failed to read request body from client", err)
		}
		if n > 0 { // handle n=0, err=io.EOF case
			sendErr := mkClient.Send(&grpctool.HttpRequest{
				Message: &grpctool.HttpRequest_Data_{
					Data: &grpctool.HttpRequest_Data{
						Data: buffer[:n],
					},
				},
			})
			if sendErr != nil {
				return p.handleSendError(log, "Proxy failed to send request body to agent", sendErr)
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	err = mkClient.Send(&grpctool.HttpRequest{
		Message: &grpctool.HttpRequest_Trailer_{
			Trailer: &grpctool.HttpRequest_Trailer{},
		},
	})
	if err != nil {
		return p.handleSendError(log, "Proxy failed to send trailers to agent", err)
	}
	err = mkClient.CloseSend()
	if err != nil {
		return p.handleSendError(log, "Proxy failed to send close frame to agent", err)
	}
	return nil
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
	h := r.Header.Values(authorizationHeader)
	if len(h) == 0 {
		return 0, "", fmt.Errorf("%s header: expecting token", authorizationHeader)
	}
	if len(h) > 1 {
		return 0, "", fmt.Errorf("%s header: expecting a single header, got %d", authorizationHeader, len(h))
	}
	return getAgentIdAndJobTokenFromHeader(h[0])
}

func getAgentIdAndJobTokenFromHeader(header string) (int64, string, error) {
	if !strings.HasPrefix(header, authorizationHeaderBearerPrefix) {
		// "missing" space in message - it's in the authorizationHeaderBearerPrefix constant already
		return 0, "", fmt.Errorf("%s header: expecting %stoken", authorizationHeader, authorizationHeaderBearerPrefix)
	}
	tokenValue := header[len(authorizationHeaderBearerPrefix):]
	tokenValueParts := strings.SplitN(tokenValue, tokenSeparator, 2)
	if len(tokenValueParts) != 2 {
		return 0, "", fmt.Errorf("%s header: invalid value", authorizationHeader)
	}
	switch tokenValueParts[0] {
	case tokenTypeCi:
	default:
		return 0, "", fmt.Errorf("%s header: unknown token type", authorizationHeader)
	}
	agentIdAndToken := tokenValueParts[1]
	agentIdAndTokenParts := strings.SplitN(agentIdAndToken, tokenSeparator, 2)
	if len(agentIdAndTokenParts) != 2 {
		return 0, "", fmt.Errorf("%s header: invalid value", authorizationHeader)
	}
	agentId, err := strconv.ParseInt(agentIdAndTokenParts[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("%s header: failed to parse: %w", authorizationHeader, err)
	}
	token := agentIdAndTokenParts[1]
	if token == "" {
		return 0, "", fmt.Errorf("%s header: empty token", authorizationHeader)
	}
	return agentId, token, nil
}

func headerFromHttpRequestHeader(header http.Header) map[string]*prototool.Values {
	header = header.Clone()
	header.Del(hostHeader)          // Use the destination host name
	header.Del(authorizationHeader) // Remove Authorization header - we got the CI job token in it

	// Remove hop-by-hop headers
	// 1. Remove headers listed in the Connection header
	httpz.RemoveConnectionHeaders(header)
	// 2. Remove well-known headers
	for _, name := range hopHeaders {
		header.Del(name)
	}

	return prototool.HttpHeaderToValuesMap(header)
}

func (p *kubernetesApiProxy) handleSendError(log *zap.Logger, msg string, err error) errFunc {
	if !grpctool.RequestCanceled(err) {
		log.Debug(msg, logz.Error(err))
	}
	return writeError(msg, err)
}

func (p *kubernetesApiProxy) handleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) errFunc {
	p.api.HandleProcessingError(ctx, log, agentId, msg, err)
	return writeError(msg, err)
}

func constructImpersonationConfig(aa *gapi.AllowedAgent) (*rpc.ImpersonationConfig, error) {
	as := aa.GetConfiguration().GetAccessAs().GetAs() // all these fields are optional, so handle nils.
	if as == nil {
		as = &agentcfg.CiAccessAsCF_Agent{} // default value
	}
	switch imp := as.(type) {
	case *agentcfg.CiAccessAsCF_Impersonate:
		i := imp.Impersonate
		return &rpc.ImpersonationConfig{
			Username: i.Username,
			Groups:   i.Groups,
			Uid:      i.Uid,
			Extra:    convertExtra(i.Extra),
		}, nil
	case *agentcfg.CiAccessAsCF_Agent:
		return &rpc.ImpersonationConfig{}, nil
	default:
		// Normally this should never happen
		return nil, fmt.Errorf("unexpected impersonation mode: %T", imp)
	}
}

func convertExtra(in []*agentcfg.ExtraKeyValCF) []*rpc.ExtraKeyVal {
	out := make([]*rpc.ExtraKeyVal, 0, len(in))
	for _, kv := range in {
		out = append(out, &rpc.ExtraKeyVal{
			Key: kv.Key,
			Val: kv.Val,
		})
	}
	return out
}

func writeError(msg string, err error) errFunc {
	return func(w http.ResponseWriter) {
		// See https://tools.ietf.org/html/rfc7231#section-6.6.3
		http.Error(w, fmt.Sprintf("%s: %v", msg, err), http.StatusBadGateway)
	}
}

var (
	_ error = errFunc(nil)
)

// errFunc enhances type safety.
type errFunc func(http.ResponseWriter)

func (e errFunc) Error() string {
	return "errorFunc"
}
