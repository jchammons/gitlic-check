// Code generated by protoc-gen-go. DO NOT EDIT.
// source: collector.proto

package collector

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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResultCode int32

const (
	ResultCode_OK              ResultCode = 0
	ResultCode_TRY_LATER       ResultCode = 1
	ResultCode_INVALID_API_KEY ResultCode = 2
	ResultCode_LIMIT_EXCEEDED  ResultCode = 3
	ResultCode_REDIRECT        ResultCode = 4
)

var ResultCode_name = map[int32]string{
	0: "OK",
	1: "TRY_LATER",
	2: "INVALID_API_KEY",
	3: "LIMIT_EXCEEDED",
	4: "REDIRECT",
}
var ResultCode_value = map[string]int32{
	"OK":              0,
	"TRY_LATER":       1,
	"INVALID_API_KEY": 2,
	"LIMIT_EXCEEDED":  3,
	"REDIRECT":        4,
}

func (x ResultCode) String() string {
	return proto.EnumName(ResultCode_name, int32(x))
}
func (ResultCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{0}
}

type EncodingType int32

const (
	EncodingType_BSON     EncodingType = 0
	EncodingType_PROTOBUF EncodingType = 1
)

var EncodingType_name = map[int32]string{
	0: "BSON",
	1: "PROTOBUF",
}
var EncodingType_value = map[string]int32{
	"BSON":     0,
	"PROTOBUF": 1,
}

func (x EncodingType) String() string {
	return proto.EnumName(EncodingType_name, int32(x))
}
func (EncodingType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{1}
}

type OboeSettingType int32

const (
	OboeSettingType_DEFAULT_SAMPLE_RATE        OboeSettingType = 0
	OboeSettingType_LAYER_SAMPLE_RATE          OboeSettingType = 1
	OboeSettingType_LAYER_APP_SAMPLE_RATE      OboeSettingType = 2
	OboeSettingType_LAYER_HTTPHOST_SAMPLE_RATE OboeSettingType = 3
	OboeSettingType_CONFIG_STRING              OboeSettingType = 4
	OboeSettingType_CONFIG_INT                 OboeSettingType = 5
)

var OboeSettingType_name = map[int32]string{
	0: "DEFAULT_SAMPLE_RATE",
	1: "LAYER_SAMPLE_RATE",
	2: "LAYER_APP_SAMPLE_RATE",
	3: "LAYER_HTTPHOST_SAMPLE_RATE",
	4: "CONFIG_STRING",
	5: "CONFIG_INT",
}
var OboeSettingType_value = map[string]int32{
	"DEFAULT_SAMPLE_RATE":        0,
	"LAYER_SAMPLE_RATE":          1,
	"LAYER_APP_SAMPLE_RATE":      2,
	"LAYER_HTTPHOST_SAMPLE_RATE": 3,
	"CONFIG_STRING":              4,
	"CONFIG_INT":                 5,
}

func (x OboeSettingType) String() string {
	return proto.EnumName(OboeSettingType_name, int32(x))
}
func (OboeSettingType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{2}
}

