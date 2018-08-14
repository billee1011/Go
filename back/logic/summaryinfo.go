package logic

import (
	"encoding/json"
	"steve/back/data"
	"steve/entity/db"
	"steve/entity/gamelog"
	"time"

	"github.com/Sirupsen/logrus"
)

func insertSummaryInfo(info gamelog.TGameSummary) error {
	summary := db.TGameSumary{
		Sumaryid:      info.Sumaryid,
		Deskid:        info.Deskid,
		Gameid:        info.Gameid,
		Levelid:       info.Levelid,
		Gamestarttime: info.Gamestarttime,
		Gameovertime:  info.Gameovertime,
		Createtime:    time.Now(),
		Createby:      info.Createby,
	}
	if playeridsData, err := translantPlayerids(info.Playerids); err == nil {
		summary.Playerids = playeridsData
	}
	if scoreinfoData, err := translantScoreInfo(info.Scoreinfo); err == nil {
		summary.Scoreinfo = scoreinfoData
	}
	if winneridsData, err := translantWinnerids(info.Winnerids); err == nil {
		summary.Winnerids = winneridsData
	}
	if roundCurrencyData, err := translantRoundCurrency(info.Roundcurrency); err == nil {
		summary.Roundcurrency = roundCurrencyData
	}
	return data.InsertSummary(&summary)
}

func translantPlayerids(playerIDs []uint64) (string, error) {
	data := ""
	b, err := json.Marshal(playerIDs)
	if err != nil {
		logrus.Errorf("failed to json marshal with playerids,err:%v", err)
		return data, err
	}
	data = string(b)
	return data, nil
}

func translantScoreInfo(ScoreInfo []int64) (string, error) {
	data := ""
	b, err := json.Marshal(ScoreInfo)
	if err != nil {
		logrus.Errorf("failed to json marshal with ScoreInfo,err:%v", err)
		return data, err
	}
	data = string(b)
	return data, nil
}

func translantWinnerids(winnerIDs []uint64) (string, error) {
	data := ""
	b, err := json.Marshal(winnerIDs)
	if err != nil {
		logrus.Errorf("failed to json marshal with winnerIDs,err:%v", err)
		return data, err
	}
	data = string(b)
	return data, nil
}

func translantRoundCurrency(roundCurrency []gamelog.RoundCurrency) (string, error) {
	data := ""
	b, err := json.Marshal(roundCurrency)
	if err != nil {
		logrus.Errorf("failed to json marshal with roundCurrency,err:%v", err)
		return data, err
	}
	data = string(b)
	return data, nil
}
