package db

import (
	"time"
)

type TPlayerCurrency struct {
	Id             int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid       int64     `xorm:"not null comment('玩家ID') BIGINT(20)"`
	Coins          int       `xorm:"comment('当前金币数') INT(11)"`
	Ingots         int       `xorm:"comment('当前元宝数') INT(11)"`
	Keycards       int       `xorm:"comment('当前房卡数') INT(11)"`
	Obtainingots   int       `xorm:"comment('总获得元宝数') INT(11)"`
	Obtainkeycards int       `xorm:"comment('总获得房卡数') INT(11)"`
	Costingots     int       `xorm:"comment('累计消耗元宝数') INT(11)"`
	Costkeycards   int       `xorm:"comment('累计消耗房卡数') INT(11)"`
	Remark         string    `xorm:"comment('备注') VARCHAR(256)"`
	Createtime     time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby       string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime     time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby       string    `xorm:"comment('更新人') VARCHAR(64)"`
}
