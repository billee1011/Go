// Code generated by protoc-gen-go. DO NOT EDIT.
// source: errors.proto

package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// ErrCode 错误码
type ErrCode int32

const (
	ErrCode_EC_SUCCESS ErrCode = 0
	ErrCode_EC_FAIL    ErrCode = 1
	// 匹配错误
	ErrCode_EC_MATCH_ALREADY_GAMEING ErrCode = 257
)

var ErrCode_name = map[int32]string{
	0:   "EC_SUCCESS",
	1:   "EC_FAIL",
	257: "EC_MATCH_ALREADY_GAMEING",
}
var ErrCode_value = map[string]int32{
	"EC_SUCCESS":               0,
	"EC_FAIL":                  1,
	"EC_MATCH_ALREADY_GAMEING": 257,
}

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}
func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}
func (x *ErrCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ErrCode_value, data, "ErrCode")
	if err != nil {
		return err
	}
	*x = ErrCode(value)
	return nil
}
func (ErrCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func init() {
	proto.RegisterEnum("common.ErrCode", ErrCode_name, ErrCode_value)
}

func init() { proto.RegisterFile("errors.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x2d, 0x2a, 0xca,
	0x2f, 0x2a, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4b, 0xce, 0xcf, 0xcd, 0xcd, 0xcf,
	0xd3, 0x72, 0xe5, 0x62, 0x77, 0x2d, 0x2a, 0x72, 0xce, 0x4f, 0x49, 0x15, 0xe2, 0xe3, 0xe2, 0x72,
	0x75, 0x8e, 0x0f, 0x0e, 0x75, 0x76, 0x76, 0x0d, 0x0e, 0x16, 0x60, 0x10, 0xe2, 0xe6, 0x62, 0x77,
	0x75, 0x8e, 0x77, 0x73, 0xf4, 0xf4, 0x11, 0x60, 0x14, 0x92, 0xe5, 0x92, 0x70, 0x75, 0x8e, 0xf7,
	0x75, 0x0c, 0x71, 0xf6, 0x88, 0x77, 0xf4, 0x09, 0x72, 0x75, 0x74, 0x89, 0x8c, 0x77, 0x77, 0xf4,
	0x75, 0xf5, 0xf4, 0x73, 0x17, 0x68, 0x64, 0x72, 0x92, 0x88, 0x12, 0x2b, 0x2e, 0x49, 0x2d, 0x4b,
	0xd5, 0x4f, 0xce, 0xc9, 0x4c, 0xcd, 0x2b, 0x89, 0x2f, 0x48, 0xd2, 0x87, 0x58, 0x00, 0x08, 0x00,
	0x00, 0xff, 0xff, 0x0e, 0x3f, 0xdd, 0xaa, 0x77, 0x00, 0x00, 0x00,
}