type HostID struct {
	Hostname             string   `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	IpAddresses          []string `protobuf:"bytes,2,rep,name=ip_addresses,json=ipAddresses,proto3" json:"ip_addresses,omitempty"`
	Uuid                 string   `protobuf:"bytes,3,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Pid                  int32    `protobuf:"varint,4,opt,name=pid,proto3" json:"pid,omitempty"`
	Ec2InstanceID        string   `protobuf:"bytes,5,opt,name=ec2InstanceID,proto3" json:"ec2InstanceID,omitempty"`
	Ec2AvailabilityZone  string   `protobuf:"bytes,6,opt,name=ec2AvailabilityZone,proto3" json:"ec2AvailabilityZone,omitempty"`
	DockerContainerID    string   `protobuf:"bytes,7,opt,name=dockerContainerID,proto3" json:"dockerContainerID,omitempty"`
	MacAddresses         []string `protobuf:"bytes,8,rep,name=macAddresses,proto3" json:"macAddresses,omitempty"`
	HerokuDynoID         string   `protobuf:"bytes,9,opt,name=herokuDynoID,proto3" json:"herokuDynoID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HostID) Reset()         { *m = HostID{} }
func (m *HostID) String() string { return proto.CompactTextString(m) }
func (*HostID) ProtoMessage()    {}
func (*HostID) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{0}
}
func (m *HostID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostID.Unmarshal(m, b)
}
func (m *HostID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostID.Marshal(b, m, deterministic)
}
func (dst *HostID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostID.Merge(dst, src)
}
func (m *HostID) XXX_Size() int {
	return xxx_messageInfo_HostID.Size(m)
}
func (m *HostID) XXX_DiscardUnknown() {
	xxx_messageInfo_HostID.DiscardUnknown(m)
}

var xxx_messageInfo_HostID proto.InternalMessageInfo

func (m *HostID) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *HostID) GetIpAddresses() []string {
	if m != nil {
		return m.IpAddresses
	}
	return nil
}

func (m *HostID) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *HostID) GetPid() int32 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *HostID) GetEc2InstanceID() string {
	if m != nil {
		return m.Ec2InstanceID
	}
	return ""
}

func (m *HostID) GetEc2AvailabilityZone() string {
	if m != nil {
		return m.Ec2AvailabilityZone
	}
	return ""
}

func (m *HostID) GetDockerContainerID() string {
	if m != nil {
		return m.DockerContainerID
	}
	return ""
}

func (m *HostID) GetMacAddresses() []string {
	if m != nil {
		return m.MacAddresses
	}
	return nil
}

func (m *HostID) GetHerokuDynoID() string {
	if m != nil {
		return m.HerokuDynoID
	}
	return ""
}

type OboeSetting struct {
	Type                 OboeSettingType   `protobuf:"varint,1,opt,name=type,proto3,enum=collector.OboeSettingType" json:"type,omitempty"`
	Flags                []byte            `protobuf:"bytes,2,opt,name=flags,proto3" json:"flags,omitempty"`
	Timestamp            int64             `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Value                int64             `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	Layer                []byte            `protobuf:"bytes,5,opt,name=layer,proto3" json:"layer,omitempty"`
	Arguments            map[string][]byte `protobuf:"bytes,7,rep,name=arguments,proto3" json:"arguments,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Ttl                  int64             `protobuf:"varint,8,opt,name=ttl,proto3" json:"ttl,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *OboeSetting) Reset()         { *m = OboeSetting{} }
func (m *OboeSetting) String() string { return proto.CompactTextString(m) }
func (*OboeSetting) ProtoMessage()    {}
func (*OboeSetting) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{1}
}
func (m *OboeSetting) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OboeSetting.Unmarshal(m, b)
}
func (m *OboeSetting) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OboeSetting.Marshal(b, m, deterministic)
}
func (dst *OboeSetting) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OboeSetting.Merge(dst, src)
}
func (m *OboeSetting) XXX_Size() int {
	return xxx_messageInfo_OboeSetting.Size(m)
}
func (m *OboeSetting) XXX_DiscardUnknown() {
	xxx_messageInfo_OboeSetting.DiscardUnknown(m)
}

var xxx_messageInfo_OboeSetting proto.InternalMessageInfo

func (m *OboeSetting) GetType() OboeSettingType {
	if m != nil {
		return m.Type
	}
	return OboeSettingType_DEFAULT_SAMPLE_RATE
}

func (m *OboeSetting) GetFlags() []byte {
	if m != nil {
		return m.Flags
	}
	return nil
}

func (m *OboeSetting) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *OboeSetting) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *OboeSetting) GetLayer() []byte {
	if m != nil {
		return m.Layer
	}
	return nil
}

func (m *OboeSetting) GetArguments() map[string][]byte {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *OboeSetting) GetTtl() int64 {
	if m != nil {
		return m.Ttl
	}
	return 0
}

