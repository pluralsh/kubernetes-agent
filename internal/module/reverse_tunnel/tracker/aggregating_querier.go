package tracker

import (
	"context"
	"sync"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"go.uber.org/zap"
)

// PollKasUrlsByAgentIdCallback is called periodically for each found kas URL for a particular agent id.
// newCycle is set to true on the first item of a new polling cycle i.e. after poller has slept for the polling interval.
type PollKasUrlsByAgentIdCallback func(newCycle bool, kasUrl string) bool

type PollingQuerier interface {
	PollKasUrlsByAgentId(ctx context.Context, agentId int64, cb PollKasUrlsByAgentIdCallback)
}

type pollConsumer struct {
	ctxDone   <-chan struct{}
	pollItems chan<- pollItem
}

type pollingContext struct {
	mu        sync.Mutex
	consumers map[*pollConsumer]struct{}
	cancel    context.CancelFunc
}

type pollItem struct {
	kasUrl   string
	newCycle bool
}

func (c *pollingContext) copyConsumersInto(consumers []pollConsumer) []pollConsumer {
	consumers = consumers[:0]
	c.mu.Lock()
	defer c.mu.Unlock()
	for h := range c.consumers {
		consumers = append(consumers, *h)
	}
	return consumers
}

// AggregatingQuerier groups polling requests.
type AggregatingQuerier struct {
	log        *zap.Logger
	delegate   Querier
	api        modshared.Api
	pollConfig retry.PollConfigFactory

	mu        sync.Mutex
	listeners map[int64]*pollingContext
}

func NewAggregatingQuerier(log *zap.Logger, delegate Querier, api modshared.Api, pollConfig retry.PollConfigFactory) *AggregatingQuerier {
	return &AggregatingQuerier{
		log:        log,
		delegate:   delegate,
		api:        api,
		pollConfig: pollConfig,
		listeners:  make(map[int64]*pollingContext),
	}
}

func (q *AggregatingQuerier) PollKasUrlsByAgentId(ctx context.Context, agentId int64, cb PollKasUrlsByAgentIdCallback) {
	pollItems := make(chan pollItem)
	ctxDone := ctx.Done()
	h := &pollConsumer{
		ctxDone:   ctxDone,
		pollItems: pollItems,
	}
	q.maybeStartPolling(agentId, h) // nolint: contextcheck
	defer q.maybeStopPolling(agentId, h)
	for {
		select {
		case <-ctxDone:
			return
		case item := <-pollItems:
			done := cb(item.newCycle, item.kasUrl)
			if done {
				return
			}
		}
	}
}

func (q *AggregatingQuerier) maybeStartPolling(agentId int64, h *pollConsumer) {
	q.mu.Lock()
	defer q.mu.Unlock()
	pc := q.listeners[agentId]
	if pc != nil { // already polling
		pc.mu.Lock()
		pc.consumers[h] = struct{}{} // register for notifications
		pc.mu.Unlock()
	} else { // not polling, start.
		ctx, cancel := context.WithCancel(context.Background())
		pc = &pollingContext{
			consumers: map[*pollConsumer]struct{}{
				h: {},
			},
			cancel: cancel,
		}
		q.listeners[agentId] = pc
		go q.poll(ctx, agentId, pc)
	}
}

func (q *AggregatingQuerier) maybeStopPolling(agentId int64, h *pollConsumer) {
	q.mu.Lock()
	defer q.mu.Unlock()

	pc := q.listeners[agentId]
	pc.mu.Lock()
	defer pc.mu.Unlock()

	delete(pc.consumers, h)
	if len(pc.consumers) == 0 {
		pc.cancel() // stop polling
		delete(q.listeners, agentId)
	}
}

func (q *AggregatingQuerier) poll(ctx context.Context, agentId int64, pc *pollingContext) {
	// err can only be retry.ErrWaitTimeout
	var consumers []pollConsumer // reuse slice between polls
	_ = retry.PollWithBackoff(ctx, q.pollConfig(), func(ctx context.Context) (error, retry.AttemptResult) {
		newCycle := true
		err := q.delegate.KasUrlsByAgentId(ctx, agentId, func(kasUrl string) (bool, error) {
			consumers = pc.copyConsumersInto(consumers)
			for _, h := range consumers {
				select {
				case <-h.ctxDone:
					// This PollKasUrlsByAgentId() invocation is no longer interested in being called. Ignore it.
				case h.pollItems <- pollItem{kasUrl: kasUrl, newCycle: newCycle}:
					// Data sent.
				}
			}
			newCycle = false
			return false, nil
		})
		if err != nil {
			q.api.HandleProcessingError(ctx, q.log, agentId, "KasUrlsByAgentId() failed", err)
			// fallthrough
		}
		return nil, retry.Continue
	})
}
