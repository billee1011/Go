package gateclient

import (
	"github.com/golang/protobuf/proto"
	"steve/structs"
	"steve/structs/proto/gate_rpc"
)

/*
	功能：网关GateWay client API
	作者： SkyWang
	日期： 2018-8-3
*/

// 方法：通过GateWay向指定玩家ID的client发送通知消息
// 参数: playerID=玩家ID, cmd=消息ID, body=消息体
// 返回: 错误
func SendPackageByPlayerID(playerID uint64, cmd uint32, body proto.Message) error {
	exposer := structs.GetGlobalExposer()
	head := &steve_proto_gaterpc.Header{}
	head.MsgId = cmd
	return exposer.Exchanger.SendPackageByPlayerID(playerID, head, body)
}

// 方法：通过GateWay向多个指定玩家ID的client发送通知消息
// 参数: playerIDs=玩家ID列表, cmd=消息ID, body=消息体
// 返回: 错误
func BroadcastPackageByPlayerID(playerIDs []uint64, cmd uint32, body proto.Message) error {
	exposer := structs.GetGlobalExposer()
	head := &steve_proto_gaterpc.Header{}
	head.MsgId = cmd
	return exposer.Exchanger.BroadcastPackageByPlayerID(playerIDs, head, body)
}

// 方法：通过NSQ消息队列，向GateWay的所有玩家发送广播消息,GateWay将会订阅此消息
// 参数: cmd=消息ID, body=消息体
// 返回：错误
func NsqBroadcastAllMsg(cmd uint32, body proto.Message) error {
	return nsqBroadcastMsg(steve_proto_gaterpc.BroadCastType_TO_ALL, 0, cmd, body)
}

// 方法：通过NSQ消息队列，向GateWay的指定channel玩家发送广播消息,GateWay将会订阅此消息
// 参数: channel=渠道ID, cmd=消息ID, body=消息体
// 返回：错误
func NsqBroadcastChannelMsg(channel int64, cmd uint32, body proto.Message) error {
	return nsqBroadcastMsg(steve_proto_gaterpc.BroadCastType_TO_CHANNEL, channel, cmd, body)
}
// 方法：通过NSQ消息队列，向GateWay的指定prov玩家发送广播消息,GateWay将会订阅此消息
// 参数: prov=省ID, cmd=消息ID, body=消息体
// 返回：错误
func NsqBroadcastProvMsg(prov int64, cmd uint32, body proto.Message) error {
	return nsqBroadcastMsg(steve_proto_gaterpc.BroadCastType_TO_PROV, prov, cmd, body)
}

// 方法：通过NSQ消息队列，向GateWay的指定city玩家发送广播消息,GateWay将会订阅此消息
// 参数: city=城市ID, cmd=消息ID, body=消息体
// 返回：错误
func NsqBroadcastCityMsg(city int64, cmd uint32, body proto.Message) error {
	return nsqBroadcastMsg(steve_proto_gaterpc.BroadCastType_TO_CITY, city, cmd, body)
}

// 方法：通过NSQ消息队列，向所有GateWay发送广播消息,GateWay将会订阅此消息
// 参数: channel=渠道ID, prov=省包ID, city=城市ID, cmd=消息ID, body=消息体
// 返回：错误
// 说明: channel != 0 广播给渠道所有的人,然后 prov != 0 广播给所有省包的人, 再然后 city != 0 广播给所有城市的人
func nsqBroadcastMsg(sendType steve_proto_gaterpc.BroadCastType, sendId int64, cmd uint32, body proto.Message) error {
	exposer := structs.GetGlobalExposer()
	messageData, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	head := &steve_proto_gaterpc.Header{}
	head.MsgId = cmd
	req := &steve_proto_gaterpc.BroadcastMsgRequest{}
	req.SendType = sendType
	req.Header = head
	req.SendId = sendId

	req.Data = messageData

	packData, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	return exposer.Publisher.Publish("broadcast_msg", packData)
}