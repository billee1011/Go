package utils

import (
	"errors"
	"steve/gutils"
	majongpb "steve/server_pb/majong"
)

// GetMajongPlayer 从 MajongContext 中根据玩家 ID 获取玩家
func GetMajongPlayer(playerID uint64, mjContext *majongpb.MajongContext) *majongpb.Player {
	return gutils.GetMajongPlayer(playerID, mjContext)
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

// GetPlayerIDIndex 获取玩家索引
func GetPlayerIDIndex(playerID uint64, players []uint64) (int, error) {
	for index, pid := range players {
		if pid == playerID {
			return index, nil
		}
	}
	return 0, errors.New("not exists")
}

// GetPalyerCloseFromTarget 从targets获取离玩家index最近的玩家id
func GetPalyerCloseFromTarget(index int, allPlayer, targets []uint64) uint64 {
	for i := 0; i <= len(allPlayer); i++ {
		nextIndex := (index + i) % len(allPlayer)
		for _, target := range targets {
			index, _ := GetPlayerIDIndex(target, allPlayer)
			if index == nextIndex {
				return target
			}
		}
	}
	return 0
}

// GetCardsGroup 获取玩家牌组信息
func GetCardsGroup(player *majongpb.Player, huCard *majongpb.Card) []*majongpb.CardsGroup {
	cardsGroupList := make([]*majongpb.CardsGroup, 0)
	// 手牌组
	handCards := ServerCards2Numbers(player.GetHandCards())
	cltHandCard := make([]uint32, 0)
	for _, handCard := range handCards {
		cltHandCard = append(cltHandCard, uint32(handCard))
	}
	handCardGroup := &majongpb.CardsGroup{
		Cards: cltHandCard,
		Type:  majongpb.CardsGroupType_CGT_HAND,
	}
	cardsGroupList = append(cardsGroupList, handCardGroup)
	// 吃牌组
	var chiCardGroups []*majongpb.CardsGroup
	for _, chiCard := range player.GetChiCards() {
		srcPlayerID := chiCard.GetSrcPlayer()
		card := ServerCard2Number(chiCard.GetCard())
		chiCardGroup := &majongpb.CardsGroup{
			Cards: []uint32{uint32(card), uint32(card) + 1, uint32(card) + 2},
			Type:  majongpb.CardsGroupType_CGT_PENG,
			Pid:   srcPlayerID,
		}
		chiCardGroups = append(chiCardGroups, chiCardGroup)
	}
	cardsGroupList = append(cardsGroupList, chiCardGroups...)
	// 碰牌组,每一次碰牌填1张还是三张
	var pengCardGroups []*majongpb.CardsGroup
	for _, pengCard := range player.GetPengCards() {
		srcPlayerID := pengCard.GetSrcPlayer()
		card := gutils.ServerCard2Number(pengCard.Card)
		pengCardGroup := &majongpb.CardsGroup{
			Cards: []uint32{card, card, card},
			Type:  majongpb.CardsGroupType_CGT_PENG,
			Pid:   srcPlayerID,
		}
		pengCardGroups = append(pengCardGroups, pengCardGroup)
	}
	cardsGroupList = append(cardsGroupList, pengCardGroups...)
	// 杠牌组
	var gangCardGroups []*majongpb.CardsGroup
	var groupType majongpb.CardsGroupType
	for _, gangCard := range player.GangCards {
		if gangCard.Type == majongpb.GangType_gang_angang {
			groupType = majongpb.CardsGroupType_CGT_ANGANG
		}
		if gangCard.Type == majongpb.GangType_gang_minggang {
			groupType = majongpb.CardsGroupType_CGT_MINGGANG
		}
		if gangCard.Type == majongpb.GangType_gang_bugang {
			groupType = majongpb.CardsGroupType_CGT_BUGANG
		}
		card := gutils.ServerCard2Number(gangCard.Card)
		cardsGroup := &majongpb.CardsGroup{
			Pid:   player.PalyerId,
			Type:  groupType,
			Cards: []uint32{card, card, card, card},
		}
		gangCardGroups = append(gangCardGroups, cardsGroup)
	}
	cardsGroupList = append(cardsGroupList, gangCardGroups...)

	// 花牌组
	huaCards := ServerCards2Numbers(player.GetHuaCards())
	cltHuaCard := make([]uint32, 0)
	for _, huaCard := range huaCards {
		cltHuaCard = append(cltHuaCard, uint32(huaCard))
	}
	huaCardGroup := &majongpb.CardsGroup{
		Cards: cltHuaCard,
		Type:  majongpb.CardsGroupType_CGT_HUA,
	}
	cardsGroupList = append(cardsGroupList, huaCardGroup)
	// 胡牌组
	huCardGroup := &majongpb.CardsGroup{
		Cards: []uint32{gutils.ServerCard2Number(huCard)},
		Type:  majongpb.CardsGroupType_CGT_HU,
	}
	cardsGroupList = append(cardsGroupList, huCardGroup)
	return cardsGroupList
}

// GetAllPlayers 所有玩家
func GetAllPlayers(mjContext *majongpb.MajongContext) (allPlayers []uint64) {
	allPlayers = make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}
	return
}

// GetHuPlayers 已胡牌玩家
func GetHuPlayers(mjContext *majongpb.MajongContext) (huPlayers []uint64) {
	huPlayers = make([]uint64, 0)
	for _, player := range mjContext.Players {
		if player.XpState == majongpb.XingPaiState_hu {
			huPlayers = append(huPlayers, player.GetPalyerId())
		}
	}
	return
}

// GetQuitPlayers 退出玩家
func GetQuitPlayers(mjContext *majongpb.MajongContext) (quitPalyers []uint64) {
	quitPalyers = make([]uint64, 0)
	for _, player := range mjContext.Players {
		if player.IsQuit {
			quitPalyers = append(quitPalyers, player.GetPalyerId())
		}
	}
	return
}

// GetGiveupPlayers 认输玩家
func GetGiveupPlayers(mjContext *majongpb.MajongContext) (giveupPlayers []uint64) {
	giveupPlayers = make([]uint64, 0)
	for _, player := range mjContext.Players {
		if player.XpState == majongpb.XingPaiState_give_up {
			giveupPlayers = append(giveupPlayers, player.GetPalyerId())
		}
	}
	return
}
