// Code generated by protoc-gen-go. DO NOT EDIT.
// source: agent.proto

package proto

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

// Metadata for a Check.
type Check struct {
	Input    *Check_Component   `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	Target   *Check_Component   `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Output   *Check_Component   `protobuf:"bytes,3,opt,name=output,proto3" json:"output,omitempty"`
	Payloads []*Check_Component `protobuf:"bytes,4,rep,name=payloads,proto3" json:"payloads,omitempty"`
	// Both the interval and the timeout values are in
	// seconds.
	Interval             int64    `protobuf:"varint,5,opt,name=interval,proto3" json:"interval,omitempty"`
	Timeout              int64    `protobuf:"varint,6,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Name                 string   `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Check) Reset()         { *m = Check{} }
func (m *Check) String() string { return proto.CompactTextString(m) }
func (*Check) ProtoMessage()    {}
func (*Check) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{0}
}
func (m *Check) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Check.Unmarshal(m, b)
}
func (m *Check) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Check.Marshal(b, m, deterministic)
}
func (dst *Check) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Check.Merge(dst, src)
}
func (m *Check) XXX_Size() int {
	return xxx_messageInfo_Check.Size(m)
}
func (m *Check) XXX_DiscardUnknown() {
	xxx_messageInfo_Check.DiscardUnknown(m)
}

var xxx_messageInfo_Check proto.InternalMessageInfo

func (m *Check) GetInput() *Check_Component {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *Check) GetTarget() *Check_Component {
	if m != nil {
		return m.Target
	}
	return nil
}

func (m *Check) GetOutput() *Check_Component {
	if m != nil {
		return m.Output
	}
	return nil
}

func (m *Check) GetPayloads() []*Check_Component {
	if m != nil {
		return m.Payloads
	}
	return nil
}

