package agent

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestFullSyncExecution(t *testing.T) {
	partialSyncInterval := durationpb.New(30 * time.Millisecond)
	fullSyncInterval := durationpb.New(70 * time.Millisecond)

	// the test will be configured to run for a target number of full sync cycles
	// and the number of intermittent partial syncs will be compared to an expected value
	targetFullSyncCount := uint32(3)

	// 3 full sync cycles are expected to occur at 0ms, 70ms, 140ms.
	// Between the first and the last full sync, partial syncs are expected to occur
	// at 30ms, 60ms, 90ms and 120ms i.e. 4 times
	expectedPartialSyncCount := uint32(4)

	fullSyncCallCounter := uint32(0)

	ctx, cancel := context.WithCancel(context.Background())
	mock := &mockReconciler{}

	w := &worker{
		log: zaptest.NewLogger(t),
		api: newMockApi(t),
		reconcilerFactory: func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevReconciler, error) {
			fullSyncCallCounter += 1

			if fullSyncCallCounter == targetFullSyncCount {
				cancel()
			}

			return mock, nil
		},
	}

	cfg := &agentcfg.RemoteCF{
		Enabled:             true,
		DnsZone:             "",
		PartialSyncInterval: partialSyncInterval,
		FullSyncInterval:    fullSyncInterval,
	}

	err := w.StartReconciliation(ctx, cfg)
	require.NoError(t, err)

	// mock reconciler will be invoked for every full sync cycle
	// so full sync call counter must be subtracted to get partial sync cycles
	partialSyncCallCounter := mock.timesCalled - fullSyncCallCounter

	require.EqualValues(t, expectedPartialSyncCount, partialSyncCallCounter, "partial sync call count: %d", partialSyncCallCounter)
}

func newMockApi(t *testing.T) *mock_modagent.MockApi {
	mockApi := mock_modagent.NewMockApi(gomock.NewController(t))
	mockApi.EXPECT().GetAgentId(gomock.Any()).AnyTimes()

	return mockApi
}
