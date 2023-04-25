package prototool

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func ProtoMarshal(m proto.Message) ([]byte, error) {
	any, err := anypb.New(m) // use Any to capture type information so that a value can be instantiated in protoUnmarshal()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(any)
}

func ProtoUnmarshal(data []byte) (proto.Message, error) {
	var any anypb.Any
	err := proto.Unmarshal(data, &any)
	if err != nil {
		return nil, err
	}
	return any.UnmarshalNew()
}
