// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/twap/v1beta1/twap_record.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// A TWAP record should be indexed in state by pool_id, (asset pair), timestamp
// The asset pair assets should be lexicographically sorted.
// Technically (pool_id, asset_0_denom, asset_1_denom, height) do not need to
// appear in the struct however we view this as the wrong performance tradeoff
// given SDK today. Would rather we optimize for readability and correctness,
// than an optimal state storage format. The system bottleneck is elsewhere for
// now.
type TwapRecord struct {
	PoolId uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	// Lexicographically smaller denom of the pair
	Asset0Denom string `protobuf:"bytes,2,opt,name=asset0_denom,json=asset0Denom,proto3" json:"asset0_denom,omitempty"`
	// Lexicographically larger denom of the pair
	Asset1Denom string `protobuf:"bytes,3,opt,name=asset1_denom,json=asset1Denom,proto3" json:"asset1_denom,omitempty"`
	// height this record corresponds to, for debugging purposes
	Height int64 `protobuf:"varint,4,opt,name=height,proto3" json:"record_height" yaml:"record_height"`
	// This field should only exist until we have a global registry in the state
	// machine, mapping prior block heights within {TIME RANGE} to times.
	Time time.Time `protobuf:"bytes,5,opt,name=time,proto3,stdtime" json:"time" yaml:"record_time"`
	// We store the last spot prices in the struct, so that we can interpolate
	// accumulator values for times between when accumulator records are stored.
	P0LastSpotPrice             cosmossdk_io_math.LegacyDec `protobuf:"bytes,6,opt,name=p0_last_spot_price,json=p0LastSpotPrice,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"p0_last_spot_price"`
	P1LastSpotPrice             cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=p1_last_spot_price,json=p1LastSpotPrice,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"p1_last_spot_price"`
	P0ArithmeticTwapAccumulator cosmossdk_io_math.LegacyDec `protobuf:"bytes,8,opt,name=p0_arithmetic_twap_accumulator,json=p0ArithmeticTwapAccumulator,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"p0_arithmetic_twap_accumulator"`
	P1ArithmeticTwapAccumulator cosmossdk_io_math.LegacyDec `protobuf:"bytes,9,opt,name=p1_arithmetic_twap_accumulator,json=p1ArithmeticTwapAccumulator,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"p1_arithmetic_twap_accumulator"`
	GeometricTwapAccumulator    cosmossdk_io_math.LegacyDec `protobuf:"bytes,10,opt,name=geometric_twap_accumulator,json=geometricTwapAccumulator,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"geometric_twap_accumulator"`
	// This field contains the time in which the last spot price error occurred.
	// It is used to alert the caller if they are getting a potentially erroneous
	// TWAP, due to an unforeseen underlying error.
	LastErrorTime time.Time `protobuf:"bytes,11,opt,name=last_error_time,json=lastErrorTime,proto3,stdtime" json:"last_error_time" yaml:"last_error_time"`
}

func (m *TwapRecord) Reset()         { *m = TwapRecord{} }
func (m *TwapRecord) String() string { return proto.CompactTextString(m) }
func (*TwapRecord) ProtoMessage()    {}
func (*TwapRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbf5c78678e601aa, []int{0}
}
func (m *TwapRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TwapRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TwapRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TwapRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TwapRecord.Merge(m, src)
}
func (m *TwapRecord) XXX_Size() int {
	return m.Size()
}
func (m *TwapRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_TwapRecord.DiscardUnknown(m)
}

var xxx_messageInfo_TwapRecord proto.InternalMessageInfo

func (m *TwapRecord) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *TwapRecord) GetAsset0Denom() string {
	if m != nil {
		return m.Asset0Denom
	}
	return ""
}

func (m *TwapRecord) GetAsset1Denom() string {
	if m != nil {
		return m.Asset1Denom
	}
	return ""
}

func (m *TwapRecord) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *TwapRecord) GetTime() time.Time {
	if m != nil {
		return m.Time
	}
	return time.Time{}
}

