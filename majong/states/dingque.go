package states

import (
	"fmt"
	"steve/gutils"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"steve/client_pb/room"

	"steve/client_pb/msgId"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

//DingqueState 定缺状态 @Author:wuhongwei
type DingqueState struct {
}

var _ interfaces.MajongState = new(DingqueState)

// ProcessEvent 处理事件
func (s *DingqueState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_dingque_request {
		isFinish, err := s.dingque(eventContext, flow)
		if err != nil || !isFinish {
			return majongpb.StateID_state_dingque, err
		}
		return majongpb.StateID_state_zixun, nil
	}
	return majongpb.StateID_state_dingque, nil
}

//定缺操作
func (s *DingqueState) dingque(eventContext []byte, flow interfaces.MajongFlow) (bool, error) {
	// 序列化utils
	dinqueEvent := new(majongpb.DingqueRequestEvent)
	if err := proto.Unmarshal(eventContext, dinqueEvent); err != nil {
		return false, fmt.Errorf("定缺 ： %v", errUnmarshalEvent)
	}
	//麻将牌局现场
	mjContext := flow.GetMajongContext()
	// 所有玩家
	players := mjContext.Players
	// 获取定缺玩家ID
	playerID := dinqueEvent.GetHead().GetPlayerId()
	// 获取定缺玩家
	dqPlayer := utils.GetPlayerByID(players, playerID)
	if dqPlayer == nil {
		return false, fmt.Errorf("定缺事件失败-定缺玩家ID不存在: %v ", playerID)
	}
	// 错误码-成功
	err := room.RoomError_Success
	// 定缺应答 请求-响应
	toClientRsq := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_DINGQUE_RSP),
		Msg: &room.RoomDingqueRsp{
			ErrCode: &err,
		},
	}
	// 推送消息应答
	flow.PushMessages([]uint64{playerID}, toClientRsq)
	// 获取定缺颜色
	dqColor := dinqueEvent.GetColor()

	// 校验颜色是否合法
	sichuangxueliuDingQueColor := map[majongpb.CardColor]string{
		majongpb.CardColor_ColorWan:  "万",
		majongpb.CardColor_ColorTong: "筒",
		majongpb.CardColor_ColorTiao: "条",
	}
	if _, ok := sichuangxueliuDingQueColor[dqColor]; !ok {
		return false, fmt.Errorf("定缺事件失败-定缺花色不存在: %v ", dqColor)
	}
	// 设置玩家定缺颜色
	dqPlayer.DingqueColor = dqColor
	// 设置已经定缺
	dqPlayer.HasDingque = true

	// 日志
	logrus.WithFields(logrus.Fields{
		"msgID":           msgid.MsgID_ROOM_DINGQUE_RSP,
		"toClientRsq":     toClientRsq,
		"dinqueEvent":     dinqueEvent,
		"dingQuePlayerID": dqPlayer.PalyerId,
		"dingQueColor":    dqPlayer.DingqueColor,
		"isDingQue":       dqPlayer.HasDingque,
	}).Info("-----定缺中")

	// 遍历其他玩家是否都已经定缺
	for i := 0; i < len(players); i++ {
		if !players[i].HasDingque {
			return false, nil
		}
	}
	return true, nil
}

// OnEntry 进入状态，进入定缺状态，发送到客户端，进入定缺
func (s *DingqueState) OnEntry(flow interfaces.MajongFlow) {
	// 定缺消息NTF被注释了 // 广播通知客户端进入定缺
	// dingQueNtf := room.RoomDingqueNtf{}
	// facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_DINGQUE_NTF, &dingQueNtf)
	// // 日志
	// logrus.WithFields(logrus.Fields{
	// 	"msgID":      msgid.MsgID_ROOM_DINGQUE_NTF,
	// 	"dingQueNtf": dingQueNtf,
	// }).Info("-----定缺开始-进入定缺状态")
}

// OnExit 退出状态，定缺完成，发送定缺完成通知，进入下一个状态
func (s *DingqueState) OnExit(flow interfaces.MajongFlow) {
	players := flow.GetMajongContext().Players
	// 所有定缺玩家消息通知
	playerDingQueMsg := make([]*room.PlayerDingqueColor, 0)
	// 设置每个玩家定缺颜色消息
	for _, player := range players {
		// 房间定缺完成通知的玩家定缺消息
		playerDingQue := &room.PlayerDingqueColor{
			PlayerId: proto.Uint64(player.PalyerId),
			Color:    gutils.ServerColor2ClientColor(player.DingqueColor).Enum(),
		}
		playerDingQueMsg = append(playerDingQueMsg, playerDingQue)
	}
	dingQueFinishNtf := room.RoomDingqueFinishNtf{
		PlayerDingqueColor: playerDingQueMsg,
	}
	// 广播定缺完成消息
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_DINGQUE_FINISH_NTF, &dingQueFinishNtf)
	// 日志
	logrus.WithFields(logrus.Fields{
		"msgID":            msgid.MsgID_ROOM_DINGQUE_FINISH_NTF,
		"dingQueFinishNtf": dingQueFinishNtf,
	}).Info("-----定缺完成-退出定缺状态")
}
