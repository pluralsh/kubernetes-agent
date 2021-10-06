package server

import (
	"context"
	"fmt"
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
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
	"gitlab.com/gitlab-org/labkit/correlation"
	"gitlab.com/gitlab-org/labkit/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	defaultMaxRequestDuration = 15 * time.Second
	shutdownTimeout           = defaultMaxRequestDuration
	readTimeout               = 1 * time.Second
	writeTimeout              = defaultMaxRequestDuration
	idleTimeout               = 1 * time.Minute

	authorizationHeader             = "Authorization"
	serverHeader                    = "Server"
	viaHeader                       = "Via"
	authorizationHeaderBearerPrefix = "Bearer " // must end with a space
	tokenSeparator                  = ":"
	tokenTypeCi                     = "ci"
)

type kubernetesApiProxy struct {
	log                       *zap.Logger
	api                       modserver.Api
	kubernetesApiClient       rpc.KubernetesApiClient
	gitLabClient              gitlab.ClientInterface
	allowedAgentsCache        *cache.CacheWithErr
	requestCount              usage_metrics.Counter
	metricsHttpHandlerFactory metrics.HandlerFactory
	serverName                string
	// urlPathPrefix is guaranteed to end with / by defaulting.
	urlPathPrefix string
}

func (p *kubernetesApiProxy) Run(ctx context.Context, listener net.Listener) error {
	var handler http.Handler
	handler = http.HandlerFunc(p.proxy)
	handler = correlation.InjectCorrelationID(handler, correlation.WithSetResponseHeader())
	handler = p.metricsHttpHandlerFactory(handler)
	srv := &http.Server{
		Handler:      handler,
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

	md := metadata.Pairs(modserver.RoutingAgentIdMetadataKey, strconv.FormatInt(agentId, 10))
	mkClient, err := p.kubernetesApiClient.MakeRequest(metadata.NewOutgoingContext(ctx, md))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.api.HandleProcessingError(ctx, log, agentId, "Proxy failed to make outbound request", err)
		return
	}

	p.pipeStreams(log, agentId, w, r, mkClient, impConfig)
}

func (p *kubernetesApiProxy) getAllowedAgentsForJob(ctx context.Context, jobToken string) (*gapi.AllowedAgentsForJob, error) {
	allowedForJob, err := p.allowedAgentsCache.GetItem(ctx, jobToken, func() (interface{}, error) {
		return gapi.GetAllowedAgentsForJob(ctx, p.gitLabClient, jobToken)
	})
	if err != nil {
		return nil, err
	}
	return allowedForJob.(*gapi.AllowedAgentsForJob), nil
}

func (p *kubernetesApiProxy) pipeStreams(log *zap.Logger, agentId int64, w http.ResponseWriter, r *http.Request,
	client rpc.KubernetesApi_MakeRequestClient, impConfig *rpc.ImpersonationConfig) {

	// urlPathPrefix is guaranteed to end with / by defaulting. That means / will be removed here.
	// Put it back by -1 on length.
	r.URL.Path = r.URL.Path[len(p.urlPathPrefix)-1:]
	r.Header.Del(authorizationHeader) // Remove Authorization header - we got the CI job token in it
	serverProto := "gRPC/1.0 " + p.serverName
	r.Header.Add(viaHeader, serverProto)

	http2grpc := grpctool.InboundHttpToOutboundGrpc{
		Log: log,
		HandleProcessingError: func(msg string, err error) {
			p.api.HandleProcessingError(r.Context(), log, agentId, msg, err)
		},
		MergeHeaders: func(fromOutbound, toInbound http.Header) {
			toInbound.Del(serverHeader) // remove the header we've added above. We use Via instead.
			for k, vals := range fromOutbound {
				toInbound[k] = vals
			}
			toInbound.Add(viaHeader, serverProto)
		},
	}
	http2grpc.Pipe(client, w, r, &rpc.HeaderExtra{
		ImpConfig: impConfig,
	})
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