func (m *TwapRecord) GetLastErrorTime() time.Time {
	if m != nil {
		return m.LastErrorTime
	}
	return time.Time{}
}

// PruningState allows us to spread out the pruning of TWAP records over time,
// instead of pruning all at once at the end of the epoch.
type PruningState struct {
	// is_pruning is true if the pruning process is ongoing.
	// This tells the module to continue pruning the TWAP records
	// at the EndBlock.
	IsPruning bool `protobuf:"varint,1,opt,name=is_pruning,json=isPruning,proto3" json:"is_pruning,omitempty"`
	// last_kept_time is the time of the last kept TWAP record.
	// This is used to determine all TWAP records that are older than
	// last_kept_time and should be pruned.
	LastKeptTime time.Time `protobuf:"bytes,2,opt,name=last_kept_time,json=lastKeptTime,proto3,stdtime" json:"last_kept_time" yaml:"last_kept_time"`
	// last_key_seen is the last key of the TWAP records that were pruned
	// before reaching the block's prune limit
	LastKeySeen []byte `protobuf:"bytes,3,opt,name=last_key_seen,json=lastKeySeen,proto3" json:"last_key_seen,omitempty"`
}

func (m *PruningState) Reset()         { *m = PruningState{} }
func (m *PruningState) String() string { return proto.CompactTextString(m) }
func (*PruningState) ProtoMessage()    {}
func (*PruningState) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbf5c78678e601aa, []int{1}
}
func (m *PruningState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PruningState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PruningState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PruningState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PruningState.Merge(m, src)
}
func (m *PruningState) XXX_Size() int {
	return m.Size()
}
func (m *PruningState) XXX_DiscardUnknown() {
	xxx_messageInfo_PruningState.DiscardUnknown(m)
}

var xxx_messageInfo_PruningState proto.InternalMessageInfo

func (m *PruningState) GetIsPruning() bool {
	if m != nil {
		return m.IsPruning
	}
	return false
}

func (m *PruningState) GetLastKeptTime() time.Time {
	if m != nil {
		return m.LastKeptTime
	}
	return time.Time{}
}

func (m *PruningState) GetLastKeySeen() []byte {
	if m != nil {
		return m.LastKeySeen
	}
	return nil
}

func init() {
	proto.RegisterType((*TwapRecord)(nil), "osmosis.twap.v1beta1.TwapRecord")
	proto.RegisterType((*PruningState)(nil), "osmosis.twap.v1beta1.PruningState")
}

func init() {
	proto.RegisterFile("osmosis/twap/v1beta1/twap_record.proto", fileDescriptor_dbf5c78678e601aa)
}

