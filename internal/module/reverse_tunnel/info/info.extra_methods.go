package info

func (x *AgentDescriptor) SupportsServiceAndMethod(service, method string) bool {
	for _, s := range x.GetServices() {
		if s.GetName() != service {
			continue
		}
		// Service found, looking for method.
		for _, m := range s.GetMethods() {
			if m.GetName() == method {
				return true
			}
		}
		break // service checked, no need to continue
	}
	return false
}
