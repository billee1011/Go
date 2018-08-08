package db

import (
	"time"
)

type TPlayerCurrency struct {
	Id             int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid       int64     `xorm:"not null BIGINT(20)"`
	Coins          int       `xorm:"INT(11)"`
	Ingots         int       `xorm:"INT(11)"`
	Keycards       int       `xorm:"INT(11)"`
	Obtainingots   int       `xorm:"INT(11)"`
	Obtainkeycards int       `xorm:"INT(11)"`
	Costingots     int       `xorm:"INT(11)"`
	Costkeycards   int       `xorm:"INT(11)"`
	Remark         string    `xorm:"VARCHAR(256)"`
	Createtime     time.Time `xorm:"DATETIME"`
	Createby       string    `xorm:"VARCHAR(64)"`
	Updatetime     time.Time `xorm:"DATETIME"`
	Updateby       string    `xorm:"VARCHAR(64)"`
}
