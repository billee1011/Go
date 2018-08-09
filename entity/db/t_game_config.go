package db

import (
	"time"
)

type TGameConfig struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)"`
	Gameid     int       `xorm:"comment('游戏ID') INT(11)"`
	Name       string    `xorm:"comment('游戏名称') VARCHAR(128)"`
	Type       int       `xorm:"comment('游戏类型') INT(11)"`
	Minpeople  int       `xorm:"comment('允许最少人数') INT(11)"`
	Maxpeople  int       `xorm:"comment('允许最多人数') INT(11)"`
	Playform   int       `xorm:"comment('平台,1:安卓;2:ios') INT(11)"`
	Countryid  int       `xorm:"comment('国区（默认中国）') INT(11)"`
	Provinceid int       `xorm:"comment('省ID') INT(11)"`
	Cityid     int       `xorm:"comment('市ID') INT(11)"`
	Channelid  int       `xorm:"comment('渠道ID') INT(11)"`
	Createtime time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby   string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby   string    `xorm:"comment('更新人') VARCHAR(64)"`
}
