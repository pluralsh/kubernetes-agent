// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: pkg/module/agent_configuration/rpc/rpc.proto

// If you make any changes make sure you run: make regenerate-proto

package rpc

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	agentcfg "github.com/pluralsh/kuberentes-agent/pkg/agentcfg"
	entity "github.com/pluralsh/kuberentes-agent/pkg/entity"
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

type ConfigurationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Last processed commit id. Optional.
	// Server will only send configuration if the last commit on the branch
	// is a different one. If a connection breaks, this allows to resume
	// the stream without sending the same data again.
	CommitId string `protobuf:"bytes,1,opt,name=commit_id,json=commitId,proto3" json:"commit_id,omitempty"`
	// Information about the agent.
	AgentMeta *entity.AgentMeta `protobuf:"bytes,2,opt,name=agent_meta,json=agentMeta,proto3" json:"agent_meta,omitempty"`
	// Flag to skip agent registration.
	SkipRegister bool `protobuf:"varint,3,opt,name=skip_register,json=skipRegister,proto3" json:"skip_register,omitempty"`
}

func (x *ConfigurationRequest) Reset() {
	*x = ConfigurationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigurationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigurationRequest) ProtoMessage() {}

func (x *ConfigurationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigurationRequest.ProtoReflect.Descriptor instead.
func (*ConfigurationRequest) Descriptor() ([]byte, []int) {
	return file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigurationRequest) GetCommitId() string {
	if x != nil {
		return x.CommitId
	}
	return ""
}

func (x *ConfigurationRequest) GetAgentMeta() *entity.AgentMeta {
	if x != nil {
		return x.AgentMeta
	}
	return nil
}

func (x *ConfigurationRequest) GetSkipRegister() bool {
	if x != nil {
		return x.SkipRegister
	}
	return false
}

type ConfigurationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Configuration *agentcfg.AgentConfiguration `protobuf:"bytes,1,opt,name=configuration,proto3" json:"configuration,omitempty"`
	// Commit id of the configuration repository.
	// Can be used to resume connection from where it dropped.
	CommitId string `protobuf:"bytes,2,opt,name=commit_id,json=commitId,proto3" json:"commit_id,omitempty"`
}

func (x *ConfigurationResponse) Reset() {
	*x = ConfigurationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigurationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigurationResponse) ProtoMessage() {}

func (x *ConfigurationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigurationResponse.ProtoReflect.Descriptor instead.
func (*ConfigurationResponse) Descriptor() ([]byte, []int) {
	return file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *ConfigurationResponse) GetConfiguration() *agentcfg.AgentConfiguration {
	if x != nil {
		return x.Configuration
	}
	return nil
}

func (x *ConfigurationResponse) GetCommitId() string {
	if x != nil {
		return x.CommitId
	}
	return ""
}

var File_pkg_module_agent_configuration_rpc_rpc_proto protoreflect.FileDescriptor

var file_pkg_module_agent_configuration_rpc_rpc_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x24,
	0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x72, 0x70, 0x63, 0x1a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x63,
	0x66, 0x67, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x63, 0x66, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x97, 0x01, 0x0a, 0x14, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x09, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6b, 0x69, 0x70,
	0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0c, 0x73, 0x6b, 0x69, 0x70, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x22, 0x8e, 0x01,
	0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29,
	0x2e, 0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x63, 0x66, 0x67, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x20, 0x01, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49, 0x64, 0x32, 0xa6,
	0x01, 0x0a, 0x12, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x8f, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x2e, 0x70, 0x6c, 0x75,
	0x72, 0x61, 0x6c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3b, 0x2e, 0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x73, 0x68, 0x2f, 0x6b,
	0x75, 0x62, 0x65, 0x72, 0x65, 0x6e, 0x74, 0x65, 0x73, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescOnce sync.Once
	file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescData = file_pkg_module_agent_configuration_rpc_rpc_proto_rawDesc
)

func file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescGZIP() []byte {
	file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescOnce.Do(func() {
		file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescData)
	})
	return file_pkg_module_agent_configuration_rpc_rpc_proto_rawDescData
}

var file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_module_agent_configuration_rpc_rpc_proto_goTypes = []interface{}{
	(*ConfigurationRequest)(nil),        // 0: plural.agent.agent_configuration.rpc.ConfigurationRequest
	(*ConfigurationResponse)(nil),       // 1: plural.agent.agent_configuration.rpc.ConfigurationResponse
	(*entity.AgentMeta)(nil),            // 2: plural.agent.entity.AgentMeta
	(*agentcfg.AgentConfiguration)(nil), // 3: plural.agent.agentcfg.AgentConfiguration
}
var file_pkg_module_agent_configuration_rpc_rpc_proto_depIdxs = []int32{
	2, // 0: plural.agent.agent_configuration.rpc.ConfigurationRequest.agent_meta:type_name -> plural.agent.entity.AgentMeta
	3, // 1: plural.agent.agent_configuration.rpc.ConfigurationResponse.configuration:type_name -> plural.agent.agentcfg.AgentConfiguration
	0, // 2: plural.agent.agent_configuration.rpc.AgentConfiguration.GetConfiguration:input_type -> plural.agent.agent_configuration.rpc.ConfigurationRequest
	1, // 3: plural.agent.agent_configuration.rpc.AgentConfiguration.GetConfiguration:output_type -> plural.agent.agent_configuration.rpc.ConfigurationResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_module_agent_configuration_rpc_rpc_proto_init() }
func file_pkg_module_agent_configuration_rpc_rpc_proto_init() {
	if File_pkg_module_agent_configuration_rpc_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigurationRequest); i {
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
		file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigurationResponse); i {
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
			RawDescriptor: file_pkg_module_agent_configuration_rpc_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_module_agent_configuration_rpc_rpc_proto_goTypes,
		DependencyIndexes: file_pkg_module_agent_configuration_rpc_rpc_proto_depIdxs,
		MessageInfos:      file_pkg_module_agent_configuration_rpc_rpc_proto_msgTypes,
	}.Build()
	File_pkg_module_agent_configuration_rpc_rpc_proto = out.File
	file_pkg_module_agent_configuration_rpc_rpc_proto_rawDesc = nil
	file_pkg_module_agent_configuration_rpc_rpc_proto_goTypes = nil
	file_pkg_module_agent_configuration_rpc_rpc_proto_depIdxs = nil
}