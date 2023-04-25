package rpc

import "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver/notifications"

func (x *Project) ToNotificationsProject() *notifications.Project {
	return &notifications.Project{
		Id:       x.Id,
		FullPath: x.FullPath,
	}
}