var fileDescriptor_dbf5c78678e601aa = []byte{
	// 624 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4d, 0x4f, 0xdb, 0x40,
	0x10, 0x8d, 0x81, 0x06, 0xd8, 0x84, 0x22, 0x59, 0xb4, 0x75, 0x83, 0x6a, 0x07, 0x57, 0xaa, 0xd2,
	0x43, 0xfd, 0x01, 0xb7, 0xde, 0x88, 0xe8, 0xa1, 0x2d, 0xaa, 0x90, 0xe1, 0xd4, 0x8b, 0xb5, 0x76,
	0x16, 0x7b, 0x85, 0xed, 0x5d, 0x79, 0x37, 0x50, 0xff, 0x0b, 0x7e, 0x4d, 0xef, 0xbd, 0x71, 0xe4,
	0x58, 0xf5, 0x90, 0x56, 0xe4, 0xd6, 0x23, 0xbf, 0xa0, 0xda, 0x5d, 0x27, 0x25, 0xe9, 0x07, 0x70,
	0xf3, 0xcc, 0xbc, 0x79, 0xef, 0x8d, 0x77, 0x76, 0xc1, 0x0b, 0xc2, 0x72, 0xc2, 0x30, 0x73, 0xf9,
	0x19, 0xa4, 0xee, 0xa9, 0x1f, 0x21, 0x0e, 0x7d, 0x19, 0x84, 0x25, 0x8a, 0x49, 0x39, 0x70, 0x68,
	0x49, 0x38, 0xd1, 0x37, 0x6a, 0x9c, 0x23, 0x4a, 0x4e, 0x8d, 0xeb, 0x6c, 0x24, 0x24, 0x21, 0x12,
	0xe0, 0x8a, 0x2f, 0x85, 0xed, 0x3c, 0x4d, 0x08, 0x49, 0x32, 0xe4, 0xca, 0x28, 0x1a, 0x1e, 0xbb,
	0xb0, 0xa8, 0x26, 0xa5, 0x58, 0xf2, 0x84, 0xaa, 0x47, 0x05, 0x75, 0xc9, 0x54, 0x91, 0x1b, 0x41,
	0x86, 0xa6, 0x46, 0x62, 0x82, 0x8b, 0xba, 0x6e, 0xcd, 0xb3, 0x72, 0x9c, 0x23, 0xc6, 0x61, 0x4e,
	0x15, 0xc0, 0xfe, 0xd2, 0x04, 0xe0, 0xe8, 0x0c, 0xd2, 0x40, 0xfa, 0xd6, 0x9f, 0x80, 0x65, 0x4a,
	0x48, 0x16, 0xe2, 0x81, 0xa1, 0x75, 0xb5, 0xde, 0x52, 0xd0, 0x14, 0xe1, 0xdb, 0x81, 0xbe, 0x05,
	0xda, 0x90, 0x31, 0xc4, 0xbd, 0x70, 0x80, 0x0a, 0x92, 0x1b, 0x0b, 0x5d, 0xad, 0xb7, 0x1a, 0xb4,
	0x54, 0x6e, 0x4f, 0xa4, 0xa6, 0x10, 0xbf, 0x86, 0x2c, 0xde, 0x80, 0xf8, 0x0a, 0xb2, 0x0b, 0x9a,
	0x29, 0xc2, 0x49, 0xca, 0x8d, 0xa5, 0xae, 0xd6, 0x5b, 0xec, 0xbf, 0xfc, 0x39, 0xb2, 0xd6, 0xd4,
	0x2f, 0x0b, 0x55, 0xe1, 0x7a, 0x64, 0x6d, 0x54, 0x30, 0xcf, 0x5e, 0xdb, 0x33, 0x69, 0x3b, 0xa8,
	0x1b, 0xf5, 0x0f, 0x60, 0x49, 0xcc, 0x60, 0x3c, 0xe8, 0x6a, 0xbd, 0xd6, 0x76, 0xc7, 0x51, 0x03,
	0x3a, 0x93, 0x01, 0x9d, 0xa3, 0xc9, 0x80, 0x7d, 0xf3, 0x62, 0x64, 0x35, 0xae, 0x47, 0x96, 0x3e,
	0xc3, 0x27, 0x9a, 0xed, 0xf3, 0xef, 0x96, 0x16, 0x48, 0x1e, 0xfd, 0x00, 0xe8, 0xd4, 0x0b, 0x33,
	0xc8, 0x78, 0xc8, 0x28, 0xe1, 0x21, 0x2d, 0x71, 0x8c, 0x8c, 0xa6, 0xf0, 0xde, 0x7f, 0x2e, 0x18,
	0xbe, 0x8d, 0xac, 0x4d, 0xf5, 0x97, 0xd9, 0xe0, 0xc4, 0xc1, 0xc4, 0xcd, 0x21, 0x4f, 0x9d, 0x7d,
	0x94, 0xc0, 0xb8, 0xda, 0x43, 0x71, 0xb0, 0x4e, 0xbd, 0x7d, 0xc8, 0xf8, 0x21, 0x25, 0xfc, 0x40,
	0xf4, 0x4a, 0x46, 0xff, 0x0f, 0xc6, 0xe5, 0xfb, 0x30, 0xfa, 0xb3, 0x8c, 0x29, 0x30, 0xa9, 0x17,
	0xc2, 0x12, 0xf3, 0x34, 0x47, 0x1c, 0xc7, 0xa1, 0x5c, 0x35, 0x18, 0xc7, 0xc3, 0x7c, 0x98, 0x41,
	0x4e, 0x4a, 0x63, 0xe5, 0xee, 0xec, 0x9b, 0xd4, 0xdb, 0x9d, 0x32, 0x89, 0xa3, 0xdf, 0xfd, 0xcd,
	0x23, 0x95, 0xfc, 0xff, 0x2a, 0xad, 0xde, 0x47, 0xc9, 0xff, 0xb7, 0x12, 0x04, 0x9d, 0x04, 0x91,
	0x1c, 0xf1, 0xf2, 0x6f, 0x2a, 0xe0, 0xee, 0x2a, 0xc6, 0x94, 0x66, 0x5e, 0xe2, 0x18, 0xac, 0xcb,
	0x53, 0x40, 0x65, 0x49, 0x4a, 0x79, 0xf0, 0x46, 0xeb, 0xd6, 0xad, 0xb1, 0xeb, 0xad, 0x79, 0xac,
	0xb6, 0x66, 0x8e, 0x40, 0x6d, 0xce, 0x9a, 0xc8, 0xbe, 0x11, 0x49, 0xd1, 0x67, 0x7f, 0xd6, 0x40,
	0xfb, 0xa0, 0x1c, 0x16, 0xb8, 0x48, 0x0e, 0x39, 0xe4, 0x48, 0x7f, 0x06, 0x00, 0x16, 0xd7, 0x55,
	0xa6, 0xe4, 0x45, 0x5a, 0x09, 0x56, 0x31, 0xab, 0x31, 0x7a, 0x0c, 0x1e, 0x4a, 0xda, 0x13, 0x44,
	0xb9, 0xb2, 0xb5, 0x70, 0xab, 0xad, 0xad, 0xda, 0xd6, 0xa3, 0x1b, 0xb6, 0xa6, 0xfd, 0xca, 0x55,
	0x5b, 0x24, 0xdf, 0x23, 0xca, 0x45, 0x97, 0x6e, 0x83, 0xb5, 0x1a, 0x54, 0x85, 0x0c, 0xa1, 0x42,
	0x5e, 0xc7, 0x76, 0xd0, 0x52, 0xa0, 0xea, 0x10, 0xa1, 0xa2, 0xff, 0xee, 0xe2, 0xca, 0xd4, 0x2e,
	0xaf, 0x4c, 0xed, 0xc7, 0x95, 0xa9, 0x9d, 0x8f, 0xcd, 0xc6, 0xe5, 0xd8, 0x6c, 0x7c, 0x1d, 0x9b,
	0x8d, 0x8f, 0x5e, 0x82, 0x79, 0x3a, 0x8c, 0x9c, 0x98, 0xe4, 0x6e, 0xfd, 0x88, 0xbd, 0xca, 0x60,
	0xc4, 0x26, 0x81, 0x7b, 0xba, 0xbd, 0xe3, 0x7e, 0x52, 0xef, 0x1f, 0xaf, 0x28, 0x62, 0x51, 0x53,
	0x9a, 0xde, 0xf9, 0x15, 0x00, 0x00, 0xff, 0xff, 0x59, 0x03, 0x4c, 0x93, 0x1c, 0x05, 0x00, 0x00,
}

