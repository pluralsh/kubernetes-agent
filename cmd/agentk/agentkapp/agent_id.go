package agentkapp

import (
	"context"
	"fmt"
)

// AgentIdHolder holds agent id of this agentk.
type AgentIdHolder struct {
	agentId    int64
	agentIdSet chan struct{}
}

func NewAgentIdHolder() *AgentIdHolder {
	return &AgentIdHolder{
		agentIdSet: make(chan struct{}),
	}
}

// set is not safe for concurrent use. It's ok since we don't need that.
func (a *AgentIdHolder) set(agentId int64) error {
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

func (a *AgentIdHolder) get(ctx context.Context) (int64, error) {
	select {
	case <-a.agentIdSet:
		return a.agentId, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (a *AgentIdHolder) tryGet() (int64, bool) {
	select {
	case <-a.agentIdSet:
		return a.agentId, true
	default:
		return 0, false
	}
}
