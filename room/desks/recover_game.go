package desks

import (
	"steve/client_pb/room"
	"steve/gutils"
	server_pb "steve/server_pb/majong"
	"time"

	"github.com/golang/protobuf/proto"
)

func getStateCostTime(entryTime int64) (costTime uint32) {
	nowTime := time.Now().Unix()
	if nowTime > entryTime {
		costTime = uint32(nowTime - entryTime)
	}
	return
}

func getOperatePlayerID(mjContext *server_pb.MajongContext) *uint64 {
	state := mjContext.GetCurState()
	var playerID uint64
	switch state {
	case server_pb.StateID_state_chupai, server_pb.StateID_state_hu, server_pb.StateID_state_chupaiwenxun:
		playerID = mjContext.GetLastChupaiPlayer()
	case server_pb.StateID_state_angang, server_pb.StateID_state_gang, server_pb.StateID_state_waitqiangganghu:
		playerID = mjContext.GetLastGangPlayer()
	case server_pb.StateID_state_zimo, server_pb.StateID_state_zixun, server_pb.StateID_state_bugang:
		playerID = mjContext.GetLastMopaiPlayer()
	case server_pb.StateID_state_peng:
		playerID = mjContext.GetLastPengPlayer()
	case server_pb.StateID_state_mopai:
		playerID = mjContext.GetMopaiPlayer()
	case server_pb.StateID_state_qiangganghu:
		if players := mjContext.GetLastHuPlayers(); len(players) > 0 {
			playerID = players[0]
		}
	default: // 前端要求必须有一个有效值
		playerID = mjContext.GetLastMopaiPlayer()
	}
	return &playerID
}

func getGameStage(curState server_pb.StateID) (stage room.GameStage) {
	switch curState {
	case server_pb.StateID_state_huansanzhang:
		stage = room.GameStage_GAMESTAGE_HUANSANZHANG
	case server_pb.StateID_state_dingque:
		stage = room.GameStage_GAMESTAGE_DINGQUE
	default:
		stage = room.GameStage_GAMESTAGE_PLAYCARD
	}
	return
}

func getDoorCard(mjContext *server_pb.MajongContext) *uint32 {
	if mjContext.GetCurState() == server_pb.StateID_state_zixun {
		DoorCard := uint32(mjContext.GetLastMopaiCard().GetPoint())
		return &DoorCard
	}
	return nil
}

func getRecoverPlayerInfo(d *desk) (recoverPlayerInfo []*room.GamePlayerInfo) {
	mjContext := &d.dContext.mjContext
	roomPlayerInfos := d.GetPlayers()
	for _, roomPlayerInfo := range roomPlayerInfos {
		var player *server_pb.Player
		// 这里假设总能找到一个对应玩家
		for _, player = range mjContext.GetPlayers() {
			if player.GetPalyerId() == roomPlayerInfo.GetPlayerId() {
				break
			}
		}
		playerID := player.GetPalyerId()
		svrHandCard := player.GetHandCards()
		handCardCount := uint32(len(svrHandCard))
		gamePlayerInfo := &room.GamePlayerInfo{
			PlayerInfo:    roomPlayerInfo,
			Color:         gutils.ServerColor2ClientColor(player.DingqueColor).Enum(),
			HandCardCount: &handCardCount,
		}

		// 手牌组
		cltHandCard := gutils.ServerCards2Numbers(svrHandCard)
		handCardGroup := &room.CardsGroup{
			Cards: cltHandCard,
			Type:  room.CardsGroupType_CGT_HAND.Enum(),
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, handCardGroup)
		// 吃牌组

		// 碰牌组,每一次碰牌填1张还是三张
		var pengCardGroups []*room.CardsGroup
		for _, pengCard := range player.GetPengCards() {
			srcPlayerID := pengCard.GetSrcPlayer()
			cards := []uint32{gutils.ServerCard2Number(pengCard.GetCard())}
			pengCardGroup := &room.CardsGroup{
				Cards: append(cards, cards[0], cards[0]),
				Type:  room.CardsGroupType_CGT_PENG.Enum(),
				Pid:   &srcPlayerID,
			}
			pengCardGroups = append(pengCardGroups, pengCardGroup)
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, pengCardGroups...)
		// 杠牌组
		var gangCardGroups []*room.CardsGroup
		for _, gangCard := range player.GetGangCards() {
			groupType := gutils.GangTypeSvr2Client(gangCard.GetType())
			srcPlayerID := gangCard.GetSrcPlayer()
			cards := []uint32{gutils.ServerCard2Number(gangCard.GetCard())}
			gangCardGroup := &room.CardsGroup{
				Cards: append(cards, cards[0], cards[0], cards[0]),
				Type:  &groupType,
				Pid:   &srcPlayerID,
			}
			gangCardGroups = append(gangCardGroups, gangCardGroup)
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, gangCardGroups...)
		// 胡牌组
		var huCards []*server_pb.Card
		for _, card := range player.GetHuCards() {
			huCards = append(huCards, card.GetCard())
		}
		huCardGroups := []*room.CardsGroup{
			&room.CardsGroup{
				Cards: gutils.ServerCards2Numbers(huCards),
				Type:  room.CardsGroupType_CGT_HU.Enum(),
			},
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, huCardGroups...)
		// 花牌组

		// 出牌组
		outCardGroup := &room.CardsGroup{
			Cards: gutils.ServerCards2Numbers(player.GetOutCards()),
			Type:  room.CardsGroupType_CGT_OUT.Enum(),
			Pid:   &playerID,
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, outCardGroup)
		recoverPlayerInfo = append(recoverPlayerInfo, gamePlayerInfo)
	}
	return
}

