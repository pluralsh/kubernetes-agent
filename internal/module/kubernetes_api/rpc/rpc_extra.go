package rpc

func (x *ImpersonationConfig) GetExtraAsMap() map[string][]string {
	extra := x.GetExtra() // nil-safe
	res := make(map[string][]string, len(extra))
	for _, kv := range extra {
		res[kv.GetKey()] = kv.GetVal()
	}
	return res
}

func (x *ImpersonationConfig) IsEmpty() bool {
	if x == nil {
		return true
	}
	return x.GetUsername() == "" && len(x.GetGroups()) == 0 && x.GetUid() == "" && len(x.GetExtra()) == 0
}
