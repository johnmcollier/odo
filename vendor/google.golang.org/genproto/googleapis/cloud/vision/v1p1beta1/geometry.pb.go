// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/vision/v1p1beta1/geometry.proto

package vision // import "google.golang.org/genproto/googleapis/cloud/vision/v1p1beta1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A vertex represents a 2D point in the image.
// NOTE: the vertex coordinates are in the same scale as the original image.
type Vertex struct {
	// X coordinate.
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y                    int32    `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vertex) Reset()         { *m = Vertex{} }
func (m *Vertex) String() string { return proto.CompactTextString(m) }
func (*Vertex) ProtoMessage()    {}
func (*Vertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_9a7190aad6b30813, []int{0}
}
func (m *Vertex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vertex.Unmarshal(m, b)
}
func (m *Vertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vertex.Marshal(b, m, deterministic)
}
func (dst *Vertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vertex.Merge(dst, src)
}
func (m *Vertex) XXX_Size() int {
	return xxx_messageInfo_Vertex.Size(m)
}
func (m *Vertex) XXX_DiscardUnknown() {
	xxx_messageInfo_Vertex.DiscardUnknown(m)
}

var xxx_messageInfo_Vertex proto.InternalMessageInfo

func (m *Vertex) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Vertex) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// A bounding polygon for the detected image annotation.
type BoundingPoly struct {
	// The bounding polygon vertices.
	Vertices             []*Vertex `protobuf:"bytes,1,rep,name=vertices,proto3" json:"vertices,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *BoundingPoly) Reset()         { *m = BoundingPoly{} }
func (m *BoundingPoly) String() string { return proto.CompactTextString(m) }
func (*BoundingPoly) ProtoMessage()    {}
func (*BoundingPoly) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_9a7190aad6b30813, []int{1}
}
func (m *BoundingPoly) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BoundingPoly.Unmarshal(m, b)
}
func (m *BoundingPoly) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BoundingPoly.Marshal(b, m, deterministic)
}
func (dst *BoundingPoly) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BoundingPoly.Merge(dst, src)
}
func (m *BoundingPoly) XXX_Size() int {
	return xxx_messageInfo_BoundingPoly.Size(m)
}
func (m *BoundingPoly) XXX_DiscardUnknown() {
	xxx_messageInfo_BoundingPoly.DiscardUnknown(m)
}

var xxx_messageInfo_BoundingPoly proto.InternalMessageInfo

func (m *BoundingPoly) GetVertices() []*Vertex {
	if m != nil {
		return m.Vertices
	}
	return nil
}

// A 3D position in the image, used primarily for Face detection landmarks.
// A valid Position must have both x and y coordinates.
// The position coordinates are in the same scale as the original image.
type Position struct {
	// X coordinate.
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	// Z coordinate (or depth).
	Z                    float32  `protobuf:"fixed32,3,opt,name=z,proto3" json:"z,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_9a7190aad6b30813, []int{2}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Position.Unmarshal(m, b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Position.Marshal(b, m, deterministic)
}
func (dst *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(dst, src)
}
func (m *Position) XXX_Size() int {
	return xxx_messageInfo_Position.Size(m)
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Position) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Position) GetZ() float32 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
	proto.RegisterType((*Vertex)(nil), "google.cloud.vision.v1p1beta1.Vertex")
	proto.RegisterType((*BoundingPoly)(nil), "google.cloud.vision.v1p1beta1.BoundingPoly")
	proto.RegisterType((*Position)(nil), "google.cloud.vision.v1p1beta1.Position")
}

func init() {
	proto.RegisterFile("google/cloud/vision/v1p1beta1/geometry.proto", fileDescriptor_geometry_9a7190aad6b30813)
}

var fileDescriptor_geometry_9a7190aad6b30813 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xb1, 0x4b, 0xc3, 0x40,
	0x14, 0x87, 0x79, 0x29, 0x96, 0x72, 0xd6, 0x25, 0x53, 0x16, 0xa1, 0x06, 0x85, 0x0e, 0x72, 0x47,
	0xd4, 0xcd, 0xc9, 0x38, 0xb8, 0xc6, 0x0c, 0x0e, 0x6e, 0x69, 0xfa, 0x78, 0x1c, 0xa4, 0xf7, 0xc2,
	0xe5, 0x1a, 0x7a, 0xc5, 0x3f, 0xdc, 0x51, 0x7a, 0x57, 0x2a, 0x0e, 0x76, 0xfc, 0xdd, 0x7d, 0x8f,
	0x0f, 0x3e, 0x71, 0x4f, 0xcc, 0xd4, 0xa1, 0x6a, 0x3b, 0xde, 0xae, 0xd5, 0xa8, 0x07, 0xcd, 0x46,
	0x8d, 0x45, 0x5f, 0xac, 0xd0, 0x35, 0x85, 0x22, 0xe4, 0x0d, 0x3a, 0xeb, 0x65, 0x6f, 0xd9, 0x71,
	0x7a, 0x1d, 0x69, 0x19, 0x68, 0x19, 0x69, 0x79, 0xa2, 0xf3, 0x5b, 0x31, 0xfd, 0x40, 0xeb, 0x70,
	0x97, 0xce, 0x05, 0xec, 0x32, 0x58, 0xc0, 0xf2, 0xa2, 0x86, 0xb0, 0x7c, 0x96, 0xc4, 0xe5, 0xf3,
	0x77, 0x31, 0x2f, 0x79, 0x6b, 0xd6, 0xda, 0x50, 0xc5, 0x9d, 0x4f, 0x5f, 0xc4, 0x6c, 0x44, 0xeb,
	0x74, 0x8b, 0x43, 0x06, 0x8b, 0xc9, 0xf2, 0xf2, 0xe1, 0x4e, 0x9e, 0xf5, 0xc8, 0x28, 0xa9, 0x4f,
	0x67, 0xf9, 0x93, 0x98, 0x55, 0x3c, 0x68, 0xa7, 0xd9, 0xfc, 0xaa, 0x93, 0x3f, 0xea, 0xa4, 0x06,
	0x7f, 0x58, 0xfb, 0x6c, 0x12, 0xd7, 0xbe, 0xfc, 0x12, 0x37, 0x2d, 0x6f, 0xce, 0xbb, 0xca, 0xab,
	0xb7, 0x63, 0x82, 0xea, 0x50, 0xa0, 0x82, 0xcf, 0xd7, 0x23, 0x4f, 0xdc, 0x35, 0x86, 0x24, 0x5b,
	0x52, 0x84, 0x26, 0xf4, 0x51, 0xf1, 0xab, 0xe9, 0xf5, 0xf0, 0x4f, 0xd0, 0xe7, 0xf8, 0xf0, 0x0d,
	0xb0, 0x9a, 0x86, 0x93, 0xc7, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x91, 0xa5, 0x86, 0xce, 0x82,
	0x01, 0x00, 0x00,
}
