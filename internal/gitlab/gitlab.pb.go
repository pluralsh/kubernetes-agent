// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: internal/gitlab/gitlab.proto

// If you make any changes make sure you run: make regenerate-proto

package gitlab

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ClientError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Path       string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Reason     string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *ClientError) Reset() {
	*x = ClientError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_gitlab_gitlab_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientError) ProtoMessage() {}

func (x *ClientError) ProtoReflect() protoreflect.Message {
	mi := &file_internal_gitlab_gitlab_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientError.ProtoReflect.Descriptor instead.
func (*ClientError) Descriptor() ([]byte, []int) {
	return file_internal_gitlab_gitlab_proto_rawDescGZIP(), []int{0}
}

func (x *ClientError) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *ClientError) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ClientError) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

// see https://gitlab.com/gitlab-org/gitlab/blob/2864126a72835bd0b29f670ffc36828014850f5f/lib/api/helpers.rb#L534-534
type DefaultApiError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DefaultApiError) Reset() {
	*x = DefaultApiError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_gitlab_gitlab_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefaultApiError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefaultApiError) ProtoMessage() {}

func (x *DefaultApiError) ProtoReflect() protoreflect.Message {
	mi := &file_internal_gitlab_gitlab_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefaultApiError.ProtoReflect.Descriptor instead.
func (*DefaultApiError) Descriptor() ([]byte, []int) {
	return file_internal_gitlab_gitlab_proto_rawDescGZIP(), []int{1}
}

func (x *DefaultApiError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_internal_gitlab_gitlab_proto protoreflect.FileDescriptor

var file_internal_gitlab_gitlab_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x67, 0x69, 0x74,
	0x6c, 0x61, 0x62, 0x22, 0x5a, 0x0a, 0x0b, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22,
	0x2b, 0x0a, 0x0f, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x41, 0x70, 0x69, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x36, 0x5a, 0x34,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x75, 0x72, 0x61,
	0x6c, 0x73, 0x68, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x65, 0x6e, 0x74, 0x65, 0x73, 0x2d, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_gitlab_gitlab_proto_rawDescOnce sync.Once
	file_internal_gitlab_gitlab_proto_rawDescData = file_internal_gitlab_gitlab_proto_rawDesc
)

func file_internal_gitlab_gitlab_proto_rawDescGZIP() []byte {
	file_internal_gitlab_gitlab_proto_rawDescOnce.Do(func() {
		file_internal_gitlab_gitlab_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_gitlab_gitlab_proto_rawDescData)
	})
	return file_internal_gitlab_gitlab_proto_rawDescData
}

var file_internal_gitlab_gitlab_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_gitlab_gitlab_proto_goTypes = []interface{}{
	(*ClientError)(nil),     // 0: gitlab.agent.gitlab.ClientError
	(*DefaultApiError)(nil), // 1: gitlab.agent.gitlab.DefaultApiError
}
var file_internal_gitlab_gitlab_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_gitlab_gitlab_proto_init() }
func file_internal_gitlab_gitlab_proto_init() {
	if File_internal_gitlab_gitlab_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_gitlab_gitlab_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientError); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_gitlab_gitlab_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefaultApiError); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_gitlab_gitlab_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_gitlab_gitlab_proto_goTypes,
		DependencyIndexes: file_internal_gitlab_gitlab_proto_depIdxs,
		MessageInfos:      file_internal_gitlab_gitlab_proto_msgTypes,
	}.Build()
	File_internal_gitlab_gitlab_proto = out.File
	file_internal_gitlab_gitlab_proto_rawDesc = nil
	file_internal_gitlab_gitlab_proto_goTypes = nil
	file_internal_gitlab_gitlab_proto_depIdxs = nil
}
