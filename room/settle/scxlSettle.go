package settle

import (
	"steve/room/interfaces"
	server_pb "steve/server_pb/majong"
)

type scxlSettle struct{}

// Settle Desk玩家对结算信息立即扣分
func (s *scxlSettle) Settle(desk interfaces.Desk, mjContext server_pb.MajongContext) {
	lastSettleInfo := mjContext.SettleInfos[len(mjContext.SettleInfos)-1]
	for _, player := range desk.GetPlayers() {
		score := lastSettleInfo.Scores[*player.PlayerId]
		coin := int64(*player.Coin)
		if score <= 0 {
			*player.Coin = uint64(coin + score)
		} else if coin >= score {
			*player.Coin = uint64(coin - score)
		} else if coin < score {
			*player.Coin = uint64(coin - coin)
			lastSettleInfo.Scores[*player.PlayerId] = coin
		}
	}
}

func (s *scxlSettle) RoundSettle(desk interfaces.Desk, mjContext server_pb.MajongContext) {
}