type MessageRequest struct {
	ApiKey               string       `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Messages             [][]byte     `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages,omitempty"`
	Encoding             EncodingType `protobuf:"varint,3,opt,name=encoding,proto3,enum=collector.EncodingType" json:"encoding,omitempty"`
	Identity             *HostID      `protobuf:"bytes,4,opt,name=identity,proto3" json:"identity,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *MessageRequest) Reset()         { *m = MessageRequest{} }
func (m *MessageRequest) String() string { return proto.CompactTextString(m) }
func (*MessageRequest) ProtoMessage()    {}
func (*MessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{2}
}
func (m *MessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageRequest.Unmarshal(m, b)
}
func (m *MessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageRequest.Marshal(b, m, deterministic)
}
func (dst *MessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageRequest.Merge(dst, src)
}
func (m *MessageRequest) XXX_Size() int {
	return xxx_messageInfo_MessageRequest.Size(m)
}
func (m *MessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessageRequest proto.InternalMessageInfo

func (m *MessageRequest) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

func (m *MessageRequest) GetMessages() [][]byte {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *MessageRequest) GetEncoding() EncodingType {
	if m != nil {
		return m.Encoding
	}
	return EncodingType_BSON
}

func (m *MessageRequest) GetIdentity() *HostID {
	if m != nil {
		return m.Identity
	}
	return nil
}

type MessageResult struct {
	Result               ResultCode `protobuf:"varint,1,opt,name=result,proto3,enum=collector.ResultCode" json:"result,omitempty"`
	Arg                  string     `protobuf:"bytes,2,opt,name=arg,proto3" json:"arg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *MessageResult) Reset()         { *m = MessageResult{} }
func (m *MessageResult) String() string { return proto.CompactTextString(m) }
func (*MessageResult) ProtoMessage()    {}
func (*MessageResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{3}
}
func (m *MessageResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageResult.Unmarshal(m, b)
}
func (m *MessageResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageResult.Marshal(b, m, deterministic)
}
func (dst *MessageResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageResult.Merge(dst, src)
}
func (m *MessageResult) XXX_Size() int {
	return xxx_messageInfo_MessageResult.Size(m)
}
func (m *MessageResult) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageResult.DiscardUnknown(m)
}

var xxx_messageInfo_MessageResult proto.InternalMessageInfo

func (m *MessageResult) GetResult() ResultCode {
	if m != nil {
		return m.Result
	}
	return ResultCode_OK
}

func (m *MessageResult) GetArg() string {
	if m != nil {
		return m.Arg
	}
	return ""
}

