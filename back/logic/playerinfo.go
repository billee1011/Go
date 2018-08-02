package logic

import (
	"steve/back/data"
	"steve/entity/gamelog"

	"github.com/Sirupsen/logrus"
)

// 更新总局数
func updateTotalBuearu(detailInfo gamelog.TGameDetail) {
	logEntry := logrus.WithFields(logrus.Fields{})
	totalBureau, err := data.GetTotalBureau(detailInfo.Playerid, detailInfo.Gameid)
	if err != nil {
		logEntry.Errorf("failed to get totalBureau,err：%v", err)
	}
	totalBureau++
	if err = data.UpdateTotalBureau(detailInfo.Playerid, detailInfo.Gameid, totalBureau); err != nil {
		logEntry.Errorf("failed to update totalBureau,err：%v", err)
	}
}
