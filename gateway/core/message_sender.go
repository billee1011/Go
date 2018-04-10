package core

import (
	"context"
	"steve/structs/proto/gate_rpc"
)

// type MessageSenderServer interface {
// 	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResult, error)
// }

// 消息发送请求
// type SendMessageRequest struct {
// 	// 客户端 ID 列表
// 	ClientId []uint64 `protobuf:"varint,1,rep,packed,name=client_id,json=clientId" json:"client_id,omitempty"`
// 	// 消息头
// 	Header *Header `protobuf:"bytes,2,opt,name=header" json:"header,omitempty"`
// 	// 消息内容
// 	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
// }

type messageSenderServer struct {
	core *gatewayCore
}

/*
	MsgId            *uint32 `protobuf:"varint,1,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	SendSeq          *uint64 `protobuf:"varint,2,opt,name=send_seq,json=sendSeq" json:"send_seq,omitempty"`
	RecvSeq          *uint64 `protobuf:"varint,3,opt,name=recv_seq,json=recvSeq" json:"recv_seq,omitempty"`
	StampTime        *uint64 `protobuf:"varint,4,opt,name=stamp_time,json=stampTime" json:"stamp_time,omitempty"`
	BodyLength       *uint32 `protobuf:"varint,5,opt,name=body_length,json=bodyLength" json:"body_length,omitempty"`
	RspSeq           *uint64 `protobuf:"varint,6,opt,name=rsp_seq,json=rspSeq" json:"rsp_seq,omitempty"`
	Version          *string `protobuf:"bytes,7,opt,name=version" json:"version,omitempty"`
	XXX_unrecognized []byte  `json:"-"`

*/

func (mss *messageSenderServer) SendMessage(ctx context.Context, req *steve_proto_gaterpc.SendMessageRequest) (*steve_proto_gaterpc.SendMessageResult, error) {

	// header := steve_proto_base.Header{
	// 	MsgId:   proto.Uint32(req.GetHeader().GetMsgId()),
	// 	Version: proto.String("1.0"), // TODO
	// }
	// mss.core.dog.BroadPackage(req.GetClientId(), &header, req.GetData())
	return nil, nil
}
