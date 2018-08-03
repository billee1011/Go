package logic

import (
	"steve/back/data"
	"steve/entity/db"
	"steve/entity/gamelog"
)

func insertDetailInfo(info gamelog.TGameDetail) error {
	detail := db.TGameDetail{
		Detailid:   info.Detailid,
		Sumaryid:   info.Sumaryid,
		Playerid:   int64(info.Playerid),
		Deskid:     info.Deskid,
		Gameid:     info.Gameid,
		Amount:     info.Amount,
		Iswinner:   info.Iswinner,
		Createtime: info.Createtime,
		Createby:   info.Createby,
		Updatetime: info.Updatetime,
		Updateby:   info.Updateby,
	}
	return data.InsertDetail(&detail)
}
