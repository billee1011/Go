package login

import (
	"steve/structs"
	"steve/structs/proto/gate_rpc"
	"steve/structs/proto/msg"

	"github.com/golang/protobuf/proto"

	"github.com/Sirupsen/logrus"
)

// HandleLogin 处理登录请求
// 返回 LoginRsp 作为登录回复
func HandleLogin(clientID uint64, head *steve_proto_gaterpc.Header, req steve_proto_msg.LoginReq) []proto.Message {
	entry := logrus.WithFields(logrus.Fields{
		"name":      "HandleLogin",
		"client_id": clientID,
		"user_name": req.GetUserName(),
	})
	entry.Info("用户登录")

	resp := &steve_proto_msg.LoginRsp{
		Result: steve_proto_msg.ErrorCode_err_OK.Enum(),
		UserId: proto.Uint64(1),
	}
	e := structs.GetGlobalExposer()
	e.Exchanger.SendPackage(clientID, head, resp)
	return nil
	// return []proto.Message{resp}  // 使用这种方式可以原序返回
}
