package db

import (
	"time"
)

type TLoginRecord struct {
	Id             int64     `xorm:"pk autoincr BIGINT(20)"`
	Recordid       int64     `xorm:"not null BIGINT(20)"`
	Playerid       int64     `xorm:"not null BIGINT(20)"`
	Onlineduration int       `xorm:"default 0 comment('在线时长') INT(11)"`
	Gamingduration int       `xorm:"default 0 comment('游戏时长') INT(11)"`
	Area           string    `xorm:"VARCHAR(64)"`
	Loginchannel   int       `xorm:"comment('上一次登录游戏的渠道号：省ID + 渠道ID') INT(11)"`
	Logintype      int       `xorm:"comment('玩家上一次登陆游戏时，所选方式。') INT(11)"`
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
