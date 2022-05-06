package mock_modserver

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
)

func NewMockAgentRpcApiWithMockPoller(ctrl *gomock.Controller, pollTimes int) *MockAgentRpcApi {
	mockRpcApi := NewMockAgentRpcApi(ctrl)
	if pollTimes > 0 {
		mockRpcApi.EXPECT().
			PollWithBackoff(gomock.Any(), gomock.Any()).
			DoAndReturn(func(cfg retry.PollConfig, f retry.PollWithBackoffFunc) error {
				for i := 0; i < pollTimes; i++ {
					err, res := f()
					if res == retry.Done {
						return err
					}
				}
				return nil
			})
	}
	return mockRpcApi
}
