// Code generated by protoc-gen-go. DO NOT EDIT.
// source: agent.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Metadata for a Check.
type Check struct {
	Input    *Check_Component  `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	Target   *Check_Component  `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Output   *Check_Component  `protobuf:"bytes,3,opt,name=output,proto3" json:"output,omitempty"`
	Payloads []*Check_Payloads `protobuf:"bytes,4,rep,name=payloads,proto3" json:"payloads,omitempty"`
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
	return fileDescriptor_56ede974c0020f77, []int{0}
}

func (m *Check) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Check.Unmarshal(m, b)
}
func (m *Check) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Check.Marshal(b, m, deterministic)
}
func (m *Check) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Check.Merge(m, src)
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

func (m *Check) GetPayloads() []*Check_Payloads {
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
	return fileDescriptor_56ede974c0020f77, []int{0, 0}
}

func (m *Check_Component) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Check_Component.Unmarshal(m, b)
}
func (m *Check_Component) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Check_Component.Marshal(b, m, deterministic)
}
func (m *Check_Component) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Check_Component.Merge(m, src)
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

type Check_Payloads struct {
	PayloadType          string   `protobuf:"bytes,1,opt,name=payload_type,json=payloadType,proto3" json:"payload_type,omitempty"`
	Payload              string   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Check_Payloads) Reset()         { *m = Check_Payloads{} }
func (m *Check_Payloads) String() string { return proto.CompactTextString(m) }
func (*Check_Payloads) ProtoMessage()    {}
func (*Check_Payloads) Descriptor() ([]byte, []int) {
	return fileDescriptor_56ede974c0020f77, []int{0, 1}
}

func (m *Check_Payloads) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Check_Payloads.Unmarshal(m, b)
}
func (m *Check_Payloads) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Check_Payloads.Marshal(b, m, deterministic)
}
func (m *Check_Payloads) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Check_Payloads.Merge(m, src)
}
func (m *Check_Payloads) XXX_Size() int {
	return xxx_messageInfo_Check_Payloads.Size(m)
}
func (m *Check_Payloads) XXX_DiscardUnknown() {
	xxx_messageInfo_Check_Payloads.DiscardUnknown(m)
}

var xxx_messageInfo_Check_Payloads proto.InternalMessageInfo

func (m *Check_Payloads) GetPayloadType() string {
	if m != nil {
		return m.PayloadType
	}
	return ""
}

func (m *Check_Payloads) GetPayload() string {
	if m != nil {
		return m.Payload
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
	return fileDescriptor_56ede974c0020f77, []int{1}
}

func (m *ManagerStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats.Unmarshal(m, b)
}
func (m *ManagerStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats.Marshal(b, m, deterministic)
}
func (m *ManagerStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats.Merge(m, src)
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
	return fileDescriptor_56ede974c0020f77, []int{1, 0}
}

func (m *ManagerStats_ControllerConfigurationStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Unmarshal(m, b)
}
func (m *ManagerStats_ControllerConfigurationStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Marshal(b, m, deterministic)
}
func (m *ManagerStats_ControllerConfigurationStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats_ControllerConfigurationStatus.Merge(m, src)
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
	return fileDescriptor_56ede974c0020f77, []int{1, 1}
}

func (m *ManagerStats_ControllerRunStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats_ControllerRunStatus.Unmarshal(m, b)
}
func (m *ManagerStats_ControllerRunStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats_ControllerRunStatus.Marshal(b, m, deterministic)
}
func (m *ManagerStats_ControllerRunStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats_ControllerRunStatus.Merge(m, src)
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
	return fileDescriptor_56ede974c0020f77, []int{1, 2}
}

func (m *ManagerStats_ControllerStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerStats_ControllerStatus.Unmarshal(m, b)
}
func (m *ManagerStats_ControllerStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerStats_ControllerStatus.Marshal(b, m, deterministic)
}
func (m *ManagerStats_ControllerStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerStats_ControllerStatus.Merge(m, src)
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
	return fileDescriptor_56ede974c0020f77, []int{2}
}

func (m *PushStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushStatus.Unmarshal(m, b)
}
func (m *PushStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushStatus.Marshal(b, m, deterministic)
}
func (m *PushStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushStatus.Merge(m, src)
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

type None struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *None) Reset()         { *m = None{} }
func (m *None) String() string { return proto.CompactTextString(m) }
func (*None) ProtoMessage()    {}
func (*None) Descriptor() ([]byte, []int) {
	return fileDescriptor_56ede974c0020f77, []int{3}
}

func (m *None) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_None.Unmarshal(m, b)
}
func (m *None) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_None.Marshal(b, m, deterministic)
}
func (m *None) XXX_Merge(src proto.Message) {
	xxx_messageInfo_None.Merge(m, src)
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
	proto.RegisterType((*Check_Payloads)(nil), "proto.Check.Payloads")
	proto.RegisterType((*ManagerStats)(nil), "proto.ManagerStats")
	proto.RegisterType((*ManagerStats_ControllerConfigurationStatus)(nil), "proto.ManagerStats.ControllerConfigurationStatus")
	proto.RegisterType((*ManagerStats_ControllerRunStatus)(nil), "proto.ManagerStats.ControllerRunStatus")
	proto.RegisterType((*ManagerStats_ControllerStatus)(nil), "proto.ManagerStats.ControllerStatus")
	proto.RegisterType((*PushStatus)(nil), "proto.PushStatus")
	proto.RegisterType((*None)(nil), "proto.None")
}

