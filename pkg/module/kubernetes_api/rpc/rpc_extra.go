package rpc

func (x *ImpersonationConfig) GetExtraAsMap() map[string][]string {
	extra := x.GetExtra() // nil-safe
	res := make(map[string][]string, len(extra))
	for _, kv := range extra {
		res[kv.Key] = kv.Val
	}
	return res
}

func (x *ImpersonationConfig) IsEmpty() bool {
	if x == nil {
		return true
	}
	return x.Username == "" && len(x.Groups) == 0 && x.Uid == "" && len(x.Extra) == 0
}