func (m *Check) GetInterval() int64 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *Check) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *Check) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Each component of a check has two things associated with it.
// First the type of the value and second the value of the component.
type Check_Component struct {
	Type                 string   `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Check_Component) Reset()         { *m = Check_Component{} }
func (m *Check_Component) String() string { return proto.CompactTextString(m) }
func (*Check_Component) ProtoMessage()    {}
func (*Check_Component) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{0, 0}
}
func (m *Check_Component) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Check_Component.Unmarshal(m, b)
}
func (m *Check_Component) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Check_Component.Marshal(b, m, deterministic)
}
func (dst *Check_Component) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Check_Component.Merge(dst, src)
}
func (m *Check_Component) XXX_Size() int {
	return xxx_messageInfo_Check_Component.Size(m)
}
func (m *Check_Component) XXX_DiscardUnknown() {
	xxx_messageInfo_Check_Component.DiscardUnknown(m)
}

var xxx_messageInfo_Check_Component proto.InternalMessageInfo

func (m *Check_Component) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Check_Component) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// This is the data structure storing the stats returned by
// Manager.
type ManagerStats struct {
	ControllerStatus     []*ManagerStats_ControllerStatus `protobuf:"bytes,1,rep,name=controller_status,json=controllerStatus,proto3" json:"controller_status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *ManagerStats) Reset()         { *m = ManagerStats{} }
func (m *ManagerStats) String() string { return proto.CompactTextString(m) }
func (*ManagerStats) ProtoMessage()    {}
func (*ManagerStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{1}
}
func (m *ManagerStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats.Unmarshal(m, b)
}
func (m *ManagerStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats.Marshal(b, m, deterministic)
}
func (dst *ManagerStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats.Merge(dst, src)
}
func (m *ManagerStats) XXX_Size() int {
	return xxx_messageInfo_ManagerStats.Size(m)
}
func (m *ManagerStats) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerStats.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerStats proto.InternalMessageInfo

func (m *ManagerStats) GetControllerStatus() []*ManagerStats_ControllerStatus {
	if m != nil {
		return m.ControllerStatus
	}
	return nil
}

type ManagerStats_ControllerConfigurationStatus struct {
	ErrorRetry           bool     `protobuf:"varint,1,opt,name=error_retry,json=errorRetry,proto3" json:"error_retry,omitempty"`
	ShouldBackOff        bool     `protobuf:"varint,2,opt,name=should_back_off,json=shouldBackOff,proto3" json:"should_back_off,omitempty"`
	Interval             string   `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ManagerStats_ControllerConfigurationStatus) Reset() {
	*m = ManagerStats_ControllerConfigurationStatus{}
}
func (m *ManagerStats_ControllerConfigurationStatus) String() string {
	return proto.CompactTextString(m)
}
func (*ManagerStats_ControllerConfigurationStatus) ProtoMessage() {}
func (*ManagerStats_ControllerConfigurationStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{1, 0}
}
func (m *ManagerStats_ControllerConfigurationStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Unmarshal(m, b)
}
func (m *ManagerStats_ControllerConfigurationStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Marshal(b, m, deterministic)
}
func (dst *ManagerStats_ControllerConfigurationStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Merge(dst, src)
}
func (m *ManagerStats_ControllerConfigurationStatus) XXX_Size() int {
	return xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Size(m)
}
func (m *ManagerStats_ControllerConfigurationStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerStats_ControllerConfigurationStatus proto.InternalMessageInfo

func (m *ManagerStats_ControllerConfigurationStatus) GetErrorRetry() bool {
	if m != nil {
		return m.ErrorRetry
	}
	return false
}

func (m *ManagerStats_ControllerConfigurationStatus) GetShouldBackOff() bool {
	if m != nil {
		return m.ShouldBackOff
	}
	return false
}

func (m *ManagerStats_ControllerConfigurationStatus) GetInterval() string {
	if m != nil {
		return m.Interval
	}
	return ""
}

type ManagerStats_ControllerRunStatus struct {
	SuccessCount         int64    `protobuf:"varint,1,opt,name=success_count,json=successCount,proto3" json:"success_count,omitempty"`
	FailureCount         int64    `protobuf:"varint,2,opt,name=failure_count,json=failureCount,proto3" json:"failure_count,omitempty"`
	ConsecFailureCount   int64    `protobuf:"varint,3,opt,name=consec_failure_count,json=consecFailureCount,proto3" json:"consec_failure_count,omitempty"`
	LastSuccessTime      string   `protobuf:"bytes,4,opt,name=last_success_time,json=lastSuccessTime,proto3" json:"last_success_time,omitempty"`
	LastFailureTime      string   `protobuf:"bytes,5,opt,name=last_failure_time,json=lastFailureTime,proto3" json:"last_failure_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ManagerStats_ControllerRunStatus) Reset()         { *m = ManagerStats_ControllerRunStatus{} }
func (m *ManagerStats_ControllerRunStatus) String() string { return proto.CompactTextString(m) }
func (*ManagerStats_ControllerRunStatus) ProtoMessage()    {}
func (*ManagerStats_ControllerRunStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{1, 1}
}
func (m *ManagerStats_ControllerRunStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats_ControllerRunStatus.Unmarshal(m, b)
}
func (m *ManagerStats_ControllerRunStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats_ControllerRunStatus.Marshal(b, m, deterministic)
}
func (dst *ManagerStats_ControllerRunStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats_ControllerRunStatus.Merge(dst, src)
}
func (m *ManagerStats_ControllerRunStatus) XXX_Size() int {
	return xxx_messageInfo_ManagerStats_ControllerRunStatus.Size(m)
}
func (m *ManagerStats_ControllerRunStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerStats_ControllerRunStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerStats_ControllerRunStatus proto.InternalMessageInfo

func (m *ManagerStats_ControllerRunStatus) GetSuccessCount() int64 {
	if m != nil {
		return m.SuccessCount
	}
	return 0
}

func (m *ManagerStats_ControllerRunStatus) GetFailureCount() int64 {
	if m != nil {
		return m.FailureCount
	}
	return 0
}

func (m *ManagerStats_ControllerRunStatus) GetConsecFailureCount() int64 {
	if m != nil {
		return m.ConsecFailureCount
	}
	return 0
}

func (m *ManagerStats_ControllerRunStatus) GetLastSuccessTime() string {
	if m != nil {
		return m.LastSuccessTime
	}
	return ""
}

func (m *ManagerStats_ControllerRunStatus) GetLastFailureTime() string {
	if m != nil {
		return m.LastFailureTime
	}
	return ""
}

type ManagerStats_ControllerStatus struct {
	Name                 string                                      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ConfigStatus         *ManagerStats_ControllerConfigurationStatus `protobuf:"bytes,2,opt,name=config_status,json=configStatus,proto3" json:"config_status,omitempty"`
	RunStatus            *ManagerStats_ControllerRunStatus           `protobuf:"bytes,3,opt,name=run_status,json=runStatus,proto3" json:"run_status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *ManagerStats_ControllerStatus) Reset()         { *m = ManagerStats_ControllerStatus{} }
func (m *ManagerStats_ControllerStatus) String() string { return proto.CompactTextString(m) }
func (*ManagerStats_ControllerStatus) ProtoMessage()    {}
func (*ManagerStats_ControllerStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{1, 2}
}
func (m *ManagerStats_ControllerStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats_ControllerStatus.Unmarshal(m, b)
}
func (m *ManagerStats_ControllerStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats_ControllerStatus.Marshal(b, m, deterministic)
}
func (dst *ManagerStats_ControllerStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats_ControllerStatus.Merge(dst, src)
}
func (m *ManagerStats_ControllerStatus) XXX_Size() int {
	return xxx_messageInfo_ManagerStats_ControllerStatus.Size(m)
}
func (m *ManagerStats_ControllerStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerStats_ControllerStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerStats_ControllerStatus proto.InternalMessageInfo

func (m *ManagerStats_ControllerStatus) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ManagerStats_ControllerStatus) GetConfigStatus() *ManagerStats_ControllerConfigurationStatus {
	if m != nil {
		return m.ConfigStatus
	}
	return nil
}

func (m *ManagerStats_ControllerStatus) GetRunStatus() *ManagerStats_ControllerRunStatus {
	if m != nil {
		return m.RunStatus
	}
	return nil
}

// This is to check if the pushed check is consumed by the agent
// or not.
type PushStatus struct {
	Pushed               bool     `protobuf:"varint,1,opt,name=Pushed,proto3" json:"Pushed,omitempty"`
	Reason               string   `protobuf:"bytes,2,opt,name=Reason,proto3" json:"Reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushStatus) Reset()         { *m = PushStatus{} }
func (m *PushStatus) String() string { return proto.CompactTextString(m) }
func (*PushStatus) ProtoMessage()    {}
func (*PushStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{2}
}
func (m *PushStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushStatus.Unmarshal(m, b)
}
func (m *PushStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushStatus.Marshal(b, m, deterministic)
}
func (dst *PushStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushStatus.Merge(dst, src)
}
func (m *PushStatus) XXX_Size() int {
	return xxx_messageInfo_PushStatus.Size(m)
}
func (m *PushStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PushStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PushStatus proto.InternalMessageInfo

func (m *PushStatus) GetPushed() bool {
	if m != nil {
		return m.Pushed
	}
	return false
}

func (m *PushStatus) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type RemoveStatus struct {
	Removed              bool     `protobuf:"varint,1,opt,name=removed,proto3" json:"removed,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveStatus) Reset()         { *m = RemoveStatus{} }
func (m *RemoveStatus) String() string { return proto.CompactTextString(m) }
func (*RemoveStatus) ProtoMessage()    {}
func (*RemoveStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{3}
}
func (m *RemoveStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveStatus.Unmarshal(m, b)
}
func (m *RemoveStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveStatus.Marshal(b, m, deterministic)
}
func (dst *RemoveStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveStatus.Merge(dst, src)
}
func (m *RemoveStatus) XXX_Size() int {
	return xxx_messageInfo_RemoveStatus.Size(m)
}
func (m *RemoveStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveStatus.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveStatus proto.InternalMessageInfo

func (m *RemoveStatus) GetRemoved() bool {
	if m != nil {
		return m.Removed
	}
	return false
}

func (m *RemoveStatus) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type CheckMeta struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckMeta) Reset()         { *m = CheckMeta{} }
func (m *CheckMeta) String() string { return proto.CompactTextString(m) }
func (*CheckMeta) ProtoMessage()    {}
func (*CheckMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{4}
}
func (m *CheckMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckMeta.Unmarshal(m, b)
}
func (m *CheckMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckMeta.Marshal(b, m, deterministic)
}
func (dst *CheckMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckMeta.Merge(dst, src)
}
func (m *CheckMeta) XXX_Size() int {
	return xxx_messageInfo_CheckMeta.Size(m)
}
func (m *CheckMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckMeta.DiscardUnknown(m)
}

var xxx_messageInfo_CheckMeta proto.InternalMessageInfo

func (m *CheckMeta) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ChecksList struct {
	Checks               []*CheckMeta `protobuf:"bytes,1,rep,name=checks,proto3" json:"checks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ChecksList) Reset()         { *m = ChecksList{} }
func (m *ChecksList) String() string { return proto.CompactTextString(m) }
func (*ChecksList) ProtoMessage()    {}
func (*ChecksList) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{5}
}
func (m *ChecksList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChecksList.Unmarshal(m, b)
}
func (m *ChecksList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChecksList.Marshal(b, m, deterministic)
}
func (dst *ChecksList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChecksList.Merge(dst, src)
}
func (m *ChecksList) XXX_Size() int {
	return xxx_messageInfo_ChecksList.Size(m)
}
func (m *ChecksList) XXX_DiscardUnknown() {
	xxx_messageInfo_ChecksList.DiscardUnknown(m)
}

var xxx_messageInfo_ChecksList proto.InternalMessageInfo

func (m *ChecksList) GetChecks() []*CheckMeta {
	if m != nil {
		return m.Checks
	}
	return nil
}

type None struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *None) Reset()         { *m = None{} }
func (m *None) String() string { return proto.CompactTextString(m) }
func (*None) ProtoMessage()    {}
func (*None) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_8fa96ae5205500cf, []int{6}
}
func (m *None) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_None.Unmarshal(m, b)
}
func (m *None) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_None.Marshal(b, m, deterministic)
}
func (dst *None) XXX_Merge(src proto.Message) {
	xxx_messageInfo_None.Merge(dst, src)
}
func (m *None) XXX_Size() int {
	return xxx_messageInfo_None.Size(m)
}
func (m *None) XXX_DiscardUnknown() {
	xxx_messageInfo_None.DiscardUnknown(m)
}

var xxx_messageInfo_None proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Check)(nil), "proto.Check")
	proto.RegisterType((*Check_Component)(nil), "proto.Check.Component")
	proto.RegisterType((*ManagerStats)(nil), "proto.ManagerStats")
	proto.RegisterType((*ManagerStats_ControllerConfigurationStatus)(nil), "proto.ManagerStats.ControllerConfigurationStatus")
	proto.RegisterType((*ManagerStats_ControllerRunStatus)(nil), "proto.ManagerStats.ControllerRunStatus")
	proto.RegisterType((*ManagerStats_ControllerStatus)(nil), "proto.ManagerStats.ControllerStatus")
	proto.RegisterType((*PushStatus)(nil), "proto.PushStatus")
	proto.RegisterType((*RemoveStatus)(nil), "proto.RemoveStatus")
	proto.RegisterType((*CheckMeta)(nil), "proto.CheckMeta")
	proto.RegisterType((*ChecksList)(nil), "proto.ChecksList")
	proto.RegisterType((*None)(nil), "proto.None")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AgentServiceClient is the client API for AgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AgentServiceClient interface {
	// Pushes a check to the agent service to be processed and managed
	PushCheck(ctx context.Context, in *Check, opts ...grpc.CallOption) (*PushStatus, error)
	// Returns the status of manager managing the controllers corresponding to
	// the checks.
	GetManagerStats(ctx context.Context, in *None, opts ...grpc.CallOption) (*ManagerStats, error)
	// Removes a check from the manager managing the controller for the checks.
	// Only reuired filed for the CheckMeta is `Name` rest everything can be ignored.
	RemoveCheck(ctx context.Context, in *CheckMeta, opts ...grpc.CallOption) (*RemoveStatus, error)
	// Lists all the checks currently being managed by the controller manager.
	// The only reuired field in the CheckList is `Name` rest everything can be blank.
	ListChecks(ctx context.Context, in *None, opts ...grpc.CallOption) (*ChecksList, error)
}

