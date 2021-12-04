// protoc --go_out=. --go-grpc_out=. messenger.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: messenger.proto

package messenger

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

type SendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message        string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	ConnectionId   int64  `protobuf:"varint,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	SenderFaceId   string `protobuf:"bytes,3,opt,name=sender_face_id,json=senderFaceId,proto3" json:"sender_face_id,omitempty"`
	ReceiverFaceId string `protobuf:"bytes,4,opt,name=receiver_face_id,json=receiverFaceId,proto3" json:"receiver_face_id,omitempty"`
}

func (x *SendRequest) Reset() {
	*x = SendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messenger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRequest) ProtoMessage() {}

func (x *SendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRequest.ProtoReflect.Descriptor instead.
func (*SendRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{0}
}

func (x *SendRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *SendRequest) GetConnectionId() int64 {
	if x != nil {
		return x.ConnectionId
	}
	return 0
}

func (x *SendRequest) GetSenderFaceId() string {
	if x != nil {
		return x.SenderFaceId
	}
	return ""
}

func (x *SendRequest) GetReceiverFaceId() string {
	if x != nil {
		return x.ReceiverFaceId
	}
	return ""
}

type SendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result bool `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *SendResponse) Reset() {
	*x = SendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messenger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendResponse) ProtoMessage() {}

func (x *SendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendResponse.ProtoReflect.Descriptor instead.
func (*SendResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{1}
}

func (x *SendResponse) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

type ListenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnectionId   int64  `protobuf:"varint,1,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	SenderFaceId   string `protobuf:"bytes,2,opt,name=sender_face_id,json=senderFaceId,proto3" json:"sender_face_id,omitempty"`
	ReceiverFaceId string `protobuf:"bytes,3,opt,name=receiver_face_id,json=receiverFaceId,proto3" json:"receiver_face_id,omitempty"`
}

func (x *ListenRequest) Reset() {
	*x = ListenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messenger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenRequest) ProtoMessage() {}

func (x *ListenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenRequest.ProtoReflect.Descriptor instead.
func (*ListenRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{2}
}

func (x *ListenRequest) GetConnectionId() int64 {
	if x != nil {
		return x.ConnectionId
	}
	return 0
}

func (x *ListenRequest) GetSenderFaceId() string {
	if x != nil {
		return x.SenderFaceId
	}
	return ""
}

func (x *ListenRequest) GetReceiverFaceId() string {
	if x != nil {
		return x.ReceiverFaceId
	}
	return ""
}

type ListenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content        string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	ConnectionId   int64  `protobuf:"varint,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	Sender         string `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
	SenderFaceId   string `protobuf:"bytes,4,opt,name=sender_face_id,json=senderFaceId,proto3" json:"sender_face_id,omitempty"`
	ReceiverFaceId string `protobuf:"bytes,5,opt,name=receiver_face_id,json=receiverFaceId,proto3" json:"receiver_face_id,omitempty"`
	Timestamp      int64  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ListenResponse) Reset() {
	*x = ListenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messenger_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenResponse) ProtoMessage() {}

func (x *ListenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenResponse.ProtoReflect.Descriptor instead.
func (*ListenResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{3}
}

func (x *ListenResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *ListenResponse) GetConnectionId() int64 {
	if x != nil {
		return x.ConnectionId
	}
	return 0
}

func (x *ListenResponse) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *ListenResponse) GetSenderFaceId() string {
	if x != nil {
		return x.SenderFaceId
	}
	return ""
}

func (x *ListenResponse) GetReceiverFaceId() string {
	if x != nil {
		return x.ReceiverFaceId
	}
	return ""
}

func (x *ListenResponse) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_messenger_proto protoreflect.FileDescriptor

var file_messenger_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22, 0x9c, 0x01, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x5f, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x46, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x46, 0x61, 0x63, 0x65, 0x49, 0x64, 0x22, 0x26, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x84,
	0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f,
	0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x46, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x46,
	0x61, 0x63, 0x65, 0x49, 0x64, 0x22, 0xd5, 0x01, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12,
	0x24, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x46,
	0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x5f, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x46, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x32, 0x75, 0x0a,
	0x09, 0x4d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04, 0x53, 0x65,
	0x6e, 0x64, 0x12, 0x11, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x06, 0x4c,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x12, 0x13, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x42, 0x13, 0x5a, 0x11, 0x6d, 0x61, 0x6e, 0x79, 0x66, 0x61, 0x63, 0x65,
	0x2e, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_messenger_proto_rawDescOnce sync.Once
	file_messenger_proto_rawDescData = file_messenger_proto_rawDesc
)

func file_messenger_proto_rawDescGZIP() []byte {
	file_messenger_proto_rawDescOnce.Do(func() {
		file_messenger_proto_rawDescData = protoimpl.X.CompressGZIP(file_messenger_proto_rawDescData)
	})
	return file_messenger_proto_rawDescData
}

var file_messenger_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_messenger_proto_goTypes = []interface{}{
	(*SendRequest)(nil),    // 0: grpc.SendRequest
	(*SendResponse)(nil),   // 1: grpc.SendResponse
	(*ListenRequest)(nil),  // 2: grpc.ListenRequest
	(*ListenResponse)(nil), // 3: grpc.ListenResponse
}
var file_messenger_proto_depIdxs = []int32{
	0, // 0: grpc.Messenger.Send:input_type -> grpc.SendRequest
	2, // 1: grpc.Messenger.Listen:input_type -> grpc.ListenRequest
	1, // 2: grpc.Messenger.Send:output_type -> grpc.SendResponse
	3, // 3: grpc.Messenger.Listen:output_type -> grpc.ListenResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_messenger_proto_init() }
func file_messenger_proto_init() {
	if File_messenger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messenger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRequest); i {
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
		file_messenger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendResponse); i {
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
		file_messenger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenRequest); i {
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
		file_messenger_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenResponse); i {
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
			RawDescriptor: file_messenger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_messenger_proto_goTypes,
		DependencyIndexes: file_messenger_proto_depIdxs,
		MessageInfos:      file_messenger_proto_msgTypes,
	}.Build()
	File_messenger_proto = out.File
	file_messenger_proto_rawDesc = nil
	file_messenger_proto_goTypes = nil
	file_messenger_proto_depIdxs = nil
}