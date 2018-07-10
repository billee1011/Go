package core

// type handler func(uint64, *steve_proto_base.Header, []byte)

// type receiver struct {
// 	handlers map[msgid.MsgID]handler
// }

// func (r *receiver) OnRecv(clientID uint64, header *steve_proto_base.Header, body []byte) {
// 	msg := msgid.MsgID(header.GetMsgId())
// 	handler, ok := r.handlers[msg]
// 	if !ok {
// 		return
// 	}
// 	handler(clientID, header, body)
// }

// // NewReceiver 创建消息接收器
// func NewReceiver() net.MessageObserver {
// 	return &receiver{
// 		handlers: map[msgid.MsgID]handler{
// 			msgid.MsgID_LOGIN_AUTH_REQ: auth.OnAuthRequest,
// 		},
// 	}
// }
