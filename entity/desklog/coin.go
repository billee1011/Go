package log

import "time"

// TGameSummary 游戏记录汇总
type TGameSummary struct {
	Sumaryid      int64     `json:"	Sumaryid "`
	Deskid        int64     `json:"	Deskid "`
	Gameid        int       `json:"	Gameid        "`
	Levelid       int       `json:"	Levelid       "`
	Playerids     string    `json:"	Playerids    "`
	Scoreinfo     string    `json:" Scoreinfo    "`
	Winnerids     string    `json:"	Winnerids    "`
	Roundcurrency string    `json:"	Roundcurrency"`
	Createtime    time.Time `json:"	Createtime    "`
	Createby      string    `json:"	Createby      "`
	Updatetime    time.Time `json:"	Updatetime    "`
	Updateby      string    `json:"	Updateby      "`
}

// TGameDetail 游戏明细
type TGameDetail struct {
	Sumaryid   int64     `json:"Sumaryid  "`
	Playerid   int64     `json:"Playerid  "`
	Deskid     int       `json:"Deskid    "`
	Gameid     int       `json:"Gameid    "`
	Amount     int       `json:"Amount    "`
	Iswinner   int       `json:"Iswinner  "`
	Createtime time.Time `json:"Createtime"`
	Createby   string    `json:"Createby  "`
	Updatetime time.Time `json:"Updatetime"`
	Updateby   string    `json:"Updateby  "`
}

// RoundCurrency 对局金币流水
type RoundCurrency struct {
	Settletype    int32          `json:"Settletype"`
	Settledetails []SettleDetail `json:"Settledetails"`
}

// SettleDetail 对局金币流水明细
type SettleDetail struct {
	Playerid  uint64 `json:"Playerid"`
	ChangeVal uint64 `json:"ChangeVal"`
}
