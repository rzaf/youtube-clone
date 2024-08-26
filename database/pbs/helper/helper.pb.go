// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: youtube-clone/database/pbs/helper.proto

package helper

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
	SortType_MostViewed       SortType = 0
	SortType_LeastViewed      SortType = 1
	SortType_Newest           SortType = 3
	SortType_Oldest           SortType = 4
	SortType_MostLiked        SortType = 5
	SortType_LeastLiked       SortType = 6
	SortType_MostDisiked      SortType = 7
	SortType_LeastDisliked    SortType = 8
	SortType_MostSubscribers  SortType = 9
	SortType_LeastSubscribers SortType = 10
	SortType_MostReplied      SortType = 11
	SortType_LeastReplied     SortType = 12
)

// Enum value maps for SortType.
var (
	SortType_name = map[int32]string{
		0:  "MostViewed",
		1:  "LeastViewed",
		3:  "Newest",
		4:  "Oldest",
		5:  "MostLiked",
		6:  "LeastLiked",
		7:  "MostDisiked",
		8:  "LeastDisliked",
		9:  "MostSubscribers",
		10: "LeastSubscribers",
		11: "MostReplied",
		12: "LeastReplied",
	}
	SortType_value = map[string]int32{
		"MostViewed":       0,
		"LeastViewed":      1,
		"Newest":           3,
		"Oldest":           4,
		"MostLiked":        5,
		"LeastLiked":       6,
		"MostDisiked":      7,
		"LeastDisliked":    8,
		"MostSubscribers":  9,
		"LeastSubscribers": 10,
		"MostReplied":      11,
		"LeastReplied":     12,
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
	return file_youtube_clone_database_pbs_helper_proto_enumTypes[0].Descriptor()
}

func (SortType) Type() protoreflect.EnumType {
	return &file_youtube_clone_database_pbs_helper_proto_enumTypes[0]
}

func (x SortType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SortType.Descriptor instead.
func (SortType) EnumDescriptor() ([]byte, []int) {
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{0}
}

type MediaType int32

const (
	MediaType_VIDEO MediaType = 0
	MediaType_MUSIC MediaType = 1
	MediaType_PHOTO MediaType = 2
	MediaType_ALL   MediaType = 3
)

// Enum value maps for MediaType.
var (
	MediaType_name = map[int32]string{
		0: "VIDEO",
		1: "MUSIC",
		2: "PHOTO",
		3: "ALL",
	}
	MediaType_value = map[string]int32{
		"VIDEO": 0,
		"MUSIC": 1,
		"PHOTO": 2,
		"ALL":   3,
	}
)

func (x MediaType) Enum() *MediaType {
	p := new(MediaType)
	*p = x
	return p
}

func (x MediaType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MediaType) Descriptor() protoreflect.EnumDescriptor {
	return file_youtube_clone_database_pbs_helper_proto_enumTypes[1].Descriptor()
}

func (MediaType) Type() protoreflect.EnumType {
	return &file_youtube_clone_database_pbs_helper_proto_enumTypes[1]
}

func (x MediaType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MediaType.Descriptor instead.
func (MediaType) EnumDescriptor() ([]byte, []int) {
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{1}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[0]
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
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{0}
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
		mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpError) ProtoMessage() {}

func (x *HttpError) ProtoReflect() protoreflect.Message {
	mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[1]
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
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{1}
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
		mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Paging) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Paging) ProtoMessage() {}

func (x *Paging) ProtoReflect() protoreflect.Message {
	mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[2]
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
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{2}
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
		mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PagesInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagesInfo) ProtoMessage() {}

func (x *PagesInfo) ProtoReflect() protoreflect.Message {
	mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[3]
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
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{3}
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

type LikeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsLike bool   `protobuf:"varint,1,opt,name=isLike,proto3" json:"isLike,omitempty"`
	UserId int64  `protobuf:"varint,2,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Url    string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *LikeReq) Reset() {
	*x = LikeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LikeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeReq) ProtoMessage() {}

func (x *LikeReq) ProtoReflect() protoreflect.Message {
	mi := &file_youtube_clone_database_pbs_helper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeReq.ProtoReflect.Descriptor instead.
func (*LikeReq) Descriptor() ([]byte, []int) {
	return file_youtube_clone_database_pbs_helper_proto_rawDescGZIP(), []int{4}
}

func (x *LikeReq) GetIsLike() bool {
	if x != nil {
		return x.IsLike
	}
	return false
}

func (x *LikeReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *LikeReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_youtube_clone_database_pbs_helper_proto protoreflect.FileDescriptor

var file_youtube_clone_database_pbs_helper_proto_rawDesc = []byte{
	0x0a, 0x27, 0x79, 0x6f, 0x75, 0x74, 0x75, 0x62, 0x65, 0x2d, 0x63, 0x6c, 0x6f, 0x6e, 0x65, 0x2f,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x70, 0x62, 0x73, 0x2f, 0x68, 0x65, 0x6c,
	0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x68, 0x65, 0x6c, 0x70, 0x65,
	0x72, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x45, 0x0a, 0x09, 0x48, 0x74,
	0x74, 0x70, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x22, 0x42, 0x0a, 0x06, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x50,
	0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x50, 0x65,
	0x72, 0x50, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x4d, 0x0a, 0x09, 0x50, 0x61, 0x67, 0x65, 0x73, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x50, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50,
	0x61, 0x67, 0x65, 0x73, 0x22, 0x4b, 0x0a, 0x07, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x16, 0x0a, 0x06, 0x69, 0x73, 0x4c, 0x69, 0x6b, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x69, 0x73, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x2a, 0xd4, 0x01, 0x0a, 0x08, 0x53, 0x6f, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e,
	0x0a, 0x0a, 0x4d, 0x6f, 0x73, 0x74, 0x56, 0x69, 0x65, 0x77, 0x65, 0x64, 0x10, 0x00, 0x12, 0x0f,
	0x0a, 0x0b, 0x4c, 0x65, 0x61, 0x73, 0x74, 0x56, 0x69, 0x65, 0x77, 0x65, 0x64, 0x10, 0x01, 0x12,
	0x0a, 0x0a, 0x06, 0x4e, 0x65, 0x77, 0x65, 0x73, 0x74, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x4f,
	0x6c, 0x64, 0x65, 0x73, 0x74, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x4d, 0x6f, 0x73, 0x74, 0x4c,
	0x69, 0x6b, 0x65, 0x64, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x4c, 0x65, 0x61, 0x73, 0x74, 0x4c,
	0x69, 0x6b, 0x65, 0x64, 0x10, 0x06, 0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x6f, 0x73, 0x74, 0x44, 0x69,
	0x73, 0x69, 0x6b, 0x65, 0x64, 0x10, 0x07, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x65, 0x61, 0x73, 0x74,
	0x44, 0x69, 0x73, 0x6c, 0x69, 0x6b, 0x65, 0x64, 0x10, 0x08, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x6f,
	0x73, 0x74, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x73, 0x10, 0x09, 0x12,
	0x14, 0x0a, 0x10, 0x4c, 0x65, 0x61, 0x73, 0x74, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x72, 0x73, 0x10, 0x0a, 0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x69, 0x65, 0x64, 0x10, 0x0b, 0x12, 0x10, 0x0a, 0x0c, 0x4c, 0x65, 0x61, 0x73, 0x74, 0x52,
	0x65, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x10, 0x0c, 0x2a, 0x35, 0x0a, 0x09, 0x4d, 0x65, 0x64, 0x69,
	0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x10, 0x00,
	0x12, 0x09, 0x0a, 0x05, 0x4d, 0x55, 0x53, 0x49, 0x43, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x50,
	0x48, 0x4f, 0x54, 0x4f, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x4c, 0x4c, 0x10, 0x03, 0x42,
	0x23, 0x5a, 0x21, 0x79, 0x6f, 0x75, 0x74, 0x75, 0x62, 0x65, 0x2d, 0x63, 0x6c, 0x6f, 0x6e, 0x65,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x70, 0x62, 0x73, 0x2f, 0x68, 0x65,
	0x6c, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_youtube_clone_database_pbs_helper_proto_rawDescOnce sync.Once
	file_youtube_clone_database_pbs_helper_proto_rawDescData = file_youtube_clone_database_pbs_helper_proto_rawDesc
)

func file_youtube_clone_database_pbs_helper_proto_rawDescGZIP() []byte {
	file_youtube_clone_database_pbs_helper_proto_rawDescOnce.Do(func() {
		file_youtube_clone_database_pbs_helper_proto_rawDescData = protoimpl.X.CompressGZIP(file_youtube_clone_database_pbs_helper_proto_rawDescData)
	})
	return file_youtube_clone_database_pbs_helper_proto_rawDescData
}

var file_youtube_clone_database_pbs_helper_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_youtube_clone_database_pbs_helper_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_youtube_clone_database_pbs_helper_proto_goTypes = []interface{}{
	(SortType)(0),     // 0: helper.SortType
	(MediaType)(0),    // 1: helper.MediaType
	(*Empty)(nil),     // 2: helper.Empty
	(*HttpError)(nil), // 3: helper.HttpError
	(*Paging)(nil),    // 4: helper.Paging
	(*PagesInfo)(nil), // 5: helper.PagesInfo
	(*LikeReq)(nil),   // 6: helper.LikeReq
}
var file_youtube_clone_database_pbs_helper_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_youtube_clone_database_pbs_helper_proto_init() }
func file_youtube_clone_database_pbs_helper_proto_init() {
	if File_youtube_clone_database_pbs_helper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_youtube_clone_database_pbs_helper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_youtube_clone_database_pbs_helper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_youtube_clone_database_pbs_helper_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_youtube_clone_database_pbs_helper_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_youtube_clone_database_pbs_helper_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LikeReq); i {
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
			RawDescriptor: file_youtube_clone_database_pbs_helper_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_youtube_clone_database_pbs_helper_proto_goTypes,
		DependencyIndexes: file_youtube_clone_database_pbs_helper_proto_depIdxs,
		EnumInfos:         file_youtube_clone_database_pbs_helper_proto_enumTypes,
		MessageInfos:      file_youtube_clone_database_pbs_helper_proto_msgTypes,
	}.Build()
	File_youtube_clone_database_pbs_helper_proto = out.File
	file_youtube_clone_database_pbs_helper_proto_rawDesc = nil
	file_youtube_clone_database_pbs_helper_proto_goTypes = nil
	file_youtube_clone_database_pbs_helper_proto_depIdxs = nil
}
