package agent

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
)

// worker is responsible for coordinating and executing full and
// partial syncs by managing the lifecycle of a reconciler
type worker struct {
	log               *zap.Logger
	api               modagent.Api
	reconcilerFactory func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevReconciler, error)
}

func (w *worker) StartReconciliation(ctx context.Context, cfg *agentcfg.RemoteCF) error {
	agentId, err := w.api.GetAgentId(ctx)
	if err != nil {
		return err
	}

	fullSyncInterval := cfg.GetFullSyncInterval().AsDuration()
	partialSyncInterval := cfg.GetPartialSyncInterval().AsDuration()

	// full sync should be started immediately
	// upon module start/restart
	fullSyncTimer := time.NewTimer(0)
	defer fullSyncTimer.Stop()

	partialSyncTimer := time.NewTimer(partialSyncInterval)
	defer partialSyncTimer.Stop()

	var (
		activeReconciler remoteDevReconciler
	)
	defer func() {
		// this nil check is needed in case the context is cancelled before the first
		// full sync has even been scheduled to execute. In such a case, the
		// active reconciler will still be nil and the call to Stop may be skipped
		if activeReconciler != nil {
			activeReconciler.Stop()
		}
	}()

	for {
		select {
		// this check allows the goroutine to immediately exit
		// if the context cancellation is invoked while waiting on
		// either of the timers
		case <-ctx.Done():
			return nil

		case <-fullSyncTimer.C:
			// full sync is implemented by creating a new reconciler
			// while discarding the state accrued in the previous reconciler
			w.log.Info("starting full sync")
			if activeReconciler != nil {
				activeReconciler.Stop()
			}

			// Full sync could have been alternatively implemented by re-using a reconciler vs
			// the current approach of destroying the existing reconciler and creating a new one
			// This has been done for the following reasons
			//  1. The current approach of stopping/starting a reconciler is conceptually
			//		equivalent to a module restart with minimal changes to the reconciliation logic/code.
			//		Full sync implemented with reconciler re-use would've required introducing special handling in the
			// 		reconciliation code. This means that any future updates to the reconciliation logic would've
			//		led to increased complexity due to possible impact on full & partial sync logic
			//		leading to increased maintenance cost.
			//  2. Reconciler re-use is inadequate when dealing with issues that occur due to corruption/
			//		mishandling of internal state, for example memory leaks that may occur due to bugs in
			//		reconciliation logic. The current approach will be able to deal with these as the core logic
			// 		requires teardown of the active reconciler and creation of a new one
			activeReconciler, err = w.reconcilerFactory(ctx, cfg)
			if err != nil {
				return err
			}
			execError := activeReconciler.Run(ctx)
			if execError != nil {
				w.api.HandleProcessingError(
					ctx, w.log, agentId,
					"Remote Dev - full sync cycle ended with error", execError,
				)
			}

			// Timer is reset after the work has been completed
			// If the timer were reset before reconciliation is executed, there may be a scenario
			// where the next timer tick occurs immediately after the reconciler finishes its
			// execution (because Run() takes too long for some reason)
			fullSyncTimer.Reset(fullSyncInterval)

		case <-partialSyncTimer.C:
			w.log.Info("starting partial update")
			execError := activeReconciler.Run(ctx)
			if execError != nil {
				w.api.HandleProcessingError(
					ctx, w.log, agentId,
					"Remote Dev - partial sync cycle ended with error", execError,
				)
			}

			// Timer is reset after the work has been completed
			// If the timer were reset before reconciler is executed, there may be a scenario
			// where the next timer tick occurs immediately after the reconciler finishes its
			// execution (because Run() takes too long for some reason)
			partialSyncTimer.Reset(partialSyncInterval)
		}
	}
}
