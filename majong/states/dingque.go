package states

import (
	"fmt"
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
	// 序列化
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

	// 获取定缺玩家和定缺颜色
	dqPlayer := utils.GetPlayerByID(players, playerID)
	if dqPlayer == nil {
		return false, fmt.Errorf("定缺事件失败-定缺玩家ID不存在: %v ", playerID)
	}
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
	// 所有定缺玩家通知
	playerDqColors := make([]*room.PlayerDingqueColor, 0)
	// 遍历其他玩家是否都已经定缺,并设置每个玩家定缺颜色消息
	for i := 0; i < len(players); i++ {
		if !players[i].HasDingque {
			return false, nil
		}
		// 房间定缺完成通知的玩家定缺消息
		playerDQ := &room.PlayerDingqueColor{
			PlayerId: proto.Uint64(players[i].PalyerId),
			Color:    room.CardColor(players[i].DingqueColor).Enum(),
		}
		playerDqColors = append(playerDqColors, playerDQ)
	}
	dingQueNtf := room.RoomDingqueFinishNtf{
		PlayerDingqueColor: playerDqColors,
	}
	// 广播定缺完成消息
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_DINGQUE_FINISH_NTF, &dingQueNtf)
	// 日志
	logrus.WithFields(logrus.Fields{
		"toClientRsq": toClientRsq,
		"dingQueNtf":  dingQueNtf,
	}).Info("定缺成功")
	return true, nil
}

// OnEntry 进入状态
func (s *DingqueState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *DingqueState) OnExit(flow interfaces.MajongFlow) {

}