type agentServiceClient struct {
	cc *grpc.ClientConn
}

func NewAgentServiceClient(cc *grpc.ClientConn) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) PushCheck(ctx context.Context, in *Check, opts ...grpc.CallOption) (*PushStatus, error) {
	out := new(PushStatus)
	err := c.cc.Invoke(ctx, "/proto.AgentService/PushCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) GetManagerStats(ctx context.Context, in *None, opts ...grpc.CallOption) (*ManagerStats, error) {
	out := new(ManagerStats)
	err := c.cc.Invoke(ctx, "/proto.AgentService/GetManagerStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) RemoveCheck(ctx context.Context, in *CheckMeta, opts ...grpc.CallOption) (*RemoveStatus, error) {
	out := new(RemoveStatus)
	err := c.cc.Invoke(ctx, "/proto.AgentService/RemoveCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) ListChecks(ctx context.Context, in *None, opts ...grpc.CallOption) (*ChecksList, error) {
	out := new(ChecksList)
	err := c.cc.Invoke(ctx, "/proto.AgentService/ListChecks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServiceServer is the server API for AgentService service.
type AgentServiceServer interface {
	// Pushes a check to the agent service to be processed and managed
	PushCheck(context.Context, *Check) (*PushStatus, error)
	// Returns the status of manager managing the controllers corresponding to
	// the checks.
	GetManagerStats(context.Context, *None) (*ManagerStats, error)
	// Removes a check from the manager managing the controller for the checks.
	// Only reuired filed for the CheckMeta is `Name` rest everything can be ignored.
	RemoveCheck(context.Context, *CheckMeta) (*RemoveStatus, error)
	// Lists all the checks currently being managed by the controller manager.
	// The only reuired field in the CheckList is `Name` rest everything can be blank.
	ListChecks(context.Context, *None) (*ChecksList, error)
}

func RegisterAgentServiceServer(s *grpc.Server, srv AgentServiceServer) {
	s.RegisterService(&_AgentService_serviceDesc, srv)
}

func _AgentService_PushCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Check)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).PushCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AgentService/PushCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).PushCheck(ctx, req.(*Check))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_GetManagerStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).GetManagerStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AgentService/GetManagerStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).GetManagerStats(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_RemoveCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckMeta)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).RemoveCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AgentService/RemoveCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).RemoveCheck(ctx, req.(*CheckMeta))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_ListChecks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).ListChecks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AgentService/ListChecks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).ListChecks(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

var _AgentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushCheck",
			Handler:    _AgentService_PushCheck_Handler,
		},
		{
			MethodName: "GetManagerStats",
			Handler:    _AgentService_GetManagerStats_Handler,
		},
		{
			MethodName: "RemoveCheck",
			Handler:    _AgentService_RemoveCheck_Handler,
		},
		{
			MethodName: "ListChecks",
			Handler:    _AgentService_ListChecks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent.proto",
}

