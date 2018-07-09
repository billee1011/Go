package cardtype

import (
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"

	"github.com/Sirupsen/logrus"
)

type cardTypeCalculator struct {
	calcs map[int]interfaces.CardTypeCalculator
}

func (ctc *cardTypeCalculator) Calculate(params interfaces.CardCalcParams) (cardTypes []int64, gengCount uint32) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "cardTypeCalculator.Calculate",
		"params":    params,
	})
	calc, exists := ctc.calcs[params.GameID]
	if !exists {
		logEntry.Errorln("game not found")
		return []int64{}, 0
	}
	return calc.Calculate(params)
}
func (ctc *cardTypeCalculator) CardTypeValue(gameID int, cardTypes []int64, gengCount uint32) (uint32, uint32) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "cardTypeCalculator.CardTypeValue",
		"game_id":   gameID,
	})
	calc, exists := ctc.calcs[gameID]
	if !exists {
		logEntry.Errorln("game not found")
		return 0, 0
	}
	return calc.CardTypeValue(gameID, cardTypes, gengCount)
}

func init() {
	ctc := &cardTypeCalculator{
		calcs: map[int]interfaces.CardTypeCalculator{
			gutils.SCXLGameID: &scxlCardTypeCalculator{}, // TODO game id
			gutils.SCXZGameID: &scxlCardTypeCalculator{},
		},
	}
	global.SetCardTypeCalculator(ctc)
}
