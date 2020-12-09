// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: internal/module/gitlab_access/rpc/rpc.proto

package rpc

import (
	context "context"
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	_ "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/grpctool/automata"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Values struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *Values) Reset() {
	*x = Values{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Values) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Values) ProtoMessage() {}

func (x *Values) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Values.ProtoReflect.Descriptor instead.
func (*Values) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *Values) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//	*Request_Headers_
	//	*Request_Data_
	//	*Request_Trailers_
	Message isRequest_Message `protobuf_oneof:"message"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{1}
}

func (m *Request) GetMessage() isRequest_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Request) GetHeaders() *Request_Headers {
	if x, ok := x.GetMessage().(*Request_Headers_); ok {
		return x.Headers
	}
	return nil
}

func (x *Request) GetData() *Request_Data {
	if x, ok := x.GetMessage().(*Request_Data_); ok {
		return x.Data
	}
	return nil
}

func (x *Request) GetTrailers() *Request_Trailers {
	if x, ok := x.GetMessage().(*Request_Trailers_); ok {
		return x.Trailers
	}
	return nil
}

type isRequest_Message interface {
	isRequest_Message()
}

type Request_Headers_ struct {
	Headers *Request_Headers `protobuf:"bytes,1,opt,name=headers,proto3,oneof"`
}

type Request_Data_ struct {
	Data *Request_Data `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}

type Request_Trailers_ struct {
	Trailers *Request_Trailers `protobuf:"bytes,3,opt,name=trailers,proto3,oneof"`
}

func (*Request_Headers_) isRequest_Message() {}

func (*Request_Data_) isRequest_Message() {}

func (*Request_Trailers_) isRequest_Message() {}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//	*Response_Headers_
	//	*Response_Data_
	//	*Response_Trailers_
	Message isResponse_Message `protobuf_oneof:"message"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{2}
}

func (m *Response) GetMessage() isResponse_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Response) GetHeaders() *Response_Headers {
	if x, ok := x.GetMessage().(*Response_Headers_); ok {
		return x.Headers
	}
	return nil
}

func (x *Response) GetData() *Response_Data {
	if x, ok := x.GetMessage().(*Response_Data_); ok {
		return x.Data
	}
	return nil
}

func (x *Response) GetTrailers() *Response_Trailers {
	if x, ok := x.GetMessage().(*Response_Trailers_); ok {
		return x.Trailers
	}
	return nil
}

type isResponse_Message interface {
	isResponse_Message()
}

type Response_Headers_ struct {
	Headers *Response_Headers `protobuf:"bytes,1,opt,name=headers,proto3,oneof"`
}

type Response_Data_ struct {
	Data *Response_Data `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}

type Response_Trailers_ struct {
	Trailers *Response_Trailers `protobuf:"bytes,3,opt,name=trailers,proto3,oneof"`
}

func (*Response_Headers_) isResponse_Message() {}

func (*Response_Data_) isResponse_Message() {}

func (*Response_Trailers_) isResponse_Message() {}

type Request_Headers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ModuleName string             `protobuf:"bytes,1,opt,name=module_name,json=moduleName,proto3" json:"module_name,omitempty"`
	Method     string             `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Headers    map[string]*Values `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	UrlPath    string             `protobuf:"bytes,4,opt,name=url_path,json=urlPath,proto3" json:"url_path,omitempty"`
	Query      map[string]*Values `protobuf:"bytes,5,rep,name=query,proto3" json:"query,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Request_Headers) Reset() {
	*x = Request_Headers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request_Headers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request_Headers) ProtoMessage() {}

func (x *Request_Headers) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request_Headers.ProtoReflect.Descriptor instead.
func (*Request_Headers) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Request_Headers) GetModuleName() string {
	if x != nil {
		return x.ModuleName
	}
	return ""
}

func (x *Request_Headers) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Request_Headers) GetHeaders() map[string]*Values {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Request_Headers) GetUrlPath() string {
	if x != nil {
		return x.UrlPath
	}
	return ""
}

func (x *Request_Headers) GetQuery() map[string]*Values {
	if x != nil {
		return x.Query
	}
	return nil
}

type Request_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Request_Data) Reset() {
	*x = Request_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request_Data) ProtoMessage() {}

func (x *Request_Data) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request_Data.ProtoReflect.Descriptor instead.
func (*Request_Data) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Request_Data) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Request_Trailers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Request_Trailers) Reset() {
	*x = Request_Trailers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request_Trailers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request_Trailers) ProtoMessage() {}

func (x *Request_Trailers) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request_Trailers.ProtoReflect.Descriptor instead.
func (*Request_Trailers) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{1, 2}
}

type Response_Headers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32              `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Status     string             `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Headers    map[string]*Values `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Response_Headers) Reset() {
	*x = Response_Headers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_Headers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_Headers) ProtoMessage() {}

func (x *Response_Headers) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_Headers.ProtoReflect.Descriptor instead.
func (*Response_Headers) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Response_Headers) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Response_Headers) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Response_Headers) GetHeaders() map[string]*Values {
	if x != nil {
		return x.Headers
	}
	return nil
}

