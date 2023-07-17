package mock_modserver

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"go.uber.org/mock/gomock"
)

func NewMockAgentRpcApiWithMockPoller(ctrl *gomock.Controller, pollTimes int) *MockAgentRpcApi {
	mockRpcApi := NewMockAgentRpcApi(ctrl)
	if pollTimes > 0 {
		mockRpcApi.EXPECT().
			PollWithBackoff(gomock.Any(), gomock.Any()).
			DoAndReturn(func(cfg retry.PollConfig, f retry.PollWithBackoffFunc) error {
				for i := 0; i < pollTimes; i++ {
					err, res := f()
					if err != nil || res == retry.Done {
						return err
					}
				}
				return nil
			})
	}
	return mockRpcApi
}
