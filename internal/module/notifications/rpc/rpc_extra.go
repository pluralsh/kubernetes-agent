package rpc

import "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"

func (x *Project) ToNotificationsProject() *modserver.Project {
	return &modserver.Project{
		Id:       x.Id,
		FullPath: x.FullPath,
	}
}
