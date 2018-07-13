package states

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type overState struct{}

func (s *overState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入Over状态")

	context := getDDZContext(m)

	totalGrab := GetTotalGrab(context.GetPlayers())
	totalDouble := GetTotalDouble(context.GetPlayers())
	base := uint32(1) // TODO:待确定底分从哪获取
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
		playerId := player.PalyerId
		coin := global.GetPlayerMgr().GetPlayer(playerId).GetCoin()
		s := If(playerId == lordId, score*2, score).(uint64)
		maxScores[playerId] = If(s > coin, coin, s).(uint64)
	}

	settleScores := make(map[uint64]uint64)

	lordMax := maxScores[lordId]
	score = lordMax / 2
	lordScore := uint64(0)
	for _, player := range context.GetPlayers() {
		playerId := player.PalyerId
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

	billPlayers := []*room.DDZBillPlayerInfo{}
	for _, player := range context.GetPlayers() {
		playerId := player.PalyerId
		roomPlayer := global.GetPlayerMgr().GetPlayer(playerId)
		billPlayer := room.DDZBillPlayerInfo{}
		billPlayer.PlayerId = &playerId
		billPlayer.Base = proto.Int32(int32(base))
		billPlayer.Multiple = proto.Int32(int32(multiple))
		originCoin := roomPlayer.GetCoin()
		settleScore := settleScores[playerId]
		roomPlayer.SetCoin(originCoin - settleScore) //实施扣费
		billPlayer.Score = proto.Int64(int64(settleScore))
		billPlayer.CurrentScore = proto.Int64(int64(roomPlayer.GetCoin()))
		billPlayer.Lord = &player.Lord
		billPlayer.OutCards = player.OutCards
		billPlayer.HandCards = player.HandCards
		billPlayers = append(billPlayers, &billPlayer)
	}

	antiSpring := !context.Spring && context.AntiSpring
	broadcast(m, msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF,
		&room.DDZGameOverNtf{
			WinnerId:     &context.WinnerId,
			ShowHandTime: proto.Uint32(4),
			Spring:       &context.Spring,
			AntiSpring:   &antiSpring,
			Bills:        billPlayers,
		},
	)
}

func (s *overState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开Over状态")
}

func (s *overState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	return int(ddz.StateID_state_over), nil
}
