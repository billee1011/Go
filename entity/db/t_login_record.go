package db

import (
	"time"
)

type TLoginRecord struct {
	Recordid       int64     `xorm:"not null pk comment('记录id') unique BIGINT(20)"`
	Playerid       int64     `xorm:"not null comment('玩家ID') BIGINT(20)"`
	Onlineduration int       `xorm:"default 0 comment('在线时长（分钟）') INT(11)"`
	Gamingduration int       `xorm:"default 0 comment('游戏时长（分钟）') INT(11)"`
	Area           string    `xorm:"comment('所选地区') VARCHAR(64)"`
	Loginchannel   int       `xorm:"comment('登录渠道：省ID + 渠道ID') INT(11)"`
	Logintype      int       `xorm:"comment('登录方式') INT(11)"`
	Logintime      time.Time `xorm:"comment('登录时间') DATETIME"`
	Logouttime     time.Time `xorm:"comment('登出时间') DATETIME"`
	Ip             string    `xorm:"comment('登录IP') VARCHAR(16)"`
	Logindevice    string    `xorm:"comment('登录设备') VARCHAR(32)"`
	Devicecode     string    `xorm:"comment('设备IMEI（唯一识别码）') VARCHAR(128)"`
	Createtime     time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby       string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime     time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby       string    `xorm:"comment('更新人') VARCHAR(64)"`
}
