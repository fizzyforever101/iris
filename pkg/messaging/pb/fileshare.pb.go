// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v5.29.3
// source: fileshare.proto

/*
Example File Sharing Flow
1. The client sends a request to retrieve metadata for the file.
2. The server responds with information like the total file size, chunk size, chunk count, and available chunks.
3. After the client has retrieved the file metadata, it sends a request for specific chunks.

Notes: 
1. The server will use the FileMetadata to store the chunk information, including chunk size and availability.
2. The server will respond to FileDownloadRequest by checking if the requested chunks exist in the available_chunks list. If so, it will return the chunk data, along with the appropriate chunk index.
*/


package pb

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

type FileMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata        *MetaData `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Cid             string    `protobuf:"bytes,2,opt,name=cid,proto3" json:"cid,omitempty"`
	Description     []byte    `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Rights          []string  `protobuf:"bytes,4,rep,name=rights,proto3" json:"rights,omitempty"`
	Severity        string    `protobuf:"bytes,5,opt,name=severity,proto3" json:"severity,omitempty"`
	ExpiredAt       int64     `protobuf:"varint,6,opt,name=expiredAt,proto3" json:"expiredAt,omitempty"`
	TotalSize       int64     `protobuf:"varint,7,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`                           // field to store the total file size
	ChunkSize       int32     `protobuf:"varint,8,opt,name=chunk_size,json=chunkSize,proto3" json:"chunk_size,omitempty"`                           // Size of each chunk (if chunking is used)
	ChunkCount      int32     `protobuf:"varint,9,opt,name=chunk_count,json=chunkCount,proto3" json:"chunk_count,omitempty"`                        // Total number of chunks (based on total_size / chunk_size)
	AvailableChunks []int32   `protobuf:"varint,10,rep,packed,name=available_chunks,json=availableChunks,proto3" json:"available_chunks,omitempty"` // List of chunk indices that are available for download
}

func (x *FileMetadata) Reset() {
	*x = FileMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileshare_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileMetadata) ProtoMessage() {}

