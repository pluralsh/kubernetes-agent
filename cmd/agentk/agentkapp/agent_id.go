package agentkapp

import (
	"context"
	"fmt"
)

// agentIdHolder holds agent id of this agentk.
type agentIdHolder struct {
	agentId    int64
	agentIdSet chan struct{}
}

func newAgentIdHolder() *agentIdHolder {
	return &agentIdHolder{
		agentIdSet: make(chan struct{}),
	}
}

// set is not safe for concurrent use. It's ok since we don't need that.
func (a *agentIdHolder) set(agentId int64) error {
	select {
	case <-a.agentIdSet: // already set
		if a.agentId != agentId {
			return fmt.Errorf("agentId is already set to a different value: old %d, new %d", a.agentId, agentId)
		}
	default: // not set
		a.agentId = agentId
		close(a.agentIdSet)
	}
	return nil
}

func (a *agentIdHolder) get(ctx context.Context) (int64, error) {
	select {
	case <-a.agentIdSet:
		return a.agentId, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
