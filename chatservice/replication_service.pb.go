// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: chatservice/replication_service.proto

package chatservice

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

type EphemeralState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnectedClients map[string]*User `protobuf:"bytes,1,rep,name=connected_clients,json=connectedClients,proto3" json:"connected_clients,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EphemeralState) Reset() {
	*x = EphemeralState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatservice_replication_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EphemeralState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EphemeralState) ProtoMessage() {}

func (x *EphemeralState) ProtoReflect() protoreflect.Message {
	mi := &file_chatservice_replication_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EphemeralState.ProtoReflect.Descriptor instead.
func (*EphemeralState) Descriptor() ([]byte, []int) {
	return file_chatservice_replication_service_proto_rawDescGZIP(), []int{0}
}

func (x *EphemeralState) GetConnectedClients() map[string]*User {
	if x != nil {
		return x.ConnectedClients
	}
	return nil
}

type MessageAppend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageUuid string `protobuf:"bytes,1,opt,name=message_uuid,json=messageUuid,proto3" json:"message_uuid,omitempty"`
	AuthorId    string `protobuf:"bytes,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Body        string `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *MessageAppend) Reset() {
	*x = MessageAppend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatservice_replication_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageAppend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageAppend) ProtoMessage() {}

func (x *MessageAppend) ProtoReflect() protoreflect.Message {
	mi := &file_chatservice_replication_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageAppend.ProtoReflect.Descriptor instead.
func (*MessageAppend) Descriptor() ([]byte, []int) {
	return file_chatservice_replication_service_proto_rawDescGZIP(), []int{1}
}

func (x *MessageAppend) GetMessageUuid() string {
	if x != nil {
		return x.MessageUuid
	}
	return ""
}

func (x *MessageAppend) GetAuthorId() string {
	if x != nil {
		return x.AuthorId
	}
	return ""
}

func (x *MessageAppend) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type MessageLike struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageUuid string `protobuf:"bytes,1,opt,name=message_uuid,json=messageUuid,proto3" json:"message_uuid,omitempty"`
	LikerId     string `protobuf:"bytes,2,opt,name=liker_id,json=likerId,proto3" json:"liker_id,omitempty"`
	Like        bool   `protobuf:"varint,3,opt,name=like,proto3" json:"like,omitempty"`
}

func (x *MessageLike) Reset() {
	*x = MessageLike{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatservice_replication_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageLike) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageLike) ProtoMessage() {}

