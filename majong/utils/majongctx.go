package utils

import (
	"errors"
	"steve/client_pb/room"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// GetMajongPlayer 从 MajongContext 中根据玩家 ID 获取玩家
func GetMajongPlayer(playerID uint64, mjContext *majongpb.MajongContext) *majongpb.Player {
	for _, player := range mjContext.GetPlayers() {
		if player.GetPalyerId() == playerID {
			return player
		}
	}
	return nil
}

// ExistPossibleAction 玩家是否存在指定的可能行为
func ExistPossibleAction(player *majongpb.Player, action majongpb.Action) bool {
	for _, a := range player.GetPossibleActions() {
		if a == action {
			return true
		}
	}
	return false
}

// GetPlayerIndex 获取玩家索引
func GetPlayerIndex(playerID uint64, players []*majongpb.Player) (int, error) {
	for index, player := range players {
		if player.GetPalyerId() == playerID {
			return index, nil
		}
	}
	return 0, errors.New("not exists")
}

// GetCardsGroup 获取玩家牌组信息
func GetCardsGroup(player *majongpb.Player) []*room.CardsGroup {
	cardsGroupList := make([]*room.CardsGroup, 0)
	// 碰牌
	for _, pengCard := range player.PengCards {
		card, _ := CardToInt(*pengCard.Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  room.CardsGroupType_CGT_PENG.Enum(),
			Cards: []uint32{uint32(*card)},
		}
		cardsGroupList = append(cardsGroupList, cardsGroup)
	}
	// 杠牌
	var groupType *room.CardsGroupType
	for _, gangCard := range player.GangCards {
		if gangCard.Type == majongpb.GangType_gang_angang {
			groupType = room.CardsGroupType_CGT_ANGANG.Enum()
		}
		if gangCard.Type == majongpb.GangType_gang_minggang {
			groupType = room.CardsGroupType_CGT_MINGGANG.Enum()
		}
		if gangCard.Type == majongpb.GangType_gang_bugang {
			groupType = room.CardsGroupType_CGT_BUGANG.Enum()
		}
		card, _ := CardToInt(*gangCard.Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  groupType,
			Cards: []uint32{uint32(*card)},
		}
		cardsGroupList = append(cardsGroupList, cardsGroup)
	}
	// 手牌
	handCards, _ := CardsToInt(player.HandCards)
	cards := make([]uint32, 0)
	for _, handCard := range handCards {
		cards = append(cards, uint32(handCard))
	}
	cardsGroup := &room.CardsGroup{
		Pid:   proto.Uint64(player.PalyerId),
		Type:  room.CardsGroupType_CGT_HAND.Enum(),
		Cards: cards,
	}
	cardsGroupList = append(cardsGroupList, cardsGroup)
	return cardsGroupList
}
