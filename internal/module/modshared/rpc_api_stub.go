package modshared

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/retry"
	"go.uber.org/zap"
)

type RpcApiStub struct {
	StreamCtx context.Context
}

func (a *RpcApiStub) PollWithBackoff(cfg retry.PollConfig, f retry.PollWithBackoffFunc) error {
	// this context must only be used here, not inside of f() - connection should be closed only when idle.
	ageCtx := grpctool.MaxConnectionAgeContextFromStreamContext(a.StreamCtx)
	err := retry.PollWithBackoff(ageCtx, cfg, f)
	if errors.Is(err, retry.ErrWaitTimeout) {
		return nil // all good, ctx is done
	}
	return err
}

func (a *RpcApiStub) Log() *zap.Logger {
	return grpctool.LoggerFromContext(a.StreamCtx)
}
