// Code generated by protoc-gen-go. DO NOT EDIT.
// source: egress.proto

package loggregator_v2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EgressRequest struct {
	ShardId string  `protobuf:"bytes,1,opt,name=shard_id,json=shardId" json:"shard_id,omitempty"`
	Filter  *Filter `protobuf:"bytes,2,opt,name=filter" json:"filter,omitempty"`
	// TODO: This can be removed once the envelope.deprecated_tags is removed.
	UsePreferredTags bool `protobuf:"varint,3,opt,name=use_preferred_tags,json=usePreferredTags" json:"use_preferred_tags,omitempty"`
}

func (m *EgressRequest) Reset()                    { *m = EgressRequest{} }
func (m *EgressRequest) String() string            { return proto.CompactTextString(m) }
func (*EgressRequest) ProtoMessage()               {}
func (*EgressRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *EgressRequest) GetShardId() string {
	if m != nil {
		return m.ShardId
	}
	return ""
}

func (m *EgressRequest) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *EgressRequest) GetUsePreferredTags() bool {
	if m != nil {
		return m.UsePreferredTags
	}
	return false
}

type EgressBatchRequest struct {
	ShardId string  `protobuf:"bytes,1,opt,name=shard_id,json=shardId" json:"shard_id,omitempty"`
	Filter  *Filter `protobuf:"bytes,2,opt,name=filter" json:"filter,omitempty"`
	// TODO: This can be removed once the envelope.deprecated_tags is removed.
	UsePreferredTags bool `protobuf:"varint,3,opt,name=use_preferred_tags,json=usePreferredTags" json:"use_preferred_tags,omitempty"`
}

func (m *EgressBatchRequest) Reset()                    { *m = EgressBatchRequest{} }
func (m *EgressBatchRequest) String() string            { return proto.CompactTextString(m) }
func (*EgressBatchRequest) ProtoMessage()               {}
func (*EgressBatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *EgressBatchRequest) GetShardId() string {
	if m != nil {
		return m.ShardId
	}
	return ""
}

func (m *EgressBatchRequest) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *EgressBatchRequest) GetUsePreferredTags() bool {
	if m != nil {
		return m.UsePreferredTags
	}
	return false
}

type Filter struct {
	SourceId string `protobuf:"bytes,1,opt,name=source_id,json=sourceId" json:"source_id,omitempty"`
	// Types that are valid to be assigned to Message:
	//	*Filter_Log
	Message isFilter_Message `protobuf_oneof:"Message"`
}

func (m *Filter) Reset()                    { *m = Filter{} }
func (m *Filter) String() string            { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()               {}
func (*Filter) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

type isFilter_Message interface {
	isFilter_Message()
}

type Filter_Log struct {
	Log *LogFilter `protobuf:"bytes,2,opt,name=log,oneof"`
}

func (*Filter_Log) isFilter_Message() {}

func (m *Filter) GetMessage() isFilter_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Filter) GetSourceId() string {
	if m != nil {
		return m.SourceId
	}
	return ""
}

func (m *Filter) GetLog() *LogFilter {
	if x, ok := m.GetMessage().(*Filter_Log); ok {
		return x.Log
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Filter) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Filter_OneofMarshaler, _Filter_OneofUnmarshaler, _Filter_OneofSizer, []interface{}{
		(*Filter_Log)(nil),
	}
}

func _Filter_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Filter)
	// Message
	switch x := m.Message.(type) {
	case *Filter_Log:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Log); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Filter.Message has unexpected type %T", x)
	}
	return nil
}

