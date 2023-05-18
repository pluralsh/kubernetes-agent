// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.7
// source: internal/module/reverse_tunnel/info/info.proto

package info

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Method struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Method) Reset() {
	*x = Method{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_reverse_tunnel_info_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Method) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Method) ProtoMessage() {}

func (x *Method) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_reverse_tunnel_info_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Method.ProtoReflect.Descriptor instead.
func (*Method) Descriptor() ([]byte, []int) {
	return file_internal_module_reverse_tunnel_info_info_proto_rawDescGZIP(), []int{0}
}

func (x *Method) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Methods []*Method `protobuf:"bytes,2,rep,name=methods,proto3" json:"methods,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_reverse_tunnel_info_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_reverse_tunnel_info_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_internal_module_reverse_tunnel_info_info_proto_rawDescGZIP(), []int{1}
}

func (x *Service) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Service) GetMethods() []*Method {
	if x != nil {
		return x.Methods
	}
	return nil
}

type AgentDescriptor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Services []*Service `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
}

func (x *AgentDescriptor) Reset() {
	*x = AgentDescriptor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_reverse_tunnel_info_info_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentDescriptor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentDescriptor) ProtoMessage() {}

func (x *AgentDescriptor) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_reverse_tunnel_info_info_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentDescriptor.ProtoReflect.Descriptor instead.
func (*AgentDescriptor) Descriptor() ([]byte, []int) {
	return file_internal_module_reverse_tunnel_info_info_proto_rawDescGZIP(), []int{2}
}

func (x *AgentDescriptor) GetServices() []*Service {
	if x != nil {
		return x.Services
	}
	return nil
}

var File_internal_module_reverse_tunnel_info_info_proto protoreflect.FileDescriptor

var file_internal_module_reverse_tunnel_info_info_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x2f, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x5f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c,
	0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x20, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x72,
	0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x5f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x69, 0x6e,
	0x66, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x06, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0x6a, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x20, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x72, 0x65, 0x76, 0x65, 0x72,
	0x73, 0x65, 0x5f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x22, 0x58,
	0x0a, 0x0f, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x12, 0x45, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x5f, 0x74, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x08,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x42, 0x60, 0x5a, 0x5e, 0x67, 0x69, 0x74, 0x6c,
	0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x6f, 0x72,
	0x67, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2f, 0x76, 0x31, 0x36, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x5f, 0x74,
	0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_internal_module_reverse_tunnel_info_info_proto_rawDescOnce sync.Once
	file_internal_module_reverse_tunnel_info_info_proto_rawDescData = file_internal_module_reverse_tunnel_info_info_proto_rawDesc
)

func file_internal_module_reverse_tunnel_info_info_proto_rawDescGZIP() []byte {
	file_internal_module_reverse_tunnel_info_info_proto_rawDescOnce.Do(func() {
		file_internal_module_reverse_tunnel_info_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_module_reverse_tunnel_info_info_proto_rawDescData)
	})
	return file_internal_module_reverse_tunnel_info_info_proto_rawDescData
}

var file_internal_module_reverse_tunnel_info_info_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_module_reverse_tunnel_info_info_proto_goTypes = []interface{}{
	(*Method)(nil),          // 0: gitlab.agent.reverse_tunnel.info.Method
	(*Service)(nil),         // 1: gitlab.agent.reverse_tunnel.info.Service
	(*AgentDescriptor)(nil), // 2: gitlab.agent.reverse_tunnel.info.AgentDescriptor
}
var file_internal_module_reverse_tunnel_info_info_proto_depIdxs = []int32{
	0, // 0: gitlab.agent.reverse_tunnel.info.Service.methods:type_name -> gitlab.agent.reverse_tunnel.info.Method
	1, // 1: gitlab.agent.reverse_tunnel.info.AgentDescriptor.services:type_name -> gitlab.agent.reverse_tunnel.info.Service
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_module_reverse_tunnel_info_info_proto_init() }
func file_internal_module_reverse_tunnel_info_info_proto_init() {
	if File_internal_module_reverse_tunnel_info_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_module_reverse_tunnel_info_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Method); i {
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
		file_internal_module_reverse_tunnel_info_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service); i {
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
		file_internal_module_reverse_tunnel_info_info_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentDescriptor); i {
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
			RawDescriptor: file_internal_module_reverse_tunnel_info_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_module_reverse_tunnel_info_info_proto_goTypes,
		DependencyIndexes: file_internal_module_reverse_tunnel_info_info_proto_depIdxs,
		MessageInfos:      file_internal_module_reverse_tunnel_info_info_proto_msgTypes,
	}.Build()
	File_internal_module_reverse_tunnel_info_info_proto = out.File
	file_internal_module_reverse_tunnel_info_info_proto_rawDesc = nil
	file_internal_module_reverse_tunnel_info_info_proto_goTypes = nil
	file_internal_module_reverse_tunnel_info_info_proto_depIdxs = nil
}
