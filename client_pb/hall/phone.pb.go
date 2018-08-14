// Code generated by protoc-gen-go. DO NOT EDIT.
// source: phone.proto

package hall

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common1 "steve/client_pb/common"
import common "steve/client_pb/common"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 验证码发送场景
type AuthCodeSendScene int32

const (
	AuthCodeSendScene_REGISTER        AuthCodeSendScene = 1
	AuthCodeSendScene_RESET_PASSWORD  AuthCodeSendScene = 2
	AuthCodeSendScene_RESET_CELLPHONE AuthCodeSendScene = 3
	AuthCodeSendScene_NEW_HALL_LOGIN  AuthCodeSendScene = 4
	AuthCodeSendScene_BIND_PHONE      AuthCodeSendScene = 5
	AuthCodeSendScene_BIND_WECHAT     AuthCodeSendScene = 10
)

var AuthCodeSendScene_name = map[int32]string{
	1:  "REGISTER",
	2:  "RESET_PASSWORD",
	3:  "RESET_CELLPHONE",
	4:  "NEW_HALL_LOGIN",
	5:  "BIND_PHONE",
	10: "BIND_WECHAT",
}
var AuthCodeSendScene_value = map[string]int32{
	"REGISTER":        1,
	"RESET_PASSWORD":  2,
	"RESET_CELLPHONE": 3,
	"NEW_HALL_LOGIN":  4,
	"BIND_PHONE":      5,
	"BIND_WECHAT":     10,
}

func (x AuthCodeSendScene) Enum() *AuthCodeSendScene {
	p := new(AuthCodeSendScene)
	*p = x
	return p
}
func (x AuthCodeSendScene) String() string {
	return proto.EnumName(AuthCodeSendScene_name, int32(x))
}
func (x *AuthCodeSendScene) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AuthCodeSendScene_value, data, "AuthCodeSendScene")
	if err != nil {
		return err
	}
	*x = AuthCodeSendScene(value)
	return nil
}
func (AuthCodeSendScene) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

// 发送验证码
type AuthCodeReq struct {
	CellphoneNum     *uint64            `protobuf:"varint,1,opt,name=cellphone_num,json=cellphoneNum" json:"cellphone_num,omitempty"`
	SendCase         *AuthCodeSendScene `protobuf:"varint,2,opt,name=send_case,json=sendCase,enum=hall.AuthCodeSendScene" json:"send_case,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *AuthCodeReq) Reset()                    { *m = AuthCodeReq{} }
func (m *AuthCodeReq) String() string            { return proto.CompactTextString(m) }
func (*AuthCodeReq) ProtoMessage()               {}
func (*AuthCodeReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *AuthCodeReq) GetCellphoneNum() uint64 {
	if m != nil && m.CellphoneNum != nil {
		return *m.CellphoneNum
	}
	return 0
}

func (m *AuthCodeReq) GetSendCase() AuthCodeSendScene {
	if m != nil && m.SendCase != nil {
		return *m.SendCase
	}
	return AuthCodeSendScene_REGISTER
}

// 验证码接收
type AuthCodeRsp struct {
	ErrorCode        *uint64 `protobuf:"varint,1,opt,name=ErrorCode" json:"ErrorCode,omitempty"`
	ErrorMsg         *string `protobuf:"bytes,2,opt,name=ErrorMsg" json:"ErrorMsg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AuthCodeRsp) Reset()                    { *m = AuthCodeRsp{} }
func (m *AuthCodeRsp) String() string            { return proto.CompactTextString(m) }
func (*AuthCodeRsp) ProtoMessage()               {}
func (*AuthCodeRsp) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *AuthCodeRsp) GetErrorCode() uint64 {
	if m != nil && m.ErrorCode != nil {
		return *m.ErrorCode
	}
	return 0
}

func (m *AuthCodeRsp) GetErrorMsg() string {
	if m != nil && m.ErrorMsg != nil {
		return *m.ErrorMsg
	}
	return ""
}

