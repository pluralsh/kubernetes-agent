package rpc

func ReconcileProjectsFromSlice(projects []string) []*Project {
	ps := make([]*Project, 0, len(projects))
	for _, p := range projects {
		ps = append(ps, &Project{Id: p})
	}
	return ps
}

func (x *ReconcileProjectsRequest) ToProjectSet() map[string]struct{} {
	if x == nil {
		return map[string]struct{}{}
	}

	projects := make(map[string]struct{}, len(x.Project))
	for _, p := range x.Project {
		projects[p.Id] = struct{}{}
	}
	return projects
}
