// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/billing_metrics.proto

package storage

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

type BillingMetricsRecord struct {
	Ts                   *types.Timestamp                       `protobuf:"bytes,1,opt,name=ts,proto3" json:"ts,omitempty" sql:"pk"`
	Sr                   *BillingMetricsRecord_SecuredResources `protobuf:"bytes,2,opt,name=sr,proto3" json:"sr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *BillingMetricsRecord) Reset()         { *m = BillingMetricsRecord{} }
func (m *BillingMetricsRecord) String() string { return proto.CompactTextString(m) }
func (*BillingMetricsRecord) ProtoMessage()    {}
func (*BillingMetricsRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_338da3bb08fef41d, []int{0}
}
func (m *BillingMetricsRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BillingMetricsRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BillingMetricsRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BillingMetricsRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BillingMetricsRecord.Merge(m, src)
}
func (m *BillingMetricsRecord) XXX_Size() int {
	return m.Size()
}
func (m *BillingMetricsRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_BillingMetricsRecord.DiscardUnknown(m)
}

var xxx_messageInfo_BillingMetricsRecord proto.InternalMessageInfo

func (m *BillingMetricsRecord) GetTs() *types.Timestamp {
	if m != nil {
		return m.Ts
	}
	return nil
}

func (m *BillingMetricsRecord) GetSr() *BillingMetricsRecord_SecuredResources {
	if m != nil {
		return m.Sr
	}
	return nil
}

func (m *BillingMetricsRecord) MessageClone() proto.Message {
	return m.Clone()
}
func (m *BillingMetricsRecord) Clone() *BillingMetricsRecord {
	if m == nil {
		return nil
	}
	cloned := new(BillingMetricsRecord)
	*cloned = *m

	cloned.Ts = m.Ts.Clone()
	cloned.Sr = m.Sr.Clone()
	return cloned
}

type BillingMetricsRecord_SecuredResources struct {
	Nodes                int32    `protobuf:"varint,1,opt,name=nodes,proto3" json:"nodes,omitempty"`
	Millicores           int32    `protobuf:"varint,2,opt,name=millicores,proto3" json:"millicores,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BillingMetricsRecord_SecuredResources) Reset()         { *m = BillingMetricsRecord_SecuredResources{} }
func (m *BillingMetricsRecord_SecuredResources) String() string { return proto.CompactTextString(m) }
func (*BillingMetricsRecord_SecuredResources) ProtoMessage()    {}
func (*BillingMetricsRecord_SecuredResources) Descriptor() ([]byte, []int) {
	return fileDescriptor_338da3bb08fef41d, []int{0, 0}
}
func (m *BillingMetricsRecord_SecuredResources) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BillingMetricsRecord_SecuredResources) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BillingMetricsRecord_SecuredResources.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BillingMetricsRecord_SecuredResources) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BillingMetricsRecord_SecuredResources.Merge(m, src)
}
func (m *BillingMetricsRecord_SecuredResources) XXX_Size() int {
	return m.Size()
}
func (m *BillingMetricsRecord_SecuredResources) XXX_DiscardUnknown() {
	xxx_messageInfo_BillingMetricsRecord_SecuredResources.DiscardUnknown(m)
}

var xxx_messageInfo_BillingMetricsRecord_SecuredResources proto.InternalMessageInfo

func (m *BillingMetricsRecord_SecuredResources) GetNodes() int32 {
	if m != nil {
		return m.Nodes
	}
	return 0
}

func (m *BillingMetricsRecord_SecuredResources) GetMillicores() int32 {
	if m != nil {
		return m.Millicores
	}
	return 0
}

func (m *BillingMetricsRecord_SecuredResources) MessageClone() proto.Message {
	return m.Clone()
}
func (m *BillingMetricsRecord_SecuredResources) Clone() *BillingMetricsRecord_SecuredResources {
	if m == nil {
		return nil
	}
	cloned := new(BillingMetricsRecord_SecuredResources)
	*cloned = *m

	return cloned
}

func init() {
	proto.RegisterType((*BillingMetricsRecord)(nil), "storage.BillingMetricsRecord")
	proto.RegisterType((*BillingMetricsRecord_SecuredResources)(nil), "storage.BillingMetricsRecord.SecuredResources")
}

func init() { proto.RegisterFile("storage/billing_metrics.proto", fileDescriptor_338da3bb08fef41d) }

