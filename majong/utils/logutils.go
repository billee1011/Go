package utils

import (
	server_majongpb "steve/entity/majong"

	"github.com/Sirupsen/logrus"
)

// WithMajongContext 在 logEntry 中附加上麻将现场的一些常用属性
func WithMajongContext(logEntry *logrus.Entry, mjContext *server_majongpb.MajongContext) *logrus.Entry {
	logEntry = logEntry.WithFields(logrus.Fields{
		"game_id":               mjContext.GetGameId(),
		"cur_state":             mjContext.GetCurState(),
		"active_player":         mjContext.GetActivePlayer(),
		"last_out_card_color":   mjContext.GetLastOutCard().GetColor(),
		"last_out_card_point":   mjContext.GetLastOutCard().GetPoint(),
		"zhuangjia_index":       mjContext.GetZhuangjiaIndex(),
		"hu_players":            mjContext.GetLastHuPlayers(),
		"peng_player":           mjContext.GetLastPengPlayer(),
		"gang_player":           mjContext.GetLastGangPlayer(),
		"mopai_player":          mjContext.GetMopaiPlayer(),
		"last_mopai_player":     mjContext.GetLastMopaiPlayer(),
		"last_mopai_card_color": mjContext.GetLastMopaiCard().GetColor(),
		"last_mopai_card_point": mjContext.GetLastMopaiCard().GetPoint(),
		"gang_card_color":       mjContext.GetGangCard().GetColor(),
		"gang_card_point":       mjContext.GetGangCard().GetPoint(),
	})

	logEntry = WithMajongPlayer(logEntry, mjContext)
	return logEntry
}

// WithMajongPlayer 在 logEntry 中带上玩家的基础信息
func WithMajongPlayer(logEntry *logrus.Entry, mjContext *server_majongpb.MajongContext) *logrus.Entry {
	players := mjContext.GetPlayers()

	playerIDs := []uint64{}

	for _, player := range players {
		playerIDs = append(playerIDs, player.GetPalyerId())
	}
	logEntry = logEntry.WithField("player_id_list", playerIDs)
	return logEntry
}
