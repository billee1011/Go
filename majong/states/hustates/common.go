package hustates

import majongpb "steve/server_pb/majong"

// addHuCard 添加胡的牌
func addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64, huType majongpb.HuType) {
	player.HuCards = append(player.GetHuCards(), &majongpb.HuCard{
		Card:      card,
		Type:      huType,
		SrcPlayer: srcPlayerID,
	})
}
