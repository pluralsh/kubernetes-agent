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

// MarshalJSON implements json.Marshaler on JsonBox.
// It must have a value receiver to make it work on non-addressable values e.g.:
// - json.Marshal(JsonBox{...})
// - json.Marshal(SomeTypeWhereJsonBoxIsNotAPointerField{...})
// See https://golang.org/ref/spec#Address_operators.
func (b JsonBox) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(b.Message)
}

// UnmarshalJSON implements json.Unmarshaler on JsonBox.
func (b *JsonBox) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, b.Message)
}
