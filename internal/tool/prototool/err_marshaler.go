package prototool

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type ProtoErrMarshaler struct {
}

func (ProtoErrMarshaler) Marshal(err error) ([]byte, error) {
	e, ok := err.(proto.Message) // nolint:errorlint
	if !ok {
		return nil, fmt.Errorf("expected proto.Message, got %T", err) // nolint:errorlint
	}
	return ProtoMarshal(e)
}

func (ProtoErrMarshaler) Unmarshal(data []byte) (error, error) {
	e, err := ProtoUnmarshal(data)
	if err != nil {
		return nil, err
	}
	err, ok := e.(error)
	if !ok {
		return nil, fmt.Errorf("expected the proto.Message to be an error but it's not: %T", e)
	}
	return err, nil
}
