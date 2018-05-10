package states

import (
	"fmt"
	"steve/majong/interfaces"
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
		return false, fmt.Errorf("定缺事件反序列化失败: %v", err)
	}
	//麻将牌局现场
	mjContext := flow.GetMajongContext()
	// 所有玩家
	players := mjContext.Players
	// 获取定缺玩家和定缺颜色
	playerID := dinqueEvent.GetHead().GetPlayerId()
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
	// 定缺所有玩家ID
	playerAllID := []uint64{}
	// 所有定缺玩家通知
	playerDqColors := make([]*room.PlayerDingqueColor, 0)
	// 遍历其他玩家是否都已经定缺,并设置广播通知定缺完成
	for i := 0; i < len(players); i++ {
		if dqPlayer.PalyerId != players[i].PalyerId && !players[i].HasDingque {
			return false, nil
		}
		playerAllID = append(playerAllID, players[i].PalyerId)
		playerDQ := &room.PlayerDingqueColor{
			PlayerId: &players[i].PalyerId,
			Color:    room.CardColor(players[i].DingqueColor).Enum(),
		}
		playerDqColors = append(playerDqColors, playerDQ)
	}
	// 定缺完成通知
	dqNtf := &room.RoomDingqueFinishNtf{
		PlayerDingqueColor: playerDqColors,
	}
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_DINGQUE_FINISH_NTF),
		Msg:   dqNtf,
	}
	// 推送消息
	// flow.PushMessages(playerAllID, toClient)
	// 日志
	logrus.WithFields(logrus.Fields{
		"playerAllID": playerAllID,
		"toClient":    toClient,
	}).Info("定缺成功")
	return true, nil
}

// OnEntry 进入状态
func (s *DingqueState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *DingqueState) OnExit(flow interfaces.MajongFlow) {

}
