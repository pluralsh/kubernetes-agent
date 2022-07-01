package agent

import (
	"fmt"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"k8s.io/utils/clock"
)

type processResult byte

const (
	// success means there was no error and output should be consumed.
	success processResult = iota
	// backoff means there was a retriable error, so the caller should try later.
	backoff
	// done means there is no output to be consumed. There may or may not have been an error.
	done
)

// TODO remove aliases and use generics.

type inputT = rpc.ObjectsToSynchronizeData
type outputT = syncJob

type processFunc func(input inputT) (outputT, processResult)

// retryPipeline takes a channel with input, a processor function, and a channel for output.
// It reads values from input, processes them, and, if/when successful, sends the result into the output channel.
// If processing fails, it is retried with backoff.
// Results are sent to the output eventually. The old result, if it hasn't been sent already, is discarded when a new
// one becomes available. I.e. level-based rather than edge-based behavior.
type retryPipeline struct {
	inputCh      <-chan inputT
	outputCh     chan<- outputT
	retryBackoff retry.BackoffManager
	process      processFunc
}

func (p *retryPipeline) run() {
	var (
		input        inputT
		output       outputT
		outputCh     chan<- outputT
		attemptCh    <-chan time.Time
		attemptTimer clock.Timer
	)
	stopAttemptTimer := func() {
		if attemptTimer != nil {
			if !attemptTimer.Stop() {
				<-attemptCh
			}
		}
	}
	defer stopAttemptTimer()
	for {
		var ok bool
		select {
		case input, ok = <-p.inputCh:
			if !ok {
				return // nolint: govet
			}
			stopAttemptTimer()
			readyAttemptCh := make(chan time.Time, 1)
			readyAttemptCh <- time.Time{}
			attemptCh = readyAttemptCh // Enable and trigger the case below
		case <-attemptCh:
			newOutput, res := p.process(input)
			switch res {
			case success:
				output = newOutput
				outputCh = p.outputCh // Enable the 'output' select case
				attemptTimer = nil
				attemptCh = nil
			case backoff:
				attemptTimer = p.retryBackoff.Backoff()
				attemptCh = attemptTimer.C()
			case done:
				// Nothing to do.
				// If 'output' was already set, it remains set still.
				attemptTimer = nil
				attemptCh = nil
			default:
				panic(fmt.Errorf("unknown process result: %v", res))
			}
		case outputCh <- output:
			// Success!
			output = outputT{} // Erase contents to help GC
			outputCh = nil     // Disable this select case (send to nil channel blocks forever)
		}
	}
}