func getZixunInfo(playerID uint64, mjContext *server_pb.MajongContext) (*bool, *room.RoomZixunNtf) {
	if mjContext.GetCurState() != server_pb.StateID_state_zixun {
		return proto.Bool(false), nil
	}

	if mjContext.GetLastMopaiPlayer() != playerID {
		return proto.Bool(false), nil
	}
	player := gutils.GetMajongPlayer(playerID, mjContext)
	return proto.Bool(true), zixunTransform(player.GetZixunRecord())
}

func getWenxunInfo(playerID uint64, mjContext *server_pb.MajongContext) (*bool, *room.RoomChupaiWenxunNtf) {
	if mjContext.GetCurState() != server_pb.StateID_state_chupaiwenxun {
		return proto.Bool(false), nil
	}

	player := gutils.GetMajongPlayer(playerID, mjContext)
	enableActions := player.GetPossibleActions()
	if len(enableActions) == 0 || player.GetHasSelected() {
		return proto.Bool(false), nil
	}

	outCard := gutils.ServerCard2Number(mjContext.GetLastOutCard())
	wenXunInfo := &room.RoomChupaiWenxunNtf{
		Card: &outCard,
	}
	for _, action := range enableActions {
		switch action {
		case server_pb.Action_action_peng:
			wenXunInfo.EnablePeng = proto.Bool(true)
		case server_pb.Action_action_gang:
			wenXunInfo.EnableMinggang = proto.Bool(true)
		case server_pb.Action_action_hu:
			wenXunInfo.EnableDianpao = proto.Bool(true)
		case server_pb.Action_action_qi:
			wenXunInfo.EnableQi = proto.Bool(true)
		}
	}
	return proto.Bool(true), wenXunInfo
}

func getQghInfo(playerID uint64, mjContext *server_pb.MajongContext) (*bool, *room.RoomWaitQianggangHuNtf) {
	if mjContext.GetCurState() != server_pb.StateID_state_waitqiangganghu {
		return proto.Bool(false), nil
	}

	player := gutils.GetMajongPlayer(playerID, mjContext)
	enableActions := player.GetPossibleActions()
	if len(enableActions) == 0 || player.GetHasSelected() {
		return proto.Bool(false), nil
	}

	outCard := gutils.ServerCard2Number(mjContext.GetLastOutCard())
	gangPlayerID := mjContext.GetLastGangPlayer()
	qghInfo := &room.RoomWaitQianggangHuNtf{
		Card:         &outCard,
		SelfCan:      proto.Bool(len(player.GetPossibleActions()) != 0),
		FromPlayerId: &gangPlayerID,
	}
	return proto.Bool(true), qghInfo
}

func zixunTransform(record *server_pb.ZixunRecord) *room.RoomZixunNtf {
	zixunNtf := &room.RoomZixunNtf{}
	zixunNtf.EnableAngangCards = record.GetEnableAngangCards()
	zixunNtf.EnableBugangCards = record.GetEnableBugangCards()
	zixunNtf.EnableChupaiCards = record.GetEnableChupaiCards()
	zixunNtf.EnableQi = proto.Bool(record.GetEnableQi())
	zixunNtf.EnableZimo = proto.Bool(record.GetEnableZimo())
	huType := gutils.HuTypeSvr2Client(record.GetHuType())
	if huType != nil {
		zixunNtf.HuType = huType
	}
	zixunNtf.CanTingCardInfo = gutils.CanTingCardInfoSvr2Client(record.GetCanTingCardInfo())

	return zixunNtf
}