func (x *MessageLike) ProtoReflect() protoreflect.Message {
	mi := &file_chatservice_replication_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageLike.ProtoReflect.Descriptor instead.
func (*MessageLike) Descriptor() ([]byte, []int) {
	return file_chatservice_replication_service_proto_rawDescGZIP(), []int{2}
}

func (x *MessageLike) GetMessageUuid() string {
	if x != nil {
		return x.MessageUuid
	}
	return ""
}

func (x *MessageLike) GetLikerId() string {
	if x != nil {
		return x.LikerId
	}
	return ""
}

func (x *MessageLike) GetLike() bool {
	if x != nil {
		return x.Like
	}
	return false
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pid              string `protobuf:"bytes,1,opt,name=pid,proto3" json:"pid,omitempty"`
	SequenceNumber   int64  `protobuf:"varint,2,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"`
	LamportTimestamp int64  `protobuf:"varint,3,opt,name=lamport_timestamp,json=lamportTimestamp,proto3" json:"lamport_timestamp,omitempty"`
	ChatroomId       string `protobuf:"bytes,4,opt,name=chatroom_id,json=chatroomId,proto3" json:"chatroom_id,omitempty"`
	// Types that are assignable to Event:
	//
	//	*Event_MessageAppend
	//	*Event_MessageLike
	Event isEvent_Event `protobuf_oneof:"event"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatservice_replication_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_chatservice_replication_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_chatservice_replication_service_proto_rawDescGZIP(), []int{3}
}

func (x *Event) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *Event) GetSequenceNumber() int64 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *Event) GetLamportTimestamp() int64 {
	if x != nil {
		return x.LamportTimestamp
	}
	return 0
}

func (x *Event) GetChatroomId() string {
	if x != nil {
		return x.ChatroomId
	}
	return ""
}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *Event) GetMessageAppend() *MessageAppend {
	if x, ok := x.GetEvent().(*Event_MessageAppend); ok {
		return x.MessageAppend
	}
	return nil
}

func (x *Event) GetMessageLike() *MessageLike {
	if x, ok := x.GetEvent().(*Event_MessageLike); ok {
		return x.MessageLike
	}
	return nil
}

type isEvent_Event interface {
	isEvent_Event()
}

type Event_MessageAppend struct {
	MessageAppend *MessageAppend `protobuf:"bytes,5,opt,name=message_append,json=messageAppend,proto3,oneof"`
}

type Event_MessageLike struct {
	MessageLike *MessageLike `protobuf:"bytes,6,opt,name=message_like,json=messageLike,proto3,oneof"`
}

func (*Event_MessageAppend) isEvent_Event() {}

func (*Event_MessageLike) isEvent_Event() {}

type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The sequence number vector is a map from a process' PID to the next-expected
	// event sequence number from that process.
	SequenceNumberVector map[string]int64 `protobuf:"bytes,1,rep,name=sequence_number_vector,json=sequenceNumberVector,proto3" json:"sequence_number_vector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatservice_replication_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chatservice_replication_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_chatservice_replication_service_proto_rawDescGZIP(), []int{4}
}

func (x *SubscribeRequest) GetSequenceNumberVector() map[string]int64 {
	if x != nil {
		return x.SequenceNumberVector
	}
	return nil
}

type SubscriptionUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EphemeralState *EphemeralState `protobuf:"bytes,1,opt,name=ephemeral_state,json=ephemeralState,proto3" json:"ephemeral_state,omitempty"`
	Events         []*Event        `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
	// The garbage collected to vector is a map from a process' PID to the maximum
	// sequence number from this PID where garbage collection ran. This is useful
	// to assure the subscriber the holes in the sequences numbers below this point
	// are because of garbage collection.
	GarbageCollectedToVector map[string]int64 `protobuf:"bytes,3,rep,name=garbage_collected_to_vector,json=garbageCollectedToVector,proto3" json:"garbage_collected_to_vector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *SubscriptionUpdate) Reset() {
	*x = SubscriptionUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatservice_replication_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscriptionUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscriptionUpdate) ProtoMessage() {}

func (x *SubscriptionUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_chatservice_replication_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscriptionUpdate.ProtoReflect.Descriptor instead.
func (*SubscriptionUpdate) Descriptor() ([]byte, []int) {
	return file_chatservice_replication_service_proto_rawDescGZIP(), []int{5}
}

func (x *SubscriptionUpdate) GetEphemeralState() *EphemeralState {
	if x != nil {
		return x.EphemeralState
	}
	return nil
}

func (x *SubscriptionUpdate) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

func (x *SubscriptionUpdate) GetGarbageCollectedToVector() map[string]int64 {
	if x != nil {
		return x.GarbageCollectedToVector
	}
	return nil
}

var File_chatservice_replication_service_proto protoreflect.FileDescriptor

var file_chatservice_replication_service_proto_rawDesc = []byte{
	0x0a, 0x25, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x72, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x1a, 0x1e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x01, 0x0a, 0x0e, 0x45, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72,
	0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x5e, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x31, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x45, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x56, 0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x63, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x55,
	0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x22, 0x5f, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c,
	0x69, 0x6b, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x55, 0x75, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x69, 0x6b, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x69, 0x6b, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6b, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x04, 0x6c, 0x69, 0x6b, 0x65, 0x22, 0x9d, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69,
	0x64, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x73, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x11, 0x6c, 0x61,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x68, 0x61, 0x74, 0x72,
	0x6f, 0x6f, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x68,
	0x61, 0x74, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x43, 0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x0d,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x12, 0x3d, 0x0a,
	0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x6c, 0x69, 0x6b, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x48, 0x00, 0x52,
	0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x42, 0x07, 0x0a, 0x05,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0xca, 0x01, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x6d, 0x0a, 0x16, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x14, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a, 0x47, 0x0a, 0x19, 0x53, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x56, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0xd1, 0x02, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x44, 0x0a, 0x0f, 0x65, 0x70, 0x68,
	0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x45, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x0e, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x2a, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x7c, 0x0a, 0x1b, 0x67,
	0x61, 0x72, 0x62, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64,
	0x5f, 0x74, 0x6f, 0x5f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x3d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x47, 0x61, 0x72, 0x62, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x54, 0x6f, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x18, 0x67, 0x61, 0x72, 0x62, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65,
	0x64, 0x54, 0x6f, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a, 0x4b, 0x0a, 0x1d, 0x47, 0x61, 0x72,
	0x62, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x54, 0x6f, 0x56,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x6d, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x11,
	0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x73, 0x12, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x52, 0x69, 0x63, 0x68, 0x69, 0x65, 0x37, 0x38, 0x33, 0x32, 0x31, 0x2f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chatservice_replication_service_proto_rawDescOnce sync.Once
	file_chatservice_replication_service_proto_rawDescData = file_chatservice_replication_service_proto_rawDesc
)

func file_chatservice_replication_service_proto_rawDescGZIP() []byte {
	file_chatservice_replication_service_proto_rawDescOnce.Do(func() {
		file_chatservice_replication_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_chatservice_replication_service_proto_rawDescData)
	})
	return file_chatservice_replication_service_proto_rawDescData
}

