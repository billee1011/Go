package states

import (
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"
)

type settleState struct{}

func (s *settleState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入Settle状态")
}

func (s *settleState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开Settle状态")
}

func (s *settleState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID == int(ddz.EventID_event_showhand_finish) {
		s.settle(m)
		return int(ddz.StateID_state_over), nil
	}
	return int(ddz.StateID_state_settle), nil
}

func (s *settleState) settle(m machine.Machine) {
	context := getDDZContext(m)
	context.CurStage = ddz.DDZStage_DDZ_STAGE_OVER

	totalGrab := context.TotalGrab
	totalDouble := context.TotalDouble
	base := uint32(1000) // TODO:待确定底分从哪获取
	multiple := uint32(totalGrab * totalDouble * context.TotalBomb)
	if context.Spring || context.AntiSpring {
		multiple = multiple * 2
	}
	score := uint64(base * multiple)
	winnerId := context.WinnerId
	lordId := context.LordPlayerId
	lordWin := winnerId == lordId

	//找出每个人最大输赢金币
	maxScores := make(map[uint64]uint64)
	for _, player := range context.GetPlayers() {
		playerId := player.PlayerId
		coin := global.GetPlayerMgr().GetPlayer(playerId).GetCoin()
		s := If(playerId == lordId, score*2, score).(uint64)
		maxScores[playerId] = If(s > coin, coin, s).(uint64)
	}

	settleScores := make(map[uint64]uint64)

	lordMax := maxScores[lordId]
	score = lordMax / 2
	lordScore := uint64(0)
	for _, player := range context.GetPlayers() {
		playerId := player.PlayerId
		if playerId == lordId {
			continue
		}
		maxScore := maxScores[playerId]
		settleScore := If(score > maxScore, maxScore, score).(uint64)
		settleScores[playerId] = settleScore
		lordScore += settleScore
	}
	if lordMax%2 == 1 { //如果地主金豆数是奇数
		if lordWin {
			lordScore-- //少赢一分(赢的分不能超过本人金豆上限)
		} else {
			lordScore++ //多扣一分，系统扣成偶数(华华的要求)
		}
	}
	settleScores[lordId] = lordScore

	var billPlayers []*room.DDZBillPlayerInfo
	for _, player := range context.GetPlayers() {
		playerId := player.PlayerId
		roomPlayer := global.GetPlayerMgr().GetPlayer(playerId)
		billPlayer := room.DDZBillPlayerInfo{}
		billPlayer.PlayerId = &playerId
		isWin := winnerId == playerId
		billPlayer.Win = &isWin
		billPlayer.Base = proto.Int32(int32(base))
		billPlayer.Multiple = proto.Int32(int32(multiple))
		originCoin := roomPlayer.GetCoin()
		settleScore := settleScores[playerId]
		if isWin {
			roomPlayer.SetCoin(originCoin + settleScore) //赢钱
		} else {
			roomPlayer.SetCoin(originCoin - settleScore) //输钱
		}
		billPlayer.Score = proto.Int64(int64(settleScore))
		billPlayer.CurrentScore = proto.Int64(int64(roomPlayer.GetCoin()))
		billPlayer.Lord = &player.Lord
		billPlayer.OutCards = player.OutCards
		billPlayer.HandCards = player.HandCards
		billPlayers = append(billPlayers, &billPlayer)
	}

	antiSpring := !context.Spring && context.AntiSpring

	settleMsg := &room.DDZGameOverNtf{
		WinnerId:     &context.WinnerId,
		ShowHandTime: proto.Uint32(4),
		Spring:       &context.Spring,
		AntiSpring:   &antiSpring,
		Bills:        billPlayers,
	}
	broadcast(m, msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF,
		settleMsg,
	)
	logrus.WithField("settle msg", settleMsg).Debugln("斗地主广播结算消息")
}
