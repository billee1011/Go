package db

import (
	"time"
)

type TGameDetail struct {
	Detailid   string    `xorm:"not null pk VARCHAR(64)"`
	Playerid   int64     `xorm:"not null BIGINT(20)"`
	Deskid     int       `xorm:"INT(11)"`
	Gameid     int       `xorm:"INT(11)"`
	Amount     int       `xorm:"INT(11)"`
	Iswinner   int       `xorm:"TINYINT(1)"`
	Createtime time.Time `xorm:"DATETIME"`
	Createby   string    `xorm:"VARCHAR(64)"`
	Updatetime time.Time `xorm:"DATETIME"`
	Updateby   string    `xorm:"VARCHAR(64)"`
}