type SettingsRequest struct {
	ApiKey               string   `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Identity             *HostID  `protobuf:"bytes,2,opt,name=identity,proto3" json:"identity,omitempty"`
	ClientVersion        string   `protobuf:"bytes,3,opt,name=clientVersion,proto3" json:"clientVersion,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SettingsRequest) Reset()         { *m = SettingsRequest{} }
func (m *SettingsRequest) String() string { return proto.CompactTextString(m) }
func (*SettingsRequest) ProtoMessage()    {}
func (*SettingsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{4}
}
func (m *SettingsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SettingsRequest.Unmarshal(m, b)
}
func (m *SettingsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SettingsRequest.Marshal(b, m, deterministic)
}
func (dst *SettingsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SettingsRequest.Merge(dst, src)
}
func (m *SettingsRequest) XXX_Size() int {
	return xxx_messageInfo_SettingsRequest.Size(m)
}
func (m *SettingsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SettingsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SettingsRequest proto.InternalMessageInfo

func (m *SettingsRequest) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

func (m *SettingsRequest) GetIdentity() *HostID {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *SettingsRequest) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

type SettingsResult struct {
	Result               ResultCode     `protobuf:"varint,1,opt,name=result,proto3,enum=collector.ResultCode" json:"result,omitempty"`
	Arg                  string         `protobuf:"bytes,2,opt,name=arg,proto3" json:"arg,omitempty"`
	Settings             []*OboeSetting `protobuf:"bytes,3,rep,name=settings,proto3" json:"settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SettingsResult) Reset()         { *m = SettingsResult{} }
func (m *SettingsResult) String() string { return proto.CompactTextString(m) }
func (*SettingsResult) ProtoMessage()    {}
func (*SettingsResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{5}
}
func (m *SettingsResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SettingsResult.Unmarshal(m, b)
}
func (m *SettingsResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SettingsResult.Marshal(b, m, deterministic)
}
func (dst *SettingsResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SettingsResult.Merge(dst, src)
}
func (m *SettingsResult) XXX_Size() int {
	return xxx_messageInfo_SettingsResult.Size(m)
}
func (m *SettingsResult) XXX_DiscardUnknown() {
	xxx_messageInfo_SettingsResult.DiscardUnknown(m)
}

var xxx_messageInfo_SettingsResult proto.InternalMessageInfo

func (m *SettingsResult) GetResult() ResultCode {
	if m != nil {
		return m.Result
	}
	return ResultCode_OK
}

func (m *SettingsResult) GetArg() string {
	if m != nil {
		return m.Arg
	}
	return ""
}

func (m *SettingsResult) GetSettings() []*OboeSetting {
	if m != nil {
		return m.Settings
	}
	return nil
}

type PingRequest struct {
	ApiKey               string   `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_collector_65775f1a4ec76cc7, []int{6}
}
func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (dst *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(dst, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

func init() {
	proto.RegisterType((*HostID)(nil), "collector.HostID")
	proto.RegisterType((*OboeSetting)(nil), "collector.OboeSetting")
	proto.RegisterMapType((map[string][]byte)(nil), "collector.OboeSetting.ArgumentsEntry")
	proto.RegisterType((*MessageRequest)(nil), "collector.MessageRequest")
	proto.RegisterType((*MessageResult)(nil), "collector.MessageResult")
	proto.RegisterType((*SettingsRequest)(nil), "collector.SettingsRequest")
	proto.RegisterType((*SettingsResult)(nil), "collector.SettingsResult")
	proto.RegisterType((*PingRequest)(nil), "collector.PingRequest")
	proto.RegisterEnum("collector.ResultCode", ResultCode_name, ResultCode_value)
	proto.RegisterEnum("collector.EncodingType", EncodingType_name, EncodingType_value)
	proto.RegisterEnum("collector.OboeSettingType", OboeSettingType_name, OboeSettingType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TraceCollectorClient is the client API for TraceCollector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TraceCollectorClient interface {
	PostEvents(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error)
	PostMetrics(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error)
	PostStatus(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error)
	GetSettings(ctx context.Context, in *SettingsRequest, opts ...grpc.CallOption) (*SettingsResult, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MessageResult, error)
}

type traceCollectorClient struct {
	cc *grpc.ClientConn
}

func NewTraceCollectorClient(cc *grpc.ClientConn) TraceCollectorClient {
	return &traceCollectorClient{cc}
}

func (c *traceCollectorClient) PostEvents(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/postEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) PostMetrics(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/postMetrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) PostStatus(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/postStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) GetSettings(ctx context.Context, in *SettingsRequest, opts ...grpc.CallOption) (*SettingsResult, error) {
	out := new(SettingsResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/getSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TraceCollectorServer is the server API for TraceCollector service.
type TraceCollectorServer interface {
	PostEvents(context.Context, *MessageRequest) (*MessageResult, error)
	PostMetrics(context.Context, *MessageRequest) (*MessageResult, error)
	PostStatus(context.Context, *MessageRequest) (*MessageResult, error)
	GetSettings(context.Context, *SettingsRequest) (*SettingsResult, error)
	Ping(context.Context, *PingRequest) (*MessageResult, error)
}

func RegisterTraceCollectorServer(s *grpc.Server, srv TraceCollectorServer) {
	s.RegisterService(&_TraceCollector_serviceDesc, srv)
}

func _TraceCollector_PostEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).PostEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/PostEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).PostEvents(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_PostMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).PostMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/PostMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).PostMetrics(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_PostStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).PostStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/PostStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).PostStatus(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_GetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).GetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/GetSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).GetSettings(ctx, req.(*SettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TraceCollector_serviceDesc = grpc.ServiceDesc{
	ServiceName: "collector.TraceCollector",
	HandlerType: (*TraceCollectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "postEvents",
			Handler:    _TraceCollector_PostEvents_Handler,
		},
		{
			MethodName: "postMetrics",
			Handler:    _TraceCollector_PostMetrics_Handler,
		},
		{
			MethodName: "postStatus",
			Handler:    _TraceCollector_PostStatus_Handler,
		},
		{
			MethodName: "getSettings",
			Handler:    _TraceCollector_GetSettings_Handler,
		},
		{
			MethodName: "ping",
			Handler:    _TraceCollector_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "collector.proto",
}

func init() { proto.RegisterFile("collector.proto", fileDescriptor_collector_65775f1a4ec76cc7) }

var fileDescriptor_collector_65775f1a4ec76cc7 = []byte{
	// 864 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x41, 0x6f, 0xe2, 0x46,
	0x14, 0x8e, 0x31, 0x21, 0xf0, 0x20, 0xc4, 0x99, 0x34, 0x8d, 0x83, 0xaa, 0x2a, 0xb5, 0xda, 0x15,
	0x8a, 0xba, 0x51, 0xc5, 0x5e, 0xaa, 0x55, 0x2f, 0x5e, 0xec, 0x6c, 0xac, 0x10, 0x40, 0x83, 0x77,
	0xd5, 0xec, 0xc5, 0x9a, 0x98, 0x29, 0x3b, 0x8a, 0xb1, 0x5d, 0x7b, 0x88, 0xe4, 0x53, 0x4f, 0xfd,
	0x1d, 0x3d, 0xf7, 0xda, 0x6b, 0xff, 0x4e, 0x7f, 0x48, 0x35, 0x63, 0x07, 0xec, 0x2d, 0x6a, 0xa4,
	0x68, 0x6f, 0xef, 0x7d, 0xef, 0x9b, 0xc7, 0x37, 0xdf, 0xbc, 0x67, 0xe0, 0xc0, 0x8f, 0x82, 0x80,
	0xfa, 0x3c, 0x4a, 0x2e, 0xe2, 0x24, 0xe2, 0x11, 0x6a, 0xad, 0x01, 0xe3, 0xef, 0x1a, 0x34, 0xae,
	0xa2, 0x94, 0x3b, 0x16, 0xea, 0x41, 0xf3, 0x63, 0x94, 0xf2, 0x90, 0x2c, 0xa9, 0xae, 0x9c, 0x29,
	0xfd, 0x16, 0x5e, 0xe7, 0xe8, 0x1b, 0xe8, 0xb0, 0xd8, 0x23, 0xf3, 0x79, 0x42, 0xd3, 0x94, 0xa6,
	0x7a, 0xed, 0x4c, 0xed, 0xb7, 0x70, 0x9b, 0xc5, 0xe6, 0x23, 0x84, 0x10, 0xd4, 0x57, 0x2b, 0x36,
	0xd7, 0x55, 0x79, 0x54, 0xc6, 0x48, 0x03, 0x35, 0x66, 0x73, 0xbd, 0x7e, 0xa6, 0xf4, 0x77, 0xb1,
	0x08, 0xd1, 0xb7, 0xb0, 0x4f, 0xfd, 0x81, 0x13, 0xa6, 0x9c, 0x84, 0x3e, 0x75, 0x2c, 0x7d, 0x57,
	0xd2, 0xab, 0x20, 0xfa, 0x01, 0x8e, 0xa8, 0x3f, 0x30, 0x1f, 0x08, 0x0b, 0xc8, 0x1d, 0x0b, 0x18,
	0xcf, 0x3e, 0x44, 0x21, 0xd5, 0x1b, 0x92, 0xbb, 0xad, 0x84, 0xbe, 0x87, 0xc3, 0x79, 0xe4, 0xdf,
	0xd3, 0x64, 0x18, 0x85, 0x9c, 0xb0, 0x90, 0x26, 0x8e, 0xa5, 0xef, 0x49, 0xfe, 0x7f, 0x0b, 0xc8,
	0x80, 0xce, 0x92, 0xf8, 0x6b, 0xed, 0x7a, 0x53, 0x5e, 0xa7, 0x82, 0x09, 0xce, 0x47, 0x9a, 0x44,
	0xf7, 0x2b, 0x2b, 0x0b, 0x23, 0xc7, 0xd2, 0x5b, 0xb2, 0x59, 0x05, 0x33, 0xfe, 0xaa, 0x41, 0x7b,
	0x72, 0x17, 0xd1, 0x19, 0xe5, 0x9c, 0x85, 0x0b, 0x74, 0x01, 0x75, 0x9e, 0xc5, 0xb9, 0x7d, 0xdd,
	0x41, 0xef, 0x62, 0x63, 0x7c, 0x89, 0xe5, 0x66, 0x31, 0xc5, 0x92, 0x87, 0xbe, 0x80, 0xdd, 0x5f,
	0x02, 0xb2, 0x10, 0x7e, 0x2a, 0xfd, 0x0e, 0xce, 0x13, 0xf4, 0x15, 0xb4, 0x38, 0x5b, 0xd2, 0x94,
	0x93, 0x65, 0x2c, 0xed, 0x54, 0xf1, 0x06, 0x10, 0x67, 0x1e, 0x48, 0xb0, 0xa2, 0xd2, 0x55, 0x15,
	0xe7, 0x89, 0x40, 0x03, 0x92, 0xd1, 0x44, 0xfa, 0xd9, 0xc1, 0x79, 0x82, 0x86, 0xd0, 0x22, 0xc9,
	0x62, 0xb5, 0xa4, 0x21, 0x4f, 0xf5, 0xbd, 0x33, 0xb5, 0xdf, 0x1e, 0x7c, 0xb7, 0x5d, 0xd4, 0x85,
	0xf9, 0xc8, 0xb3, 0x43, 0x9e, 0x64, 0x78, 0x73, 0x4e, 0x3c, 0x22, 0xe7, 0x81, 0xde, 0x94, 0x3f,
	0x27, 0xc2, 0xde, 0x4f, 0xd0, 0xad, 0xd2, 0x05, 0xe7, 0x9e, 0x66, 0xc5, 0xd8, 0x88, 0x70, 0x23,
	0xb3, 0xb8, 0x9a, 0x4c, 0x5e, 0xd7, 0x7e, 0x54, 0x8c, 0x3f, 0x15, 0xe8, 0xde, 0xd0, 0x34, 0x25,
	0x0b, 0x8a, 0xe9, 0xaf, 0x2b, 0x9a, 0x72, 0x74, 0x02, 0x7b, 0x24, 0x66, 0xde, 0xa6, 0x45, 0x83,
	0xc4, 0xec, 0x9a, 0x66, 0x62, 0x26, 0x97, 0x39, 0x35, 0x9f, 0xb9, 0x0e, 0x5e, 0xe7, 0xe8, 0x15,
	0x34, 0x69, 0xe8, 0x47, 0x73, 0x16, 0x2e, 0xa4, 0x4b, 0xdd, 0xc1, 0x49, 0xe9, 0x6e, 0x76, 0x51,
	0x92, 0x6e, 0xaf, 0x89, 0xe8, 0x25, 0x34, 0xd9, 0x9c, 0x86, 0x9c, 0xf1, 0x4c, 0x1a, 0xd8, 0x1e,
	0x1c, 0x96, 0x0e, 0xe5, 0x9b, 0x80, 0xd7, 0x14, 0x63, 0x0a, 0xfb, 0x6b, 0xa9, 0xe9, 0x2a, 0xe0,
	0xe8, 0x25, 0x34, 0x12, 0x19, 0x15, 0x6f, 0x7c, 0x5c, 0x3a, 0x9d, 0x53, 0x86, 0xd1, 0x9c, 0xe2,
	0x82, 0x24, 0x7c, 0x21, 0xc9, 0x42, 0x7a, 0xd0, 0xc2, 0x22, 0x34, 0x7e, 0x83, 0x83, 0xc2, 0xf2,
	0xf4, 0xc9, 0xdb, 0x97, 0xc5, 0xd6, 0x9e, 0x14, 0x2b, 0x76, 0xcb, 0x0f, 0x18, 0x0d, 0xf9, 0x7b,
	0x9a, 0xa4, 0x2c, 0x0a, 0x8b, 0x55, 0xac, 0x82, 0xc6, 0xef, 0x0a, 0x74, 0x37, 0x0a, 0x3e, 0xcb,
	0xa5, 0xd0, 0x00, 0x9a, 0x69, 0xd1, 0x52, 0x57, 0xe5, 0x98, 0x7d, 0xb9, 0x7d, 0xcc, 0xf0, 0x9a,
	0x67, 0xbc, 0x80, 0xf6, 0x54, 0x20, 0x4f, 0x98, 0x70, 0xfe, 0x01, 0x60, 0xa3, 0x01, 0x35, 0xa0,
	0x36, 0xb9, 0xd6, 0x76, 0xd0, 0x3e, 0xb4, 0x5c, 0x7c, 0xeb, 0x8d, 0x4c, 0xd7, 0xc6, 0x9a, 0x82,
	0x8e, 0xe0, 0xc0, 0x19, 0xbf, 0x37, 0x47, 0x8e, 0xe5, 0x99, 0x53, 0xc7, 0xbb, 0xb6, 0x6f, 0xb5,
	0x1a, 0x42, 0xd0, 0x1d, 0x39, 0x37, 0x8e, 0xeb, 0xd9, 0x3f, 0x0f, 0x6d, 0xdb, 0xb2, 0x2d, 0x4d,
	0x45, 0x1d, 0x68, 0x62, 0xdb, 0x72, 0xb0, 0x3d, 0x74, 0xb5, 0xfa, 0xf9, 0x0b, 0xe8, 0x94, 0xe7,
	0x04, 0x35, 0xa1, 0xfe, 0x66, 0x36, 0x19, 0x6b, 0x3b, 0x82, 0x37, 0xc5, 0x13, 0x77, 0xf2, 0xe6,
	0xdd, 0xa5, 0xa6, 0x9c, 0xff, 0xa1, 0xc0, 0xc1, 0x27, 0x1b, 0x8c, 0x4e, 0xe0, 0xc8, 0xb2, 0x2f,
	0xcd, 0x77, 0x23, 0xd7, 0x9b, 0x99, 0x37, 0xd3, 0x91, 0xed, 0x61, 0xd3, 0xb5, 0xb5, 0x1d, 0x74,
	0x0c, 0x87, 0x23, 0xf3, 0xd6, 0xc6, 0x15, 0x58, 0x41, 0xa7, 0x70, 0x9c, 0xc3, 0xe6, 0x74, 0x5a,
	0x29, 0xd5, 0xd0, 0xd7, 0xd0, 0xcb, 0x4b, 0x57, 0xae, 0x3b, 0xbd, 0x9a, 0xcc, 0xaa, 0x1d, 0x55,
	0x74, 0x08, 0xfb, 0xc3, 0xc9, 0xf8, 0xd2, 0x79, 0xeb, 0xcd, 0x5c, 0xec, 0x8c, 0xdf, 0x6a, 0x75,
	0xd4, 0x05, 0x28, 0x20, 0x67, 0xec, 0x6a, 0xbb, 0x83, 0x7f, 0x6a, 0xd0, 0x75, 0x13, 0xe2, 0xd3,
	0xe1, 0xa3, 0xed, 0x68, 0x08, 0x10, 0x47, 0x29, 0xb7, 0x1f, 0xe4, 0x16, 0x9f, 0x96, 0x1e, 0xa4,
	0xba, 0x7d, 0x3d, 0x7d, 0x5b, 0x49, 0x38, 0x6e, 0xec, 0x20, 0x0b, 0xda, 0xa2, 0xc9, 0x0d, 0xe5,
	0x09, 0xf3, 0x9f, 0xdd, 0xa5, 0x90, 0x32, 0xe3, 0x84, 0xaf, 0x9e, 0xdd, 0xe4, 0x12, 0xda, 0x0b,
	0xca, 0x1f, 0x47, 0x17, 0x95, 0xbf, 0xae, 0x9f, 0x6c, 0x54, 0xef, 0x74, 0x6b, 0xad, 0xe8, 0xf3,
	0x1a, 0xea, 0xb1, 0xf8, 0x14, 0x94, 0x47, 0xb4, 0x34, 0x89, 0xff, 0xa7, 0xe1, 0xae, 0x21, 0xff,
	0x40, 0x5f, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xa3, 0x9f, 0xf4, 0x08, 0x53, 0x07, 0x00, 0x00,
}