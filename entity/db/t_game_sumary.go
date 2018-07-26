package db

import (
	"time"
)

type TGameSumary struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)"`
	Sumaryid   int64     `xorm:"BIGINT(20)"`
	Deskid     int64     `xorm:"BIGINT(20)"`
	Gameid     int       `xorm:"INT(11)"`
	Levelid    int       `xorm:"comment('场次ID') INT(11)"`
	Playerids  string    `xorm:"comment('桌子内玩家，多个玩家用|分割') VARCHAR(256)"`
	Winnerids  string    `xorm:"comment('赢家ID，多个赢家用|分割') VARCHAR(256)"`
	Createtime time.Time `xorm:"DATETIME"`
	Createby   string    `xorm:"VARCHAR(64)"`
	Updatetime time.Time `xorm:"DATETIME"`
	Updateby   string    `xorm:"VARCHAR(64)"`
}
