package mjdesk

import (
	"errors"
	"steve/common/mjoption"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

var errNoGameOption = errors.New("没有该游戏的游戏选项")

// fillContextOptions 填充麻将现场的 options
func fillContextOptions(gameID int, mjContext *majong.MajongContext) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "fillContextOptions",
		"game_id":   gameID,
	})
	gameOption := mjoption.GetGameOptions(gameID)
	if gameOption == nil {
		entry.Errorln(errNoGameOption)
		return errNoGameOption
	}
	mjContext.SettleOptionId = uint32(gameOption.SettleOptionID)
	mjContext.CardtypeOptionId = uint32(gameOption.CardTypeOptionID)
	mjContext.XingpaiOptionId = uint32(gameOption.XingPaiOptionID)
	return nil
}
