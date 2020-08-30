// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: credential.proto

package credential_pb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ApiCredential struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce string   `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Roles []string `protobuf:"bytes,2,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *ApiCredential) Reset() {
	*x = ApiCredential{}
	if protoimpl.UnsafeEnabled {
		mi := &file_credential_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiCredential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiCredential) ProtoMessage() {}

func (x *ApiCredential) ProtoReflect() protoreflect.Message {
	mi := &file_credential_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiCredential.ProtoReflect.Descriptor instead.
func (*ApiCredential) Descriptor() ([]byte, []int) {
	return file_credential_proto_rawDescGZIP(), []int{0}
}

func (x *ApiCredential) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *ApiCredential) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

var File_credential_proto protoreflect.FileDescriptor

var file_credential_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x22, 0x3b,
	0x0a, 0x0d, 0x41, 0x70, 0x69, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x12,
	0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x20, 0x5a, 0x1e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f,
	0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_credential_proto_rawDescOnce sync.Once
	file_credential_proto_rawDescData = file_credential_proto_rawDesc
)

func file_credential_proto_rawDescGZIP() []byte {
	file_credential_proto_rawDescOnce.Do(func() {
		file_credential_proto_rawDescData = protoimpl.X.CompressGZIP(file_credential_proto_rawDescData)
	})
	return file_credential_proto_rawDescData
}

var file_credential_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_credential_proto_goTypes = []interface{}{
	(*ApiCredential)(nil), // 0: credential.ApiCredential
}
var file_credential_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_credential_proto_init() }
func file_credential_proto_init() {
	if File_credential_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_credential_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiCredential); i {
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
			RawDescriptor: file_credential_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_credential_proto_goTypes,
		DependencyIndexes: file_credential_proto_depIdxs,
		MessageInfos:      file_credential_proto_msgTypes,
	}.Build()
	File_credential_proto = out.File
	file_credential_proto_rawDesc = nil
	file_credential_proto_goTypes = nil
	file_credential_proto_depIdxs = nil
}
