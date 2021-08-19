package prototool

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JsonBox ensures the protojson package is used for to/from JSON marshaling.
// See https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson.
type JsonBox struct {
	Message proto.Message
}

func (b *JsonBox) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(b.Message)
}

func (b *JsonBox) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, b.Message)
}