func init() { proto.RegisterFile("agent.proto", fileDescriptor_agent_8fa96ae5205500cf) }

var fileDescriptor_agent_8fa96ae5205500cf = []byte{
	// 681 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xdf, 0x4e, 0xdb, 0x3e,
	0x14, 0xc7, 0x1b, 0xda, 0x94, 0xe6, 0xb4, 0x08, 0x6a, 0x10, 0x8a, 0x22, 0xfd, 0x44, 0x95, 0xdf,
	0xb4, 0x55, 0xd3, 0x14, 0x6d, 0x4c, 0x70, 0xb5, 0x9b, 0x51, 0x89, 0xdd, 0x0c, 0xb6, 0xb9, 0x88,
	0xdb, 0xc8, 0x04, 0xb7, 0x44, 0xa4, 0x76, 0x65, 0x3b, 0x48, 0xdc, 0xef, 0x35, 0xf6, 0x3a, 0xbb,
	0xd8, 0x3b, 0xec, 0x01, 0xf6, 0x16, 0x93, 0xff, 0x24, 0x4d, 0x19, 0xea, 0x55, 0xfc, 0x3d, 0xfe,
	0xf8, 0x1c, 0x9f, 0x73, 0x7c, 0x02, 0x7d, 0x32, 0xa7, 0x4c, 0x25, 0x4b, 0xc1, 0x15, 0x47, 0xbe,
	0xf9, 0xc4, 0xbf, 0xb6, 0xc0, 0x9f, 0xdc, 0xd1, 0xec, 0x1e, 0xbd, 0x01, 0x3f, 0x67, 0xcb, 0x52,
	0x85, 0xde, 0xc8, 0x1b, 0xf7, 0x8f, 0x0f, 0x2d, 0x97, 0x98, 0xcd, 0x64, 0xc2, 0x17, 0x4b, 0xce,
	0x28, 0x53, 0xd8, 0x42, 0x28, 0x81, 0xae, 0x22, 0x62, 0x4e, 0x55, 0xb8, 0xb5, 0x11, 0x77, 0x94,
	0xe6, 0x79, 0xa9, 0xb4, 0xfb, 0xf6, 0x66, 0xde, 0x52, 0xe8, 0x18, 0x7a, 0x4b, 0xf2, 0x58, 0x70,
	0x72, 0x2b, 0xc3, 0xce, 0xa8, 0xbd, 0xe1, 0x44, 0xcd, 0xa1, 0x08, 0x7a, 0x39, 0x53, 0x54, 0x3c,
	0x90, 0x22, 0xf4, 0x47, 0xde, 0xb8, 0x8d, 0x6b, 0x8d, 0x42, 0xd8, 0x56, 0xf9, 0x82, 0xf2, 0x52,
	0x85, 0x5d, 0xb3, 0x55, 0x49, 0x84, 0xa0, 0xc3, 0xc8, 0x82, 0x86, 0xdb, 0x23, 0x6f, 0x1c, 0x60,
	0xb3, 0x8e, 0x4e, 0x20, 0xa8, 0x03, 0x68, 0xe0, 0xea, 0x71, 0x49, 0x4d, 0x5d, 0x02, 0x6c, 0xd6,
	0xe8, 0x00, 0xfc, 0x6b, 0x52, 0x94, 0xd4, 0x64, 0x1f, 0x60, 0x2b, 0xe2, 0x1f, 0x3e, 0x0c, 0x2e,
	0x08, 0x23, 0x73, 0x2a, 0xa6, 0x8a, 0x28, 0x89, 0xbe, 0xc1, 0x30, 0xe3, 0x4c, 0x09, 0x5e, 0x14,
	0x54, 0xa4, 0x52, 0x11, 0x55, 0xca, 0xd0, 0x33, 0xe9, 0xbc, 0x70, 0xe9, 0x34, 0xf9, 0x64, 0x52,
	0xc3, 0x53, 0xc3, 0xe2, 0xbd, 0xec, 0x89, 0x25, 0xfa, 0xee, 0xc1, 0x7f, 0x2b, 0x6c, 0xc2, 0xd9,
	0x2c, 0x9f, 0x97, 0x82, 0xa8, 0x9c, 0x33, 0x4b, 0xa0, 0x23, 0xe8, 0x53, 0x21, 0xb8, 0x48, 0x05,
	0x55, 0xe2, 0xd1, 0x5c, 0xbb, 0x87, 0xc1, 0x98, 0xb0, 0xb6, 0xa0, 0x97, 0xb0, 0x2b, 0xef, 0x78,
	0x59, 0xdc, 0xa6, 0x37, 0x24, 0xbb, 0x4f, 0xf9, 0x6c, 0x66, 0xd2, 0xe8, 0xe1, 0x1d, 0x6b, 0x3e,
	0x23, 0xd9, 0xfd, 0x97, 0xd9, 0x6c, 0xad, 0x9e, 0x6d, 0x93, 0x67, 0xad, 0xa3, 0x3f, 0x1e, 0xec,
	0xaf, 0xae, 0x81, 0xcb, 0x2a, 0xf8, 0xff, 0xb0, 0x23, 0xcb, 0x2c, 0xa3, 0x52, 0xa6, 0x19, 0x2f,
	0x99, 0x7d, 0x4d, 0x6d, 0x3c, 0x70, 0xc6, 0x89, 0xb6, 0x69, 0x68, 0x46, 0xf2, 0xa2, 0x14, 0xd4,
	0x41, 0x5b, 0x16, 0x72, 0x46, 0x0b, 0xbd, 0x85, 0x83, 0x8c, 0x33, 0x49, 0xb3, 0x74, 0x9d, 0x6d,
	0x1b, 0x16, 0xd9, 0xbd, 0xf3, 0xe6, 0x89, 0xd7, 0x30, 0x2c, 0x88, 0x54, 0x69, 0x75, 0x01, 0xdd,
	0xe1, 0xb0, 0x63, 0x2e, 0xbe, 0xab, 0x37, 0xa6, 0xd6, 0x7e, 0x95, 0x2f, 0x68, 0xcd, 0x56, 0xbe,
	0x0d, 0xeb, 0xaf, 0x58, 0xe7, 0x58, 0xb3, 0xd1, 0x4f, 0x0f, 0xf6, 0x9e, 0x76, 0xa6, 0x7e, 0x36,
	0xde, 0xea, 0xd9, 0xa0, 0x6b, 0xd8, 0xc9, 0x4c, 0x43, 0xaa, 0x56, 0xdb, 0xd9, 0x78, 0xb7, 0xb9,
	0xd5, 0xcf, 0xf4, 0x10, 0x0f, 0xac, 0x1f, 0x17, 0xeb, 0x1c, 0x40, 0x94, 0xac, 0x72, 0x6a, 0x07,
	0xe8, 0xd5, 0x66, 0xa7, 0x75, 0x47, 0x70, 0x20, 0xaa, 0x65, 0xfc, 0x01, 0xe0, 0x6b, 0x29, 0xef,
	0x9c, 0xd7, 0x43, 0xe8, 0x6a, 0x45, 0x6f, 0xdd, 0x13, 0x71, 0x4a, 0xdb, 0x31, 0x25, 0x92, 0x33,
	0xf7, 0xb8, 0x9d, 0x8a, 0xcf, 0x60, 0x80, 0xe9, 0x82, 0x3f, 0x50, 0x77, 0x3e, 0x84, 0x6d, 0x61,
	0x74, 0xe5, 0xa0, 0x92, 0x7a, 0x67, 0x41, 0xa5, 0x24, 0xf3, 0x6a, 0x3e, 0x2a, 0x19, 0x1f, 0x41,
	0x60, 0xe6, 0xf7, 0x82, 0x2a, 0xa2, 0x4b, 0x78, 0xd9, 0x28, 0xa1, 0x5e, 0xc7, 0xa7, 0x00, 0x06,
	0x90, 0x9f, 0x73, 0xa9, 0xd0, 0x18, 0xba, 0x99, 0x51, 0x6e, 0x68, 0xf6, 0x9a, 0xff, 0x00, 0xed,
	0x03, 0xbb, 0xfd, 0xb8, 0x0b, 0x9d, 0x4b, 0xce, 0xe8, 0xf1, 0x6f, 0x0f, 0x06, 0x1f, 0xf5, 0x6f,
	0x6e, 0x4a, 0xc5, 0x43, 0x9e, 0x51, 0x94, 0x40, 0xa0, 0xf3, 0xb2, 0xff, 0xb8, 0x41, 0xf3, 0x7c,
	0x34, 0x74, 0x6a, 0x55, 0x93, 0xb8, 0x85, 0x4e, 0x60, 0xf7, 0x13, 0x55, 0x6b, 0x53, 0xdc, 0x77,
	0x9c, 0x0e, 0x10, 0xed, 0x3f, 0x53, 0xf7, 0xb8, 0x85, 0x4e, 0xa1, 0x6f, 0x8b, 0x63, 0x03, 0xfd,
	0x73, 0xd1, 0xfa, 0x5c, 0xb3, 0x84, 0x71, 0x0b, 0x25, 0x00, 0x3a, 0x53, 0x9b, 0xf3, 0x7a, 0xa4,
	0x61, 0xd3, 0x87, 0xa9, 0x47, 0xdc, 0xba, 0xe9, 0x1a, 0xdb, 0xfb, 0xbf, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x9d, 0xec, 0xdd, 0xb6, 0xcc, 0x05, 0x00, 0x00,
}