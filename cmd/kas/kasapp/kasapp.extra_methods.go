package kasapp

import (
	"google.golang.org/grpc/metadata"

	"github.com/pluralsh/kuberentes-agent/pkg/tool/grpctool"
)

func (x *GatewayKasResponse_Header) Metadata() metadata.MD {
	return grpctool.ValuesMapToMeta(x.Meta)
}

func (x *GatewayKasResponse_Trailer) Metadata() metadata.MD {
	return grpctool.ValuesMapToMeta(x.Meta)
}
