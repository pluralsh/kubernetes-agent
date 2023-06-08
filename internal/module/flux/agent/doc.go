package agent

//go:generate mockgen.sh -source "reconcile_trigger.go" -destination "mock_reconciler_trigger_for_test.go" -package "agent" "reconcileTrigger"

//go:generate mockgen.sh -source "client.go" -destination "mock_client_for_test.go" -package "agent" "projectReconciler"

//go:generate mockgen.sh -destination "mock_rpc_for_test.go" -package "agent" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc" "GitLabFluxClient,GitLabFlux_ReconcileProjectsClient"
