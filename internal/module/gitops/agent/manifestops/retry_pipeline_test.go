package manifestops

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
)

func TestRetryPipeline_LastInputOnly(t *testing.T) {
	inputCh := make(chan inputT)
	outputCh := make(chan outputT)
	in1 := inputT{
		ProjectId: 1,
	}
	in2 := inputT{
		ProjectId: 2,
	}
	out2 := outputT{
		commitId: "2",
	}
	p := retryPipeline{
		inputCh:      inputCh,
		outputCh:     outputCh,
		retryBackoff: backoffMgr(),
		process: func(input inputT) (outputT, processResult) {
			switch input.ProjectId { // we can receive either value because `select` is not deterministic.
			case in1.ProjectId:
				return outputT{}, backoff // pretend there was an issue
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
	inputCh := make(chan inputT)
	outputCh := make(chan outputT)
	in2wait := make(chan struct{})
	in1 := inputT{
		ProjectId: 1,
	}
	in2 := inputT{
		ProjectId: 2,
	}
	out1 := outputT{
		commitId: "1",
	}
	out2 := outputT{
		commitId: "2",
	}
	p := retryPipeline{
		inputCh:      inputCh,
		outputCh:     outputCh,
		retryBackoff: backoffMgr(),
		process: func(input inputT) (outputT, processResult) {
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
	return wait.NewExponentialBackoffManager(time.Minute, time.Minute, time.Minute, 2, 1, clock.RealClock{})
}
