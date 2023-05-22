// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.7
// source: internal/module/flux/rpc/rpc.proto

package rpc

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

type ReconcileProjectsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project []*Project `protobuf:"bytes,1,rep,name=project,proto3" json:"project,omitempty"`
}

func (x *ReconcileProjectsRequest) Reset() {
	*x = ReconcileProjectsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_flux_rpc_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconcileProjectsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconcileProjectsRequest) ProtoMessage() {}

func (x *ReconcileProjectsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_flux_rpc_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconcileProjectsRequest.ProtoReflect.Descriptor instead.
func (*ReconcileProjectsRequest) Descriptor() ([]byte, []int) {
	return file_internal_module_flux_rpc_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *ReconcileProjectsRequest) GetProject() []*Project {
	if x != nil {
		return x.Project
	}
	return nil
}

type ReconcileProjectsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project *Project `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
}

func (x *ReconcileProjectsResponse) Reset() {
	*x = ReconcileProjectsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_flux_rpc_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconcileProjectsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconcileProjectsResponse) ProtoMessage() {}

func (x *ReconcileProjectsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_flux_rpc_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconcileProjectsResponse.ProtoReflect.Descriptor instead.
func (*ReconcileProjectsResponse) Descriptor() ([]byte, []int) {
	return file_internal_module_flux_rpc_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *ReconcileProjectsResponse) GetProject() *Project {
	if x != nil {
		return x.Project
	}
	return nil
}

type Project struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Project) Reset() {
	*x = Project{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_flux_rpc_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Project) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Project) ProtoMessage() {}

func (x *Project) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_flux_rpc_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Project.ProtoReflect.Descriptor instead.
func (*Project) Descriptor() ([]byte, []int) {
	return file_internal_module_flux_rpc_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *Project) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_internal_module_flux_rpc_rpc_proto protoreflect.FileDescriptor

var file_internal_module_flux_rpc_rpc_proto_rawDesc = []byte{
	0x0a, 0x22, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x2e, 0x72, 0x70, 0x63, 0x1a, 0x17, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x18, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c,
	0x65, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x38, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x66, 0x6c, 0x75, 0x78, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x5f, 0x0a, 0x19, 0x52, 0x65,
	0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02,
	0x10, 0x01, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x22, 0x0a, 0x07, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x02, 0x69, 0x64, 0x32,
	0x88, 0x01, 0x0a, 0x0a, 0x47, 0x69, 0x74, 0x4c, 0x61, 0x62, 0x46, 0x6c, 0x75, 0x78, 0x12, 0x7a,
	0x0a, 0x11, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x12, 0x2f, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x63, 0x6f,
	0x6e, 0x63, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x63,
	0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x55, 0x5a, 0x53, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d,
	0x6f, 0x72, 0x67, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x74, 0x65,
	0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x36, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x2f, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_module_flux_rpc_rpc_proto_rawDescOnce sync.Once
	file_internal_module_flux_rpc_rpc_proto_rawDescData = file_internal_module_flux_rpc_rpc_proto_rawDesc
)

func file_internal_module_flux_rpc_rpc_proto_rawDescGZIP() []byte {
	file_internal_module_flux_rpc_rpc_proto_rawDescOnce.Do(func() {
		file_internal_module_flux_rpc_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_module_flux_rpc_rpc_proto_rawDescData)
	})
	return file_internal_module_flux_rpc_rpc_proto_rawDescData
}

var file_internal_module_flux_rpc_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_module_flux_rpc_rpc_proto_goTypes = []interface{}{
	(*ReconcileProjectsRequest)(nil),  // 0: gitlab.agent.flux.rpc.ReconcileProjectsRequest
	(*ReconcileProjectsResponse)(nil), // 1: gitlab.agent.flux.rpc.ReconcileProjectsResponse
	(*Project)(nil),                   // 2: gitlab.agent.flux.rpc.Project
}
var file_internal_module_flux_rpc_rpc_proto_depIdxs = []int32{
	2, // 0: gitlab.agent.flux.rpc.ReconcileProjectsRequest.project:type_name -> gitlab.agent.flux.rpc.Project
	2, // 1: gitlab.agent.flux.rpc.ReconcileProjectsResponse.project:type_name -> gitlab.agent.flux.rpc.Project
	0, // 2: gitlab.agent.flux.rpc.GitLabFlux.ReconcileProjects:input_type -> gitlab.agent.flux.rpc.ReconcileProjectsRequest
	1, // 3: gitlab.agent.flux.rpc.GitLabFlux.ReconcileProjects:output_type -> gitlab.agent.flux.rpc.ReconcileProjectsResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_module_flux_rpc_rpc_proto_init() }
func file_internal_module_flux_rpc_rpc_proto_init() {
	if File_internal_module_flux_rpc_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_module_flux_rpc_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconcileProjectsRequest); i {
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
		file_internal_module_flux_rpc_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconcileProjectsResponse); i {
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
		file_internal_module_flux_rpc_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Project); i {
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
			RawDescriptor: file_internal_module_flux_rpc_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_module_flux_rpc_rpc_proto_goTypes,
		DependencyIndexes: file_internal_module_flux_rpc_rpc_proto_depIdxs,
		MessageInfos:      file_internal_module_flux_rpc_rpc_proto_msgTypes,
	}.Build()
	File_internal_module_flux_rpc_rpc_proto = out.File
	file_internal_module_flux_rpc_rpc_proto_rawDesc = nil
	file_internal_module_flux_rpc_rpc_proto_goTypes = nil
	file_internal_module_flux_rpc_rpc_proto_depIdxs = nil
}