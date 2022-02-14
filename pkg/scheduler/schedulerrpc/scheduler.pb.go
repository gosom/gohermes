// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: pkg/scheduler/schedulerrpc/scheduler.proto

package schedulerrpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ScheduledJob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Endpoint    string                 `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Data        []byte                 `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	ScheduledAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=scheduled_at,json=scheduledAt,proto3" json:"scheduled_at,omitempty"`
}

func (x *ScheduledJob) Reset() {
	*x = ScheduledJob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScheduledJob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduledJob) ProtoMessage() {}

func (x *ScheduledJob) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduledJob.ProtoReflect.Descriptor instead.
func (*ScheduledJob) Descriptor() ([]byte, []int) {
	return file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescGZIP(), []int{0}
}

func (x *ScheduledJob) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ScheduledJob) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *ScheduledJob) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ScheduledJob) GetScheduledAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ScheduledAt
	}
	return nil
}

type CreateScheduledJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job *ScheduledJob `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
}

func (x *CreateScheduledJobRequest) Reset() {
	*x = CreateScheduledJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateScheduledJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScheduledJobRequest) ProtoMessage() {}

func (x *CreateScheduledJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScheduledJobRequest.ProtoReflect.Descriptor instead.
func (*CreateScheduledJobRequest) Descriptor() ([]byte, []int) {
	return file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescGZIP(), []int{1}
}

func (x *CreateScheduledJobRequest) GetJob() *ScheduledJob {
	if x != nil {
		return x.Job
	}
	return nil
}

type CreateScheduledJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job *ScheduledJob `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
}

func (x *CreateScheduledJobResponse) Reset() {
	*x = CreateScheduledJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateScheduledJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScheduledJobResponse) ProtoMessage() {}

func (x *CreateScheduledJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScheduledJobResponse.ProtoReflect.Descriptor instead.
func (*CreateScheduledJobResponse) Descriptor() ([]byte, []int) {
	return file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescGZIP(), []int{2}
}

func (x *CreateScheduledJobResponse) GetJob() *ScheduledJob {
	if x != nil {
		return x.Job
	}
	return nil
}

var File_pkg_scheduler_schedulerrpc_scheduler_proto protoreflect.FileDescriptor

var file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2f,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8d, 0x01, 0x0a, 0x0c, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3d, 0x0a, 0x0c, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x41, 0x74, 0x22, 0x46, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x52, 0x03, 0x6a, 0x6f, 0x62,
	0x22, 0x47, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29,
	0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x64, 0x4a, 0x6f, 0x62, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x32, 0x78, 0x0a, 0x13, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x61, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x12, 0x24, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x64, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x73,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x67, 0x6f, 0x73, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73,
	0x2f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x3b, 0x73, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x72, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescOnce sync.Once
	file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescData = file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDesc
)

func file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescGZIP() []byte {
	file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescOnce.Do(func() {
		file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescData)
	})
	return file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDescData
}

var file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_scheduler_schedulerrpc_scheduler_proto_goTypes = []interface{}{
	(*ScheduledJob)(nil),               // 0: scheduler.ScheduledJob
	(*CreateScheduledJobRequest)(nil),  // 1: scheduler.CreateScheduledJobRequest
	(*CreateScheduledJobResponse)(nil), // 2: scheduler.CreateScheduledJobResponse
	(*timestamppb.Timestamp)(nil),      // 3: google.protobuf.Timestamp
}
var file_pkg_scheduler_schedulerrpc_scheduler_proto_depIdxs = []int32{
	3, // 0: scheduler.ScheduledJob.scheduled_at:type_name -> google.protobuf.Timestamp
	0, // 1: scheduler.CreateScheduledJobRequest.job:type_name -> scheduler.ScheduledJob
	0, // 2: scheduler.CreateScheduledJobResponse.job:type_name -> scheduler.ScheduledJob
	1, // 3: scheduler.ScheduledJobService.CreateScheduledJob:input_type -> scheduler.CreateScheduledJobRequest
	2, // 4: scheduler.ScheduledJobService.CreateScheduledJob:output_type -> scheduler.CreateScheduledJobResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pkg_scheduler_schedulerrpc_scheduler_proto_init() }
func file_pkg_scheduler_schedulerrpc_scheduler_proto_init() {
	if File_pkg_scheduler_schedulerrpc_scheduler_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScheduledJob); i {
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
		file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateScheduledJobRequest); i {
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
		file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateScheduledJobResponse); i {
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
			RawDescriptor: file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_scheduler_schedulerrpc_scheduler_proto_goTypes,
		DependencyIndexes: file_pkg_scheduler_schedulerrpc_scheduler_proto_depIdxs,
		MessageInfos:      file_pkg_scheduler_schedulerrpc_scheduler_proto_msgTypes,
	}.Build()
	File_pkg_scheduler_schedulerrpc_scheduler_proto = out.File
	file_pkg_scheduler_schedulerrpc_scheduler_proto_rawDesc = nil
	file_pkg_scheduler_schedulerrpc_scheduler_proto_goTypes = nil
	file_pkg_scheduler_schedulerrpc_scheduler_proto_depIdxs = nil
}