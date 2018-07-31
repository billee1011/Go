package db

import (
	"time"
)

type TPlayerCurrency struct {
	Id             int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid       int64     `xorm:"not null unique BIGINT(20)"`
	Coins          int       `xorm:"comment('当前金币数') INT(11)"`
	Ingots         int       `xorm:"comment('当前面元宝数') INT(11)"`
	Keycards       int       `xorm:"comment('当前房卡') INT(11)"`
	Obtainingots   int       `xorm:"comment('总获得元宝') INT(11)"`
	Obtainkeycards int       `xorm:"comment('总获得房卡') INT(11)"`
	Costingots     int       `xorm:"comment('累计消耗元宝数') INT(11)"`
	Costkeycards   int       `xorm:"comment('累计消耗房卡数') INT(11)"`
	Remark         string    `xorm:"VARCHAR(256)"`
	Createtime     time.Time `xorm:"DATETIME"`
	Createby       string    `xorm:"VARCHAR(64)"`
	Updatetime     time.Time `xorm:"DATETIME"`
	Updateby       string    `xorm:"VARCHAR(64)"`
}
