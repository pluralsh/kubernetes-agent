package kasapp

import (
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"google.golang.org/grpc/metadata"
)

func (x *GatewayKasResponse_Header) Metadata() metadata.MD {
	return grpctool.ValuesMapToMeta(x.Meta)
}

func (x *GatewayKasResponse_Trailer) Metadata() metadata.MD {
	return grpctool.ValuesMapToMeta(x.Meta)
}