func init() { proto.RegisterFile("agent.proto", fileDescriptor_56ede974c0020f77) }

var fileDescriptor_56ede974c0020f77 = []byte{
	// 613 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xcd, 0x4e, 0x1b, 0x3d,
	0x14, 0x86, 0x19, 0x26, 0x13, 0x92, 0x93, 0x20, 0xc0, 0xf0, 0xa1, 0xd1, 0x48, 0x9f, 0x9a, 0xa6,
	0x55, 0x1b, 0x55, 0x55, 0x54, 0xa8, 0xd8, 0x75, 0xd3, 0x46, 0x82, 0x55, 0x5b, 0xea, 0x20, 0xb6,
	0x23, 0x63, 0x9c, 0x30, 0x62, 0x62, 0x47, 0xfe, 0x41, 0x62, 0xdf, 0x2b, 0xe8, 0xbe, 0xb7, 0xd3,
	0x7b, 0xe9, 0x5d, 0x54, 0xfe, 0x99, 0x99, 0x04, 0xa1, 0xac, 0xc6, 0xe7, 0xf5, 0xe3, 0xd7, 0xc7,
	0x3e, 0xc7, 0x03, 0x3d, 0x32, 0x67, 0x5c, 0x8f, 0x97, 0x52, 0x68, 0x81, 0x12, 0xf7, 0x19, 0xfe,
	0x8a, 0x21, 0x99, 0xdc, 0x31, 0x7a, 0x8f, 0xde, 0x43, 0x52, 0xf0, 0xa5, 0xd1, 0x69, 0x34, 0x88,
	0x46, 0xbd, 0xd3, 0x63, 0xcf, 0x8d, 0xdd, 0xe4, 0x78, 0x22, 0x16, 0x4b, 0xc1, 0x19, 0xd7, 0xd8,
	0x43, 0x68, 0x0c, 0x6d, 0x4d, 0xe4, 0x9c, 0xe9, 0x74, 0x7b, 0x23, 0x1e, 0x28, 0xcb, 0x0b, 0xa3,
	0xad, 0x7d, 0xbc, 0x99, 0xf7, 0x14, 0x3a, 0x81, 0xce, 0x92, 0x3c, 0x96, 0x82, 0xdc, 0xaa, 0xb4,
	0x35, 0x88, 0x47, 0xbd, 0xd3, 0xff, 0xd6, 0x56, 0x5c, 0x86, 0x49, 0x5c, 0x63, 0x28, 0x83, 0x4e,
	0xc1, 0x35, 0x93, 0x0f, 0xa4, 0x4c, 0x93, 0x41, 0x34, 0x8a, 0x71, 0x1d, 0xa3, 0x14, 0x76, 0x74,
	0xb1, 0x60, 0xc2, 0xe8, 0xb4, 0xed, 0xa6, 0xaa, 0x10, 0x21, 0x68, 0x71, 0xb2, 0x60, 0xe9, 0xce,
	0x20, 0x1a, 0x75, 0xb1, 0x1b, 0x67, 0x67, 0xd0, 0xad, 0x33, 0xb2, 0xc0, 0xd5, 0xe3, 0x92, 0xb9,
	0x6b, 0xe9, 0x62, 0x37, 0x46, 0x47, 0x90, 0x5c, 0x93, 0xd2, 0x30, 0x77, 0xf8, 0x2e, 0xf6, 0x41,
	0x76, 0x01, 0x9d, 0x2a, 0x2d, 0xf4, 0x12, 0xfa, 0x21, 0xb1, 0x5c, 0x37, 0xab, 0x7b, 0x41, 0x73,
	0x26, 0x29, 0xec, 0x84, 0x30, 0xd8, 0x54, 0xe1, 0xf0, 0x77, 0x02, 0xfd, 0xaf, 0x84, 0x93, 0x39,
	0x93, 0x53, 0x4d, 0xb4, 0x42, 0x3f, 0xe0, 0x80, 0x0a, 0xae, 0xa5, 0x28, 0x4b, 0x26, 0x73, 0xa5,
	0x89, 0x36, 0x2a, 0x8d, 0xdc, 0xb5, 0xbc, 0x0e, 0xd7, 0xb2, 0xca, 0x8f, 0x27, 0x35, 0x3c, 0x75,
	0x2c, 0xde, 0xa7, 0x4f, 0x94, 0xec, 0x67, 0x04, 0xff, 0x37, 0xd8, 0x44, 0xf0, 0x59, 0x31, 0x37,
	0x92, 0xe8, 0x42, 0x70, 0x4f, 0xa0, 0x17, 0xd0, 0x63, 0x52, 0x0a, 0x99, 0x4b, 0xa6, 0xe5, 0xa3,
	0x3b, 0x41, 0x07, 0x83, 0x93, 0xb0, 0x55, 0xd0, 0x1b, 0xd8, 0x53, 0x77, 0xc2, 0x94, 0xb7, 0xf9,
	0x0d, 0xa1, 0xf7, 0xb9, 0x98, 0xcd, 0xdc, 0x41, 0x3a, 0x78, 0xd7, 0xcb, 0x5f, 0x08, 0xbd, 0xff,
	0x3e, 0x9b, 0xad, 0x15, 0x26, 0x76, 0x27, 0xad, 0xe3, 0xec, 0x6f, 0x04, 0x87, 0x4d, 0x1a, 0xd8,
	0x54, 0x9b, 0xbf, 0x82, 0x5d, 0x65, 0x28, 0x65, 0x4a, 0xe5, 0x54, 0x18, 0xee, 0xbb, 0x32, 0xc6,
	0xfd, 0x20, 0x4e, 0xac, 0x66, 0xa1, 0x19, 0x29, 0x4a, 0x23, 0x59, 0x80, 0xb6, 0x3d, 0x14, 0x44,
	0x0f, 0x7d, 0x80, 0x23, 0x2a, 0xb8, 0x62, 0x34, 0x5f, 0x67, 0x63, 0xc7, 0x22, 0x3f, 0x77, 0xbe,
	0xba, 0xe2, 0x1d, 0x1c, 0x94, 0x44, 0xe9, 0xbc, 0x4a, 0xc0, 0xb6, 0x4a, 0xda, 0x72, 0x89, 0xef,
	0xd9, 0x89, 0xa9, 0xd7, 0xaf, 0x8a, 0x05, 0xab, 0xd9, 0xca, 0xdb, 0xb1, 0x49, 0xc3, 0x06, 0x63,
	0xcb, 0x66, 0x7f, 0x22, 0xd8, 0x7f, 0x5a, 0x99, 0xba, 0xff, 0xa2, 0xa6, 0xff, 0xd0, 0x35, 0xec,
	0x52, 0x57, 0x90, 0xaa, 0xd4, 0xfe, 0x8d, 0x9d, 0x6c, 0x2e, 0xf5, 0x33, 0x35, 0xc4, 0x7d, 0xef,
	0x13, 0xf6, 0x3a, 0x07, 0x90, 0x86, 0x57, 0xa6, 0xfe, 0x21, 0xbe, 0xdd, 0x6c, 0x5a, 0x57, 0x04,
	0x77, 0x65, 0x35, 0x1c, 0x7e, 0x02, 0xb8, 0x34, 0xea, 0x2e, 0xb8, 0x1e, 0x43, 0xdb, 0x46, 0xec,
	0x36, 0xb4, 0x48, 0x88, 0xac, 0x8e, 0x19, 0x51, 0x82, 0x87, 0xf6, 0x0e, 0xd1, 0xb0, 0x0d, 0xad,
	0x6f, 0x82, 0xb3, 0x53, 0x03, 0xfd, 0xcf, 0xf6, 0x87, 0x34, 0x65, 0xf2, 0xa1, 0xa0, 0x0c, 0x8d,
	0xa1, 0x6b, 0x57, 0xfa, 0xbf, 0x51, 0x7f, 0xf5, 0xb5, 0x67, 0x07, 0x21, 0x6a, 0x76, 0x1d, 0x6e,
	0xa1, 0x33, 0xd8, 0xbb, 0x60, 0x7a, 0xed, 0x9d, 0xf4, 0x02, 0x67, 0xfd, 0xb3, 0xc3, 0x67, 0x4e,
	0x36, 0xdc, 0xba, 0x69, 0x3b, 0xf5, 0xe3, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa1, 0xa1, 0xe4,
	0xd5, 0x0e, 0x05, 0x00, 0x00,
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

// AgentServiceServer is the server API for AgentService service.
type AgentServiceServer interface {
	// Pushes a check to the agent service to be processed and managed
	PushCheck(context.Context, *Check) (*PushStatus, error)
	// Returns the status of manager managing the controllers corresponding to
	// the checks.
	GetManagerStats(context.Context, *None) (*ManagerStats, error)
}

// UnimplementedAgentServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (*UnimplementedAgentServiceServer) PushCheck(ctx context.Context, req *Check) (*PushStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushCheck not implemented")
}
func (*UnimplementedAgentServiceServer) GetManagerStats(ctx context.Context, req *None) (*ManagerStats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManagerStats not implemented")
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent.proto",
}