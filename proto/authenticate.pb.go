// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/authenticate.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TokenValidateRequest is the request message for the Authenticate method.
type TokenValidateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TokenValidateRequest) Reset() {
	*x = TokenValidateRequest{}
	mi := &file_proto_authenticate_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenValidateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenValidateRequest) ProtoMessage() {}

func (x *TokenValidateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_authenticate_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenValidateRequest.ProtoReflect.Descriptor instead.
func (*TokenValidateRequest) Descriptor() ([]byte, []int) {
	return file_proto_authenticate_proto_rawDescGZIP(), []int{0}
}

func (x *TokenValidateRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// TokenValidateResponse is the response message for the Authenticate method.
type TokenValidateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Valid         bool                   `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	Id            string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Role          string                 `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
	Email         string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	UserName      string                 `protobuf:"bytes,5,opt,name=userName,proto3" json:"userName,omitempty"`
	FullName      string                 `protobuf:"bytes,6,opt,name=fullName,proto3" json:"fullName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TokenValidateResponse) Reset() {
	*x = TokenValidateResponse{}
	mi := &file_proto_authenticate_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenValidateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenValidateResponse) ProtoMessage() {}

func (x *TokenValidateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_authenticate_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenValidateResponse.ProtoReflect.Descriptor instead.
func (*TokenValidateResponse) Descriptor() ([]byte, []int) {
	return file_proto_authenticate_proto_rawDescGZIP(), []int{1}
}

func (x *TokenValidateResponse) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

func (x *TokenValidateResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TokenValidateResponse) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *TokenValidateResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *TokenValidateResponse) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *TokenValidateResponse) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

var File_proto_authenticate_proto protoreflect.FileDescriptor

var file_proto_authenticate_proto_rawDesc = string([]byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x2c, 0x0a, 0x14, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x9f, 0x01, 0x0a, 0x15, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d,
	0x65, 0x32, 0x53, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x12, 0x43, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_authenticate_proto_rawDescOnce sync.Once
	file_proto_authenticate_proto_rawDescData []byte
)

func file_proto_authenticate_proto_rawDescGZIP() []byte {
	file_proto_authenticate_proto_rawDescOnce.Do(func() {
		file_proto_authenticate_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_authenticate_proto_rawDesc), len(file_proto_authenticate_proto_rawDesc)))
	})
	return file_proto_authenticate_proto_rawDescData
}

var file_proto_authenticate_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_authenticate_proto_goTypes = []any{
	(*TokenValidateRequest)(nil),  // 0: proto.TokenValidateRequest
	(*TokenValidateResponse)(nil), // 1: proto.TokenValidateResponse
}
var file_proto_authenticate_proto_depIdxs = []int32{
	0, // 0: proto.Authenticate.Auth:input_type -> proto.TokenValidateRequest
	1, // 1: proto.Authenticate.Auth:output_type -> proto.TokenValidateResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_authenticate_proto_init() }
func file_proto_authenticate_proto_init() {
	if File_proto_authenticate_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_authenticate_proto_rawDesc), len(file_proto_authenticate_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_authenticate_proto_goTypes,
		DependencyIndexes: file_proto_authenticate_proto_depIdxs,
		MessageInfos:      file_proto_authenticate_proto_msgTypes,
	}.Build()
	File_proto_authenticate_proto = out.File
	file_proto_authenticate_proto_goTypes = nil
	file_proto_authenticate_proto_depIdxs = nil
}
