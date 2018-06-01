//Package states implements a simple states
package states

//适用麻将：四川血流
//前置条件：无
//处理的事件请求：定缺请求
//处理请求的过程：设置定缺玩家的定缺颜色和设置定缺玩家是否已经定缺,所有玩家定缺完设置麻将现场最后摸牌玩家为庄家，
//请求成功处理，设置应答错误码成功消息通知，给当前玩家客户端
//处理请求的结果：所有玩家都定缺则返回自询状态ID，否则返回定缺状态ID
//状态退出行为：定缺完成，广播通知客户端定缺完成消息通知，该通知包含每个玩家的ID和定缺的颜色
//状态进入行为：无
//约束条件：定缺的颜色，必需是万或条或筒
import (
	"fmt"
	"steve/gutils"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/states/tingtips"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"steve/client_pb/room"

	"steve/client_pb/msgId"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

//DingqueState 定缺状态
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
		s.notifyTingCards(flow)
		mjContext := flow.GetMajongContext()
		mjContext.ZixunType = majongpb.ZixunType_ZXT_NORMAL
		return majongpb.StateID_state_zixun, nil
	}
	return majongpb.StateID_state_dingque, nil
}

// notifyTingCards 通知玩家听牌信息
func (s *DingqueState) notifyTingCards(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	for seat, player := range mjContext.GetPlayers() {
		if seat == int(mjContext.GetZhuangjiaIndex()) {
			continue
		}
		tingtips.NotifyTingCards(flow, player.GetPalyerId())
	}
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
	// 获取定缺颜色
	dqColor := dinqueEvent.GetColor()
	// 检测定缺颜色是否合法
	if ok := checkDingQueReq(dqColor); !ok {
		return false, fmt.Errorf("定缺事件失败-定缺花色不存在: %v ", dqColor)
	}
	// 设置玩家定缺颜色
	dqPlayer.DingqueColor = dqColor
	// 设置已经定缺
	dqPlayer.HasDingque = true
	// 应答
	onDingQueRsq(playerID, flow)

	// 日志
	logrus.WithFields(logrus.Fields{
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

// OnEntry 进入状态，进入定缺状态，发送到客户端，进入定缺tiao jian
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

//checkDingQueReq 检测定缺请求是否合法
func checkDingQueReq(dingQueColor majongpb.CardColor) bool {
	sichuangxueliuDingQueColor := map[majongpb.CardColor]string{
		majongpb.CardColor_ColorWan:  "万",
		majongpb.CardColor_ColorTong: "筒",
		majongpb.CardColor_ColorTiao: "条",
	}
	colorValue, ok := sichuangxueliuDingQueColor[dingQueColor]
	if !ok {
		return false
	}
	logrus.WithFields(logrus.Fields{
		"DingQueColor": colorValue,
	}).Info("--定缺颜色")
	return true
}

//onRsq 定缺应答
func onDingQueRsq(playerID uint64, flow interfaces.MajongFlow) {
	// 错误码-成功
	errCode := room.RoomError_SUCCESS
	// 定缺应答 请求-响应
	toClientRsq := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_DINGQUE_RSP),
		Msg: &room.RoomDingqueRsp{
			ErrCode: &errCode,
		},
	}
	// 推送消息应答
	flow.PushMessages([]uint64{playerID}, toClientRsq)
	logrus.WithFields(logrus.Fields{
		"msgID":      msgid.MsgID_ROOM_DINGQUE_RSP,
		"dingQueNtf": toClientRsq,
	}).Info("-----定缺成功应答")
}
