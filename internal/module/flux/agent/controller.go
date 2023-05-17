package agent

import (
	"context"

	"k8s.io/client-go/informers"
)

type controller interface {
	Run(ctx context.Context)
}

type controllerFactory func(ctx context.Context, gitRepositoryInformer informers.GenericInformer, receiverInformer informers.GenericInformer, projectReconciler projectReconciler) (controller, error)

type reconciliationResult struct {
	status reconciliationStatus
	error  error
}

type reconciliationStatus int

const (
	RetryRateLimited reconciliationStatus = iota
	Success
	Error
)