func _Filter_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Filter)
	switch tag {
	case 2: // Message.log
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LogFilter)
		err := b.DecodeMessage(msg)
		m.Message = &Filter_Log{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Filter_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Filter)
	// Message
	switch x := m.Message.(type) {
	case *Filter_Log:
		s := proto.Size(x.Log)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type LogFilter struct {
}

func (m *LogFilter) Reset()                    { *m = LogFilter{} }
func (m *LogFilter) String() string            { return proto.CompactTextString(m) }
func (*LogFilter) ProtoMessage()               {}
func (*LogFilter) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func init() {
	proto.RegisterType((*EgressRequest)(nil), "loggregator.v2.EgressRequest")
	proto.RegisterType((*EgressBatchRequest)(nil), "loggregator.v2.EgressBatchRequest")
	proto.RegisterType((*Filter)(nil), "loggregator.v2.Filter")
	proto.RegisterType((*LogFilter)(nil), "loggregator.v2.LogFilter")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Egress service

type EgressClient interface {
	Receiver(ctx context.Context, in *EgressRequest, opts ...grpc.CallOption) (Egress_ReceiverClient, error)
	BatchedReceiver(ctx context.Context, in *EgressBatchRequest, opts ...grpc.CallOption) (Egress_BatchedReceiverClient, error)
}

type egressClient struct {
	cc *grpc.ClientConn
}

func NewEgressClient(cc *grpc.ClientConn) EgressClient {
	return &egressClient{cc}
}

func (c *egressClient) Receiver(ctx context.Context, in *EgressRequest, opts ...grpc.CallOption) (Egress_ReceiverClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Egress_serviceDesc.Streams[0], c.cc, "/loggregator.v2.Egress/Receiver", opts...)
	if err != nil {
		return nil, err
	}
	x := &egressReceiverClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Egress_ReceiverClient interface {
	Recv() (*Envelope, error)
	grpc.ClientStream
}

type egressReceiverClient struct {
	grpc.ClientStream
}

func (x *egressReceiverClient) Recv() (*Envelope, error) {
	m := new(Envelope)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *egressClient) BatchedReceiver(ctx context.Context, in *EgressBatchRequest, opts ...grpc.CallOption) (Egress_BatchedReceiverClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Egress_serviceDesc.Streams[1], c.cc, "/loggregator.v2.Egress/BatchedReceiver", opts...)
	if err != nil {
		return nil, err
	}
	x := &egressBatchedReceiverClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Egress_BatchedReceiverClient interface {
	Recv() (*EnvelopeBatch, error)
	grpc.ClientStream
}

type egressBatchedReceiverClient struct {
	grpc.ClientStream
}

func (x *egressBatchedReceiverClient) Recv() (*EnvelopeBatch, error) {
	m := new(EnvelopeBatch)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Egress service

type EgressServer interface {
	Receiver(*EgressRequest, Egress_ReceiverServer) error
	BatchedReceiver(*EgressBatchRequest, Egress_BatchedReceiverServer) error
}

func RegisterEgressServer(s *grpc.Server, srv EgressServer) {
	s.RegisterService(&_Egress_serviceDesc, srv)
}

func _Egress_Receiver_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EgressRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EgressServer).Receiver(m, &egressReceiverServer{stream})
}

type Egress_ReceiverServer interface {
	Send(*Envelope) error
	grpc.ServerStream
}

type egressReceiverServer struct {
	grpc.ServerStream
}

func (x *egressReceiverServer) Send(m *Envelope) error {
	return x.ServerStream.SendMsg(m)
}

func _Egress_BatchedReceiver_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EgressBatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EgressServer).BatchedReceiver(m, &egressBatchedReceiverServer{stream})
}

type Egress_BatchedReceiverServer interface {
	Send(*EnvelopeBatch) error
	grpc.ServerStream
}

type egressBatchedReceiverServer struct {
	grpc.ServerStream
}

func (x *egressBatchedReceiverServer) Send(m *EnvelopeBatch) error {
	return x.ServerStream.SendMsg(m)
}

