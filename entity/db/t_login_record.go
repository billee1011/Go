package db

import (
	"time"
)

type TLoginRecord struct {
	Recordid       int64     `xorm:"not null pk unique BIGINT(20)"`
	Playerid       int64     `xorm:"not null BIGINT(20)"`
	Onlineduration int       `xorm:"default 0 INT(11)"`
	Gamingduration int       `xorm:"default 0 INT(11)"`
	Area           string    `xorm:"VARCHAR(64)"`
	Loginchannel   int       `xorm:"comment('ID + ID') INT(11)"`
	Logintype      int       `xorm:"INT(11)"`
	Logintime      time.Time `xorm:"DATETIME"`
	Logouttime     time.Time `xorm:"DATETIME"`
	Ip             string    `xorm:"VARCHAR(16)"`
	Logindevice    string    `xorm:"VARCHAR(32)"`
	Devicecode     string    `xorm:"VARCHAR(128)"`
	Createtime     time.Time `xorm:"DATETIME"`
	Createby       string    `xorm:"VARCHAR(64)"`
	Updatetime     time.Time `xorm:"DATETIME"`
	Updateby       string    `xorm:"VARCHAR(64)"`
}