type Response_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Response_Data) Reset() {
	*x = Response_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_Data) ProtoMessage() {}

func (x *Response_Data) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_Data.ProtoReflect.Descriptor instead.
func (*Response_Data) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Response_Data) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Response_Trailers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response_Trailers) Reset() {
	*x = Response_Trailers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_Trailers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_Trailers) ProtoMessage() {}

func (x *Response_Trailers) ProtoReflect() protoreflect.Message {
	mi := &file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_Trailers.ProtoReflect.Descriptor instead.
func (*Response_Trailers) Descriptor() ([]byte, []int) {
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP(), []int{2, 2}
}

var File_internal_module_gitlab_access_rpc_rpc_proto protoreflect.FileDescriptor

var file_internal_module_gitlab_access_rpc_rpc_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x72,
	0x70, 0x63, 0x1a, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x6f, 0x6f,
	0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x74, 0x6f, 0x6f, 0x6c, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6d,
	0x61, 0x74, 0x61, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x1e, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0xe3, 0x04, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3e,
	0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x42, 0x0c, 0x82, 0xf6, 0x2c, 0x02, 0x08, 0x02, 0x82, 0xf6, 0x2c,
	0x02, 0x08, 0x03, 0x48, 0x00, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x35,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x42,
	0x0c, 0x82, 0xf6, 0x2c, 0x02, 0x08, 0x02, 0x82, 0xf6, 0x2c, 0x02, 0x08, 0x03, 0x48, 0x00, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x44, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73, 0x42, 0x0f,
	0x82, 0xf6, 0x2c, 0x0b, 0x08, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x48,
	0x00, 0x52, 0x08, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73, 0x1a, 0xe1, 0x02, 0x0a, 0x07,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x3b, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x21, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x19, 0x0a,
	0x08, 0x75, 0x72, 0x6c, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x75, 0x72, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x12, 0x35, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x1a,
	0x47, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x21, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x45, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x21, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a,
	0x1a, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x0a, 0x0a, 0x08, 0x54,
	0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73, 0x42, 0x0f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x01, 0x22, 0xcf, 0x03, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x42, 0x0c, 0x82,
	0xf6, 0x2c, 0x02, 0x08, 0x02, 0x82, 0xf6, 0x2c, 0x02, 0x08, 0x03, 0x48, 0x00, 0x52, 0x07, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x36, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x42, 0x0c, 0x82, 0xf6, 0x2c, 0x02, 0x08, 0x02,
	0x82, 0xf6, 0x2c, 0x02, 0x08, 0x03, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x45,
	0x0a, 0x08, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x54, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73, 0x42, 0x0f, 0x82, 0xf6, 0x2c, 0x0b, 0x08, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x48, 0x00, 0x52, 0x08, 0x74, 0x72, 0x61,
	0x69, 0x6c, 0x65, 0x72, 0x73, 0x1a, 0xc9, 0x01, 0x0a, 0x07, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3c, 0x0a, 0x07, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x47, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x21, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x1a, 0x1a, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x0a, 0x0a,
	0x08, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73, 0x42, 0x0f, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x04, 0x88, 0xf6, 0x2c, 0x01, 0x32, 0x40, 0x0a, 0x0c, 0x47, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x30, 0x0a, 0x0b, 0x4d, 0x61,
	0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x5a, 0x5a, 0x58,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x69, 0x6e,
	0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x5f, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_module_gitlab_access_rpc_rpc_proto_rawDescOnce sync.Once
	file_internal_module_gitlab_access_rpc_rpc_proto_rawDescData = file_internal_module_gitlab_access_rpc_rpc_proto_rawDesc
)

func file_internal_module_gitlab_access_rpc_rpc_proto_rawDescGZIP() []byte {
	file_internal_module_gitlab_access_rpc_rpc_proto_rawDescOnce.Do(func() {
		file_internal_module_gitlab_access_rpc_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_module_gitlab_access_rpc_rpc_proto_rawDescData)
	})
	return file_internal_module_gitlab_access_rpc_rpc_proto_rawDescData
}

