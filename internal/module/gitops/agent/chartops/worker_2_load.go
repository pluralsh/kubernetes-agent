package chartops

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func (w *worker) decode(desiredState <-chan fetchedData, jobs chan<- job) {
	var (
		newJob      job
		nilableJobs chan<- job
		jobCancel   context.CancelFunc
	)
	defer func() {
		if jobCancel != nil {
			jobCancel()
		}
	}()
	for {
		select {
		case data, ok := <-desiredState:
			if !ok {
				return // nolint: govet
			}
			chart, err := loader.LoadFiles(data.files)
			if err != nil {
				data.log.Error("Failed to load chart", logz.Error(err))
				continue
			}
			if jobCancel != nil {
				jobCancel() // Cancel running/pending job ASAP
			}
			newJob = job{
				log:   data.log,
				chart: chart,
				vals:  nil, // TODO load values from somewhere
			}
			newJob.ctx, jobCancel = context.WithCancel(context.Background()) // nolint: govet
			nilableJobs = jobs
		case nilableJobs <- newJob:
			// Success!
			newJob = job{}    // Erase contents to help GC
			nilableJobs = nil // Disable this select case (send to nil channel blocks forever)
		}
	}
}