var file_chatservice_replication_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_chatservice_replication_service_proto_goTypes = []interface{}{
	(*EphemeralState)(nil),     // 0: chatservice.EphemeralState
	(*MessageAppend)(nil),      // 1: chatservice.MessageAppend
	(*MessageLike)(nil),        // 2: chatservice.MessageLike
	(*Event)(nil),              // 3: chatservice.Event
	(*SubscribeRequest)(nil),   // 4: chatservice.SubscribeRequest
	(*SubscriptionUpdate)(nil), // 5: chatservice.SubscriptionUpdate
	nil,                        // 6: chatservice.EphemeralState.ConnectedClientsEntry
	nil,                        // 7: chatservice.SubscribeRequest.SequenceNumberVectorEntry
	nil,                        // 8: chatservice.SubscriptionUpdate.GarbageCollectedToVectorEntry
	(*User)(nil),               // 9: chatservice.User
}
var file_chatservice_replication_service_proto_depIdxs = []int32{
	6, // 0: chatservice.EphemeralState.connected_clients:type_name -> chatservice.EphemeralState.ConnectedClientsEntry
	1, // 1: chatservice.Event.message_append:type_name -> chatservice.MessageAppend
	2, // 2: chatservice.Event.message_like:type_name -> chatservice.MessageLike
	7, // 3: chatservice.SubscribeRequest.sequence_number_vector:type_name -> chatservice.SubscribeRequest.SequenceNumberVectorEntry
	0, // 4: chatservice.SubscriptionUpdate.ephemeral_state:type_name -> chatservice.EphemeralState
	3, // 5: chatservice.SubscriptionUpdate.events:type_name -> chatservice.Event
	8, // 6: chatservice.SubscriptionUpdate.garbage_collected_to_vector:type_name -> chatservice.SubscriptionUpdate.GarbageCollectedToVectorEntry
	9, // 7: chatservice.EphemeralState.ConnectedClientsEntry.value:type_name -> chatservice.User
	4, // 8: chatservice.ReplicationService.subscribe_updates:input_type -> chatservice.SubscribeRequest
	5, // 9: chatservice.ReplicationService.subscribe_updates:output_type -> chatservice.SubscriptionUpdate
	9, // [9:10] is the sub-list for method output_type
	8, // [8:9] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_chatservice_replication_service_proto_init() }
func file_chatservice_replication_service_proto_init() {
	if File_chatservice_replication_service_proto != nil {
		return
	}
	file_chatservice_chat_service_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_chatservice_replication_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EphemeralState); i {
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
		file_chatservice_replication_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageAppend); i {
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
		file_chatservice_replication_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageLike); i {
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
		file_chatservice_replication_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_chatservice_replication_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeRequest); i {
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
		file_chatservice_replication_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscriptionUpdate); i {
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
	file_chatservice_replication_service_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Event_MessageAppend)(nil),
		(*Event_MessageLike)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chatservice_replication_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chatservice_replication_service_proto_goTypes,
		DependencyIndexes: file_chatservice_replication_service_proto_depIdxs,
		MessageInfos:      file_chatservice_replication_service_proto_msgTypes,
	}.Build()
	File_chatservice_replication_service_proto = out.File
	file_chatservice_replication_service_proto_rawDesc = nil
	file_chatservice_replication_service_proto_goTypes = nil
	file_chatservice_replication_service_proto_depIdxs = nil
}
