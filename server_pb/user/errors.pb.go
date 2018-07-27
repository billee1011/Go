// Code generated by protoc-gen-go. DO NOT EDIT.
// source: errors.proto

package user

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ErrCode int32

const (
	ErrCode_EC_SUCCESS ErrCode = 0
	ErrCode_EC_FAIL    ErrCode = 1
)

var ErrCode_name = map[int32]string{
	0: "EC_SUCCESS",
	1: "EC_FAIL",
}
var ErrCode_value = map[string]int32{
	"EC_SUCCESS": 0,
	"EC_FAIL":    1,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}
func (ErrCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func init() {
	proto.RegisterEnum("user.ErrCode", ErrCode_name, ErrCode_value)
}

func init() { proto.RegisterFile("errors.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 89 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x2d, 0x2a, 0xca,
	0x2f, 0x2a, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x2d, 0x4e, 0x2d, 0xd2, 0x52,
	0xe3, 0x62, 0x77, 0x2d, 0x2a, 0x72, 0xce, 0x4f, 0x49, 0x15, 0xe2, 0xe3, 0xe2, 0x72, 0x75, 0x8e,
	0x0f, 0x0e, 0x75, 0x76, 0x76, 0x0d, 0x0e, 0x16, 0x60, 0x10, 0xe2, 0xe6, 0x62, 0x77, 0x75, 0x8e,
	0x77, 0x73, 0xf4, 0xf4, 0x11, 0x60, 0x4c, 0x62, 0x03, 0x6b, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff,
	0xff, 0xbc, 0xfe, 0x34, 0x9d, 0x44, 0x00, 0x00, 0x00,
}