var fileDescriptor_338da3bb08fef41d = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0x2e, 0xc9, 0x2f,
	0x4a, 0x4c, 0x4f, 0xd5, 0x4f, 0xca, 0xcc, 0xc9, 0xc9, 0xcc, 0x4b, 0x8f, 0xcf, 0x4d, 0x2d, 0x29,
	0xca, 0x4c, 0x2e, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x4a, 0x4b, 0x89, 0xa4,
	0xe7, 0xa7, 0xe7, 0x83, 0xc5, 0xf4, 0x41, 0x2c, 0x88, 0xb4, 0x94, 0x7c, 0x7a, 0x7e, 0x7e, 0x7a,
	0x4e, 0xaa, 0x3e, 0x98, 0x97, 0x54, 0x9a, 0xa6, 0x5f, 0x92, 0x99, 0x9b, 0x5a, 0x5c, 0x92, 0x98,
	0x5b, 0x00, 0x51, 0xa0, 0x74, 0x8b, 0x91, 0x4b, 0xc4, 0x09, 0x62, 0xb2, 0x2f, 0xc4, 0xe0, 0xa0,
	0xd4, 0xe4, 0xfc, 0xa2, 0x14, 0x21, 0x0b, 0x2e, 0xa6, 0x92, 0x62, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0x6e, 0x23, 0x29, 0x3d, 0x88, 0x31, 0x7a, 0x30, 0x63, 0xf4, 0x42, 0x60, 0xc6, 0x38, 0xf1, 0x7c,
	0xba, 0x27, 0xcf, 0x51, 0x5c, 0x98, 0x63, 0xa5, 0x54, 0x90, 0xad, 0x14, 0xc4, 0x54, 0x52, 0x2c,
	0x64, 0xc7, 0xc5, 0x54, 0x5c, 0x24, 0xc1, 0x04, 0xd6, 0xa9, 0xa7, 0x07, 0x75, 0x9f, 0x1e, 0x36,
	0x4b, 0xf4, 0x82, 0x53, 0x93, 0x4b, 0x8b, 0x52, 0x53, 0x82, 0x52, 0x8b, 0xf3, 0x4b, 0x8b, 0x92,
	0x53, 0x8b, 0x83, 0x98, 0x8a, 0x8b, 0xa4, 0x3c, 0xb8, 0x04, 0xd0, 0xc5, 0x85, 0x44, 0xb8, 0x58,
	0xf3, 0xf2, 0x53, 0x52, 0x21, 0x0e, 0x62, 0x0d, 0x82, 0x70, 0x84, 0xe4, 0xb8, 0xb8, 0x72, 0x41,
	0xc6, 0x26, 0xe7, 0x17, 0xa5, 0x16, 0x83, 0x6d, 0x64, 0x0d, 0x42, 0x12, 0x71, 0x92, 0x3c, 0xf1,
	0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x67, 0x3c, 0x96, 0x63, 0x88,
	0x82, 0x05, 0x57, 0x12, 0x1b, 0xd8, 0x2b, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1f, 0x1e,
	0x6c, 0x1c, 0x5f, 0x01, 0x00, 0x00,
}

func (m *BillingMetricsRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BillingMetricsRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BillingMetricsRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Sr != nil {
		{
			size, err := m.Sr.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBillingMetrics(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Ts != nil {
		{
			size, err := m.Ts.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBillingMetrics(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BillingMetricsRecord_SecuredResources) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BillingMetricsRecord_SecuredResources) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BillingMetricsRecord_SecuredResources) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Millicores != 0 {
		i = encodeVarintBillingMetrics(dAtA, i, uint64(m.Millicores))
		i--
		dAtA[i] = 0x10
	}
	if m.Nodes != 0 {
		i = encodeVarintBillingMetrics(dAtA, i, uint64(m.Nodes))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintBillingMetrics(dAtA []byte, offset int, v uint64) int {
	offset -= sovBillingMetrics(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BillingMetricsRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Ts != nil {
		l = m.Ts.Size()
		n += 1 + l + sovBillingMetrics(uint64(l))
	}
	if m.Sr != nil {
		l = m.Sr.Size()
		n += 1 + l + sovBillingMetrics(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *BillingMetricsRecord_SecuredResources) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Nodes != 0 {
		n += 1 + sovBillingMetrics(uint64(m.Nodes))
	}
	if m.Millicores != 0 {
		n += 1 + sovBillingMetrics(uint64(m.Millicores))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBillingMetrics(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBillingMetrics(x uint64) (n int) {
	return sovBillingMetrics(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BillingMetricsRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBillingMetrics
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
			return fmt.Errorf("proto: BillingMetricsRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BillingMetricsRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBillingMetrics
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
				return ErrInvalidLengthBillingMetrics
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBillingMetrics
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Ts == nil {
				m.Ts = &types.Timestamp{}
			}
			if err := m.Ts.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sr", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBillingMetrics
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
				return ErrInvalidLengthBillingMetrics
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBillingMetrics
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Sr == nil {
				m.Sr = &BillingMetricsRecord_SecuredResources{}
			}
			if err := m.Sr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBillingMetrics(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBillingMetrics
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BillingMetricsRecord_SecuredResources) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBillingMetrics
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
			return fmt.Errorf("proto: SecuredResources: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SecuredResources: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nodes", wireType)
			}
			m.Nodes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBillingMetrics
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nodes |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Millicores", wireType)
			}
			m.Millicores = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBillingMetrics
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Millicores |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBillingMetrics(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBillingMetrics
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBillingMetrics(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBillingMetrics
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
					return 0, ErrIntOverflowBillingMetrics
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
					return 0, ErrIntOverflowBillingMetrics
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
				return 0, ErrInvalidLengthBillingMetrics
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBillingMetrics
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBillingMetrics
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBillingMetrics        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBillingMetrics          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBillingMetrics = fmt.Errorf("proto: unexpected end of group")
)