// CheckAuthCodeReq 校验验证码请求
type CheckAuthCodeReq struct {
	SendCase         *AuthCodeSendScene `protobuf:"varint,1,opt,name=send_case,json=sendCase,enum=hall.AuthCodeSendScene" json:"send_case,omitempty"`
	Code             *string            `protobuf:"bytes,2,opt,name=code" json:"code,omitempty"`
	Phone            *string            `protobuf:"bytes,3,opt,name=phone" json:"phone,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *CheckAuthCodeReq) Reset()                    { *m = CheckAuthCodeReq{} }
func (m *CheckAuthCodeReq) String() string            { return proto.CompactTextString(m) }
func (*CheckAuthCodeReq) ProtoMessage()               {}
func (*CheckAuthCodeReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *CheckAuthCodeReq) GetSendCase() AuthCodeSendScene {
	if m != nil && m.SendCase != nil {
		return *m.SendCase
	}
	return AuthCodeSendScene_REGISTER
}

func (m *CheckAuthCodeReq) GetCode() string {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ""
}

func (m *CheckAuthCodeReq) GetPhone() string {
	if m != nil && m.Phone != nil {
		return *m.Phone
	}
	return ""
}

// CheckAuthCodeRsp 校验验证码应答
type CheckAuthCodeRsp struct {
	Result           *common.Result `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *CheckAuthCodeRsp) Reset()                    { *m = CheckAuthCodeRsp{} }
func (m *CheckAuthCodeRsp) String() string            { return proto.CompactTextString(m) }
func (*CheckAuthCodeRsp) ProtoMessage()               {}
func (*CheckAuthCodeRsp) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *CheckAuthCodeRsp) GetResult() *common.Result {
	if m != nil {
		return m.Result
	}
	return nil
}

// GetBindPhoneRewardInfoReq 获取绑定手机可获得的奖励信息请求
type GetBindPhoneRewardInfoReq struct {
	Reserve          *uint32 `protobuf:"varint,1,opt,name=reserve" json:"reserve,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GetBindPhoneRewardInfoReq) Reset()                    { *m = GetBindPhoneRewardInfoReq{} }
func (m *GetBindPhoneRewardInfoReq) String() string            { return proto.CompactTextString(m) }
func (*GetBindPhoneRewardInfoReq) ProtoMessage()               {}
func (*GetBindPhoneRewardInfoReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *GetBindPhoneRewardInfoReq) GetReserve() uint32 {
	if m != nil && m.Reserve != nil {
		return *m.Reserve
	}
	return 0
}

// GetBindPhoneRewardInfoRsp 获取绑定手机可获得的奖励信息应答
type GetBindPhoneRewardInfoRsp struct {
	Reward           *common1.Money `protobuf:"bytes,1,opt,name=reward" json:"reward,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *GetBindPhoneRewardInfoRsp) Reset()                    { *m = GetBindPhoneRewardInfoRsp{} }
func (m *GetBindPhoneRewardInfoRsp) String() string            { return proto.CompactTextString(m) }
func (*GetBindPhoneRewardInfoRsp) ProtoMessage()               {}
func (*GetBindPhoneRewardInfoRsp) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *GetBindPhoneRewardInfoRsp) GetReward() *common1.Money {
	if m != nil {
		return m.Reward
	}
	return nil
}

