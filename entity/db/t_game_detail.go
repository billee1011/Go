package db

import (
	"time"
)

type TGameDetail struct {
	Sumaryid   int64     `xorm:"not null pk BIGINT(20)"`
	Playerid   int64     `xorm:"not null BIGINT(20)"`
	Deskid     int64     `xorm:"BIGINT(20)"`
	Gameid     int       `xorm:"INT(11)"`
	Amount     int64     `xorm:"BIGINT(20)"`
	Iswinner   int       `xorm:"TINYINT(1)"`
	Createtime time.Time `xorm:"DATETIME"`
	Createby   string    `xorm:"VARCHAR(64)"`
	Updatetime time.Time `xorm:"DATETIME"`
	Updateby   string    `xorm:"VARCHAR(64)"`
}