var file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_internal_module_gitlab_access_rpc_rpc_proto_goTypes = []interface{}{
	(*Values)(nil),            // 0: rpc.Values
	(*Request)(nil),           // 1: rpc.Request
	(*Response)(nil),          // 2: rpc.Response
	(*Request_Headers)(nil),   // 3: rpc.Request.Headers
	(*Request_Data)(nil),      // 4: rpc.Request.Data
	(*Request_Trailers)(nil),  // 5: rpc.Request.Trailers
	nil,                       // 6: rpc.Request.Headers.HeadersEntry
	nil,                       // 7: rpc.Request.Headers.QueryEntry
	(*Response_Headers)(nil),  // 8: rpc.Response.Headers
	(*Response_Data)(nil),     // 9: rpc.Response.Data
	(*Response_Trailers)(nil), // 10: rpc.Response.Trailers
	nil,                       // 11: rpc.Response.Headers.HeadersEntry
}
var file_internal_module_gitlab_access_rpc_rpc_proto_depIdxs = []int32{
	3,  // 0: rpc.Request.headers:type_name -> rpc.Request.Headers
	4,  // 1: rpc.Request.data:type_name -> rpc.Request.Data
	5,  // 2: rpc.Request.trailers:type_name -> rpc.Request.Trailers
	8,  // 3: rpc.Response.headers:type_name -> rpc.Response.Headers
	9,  // 4: rpc.Response.data:type_name -> rpc.Response.Data
	10, // 5: rpc.Response.trailers:type_name -> rpc.Response.Trailers
	6,  // 6: rpc.Request.Headers.headers:type_name -> rpc.Request.Headers.HeadersEntry
	7,  // 7: rpc.Request.Headers.query:type_name -> rpc.Request.Headers.QueryEntry
	0,  // 8: rpc.Request.Headers.HeadersEntry.value:type_name -> rpc.Values
	0,  // 9: rpc.Request.Headers.QueryEntry.value:type_name -> rpc.Values
	11, // 10: rpc.Response.Headers.headers:type_name -> rpc.Response.Headers.HeadersEntry
	0,  // 11: rpc.Response.Headers.HeadersEntry.value:type_name -> rpc.Values
	1,  // 12: rpc.GitlabAccess.MakeRequest:input_type -> rpc.Request
	2,  // 13: rpc.GitlabAccess.MakeRequest:output_type -> rpc.Response
	13, // [13:14] is the sub-list for method output_type
	12, // [12:13] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_internal_module_gitlab_access_rpc_rpc_proto_init() }
func file_internal_module_gitlab_access_rpc_rpc_proto_init() {
	if File_internal_module_gitlab_access_rpc_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Values); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request_Headers); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request_Data); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request_Trailers); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_Headers); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_Data); i {
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
		file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_Trailers); i {
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
	file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Request_Headers_)(nil),
		(*Request_Data_)(nil),
		(*Request_Trailers_)(nil),
	}
	file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Response_Headers_)(nil),
		(*Response_Data_)(nil),
		(*Response_Trailers_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_module_gitlab_access_rpc_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_module_gitlab_access_rpc_rpc_proto_goTypes,
		DependencyIndexes: file_internal_module_gitlab_access_rpc_rpc_proto_depIdxs,
		MessageInfos:      file_internal_module_gitlab_access_rpc_rpc_proto_msgTypes,
	}.Build()
	File_internal_module_gitlab_access_rpc_rpc_proto = out.File
	file_internal_module_gitlab_access_rpc_rpc_proto_rawDesc = nil
	file_internal_module_gitlab_access_rpc_rpc_proto_goTypes = nil
	file_internal_module_gitlab_access_rpc_rpc_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GitlabAccessClient is the client API for GitlabAccess service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GitlabAccessClient interface {
	MakeRequest(ctx context.Context, opts ...grpc.CallOption) (GitlabAccess_MakeRequestClient, error)
}

type gitlabAccessClient struct {
	cc grpc.ClientConnInterface
}

func NewGitlabAccessClient(cc grpc.ClientConnInterface) GitlabAccessClient {
	return &gitlabAccessClient{cc}
}

func (c *gitlabAccessClient) MakeRequest(ctx context.Context, opts ...grpc.CallOption) (GitlabAccess_MakeRequestClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GitlabAccess_serviceDesc.Streams[0], "/rpc.GitlabAccess/MakeRequest", opts...)
	if err != nil {
		return nil, err
	}
	x := &gitlabAccessMakeRequestClient{stream}
	return x, nil
}

type GitlabAccess_MakeRequestClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type gitlabAccessMakeRequestClient struct {
	grpc.ClientStream
}

func (x *gitlabAccessMakeRequestClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitlabAccessMakeRequestClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GitlabAccessServer is the server API for GitlabAccess service.
type GitlabAccessServer interface {
	MakeRequest(GitlabAccess_MakeRequestServer) error
}

// UnimplementedGitlabAccessServer can be embedded to have forward compatible implementations.
type UnimplementedGitlabAccessServer struct {
}

func (*UnimplementedGitlabAccessServer) MakeRequest(GitlabAccess_MakeRequestServer) error {
	return status.Errorf(codes.Unimplemented, "method MakeRequest not implemented")
}

func RegisterGitlabAccessServer(s *grpc.Server, srv GitlabAccessServer) {
	s.RegisterService(&_GitlabAccess_serviceDesc, srv)
}

func _GitlabAccess_MakeRequest_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitlabAccessServer).MakeRequest(&gitlabAccessMakeRequestServer{stream})
}

type GitlabAccess_MakeRequestServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type gitlabAccessMakeRequestServer struct {
	grpc.ServerStream
}

func (x *gitlabAccessMakeRequestServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitlabAccessMakeRequestServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _GitlabAccess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.GitlabAccess",
	HandlerType: (*GitlabAccessServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MakeRequest",
			Handler:       _GitlabAccess_MakeRequest_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internal/module/gitlab_access/rpc/rpc.proto",
}
