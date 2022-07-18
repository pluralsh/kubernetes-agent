package agent

import (
	"bytes"
	"context"
	"os"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
	"sigs.k8s.io/cli-utils/pkg/printers"
)

type applyJob struct {
	ctx      context.Context
	commitId string
	invInfo  inventory.Info
	objects  []*unstructured.Unstructured
}

func (w *defaultGitopsWorker) apply(jobs <-chan applyJob) {
	for job := range jobs {
		l := w.log.With(logz.CommitId(job.commitId))
		_ = retry.PollWithBackoff(job.ctx, w.applierPollConfig, func(ctx context.Context) (error, retry.AttemptResult) {
			l.Info("Synchronizing objects")
			err := w.applyJob(ctx, job)
			if err != nil {
				if errz.ContextDone(err) {
					l.Info("Synchronization was canceled", logz.Error(err))
				} else {
					l.Warn("Synchronization failed", logz.Error(err))
				}
				return nil, retry.Backoff
			}
			l.Info("Objects synchronized")
			return nil, retry.Continue
		})
	}
}

func (w *defaultGitopsWorker) applyJob(ctx context.Context, job applyJob) error {
	events := w.applier.Run(ctx, job.invInfo, job.objects, w.applyOptions)
	// The printer will print updates from the channel. It will block
	// until the channel is closed.
	printer := printers.GetPrinter(printers.JSONPrinter, genericclioptions.IOStreams{
		In:     &bytes.Buffer{}, // nothing to read
		Out:    os.Stderr,
		ErrOut: os.Stderr,
	})
	return printer.Print(events, common.DryRunNone, true)
}