func (x *FileMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_fileshare_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileMetadata.ProtoReflect.Descriptor instead.
func (*FileMetadata) Descriptor() ([]byte, []int) {
	return file_fileshare_proto_rawDescGZIP(), []int{0}
}

func (x *FileMetadata) GetMetadata() *MetaData {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *FileMetadata) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

func (x *FileMetadata) GetDescription() []byte {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *FileMetadata) GetRights() []string {
	if x != nil {
		return x.Rights
	}
	return nil
}

func (x *FileMetadata) GetSeverity() string {
	if x != nil {
		return x.Severity
	}
	return ""
}

func (x *FileMetadata) GetExpiredAt() int64 {
	if x != nil {
		return x.ExpiredAt
	}
	return 0
}

func (x *FileMetadata) GetTotalSize() int64 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

func (x *FileMetadata) GetChunkSize() int32 {
	if x != nil {
		return x.ChunkSize
	}
	return 0
}

func (x *FileMetadata) GetChunkCount() int32 {
	if x != nil {
		return x.ChunkCount
	}
	return 0
}

func (x *FileMetadata) GetAvailableChunks() []int32 {
	if x != nil {
		return x.AvailableChunks
	}
	return nil
}

type FileDownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata     *MetaData `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Cid          string    `protobuf:"bytes,2,opt,name=cid,proto3" json:"cid,omitempty"`
	ChunkSize    int32     `protobuf:"varint,3,opt,name=chunk_size,json=chunkSize,proto3" json:"chunk_size,omitempty"`                 // field for chunk size
	ChunkIndices []int32   `protobuf:"varint,4,rep,packed,name=chunk_indices,json=chunkIndices,proto3" json:"chunk_indices,omitempty"` // list of chunk indices to download
}

func (x *FileDownloadRequest) Reset() {
	*x = FileDownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileshare_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDownloadRequest) ProtoMessage() {}

func (x *FileDownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fileshare_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDownloadRequest.ProtoReflect.Descriptor instead.
func (*FileDownloadRequest) Descriptor() ([]byte, []int) {
	return file_fileshare_proto_rawDescGZIP(), []int{1}
}

func (x *FileDownloadRequest) GetMetadata() *MetaData {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *FileDownloadRequest) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

func (x *FileDownloadRequest) GetChunkSize() int32 {
	if x != nil {
		return x.ChunkSize
	}
	return 0
}

func (x *FileDownloadRequest) GetChunkIndices() []int32 {
	if x != nil {
		return x.ChunkIndices
	}
	return nil
}

type FileDownloadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata   *MetaData `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Status     string    `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Data       []byte    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	ChunkIndex int32     `protobuf:"varint,4,opt,name=chunk_index,json=chunkIndex,proto3" json:"chunk_index,omitempty"` // To indicate which chunk is being returned
}

func (x *FileDownloadResponse) Reset() {
	*x = FileDownloadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileshare_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDownloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDownloadResponse) ProtoMessage() {}

func (x *FileDownloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileshare_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDownloadResponse.ProtoReflect.Descriptor instead.
func (*FileDownloadResponse) Descriptor() ([]byte, []int) {
	return file_fileshare_proto_rawDescGZIP(), []int{2}
}

func (x *FileDownloadResponse) GetMetadata() *MetaData {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *FileDownloadResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *FileDownloadResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *FileDownloadResponse) GetChunkIndex() int32 {
	if x != nil {
		return x.ChunkIndex
	}
	return 0
}

var File_fileshare_proto protoreflect.FileDescriptor

var file_fileshare_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xc8, 0x02, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x28, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03,
	0x63, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x69, 0x67, 0x68, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x69, 0x67, 0x68, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0f, 0x61, 0x76, 0x61,
	0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22, 0x95, 0x01, 0x0a,
	0x13, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10,
	0x0a, 0x03, 0x63, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x69, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x69, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0c, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x64,
	0x69, 0x63, 0x65, 0x73, 0x22, 0x8d, 0x01, 0x0a, 0x14, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x42, 0x14, 0x5a, 0x12, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_fileshare_proto_rawDescOnce sync.Once
	file_fileshare_proto_rawDescData = file_fileshare_proto_rawDesc
)

func file_fileshare_proto_rawDescGZIP() []byte {
	file_fileshare_proto_rawDescOnce.Do(func() {
		file_fileshare_proto_rawDescData = protoimpl.X.CompressGZIP(file_fileshare_proto_rawDescData)
	})
	return file_fileshare_proto_rawDescData
}

var file_fileshare_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_fileshare_proto_goTypes = []interface{}{
	(*FileMetadata)(nil),         // 0: pb.FileMetadata
	(*FileDownloadRequest)(nil),  // 1: pb.FileDownloadRequest
	(*FileDownloadResponse)(nil), // 2: pb.FileDownloadResponse
	(*MetaData)(nil),             // 3: pb.MetaData
}
var file_fileshare_proto_depIdxs = []int32{
	3, // 0: pb.FileMetadata.metadata:type_name -> pb.MetaData
	3, // 1: pb.FileDownloadRequest.metadata:type_name -> pb.MetaData
	3, // 2: pb.FileDownloadResponse.metadata:type_name -> pb.MetaData
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_fileshare_proto_init() }
func file_fileshare_proto_init() {
	if File_fileshare_proto != nil {
		return
	}
	file_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_fileshare_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileMetadata); i {
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
		file_fileshare_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDownloadRequest); i {
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
		file_fileshare_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDownloadResponse); i {
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
			RawDescriptor: file_fileshare_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fileshare_proto_goTypes,
		DependencyIndexes: file_fileshare_proto_depIdxs,
		MessageInfos:      file_fileshare_proto_msgTypes,
	}.Build()
	File_fileshare_proto = out.File
	file_fileshare_proto_rawDesc = nil
	file_fileshare_proto_goTypes = nil
	file_fileshare_proto_depIdxs = nil
}
