// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: pagination.proto

package protoapi

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

type PageInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalCount      int64 `protobuf:"varint,1,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	HasNextPage     bool  `protobuf:"varint,2,opt,name=has_next_page,json=hasNextPage,proto3" json:"has_next_page,omitempty"`
	HasPreviousPage bool  `protobuf:"varint,3,opt,name=has_previous_page,json=hasPreviousPage,proto3" json:"has_previous_page,omitempty"`
}

func (x *PageInfo) Reset() {
	*x = PageInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PageInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageInfo) ProtoMessage() {}

func (x *PageInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pagination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageInfo.ProtoReflect.Descriptor instead.
func (*PageInfo) Descriptor() ([]byte, []int) {
	return file_pagination_proto_rawDescGZIP(), []int{0}
}

func (x *PageInfo) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *PageInfo) GetHasNextPage() bool {
	if x != nil {
		return x.HasNextPage
	}
	return false
}

func (x *PageInfo) GetHasPreviousPage() bool {
	if x != nil {
		return x.HasPreviousPage
	}
	return false
}

var File_pagination_proto protoreflect.FileDescriptor

var file_pagination_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x08, 0x50, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1f,
	0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x22, 0x0a, 0x0d, 0x68, 0x61, 0x73, 0x5f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x4e, 0x65, 0x78, 0x74, 0x50,
	0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x68, 0x61, 0x73, 0x5f, 0x70, 0x72, 0x65, 0x76, 0x69,
	0x6f, 0x75, 0x73, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f,
	0x68, 0x61, 0x73, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x50, 0x61, 0x67, 0x65, 0x42,
	0x10, 0x5a, 0x0e, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pagination_proto_rawDescOnce sync.Once
	file_pagination_proto_rawDescData = file_pagination_proto_rawDesc
)

func file_pagination_proto_rawDescGZIP() []byte {
	file_pagination_proto_rawDescOnce.Do(func() {
		file_pagination_proto_rawDescData = protoimpl.X.CompressGZIP(file_pagination_proto_rawDescData)
	})
	return file_pagination_proto_rawDescData
}

var file_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pagination_proto_goTypes = []interface{}{
	(*PageInfo)(nil), // 0: PageInfo
}
var file_pagination_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pagination_proto_init() }
func file_pagination_proto_init() {
	if File_pagination_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pagination_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PageInfo); i {
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
			RawDescriptor: file_pagination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pagination_proto_goTypes,
		DependencyIndexes: file_pagination_proto_depIdxs,
		MessageInfos:      file_pagination_proto_msgTypes,
	}.Build()
	File_pagination_proto = out.File
	file_pagination_proto_rawDesc = nil
	file_pagination_proto_goTypes = nil
	file_pagination_proto_depIdxs = nil
}
