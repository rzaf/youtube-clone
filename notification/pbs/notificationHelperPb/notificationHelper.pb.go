// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: notificationHelper.proto

package notificationHelperPb

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

type SortType int32

const (
	SortType_Newest SortType = 0
	SortType_Oldest SortType = 1
)

// Enum value maps for SortType.
var (
	SortType_name = map[int32]string{
		0: "Newest",
		1: "Oldest",
	}
	SortType_value = map[string]int32{
		"Newest": 0,
		"Oldest": 1,
	}
)

func (x SortType) Enum() *SortType {
	p := new(SortType)
	*p = x
	return p
}

func (x SortType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SortType) Descriptor() protoreflect.EnumDescriptor {
	return file_notificationHelper_proto_enumTypes[0].Descriptor()
}

func (SortType) Type() protoreflect.EnumType {
	return &file_notificationHelper_proto_enumTypes[0]
}

func (x SortType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SortType.Descriptor instead.
func (SortType) EnumDescriptor() ([]byte, []int) {
	return file_notificationHelper_proto_rawDescGZIP(), []int{0}
}

type SeenType int32

const (
	SeenType_Any     SeenType = 0
	SeenType_Seen    SeenType = 1
	SeenType_NotSeen SeenType = 2
)

// Enum value maps for SeenType.
var (
	SeenType_name = map[int32]string{
		0: "Any",
		1: "Seen",
		2: "NotSeen",
	}
	SeenType_value = map[string]int32{
		"Any":     0,
		"Seen":    1,
		"NotSeen": 2,
	}
)

func (x SeenType) Enum() *SeenType {
	p := new(SeenType)
	*p = x
	return p
}

func (x SeenType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SeenType) Descriptor() protoreflect.EnumDescriptor {
	return file_notificationHelper_proto_enumTypes[1].Descriptor()
}

func (SeenType) Type() protoreflect.EnumType {
	return &file_notificationHelper_proto_enumTypes[1]
}

func (x SeenType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SeenType.Descriptor instead.
func (SeenType) EnumDescriptor() ([]byte, []int) {
	return file_notificationHelper_proto_rawDescGZIP(), []int{1}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notificationHelper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_notificationHelper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_notificationHelper_proto_rawDescGZIP(), []int{0}
}

type HttpError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message    string `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	StatusCode int32  `protobuf:"varint,2,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
}

func (x *HttpError) Reset() {
	*x = HttpError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notificationHelper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpError) ProtoMessage() {}

func (x *HttpError) ProtoReflect() protoreflect.Message {
	mi := &file_notificationHelper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpError.ProtoReflect.Descriptor instead.
func (*HttpError) Descriptor() ([]byte, []int) {
	return file_notificationHelper_proto_rawDescGZIP(), []int{1}
}

func (x *HttpError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *HttpError) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

type Paging struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PerPage    int32 `protobuf:"varint,1,opt,name=PerPage,proto3" json:"PerPage,omitempty"`       ///limit
	PageNumber int32 `protobuf:"varint,2,opt,name=pageNumber,proto3" json:"pageNumber,omitempty"` ///offset
}

func (x *Paging) Reset() {
	*x = Paging{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notificationHelper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Paging) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Paging) ProtoMessage() {}

func (x *Paging) ProtoReflect() protoreflect.Message {
	mi := &file_notificationHelper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Paging.ProtoReflect.Descriptor instead.
func (*Paging) Descriptor() ([]byte, []int) {
	return file_notificationHelper_proto_rawDescGZIP(), []int{2}
}

func (x *Paging) GetPerPage() int32 {
	if x != nil {
		return x.PerPage
	}
	return 0
}

func (x *Paging) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

type PagesInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentPage int32 `protobuf:"varint,1,opt,name=CurrentPage,proto3" json:"CurrentPage,omitempty"`
	TotalPages  int32 `protobuf:"varint,2,opt,name=TotalPages,proto3" json:"TotalPages,omitempty"`
}

func (x *PagesInfo) Reset() {
	*x = PagesInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notificationHelper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PagesInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagesInfo) ProtoMessage() {}

func (x *PagesInfo) ProtoReflect() protoreflect.Message {
	mi := &file_notificationHelper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagesInfo.ProtoReflect.Descriptor instead.
func (*PagesInfo) Descriptor() ([]byte, []int) {
	return file_notificationHelper_proto_rawDescGZIP(), []int{3}
}

func (x *PagesInfo) GetCurrentPage() int32 {
	if x != nil {
		return x.CurrentPage
	}
	return 0
}

func (x *PagesInfo) GetTotalPages() int32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

var File_notificationHelper_proto protoreflect.FileDescriptor

var file_notificationHelper_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65,
	0x6c, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x6c, 0x70, 0x65, 0x72, 0x22, 0x07,
	0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x45, 0x0a, 0x09, 0x48, 0x74, 0x74, 0x70, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x42,
	0x0a, 0x06, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x65, 0x72, 0x50,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x50, 0x65, 0x72, 0x50, 0x61,
	0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x22, 0x4d, 0x0a, 0x09, 0x50, 0x61, 0x67, 0x65, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x20, 0x0a, 0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x67,
	0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65,
	0x73, 0x2a, 0x22, 0x0a, 0x08, 0x53, 0x6f, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a,
	0x06, 0x4e, 0x65, 0x77, 0x65, 0x73, 0x74, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x6c, 0x64,
	0x65, 0x73, 0x74, 0x10, 0x01, 0x2a, 0x2a, 0x0a, 0x08, 0x53, 0x65, 0x65, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x6e, 0x79, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x65,
	0x65, 0x6e, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x6f, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x10,
	0x02, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x72, 0x7a, 0x61, 0x66, 0x2f, 0x79, 0x6f, 0x75, 0x74, 0x75, 0x62, 0x65, 0x2d, 0x63, 0x6c, 0x6f,
	0x6e, 0x65, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x70, 0x62, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x48, 0x65, 0x6c, 0x70, 0x65, 0x72, 0x50, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_notificationHelper_proto_rawDescOnce sync.Once
	file_notificationHelper_proto_rawDescData = file_notificationHelper_proto_rawDesc
)

func file_notificationHelper_proto_rawDescGZIP() []byte {
	file_notificationHelper_proto_rawDescOnce.Do(func() {
		file_notificationHelper_proto_rawDescData = protoimpl.X.CompressGZIP(file_notificationHelper_proto_rawDescData)
	})
	return file_notificationHelper_proto_rawDescData
}

var file_notificationHelper_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_notificationHelper_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_notificationHelper_proto_goTypes = []interface{}{
	(SortType)(0),     // 0: notificationHelper.SortType
	(SeenType)(0),     // 1: notificationHelper.SeenType
	(*Empty)(nil),     // 2: notificationHelper.Empty
	(*HttpError)(nil), // 3: notificationHelper.HttpError
	(*Paging)(nil),    // 4: notificationHelper.Paging
	(*PagesInfo)(nil), // 5: notificationHelper.PagesInfo
}
var file_notificationHelper_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_notificationHelper_proto_init() }
func file_notificationHelper_proto_init() {
	if File_notificationHelper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_notificationHelper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_notificationHelper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpError); i {
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
		file_notificationHelper_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Paging); i {
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
		file_notificationHelper_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PagesInfo); i {
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
			RawDescriptor: file_notificationHelper_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_notificationHelper_proto_goTypes,
		DependencyIndexes: file_notificationHelper_proto_depIdxs,
		EnumInfos:         file_notificationHelper_proto_enumTypes,
		MessageInfos:      file_notificationHelper_proto_msgTypes,
	}.Build()
	File_notificationHelper_proto = out.File
	file_notificationHelper_proto_rawDesc = nil
	file_notificationHelper_proto_goTypes = nil
	file_notificationHelper_proto_depIdxs = nil
}