func (m *TwapRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TwapRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TwapRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.LastErrorTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.LastErrorTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintTwapRecord(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x5a
	{
		size := m.GeometricTwapAccumulator.Size()
		i -= size
		if _, err := m.GeometricTwapAccumulator.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwapRecord(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.P1ArithmeticTwapAccumulator.Size()
		i -= size
		if _, err := m.P1ArithmeticTwapAccumulator.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwapRecord(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.P0ArithmeticTwapAccumulator.Size()
		i -= size
		if _, err := m.P0ArithmeticTwapAccumulator.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwapRecord(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.P1LastSpotPrice.Size()
		i -= size
		if _, err := m.P1LastSpotPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwapRecord(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.P0LastSpotPrice.Size()
		i -= size
		if _, err := m.P0LastSpotPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwapRecord(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintTwapRecord(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x2a
	if m.Height != 0 {
		i = encodeVarintTwapRecord(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Asset1Denom) > 0 {
		i -= len(m.Asset1Denom)
		copy(dAtA[i:], m.Asset1Denom)
		i = encodeVarintTwapRecord(dAtA, i, uint64(len(m.Asset1Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Asset0Denom) > 0 {
		i -= len(m.Asset0Denom)
		copy(dAtA[i:], m.Asset0Denom)
		i = encodeVarintTwapRecord(dAtA, i, uint64(len(m.Asset0Denom)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintTwapRecord(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PruningState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PruningState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PruningState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.LastKeySeen) > 0 {
		i -= len(m.LastKeySeen)
		copy(dAtA[i:], m.LastKeySeen)
		i = encodeVarintTwapRecord(dAtA, i, uint64(len(m.LastKeySeen)))
		i--
		dAtA[i] = 0x1a
	}
	n3, err3 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.LastKeptTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.LastKeptTime):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintTwapRecord(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x12
	if m.IsPruning {
		i--
		if m.IsPruning {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTwapRecord(dAtA []byte, offset int, v uint64) int {
	offset -= sovTwapRecord(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TwapRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovTwapRecord(uint64(m.PoolId))
	}
	l = len(m.Asset0Denom)
	if l > 0 {
		n += 1 + l + sovTwapRecord(uint64(l))
	}
	l = len(m.Asset1Denom)
	if l > 0 {
		n += 1 + l + sovTwapRecord(uint64(l))
	}
	if m.Height != 0 {
		n += 1 + sovTwapRecord(uint64(m.Height))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovTwapRecord(uint64(l))
	l = m.P0LastSpotPrice.Size()
	n += 1 + l + sovTwapRecord(uint64(l))
	l = m.P1LastSpotPrice.Size()
	n += 1 + l + sovTwapRecord(uint64(l))
	l = m.P0ArithmeticTwapAccumulator.Size()
	n += 1 + l + sovTwapRecord(uint64(l))
	l = m.P1ArithmeticTwapAccumulator.Size()
	n += 1 + l + sovTwapRecord(uint64(l))
	l = m.GeometricTwapAccumulator.Size()
	n += 1 + l + sovTwapRecord(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.LastErrorTime)
	n += 1 + l + sovTwapRecord(uint64(l))
	return n
}

func (m *PruningState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.IsPruning {
		n += 2
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.LastKeptTime)
	n += 1 + l + sovTwapRecord(uint64(l))
	l = len(m.LastKeySeen)
	if l > 0 {
		n += 1 + l + sovTwapRecord(uint64(l))
	}
	return n
}

func sovTwapRecord(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTwapRecord(x uint64) (n int) {
	return sovTwapRecord(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TwapRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTwapRecord
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TwapRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TwapRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset0Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asset0Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset1Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asset1Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field P0LastSpotPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.P0LastSpotPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field P1LastSpotPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.P1LastSpotPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field P0ArithmeticTwapAccumulator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.P0ArithmeticTwapAccumulator.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field P1ArithmeticTwapAccumulator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.P1ArithmeticTwapAccumulator.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GeometricTwapAccumulator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GeometricTwapAccumulator.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastErrorTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.LastErrorTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTwapRecord(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PruningState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTwapRecord
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PruningState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PruningState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsPruning", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsPruning = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastKeptTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.LastKeptTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastKeySeen", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTwapRecord
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastKeySeen = append(m.LastKeySeen[:0], dAtA[iNdEx:postIndex]...)
			if m.LastKeySeen == nil {
				m.LastKeySeen = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTwapRecord(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTwapRecord
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTwapRecord(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTwapRecord
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTwapRecord
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTwapRecord
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTwapRecord
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTwapRecord
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTwapRecord        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTwapRecord          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTwapRecord = fmt.Errorf("proto: unexpected end of group")
)