var _Egress_serviceDesc = grpc.ServiceDesc{
	ServiceName: "loggregator.v2.Egress",
	HandlerType: (*EgressServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Receiver",
			Handler:       _Egress_Receiver_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BatchedReceiver",
			Handler:       _Egress_BatchedReceiver_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "egress.proto",
}

func init() { proto.RegisterFile("egress.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x92, 0xdf, 0x4a, 0xc3, 0x30,
	0x14, 0x87, 0x17, 0x07, 0x5d, 0x7b, 0xa6, 0x53, 0x72, 0x21, 0xdd, 0x64, 0x30, 0x7a, 0xb5, 0x0b,
	0x2d, 0x32, 0xdf, 0x60, 0xe0, 0x9f, 0x81, 0x82, 0x04, 0x2f, 0xbc, 0x2b, 0xb1, 0x39, 0xcb, 0x0a,
	0xc5, 0xd4, 0x24, 0xed, 0x33, 0x78, 0xe3, 0x93, 0xf8, 0x92, 0xb2, 0xa4, 0x4e, 0x37, 0xf5, 0xda,
	0xcb, 0xe4, 0x77, 0xce, 0xc7, 0x97, 0x73, 0x02, 0xfb, 0x28, 0x35, 0x1a, 0x93, 0x56, 0x5a, 0x59,
	0x45, 0x07, 0xa5, 0x92, 0x52, 0xa3, 0xe4, 0x56, 0xe9, 0xb4, 0x99, 0x8d, 0x06, 0xf8, 0xdc, 0x60,
	0xa9, 0x2a, 0xf4, 0x79, 0xf2, 0x4a, 0xe0, 0xe0, 0xd2, 0x35, 0x30, 0x7c, 0xa9, 0xd1, 0x58, 0x3a,
	0x84, 0xd0, 0xac, 0xb8, 0x16, 0x59, 0x21, 0x62, 0x32, 0x21, 0xd3, 0x88, 0xf5, 0xdc, 0x79, 0x21,
	0x68, 0x0a, 0xc1, 0xb2, 0x28, 0x2d, 0xea, 0x78, 0x6f, 0x42, 0xa6, 0xfd, 0xd9, 0x71, 0xba, 0x4d,
	0x4f, 0xaf, 0x5c, 0xca, 0xda, 0x2a, 0x7a, 0x0a, 0xb4, 0x36, 0x98, 0x55, 0x1a, 0x97, 0xa8, 0x35,
	0x8a, 0xcc, 0x72, 0x69, 0xe2, 0xee, 0x84, 0x4c, 0x43, 0x76, 0x54, 0x1b, 0xbc, 0xff, 0x0c, 0x1e,
	0xb8, 0x34, 0xc9, 0x1b, 0x01, 0xea, 0x55, 0xe6, 0xdc, 0xe6, 0xab, 0x7f, 0xf7, 0xc9, 0x20, 0xf0,
	0xfd, 0xf4, 0x04, 0x22, 0xa3, 0x6a, 0x9d, 0xe3, 0x97, 0x43, 0xe8, 0x2f, 0x16, 0x82, 0x9e, 0x41,
	0xb7, 0x54, 0xb2, 0x35, 0x18, 0xee, 0x1a, 0xdc, 0x2a, 0xe9, 0x21, 0x37, 0x1d, 0xb6, 0xae, 0x9b,
	0x47, 0xd0, 0xbb, 0x43, 0x63, 0xb8, 0xc4, 0xa4, 0x0f, 0xd1, 0x26, 0x9e, 0xbd, 0x13, 0x08, 0xfc,
	0xeb, 0xe9, 0x35, 0x84, 0x0c, 0x73, 0x2c, 0x1a, 0xd4, 0x74, 0xbc, 0x0b, 0xdc, 0x5a, 0xd6, 0x28,
	0xfe, 0x11, 0xb7, 0xeb, 0x4d, 0x3a, 0xe7, 0x84, 0x3e, 0xc2, 0xa1, 0x1b, 0x25, 0x8a, 0x0d, 0x2f,
	0xf9, 0x9d, 0xf7, 0x7d, 0xe2, 0xa3, 0xf1, 0x5f, 0x50, 0x57, 0xb5, 0x26, 0x3f, 0x05, 0xee, 0xf7,
	0x5c, 0x7c, 0x04, 0x00, 0x00, 0xff, 0xff, 0x25, 0xf7, 0xac, 0x98, 0x6d, 0x02, 0x00, 0x00,
}