// BindPhoneReq 绑定手机请求
type BindPhoneReq struct {
	Phone            *string `protobuf:"bytes,1,opt,name=phone" json:"phone,omitempty"`
	DymcCode         *string `protobuf:"bytes,2,opt,name=dymc_code,json=dymcCode" json:"dymc_code,omitempty"`
	Passwd           *string `protobuf:"bytes,3,opt,name=passwd" json:"passwd,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *BindPhoneReq) Reset()                    { *m = BindPhoneReq{} }
func (m *BindPhoneReq) String() string            { return proto.CompactTextString(m) }
func (*BindPhoneReq) ProtoMessage()               {}
func (*BindPhoneReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *BindPhoneReq) GetPhone() string {
	if m != nil && m.Phone != nil {
		return *m.Phone
	}
	return ""
}

func (m *BindPhoneReq) GetDymcCode() string {
	if m != nil && m.DymcCode != nil {
		return *m.DymcCode
	}
	return ""
}

func (m *BindPhoneReq) GetPasswd() string {
	if m != nil && m.Passwd != nil {
		return *m.Passwd
	}
	return ""
}

// BindPhoneRsp 绑定手机应答
type BindPhoneRsp struct {
	Result           *common.Result `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Reward           *common1.Money `protobuf:"bytes,2,opt,name=reward" json:"reward,omitempty"`
	NewMoney         *common1.Money `protobuf:"bytes,3,opt,name=new_money,json=newMoney" json:"new_money,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *BindPhoneRsp) Reset()                    { *m = BindPhoneRsp{} }
func (m *BindPhoneRsp) String() string            { return proto.CompactTextString(m) }
func (*BindPhoneRsp) ProtoMessage()               {}
func (*BindPhoneRsp) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *BindPhoneRsp) GetResult() *common.Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *BindPhoneRsp) GetReward() *common1.Money {
	if m != nil {
		return m.Reward
	}
	return nil
}

func (m *BindPhoneRsp) GetNewMoney() *common1.Money {
	if m != nil {
		return m.NewMoney
	}
	return nil
}

// ChangePhoneReq 修改手机请求
type ChangePhoneReq struct {
	OldPhone         *string `protobuf:"bytes,1,opt,name=old_phone,json=oldPhone" json:"old_phone,omitempty"`
	OldPhoneCode     *string `protobuf:"bytes,2,opt,name=old_phone_code,json=oldPhoneCode" json:"old_phone_code,omitempty"`
	NewPhone         *string `protobuf:"bytes,3,opt,name=new_phone,json=newPhone" json:"new_phone,omitempty"`
	NewPhoneCode     *string `protobuf:"bytes,4,opt,name=new_phone_code,json=newPhoneCode" json:"new_phone_code,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ChangePhoneReq) Reset()                    { *m = ChangePhoneReq{} }
func (m *ChangePhoneReq) String() string            { return proto.CompactTextString(m) }
func (*ChangePhoneReq) ProtoMessage()               {}
func (*ChangePhoneReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *ChangePhoneReq) GetOldPhone() string {
	if m != nil && m.OldPhone != nil {
		return *m.OldPhone
	}
	return ""
}

func (m *ChangePhoneReq) GetOldPhoneCode() string {
	if m != nil && m.OldPhoneCode != nil {
		return *m.OldPhoneCode
	}
	return ""
}

func (m *ChangePhoneReq) GetNewPhone() string {
	if m != nil && m.NewPhone != nil {
		return *m.NewPhone
	}
	return ""
}

func (m *ChangePhoneReq) GetNewPhoneCode() string {
	if m != nil && m.NewPhoneCode != nil {
		return *m.NewPhoneCode
	}
	return ""
}

// ChangePhoneRsp 修改手机应答
type ChangePhoneRsp struct {
	Result           *common.Result `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *ChangePhoneRsp) Reset()                    { *m = ChangePhoneRsp{} }
func (m *ChangePhoneRsp) String() string            { return proto.CompactTextString(m) }
func (*ChangePhoneRsp) ProtoMessage()               {}
func (*ChangePhoneRsp) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

func (m *ChangePhoneRsp) GetResult() *common.Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*AuthCodeReq)(nil), "hall.AuthCodeReq")
	proto.RegisterType((*AuthCodeRsp)(nil), "hall.AuthCodeRsp")
	proto.RegisterType((*CheckAuthCodeReq)(nil), "hall.CheckAuthCodeReq")
	proto.RegisterType((*CheckAuthCodeRsp)(nil), "hall.CheckAuthCodeRsp")
	proto.RegisterType((*GetBindPhoneRewardInfoReq)(nil), "hall.GetBindPhoneRewardInfoReq")
	proto.RegisterType((*GetBindPhoneRewardInfoRsp)(nil), "hall.GetBindPhoneRewardInfoRsp")
	proto.RegisterType((*BindPhoneReq)(nil), "hall.BindPhoneReq")
	proto.RegisterType((*BindPhoneRsp)(nil), "hall.BindPhoneRsp")
	proto.RegisterType((*ChangePhoneReq)(nil), "hall.ChangePhoneReq")
	proto.RegisterType((*ChangePhoneRsp)(nil), "hall.ChangePhoneRsp")
	proto.RegisterEnum("hall.AuthCodeSendScene", AuthCodeSendScene_name, AuthCodeSendScene_value)
}

func init() { proto.RegisterFile("phone.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 543 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xdf, 0x4f, 0x9b, 0x50,
	0x14, 0x0e, 0x5a, 0x1d, 0x9c, 0x22, 0xb2, 0x3b, 0xe3, 0x3a, 0xdd, 0x83, 0x61, 0x3f, 0x62, 0x7c,
	0xa8, 0x89, 0xd9, 0x92, 0x65, 0x6f, 0x2d, 0x92, 0xb6, 0x49, 0x6d, 0x9b, 0x4b, 0x93, 0x66, 0x7b,
	0x21, 0x0c, 0xce, 0xc4, 0x0c, 0x2e, 0xc8, 0xa5, 0x36, 0x3e, 0xed, 0x79, 0xff, 0xc0, 0xfe, 0xde,
	0xe5, 0x5e, 0x28, 0xe2, 0xcc, 0x96, 0xbe, 0xf1, 0x7d, 0xe7, 0xf4, 0xfb, 0x71, 0x0a, 0xd0, 0xce,
	0xa2, 0x94, 0x61, 0x37, 0xcb, 0xd3, 0x22, 0x25, 0xad, 0xc8, 0x8f, 0xe3, 0x23, 0x3d, 0x48, 0x93,
	0x24, 0x65, 0x25, 0x77, 0xa4, 0x63, 0x9e, 0xa7, 0x39, 0x2f, 0x91, 0x15, 0x41, 0xbb, 0xb7, 0x2c,
	0x22, 0x3b, 0x0d, 0x91, 0xe2, 0x2d, 0x79, 0x03, 0x7b, 0x01, 0xc6, 0xb1, 0xd4, 0xf0, 0xd8, 0x32,
	0xe9, 0x28, 0x27, 0xca, 0x69, 0x8b, 0xea, 0x35, 0x39, 0x59, 0x26, 0xe4, 0x03, 0x68, 0x1c, 0x59,
	0xe8, 0x05, 0x3e, 0xc7, 0xce, 0xd6, 0x89, 0x72, 0x6a, 0x5c, 0xbc, 0xec, 0x0a, 0xa7, 0xee, 0x5a,
	0xca, 0x45, 0x16, 0xba, 0x01, 0x32, 0xa4, 0xaa, 0xd8, 0xb4, 0x7d, 0x8e, 0xd6, 0xa0, 0xe1, 0xc4,
	0x33, 0xf2, 0x1a, 0x34, 0x47, 0x04, 0x11, 0xb8, 0x72, 0x79, 0x20, 0xc8, 0x11, 0xa8, 0x12, 0x5c,
	0xf1, 0x6b, 0xe9, 0xa0, 0xd1, 0x1a, 0x5b, 0x39, 0x98, 0x76, 0x84, 0xc1, 0x8f, 0x66, 0xee, 0x47,
	0x91, 0x94, 0x0d, 0x23, 0x11, 0x02, 0xad, 0x40, 0xd8, 0x97, 0x0e, 0xf2, 0x99, 0x1c, 0xc0, 0x8e,
	0x2c, 0xda, 0xd9, 0x96, 0x64, 0x09, 0xac, 0xcf, 0x7f, 0x7b, 0xf2, 0x8c, 0xbc, 0x87, 0xdd, 0x1c,
	0xf9, 0x32, 0x2e, 0xa4, 0x61, 0xfb, 0xc2, 0xe8, 0x56, 0x77, 0xa6, 0x92, 0xa5, 0xd5, 0xd4, 0xfa,
	0x08, 0xaf, 0x06, 0x58, 0xf4, 0x6f, 0x58, 0x38, 0x13, 0x5a, 0x14, 0x57, 0x7e, 0x1e, 0x8e, 0xd8,
	0xf7, 0x54, 0x04, 0xef, 0xc0, 0xb3, 0x1c, 0x39, 0xe6, 0x77, 0x65, 0xec, 0x3d, 0xba, 0x86, 0x56,
	0xff, 0x9f, 0x3f, 0xe3, 0x19, 0x79, 0x27, 0xbc, 0x05, 0x51, 0x79, 0xef, 0xad, 0xbd, 0xaf, 0x52,
	0x86, 0xf7, 0xb4, 0x1a, 0x5a, 0x5f, 0x40, 0x6f, 0x08, 0xdc, 0x3e, 0x94, 0x53, 0x1a, 0xe5, 0xc8,
	0x31, 0x68, 0xe1, 0x7d, 0x12, 0x78, 0x8d, 0x5b, 0xa8, 0x82, 0x90, 0xff, 0xc4, 0x21, 0xec, 0x66,
	0x3e, 0xe7, 0xab, 0xb0, 0x3a, 0x48, 0x85, 0xac, 0x5f, 0x4a, 0x53, 0x7b, 0xf3, 0x73, 0x34, 0xa2,
	0x6f, 0xfd, 0x27, 0x3a, 0x39, 0x03, 0x8d, 0xe1, 0xca, 0x4b, 0x04, 0x29, 0xad, 0x9f, 0x6c, 0xaa,
	0x0c, 0x57, 0xf2, 0xc9, 0xfa, 0xad, 0x80, 0x61, 0x47, 0x3e, 0xbb, 0xc6, 0xba, 0xe9, 0x31, 0x68,
	0x69, 0x1c, 0x7a, 0xcd, 0xb6, 0x6a, 0x1a, 0x97, 0x69, 0xc9, 0x5b, 0x30, 0xea, 0x61, 0xb3, 0xb5,
	0xbe, 0xde, 0x90, 0xcd, 0x8f, 0xcb, 0x04, 0xcd, 0xb7, 0x41, 0x58, 0xd6, 0x12, 0xf5, 0xb0, 0x94,
	0x68, 0x95, 0x12, 0xeb, 0x0d, 0x21, 0x61, 0x7d, 0x7a, 0x9c, 0x6b, 0xf3, 0x2b, 0x9d, 0xfd, 0x84,
	0xe7, 0x4f, 0xde, 0x5c, 0xa2, 0x83, 0x4a, 0x9d, 0xc1, 0xc8, 0x9d, 0x3b, 0xd4, 0x54, 0x08, 0x01,
	0x83, 0x3a, 0xae, 0x33, 0xf7, 0x66, 0x3d, 0xd7, 0x5d, 0x4c, 0xe9, 0xa5, 0xb9, 0x45, 0x5e, 0xc0,
	0x7e, 0xc9, 0xd9, 0xce, 0x78, 0x3c, 0x1b, 0x4e, 0x27, 0x8e, 0xb9, 0x2d, 0x16, 0x27, 0xce, 0xc2,
	0x1b, 0xf6, 0xc6, 0x63, 0x6f, 0x3c, 0x1d, 0x8c, 0x26, 0x66, 0x8b, 0x18, 0x00, 0xfd, 0xd1, 0xe4,
	0xd2, 0x2b, 0x77, 0x76, 0xc8, 0x3e, 0xb4, 0x25, 0x5e, 0x38, 0xf6, 0xb0, 0x37, 0x37, 0xa1, 0x7f,
	0xf8, 0xf5, 0x80, 0x17, 0x78, 0x87, 0xe7, 0x41, 0x7c, 0x83, 0xac, 0xf0, 0xb2, 0x6f, 0xe7, 0xe2,
	0x7b, 0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x2d, 0x15, 0xef, 0x60, 0x04, 0x00, 0x00,
}
