package manifestops

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
)

func TestRetryPipeline_LastInputOnly(t *testing.T) {
	inputCh := make(chan rpc.ObjectsToSynchronizeData)
	outputCh := make(chan applyJob)
	in1 := rpc.ObjectsToSynchronizeData{
		ProjectId: 1,
	}
	in2 := rpc.ObjectsToSynchronizeData{
		ProjectId: 2,
	}
	out2 := applyJob{
		commitId: "2",
	}
	p := retryPipeline[rpc.ObjectsToSynchronizeData, applyJob]{
		inputCh:      inputCh,
		outputCh:     outputCh,
		retryBackoff: backoffMgr(),
		process: func(input rpc.ObjectsToSynchronizeData) (applyJob, processResult) {
			switch input.ProjectId { // we can receive either value because `select` is not deterministic.
			case in1.ProjectId:
				return applyJob{}, backoff // pretend there was an issue
			case in2.ProjectId:
				return out2, success
			default:
				panic(input)
			}
		},
	}
	go p.run()
	inputCh <- in1
	inputCh <- in2
	out := <-outputCh
	close(inputCh) // stops the goroutine
	assert.Equal(t, out2, out)
}

func TestRetryPipeline_LastOutputOnly(t *testing.T) {
	inputCh := make(chan rpc.ObjectsToSynchronizeData)
	outputCh := make(chan applyJob)
	in2wait := make(chan struct{})
	in1 := rpc.ObjectsToSynchronizeData{
		ProjectId: 1,
	}
	in2 := rpc.ObjectsToSynchronizeData{
		ProjectId: 2,
	}
	out1 := applyJob{
		commitId: "1",
	}
	out2 := applyJob{
		commitId: "2",
	}
	p := retryPipeline[rpc.ObjectsToSynchronizeData, applyJob]{
		inputCh:      inputCh,
		outputCh:     outputCh,
		retryBackoff: backoffMgr(),
		process: func(input rpc.ObjectsToSynchronizeData) (applyJob, processResult) {
			switch input.ProjectId { // we can receive either value because `select` is not deterministic.
			case in1.ProjectId:
				return out1, success
			case in2.ProjectId:
				close(in2wait)
				return out2, success
			default:
				panic(input)
			}
		},
	}
	go p.run()
	inputCh <- in1
	inputCh <- in2
	<-in2wait // wait for in2 to have been processed
	out := <-outputCh
	close(inputCh) // stops the goroutine
	assert.Equal(t, out2, out)
}

func backoffMgr() wait.BackoffManager {
	return wait.NewExponentialBackoffManager(time.Minute, time.Minute, time.Minute, 2, 1, clock.RealClock{}) // nolint:staticcheck
}
