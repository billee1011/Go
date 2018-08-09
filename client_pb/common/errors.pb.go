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
	// 身份认证
	ErrCode_EC_INVALID_ID_CARD_OR_NAME ErrCode = 288
	ErrCode_EC_REAL_NAME_ALREADY       ErrCode = 289
)

var ErrCode_name = map[int32]string{
	0:   "EC_SUCCESS",
	1:   "EC_FAIL",
	257: "EC_MATCH_ALREADY_GAMEING",
	288: "EC_INVALID_ID_CARD_OR_NAME",
	289: "EC_REAL_NAME_ALREADY",
}
var ErrCode_value = map[string]int32{
	"EC_SUCCESS":                 0,
	"EC_FAIL":                    1,
	"EC_MATCH_ALREADY_GAMEING":   257,
	"EC_INVALID_ID_CARD_OR_NAME": 288,
	"EC_REAL_NAME_ALREADY":       289,
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

// Result 通用处理结果
type Result struct {
	ErrCode          *ErrCode `protobuf:"varint,1,opt,name=err_code,json=errCode,enum=common.ErrCode" json:"err_code,omitempty"`
	ErrDesc          *string  `protobuf:"bytes,2,opt,name=err_desc,json=errDesc" json:"err_desc,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Result) GetErrCode() ErrCode {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return ErrCode_EC_SUCCESS
}

func (m *Result) GetErrDesc() string {
	if m != nil && m.ErrDesc != nil {
		return *m.ErrDesc
	}
	return ""
}

func init() {
	proto.RegisterType((*Result)(nil), "common.Result")
	proto.RegisterEnum("common.ErrCode", ErrCode_name, ErrCode_value)
}

func init() { proto.RegisterFile("errors.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0xcf, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x06, 0x60, 0x93, 0x43, 0xa3, 0xa3, 0xd4, 0xb0, 0x88, 0xa4, 0x82, 0x58, 0x3c, 0x95, 0x1e,
	0x52, 0xf0, 0x0d, 0xd6, 0xc9, 0x58, 0x03, 0x49, 0x0a, 0x1b, 0x15, 0xf4, 0x32, 0xe0, 0x66, 0x0e,
	0x42, 0xdb, 0x2d, 0xbb, 0xd1, 0x7b, 0xdf, 0x44, 0xdf, 0x54, 0x68, 0xf4, 0x36, 0xcc, 0xff, 0xf1,
	0xc3, 0x0f, 0x67, 0xe2, 0xbd, 0xf3, 0x21, 0xdf, 0x79, 0xd7, 0x3b, 0x35, 0xb2, 0x6e, 0xb3, 0x71,
	0xdb, 0xdb, 0x15, 0x8c, 0x8c, 0x84, 0xcf, 0x75, 0xaf, 0xe6, 0x70, 0x2c, 0xde, 0xb3, 0x75, 0x9d,
	0x64, 0xd1, 0x34, 0x9a, 0x8d, 0xef, 0xce, 0xf3, 0x01, 0xe5, 0xe4, 0x3d, 0xba, 0x4e, 0x4c, 0x22,
	0xc3, 0xa1, 0x26, 0x83, 0xed, 0x24, 0xd8, 0x2c, 0x9e, 0x46, 0xb3, 0x93, 0x43, 0x54, 0x48, 0xb0,
	0xf3, 0x7d, 0x04, 0xc9, 0x9f, 0x57, 0x63, 0x00, 0x42, 0x6e, 0x9f, 0x11, 0xa9, 0x6d, 0xd3, 0x23,
	0x75, 0x0a, 0x09, 0x21, 0x3f, 0xe8, 0xb2, 0x4a, 0x23, 0x75, 0x0d, 0x19, 0x21, 0xd7, 0xfa, 0x09,
	0x1f, 0x59, 0x57, 0x86, 0x74, 0xf1, 0xca, 0x4b, 0x5d, 0x53, 0xd9, 0x2c, 0xd3, 0x7d, 0xac, 0x6e,
	0xe0, 0x8a, 0x90, 0xcb, 0xe6, 0x45, 0x57, 0x65, 0xc1, 0x65, 0xc1, 0xa8, 0x4d, 0xc1, 0x2b, 0xc3,
	0x8d, 0xae, 0x29, 0xfd, 0x8e, 0xd5, 0x04, 0x2e, 0x08, 0xd9, 0x90, 0xae, 0x0e, 0xaf, 0xff, 0x8e,
	0xf4, 0x27, 0xbe, 0xcf, 0xde, 0x2e, 0x43, 0x2f, 0x5f, 0xb2, 0xb0, 0xeb, 0x0f, 0xd9, 0xf6, 0xbc,
	0x7b, 0x5f, 0x0c, 0x4b, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xdd, 0x95, 0x1d, 0x45, 0x05, 0x01,
	0x00, 0x00,
}
