package logic

import (
	"steve/back/data"
	"steve/entity/db"
	"steve/entity/gamelog"
	"time"
)

func insertDetailInfo(info gamelog.TGameDetail) error {
	detail := db.TGameDetail{
		Detailid:    info.Detailid,
		Sumaryid:    info.Sumaryid,
		Playerid:    int64(info.Playerid),
		Deskid:      info.Deskid,
		Gameid:      info.Gameid,
		Amount:      info.Amount,
		Iswinner:    info.Iswinner,
		Brokercount: info.BrokerCount,
		Createtime:  time.Now(),
	}
	return data.InsertDetail(&detail)
}
