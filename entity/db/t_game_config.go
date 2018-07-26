package db

import (
	"time"
)

type TGameConfig struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)"`
	Gameid     int       `xorm:"INT(11)"`
	Name       string    `xorm:"comment('游戏名称') VARCHAR(128)"`
	Type       int       `xorm:"comment('游戏类型') INT(11)"`
	Createtime time.Time `xorm:"DATETIME"`
	Createby   string    `xorm:"VARCHAR(64)"`
	Updatetime time.Time `xorm:"DATETIME"`
	Updateby   string    `xorm:"VARCHAR(64)"`
}
