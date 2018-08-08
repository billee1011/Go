package core

/*
	功能： 服务配置中心，定义RPC服务关联，Client消息分派。

 */
import (
	"steve/client_pb/msgid"
	"steve/mailserver/msg"
	"steve/mailserver/logic"
)

/////////////////////////////////////////[1.配置线程模型]////////////////////////////////////
// 是否采用单线程运行所有协程(goroutime)
var bSingleThread = true

/////////////////////////////[2.定义RPC服务实现]//////////////////////////////////////////
// PB文件中定义的RPC服务接口,如果不提供RPC服务，设置为nil
var pbService interface{} = nil
// PB定义的RPC服务的实现类,如果不提供RPC服务，设置为nil
var pbServerImp interface{} = nil

/////////////////////////////[3.处理Client消息]////////////////////////////////////////////
// 如果需要处理Client消息，需要在GateWay配置消息ID开始~ 结束区间 关联到当前服务名,GateWay才会把消息转发到此服务.

// 添加从GateWay过来的Client消息处理
var mapMsg  = map[msgid.MsgID] interface{} {
	msgid.MsgID_MAILSVR_GET_UNREAD_SUM_REQ:msg.ProcessGetUnReadSumReq,
	msgid.MsgID_MAILSVR_GET_MAIL_LIST_REQ:msg.ProcessGetMailListReq,
	msgid.MsgID_MAILSVR_GET_MAIL_DETAIL_REQ:msg.ProcessGetMailDetailReq,
	msgid.MsgID_MAILSVR_DEL_MAIL_REQ:msg.ProcessDelMailReq,
	msgid.MsgID_MAILSVR_SET_READ_TAG_REQ:msg.ProcessSetReadTagReq,
	msgid.MsgID_MAILSVR_AWARD_ATTACH_REQ:msg.ProcessAwardAttachReq,
}

/////////////////////////////[4.向client发送通知消息]////////////////////////////////////////////
// 4.1通过GateWay向指定player_id列表发送通知消息
/*
// 方法：通过GateWay向指定玩家ID的client发送通知消息
// 参数: playerID=玩家ID, cmd=消息ID, body=消息体
// 返回: 错误
req := &msgserver.MsgSvrHorseRaceChangeNtf{}
playderId := 1001
gateclient.SendPackageByPlayerID(playderId, uint32(msgid.MsgID_MSGSVR_HORSE_RACE_UPDATE_NTF), req)

// 方法：通过GateWay向多个指定玩家ID的client发送通知消息
// 参数: playerIDs=玩家ID列表, cmd=消息ID, body=消息体
// 返回: 错误
req := &msgserver.MsgSvrHorseRaceChangeNtf{}
playderId := 1001
gateclient.BroadcastPackageByPlayerID([]uint64{playderId}, uint32(msgid.MsgID_MSGSVR_HORSE_RACE_UPDATE_NTF), req)
*/

// 4.2通过GateWay向Client发送广播消息
/*
	req := &msgserver.MsgSvrHorseRaceChangeNtf{}
	gateclient.NsqBroadcastAllMsg(uint32(msgid.MsgID_MSGSVR_HORSE_RACE_UPDATE_NTF), req)
 */
/////////////////////////////[5.通过nsq发布和订阅消息]////////////////////////////////////////////
// 5.1发布消息
/*
exposer := structs.GetGlobalExposer()
if err := exposer.Publisher.Publish("player_login", messageData); err != nil {
entry.WithError(err).Errorln("发布消息失败")
}
*/
// 5.2订阅消息
/*
	exposer := structs.GetGlobalExposer()
	if err := exposer.Subscriber.Subscribe("player_login", "match", &playerLoginHandler{}); err != nil {
		logrus.WithError(err).Panicln("订阅登录消息失败")
	}
 */
/////////////////////////////[6.服务初始化配置]////////////////////////////////////////////
// 比如从DB或文件加载配置
func InitServer() error {
	err := logic.Init()
	if err != nil {
		return err
	}
	return nil
}




