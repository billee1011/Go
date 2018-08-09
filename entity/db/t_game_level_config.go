package db

import (
	"time"
)

type TGameLevelConfig struct {
	Id               int64     `xorm:"pk autoincr BIGINT(20)"`
	Gameid           int       `xorm:"INT(11)"`
	Levelid          int       `xorm:"INT(11)"`
	Name             string    `xorm:"VARCHAR(256)"`
	Fee              int       `xorm:"comment('费用') INT(11)"`
	Basescores       int       `xorm:"INT(11)"`
	Lowscores        int       `xorm:"INT(11)"`
	Highscores       int       `xorm:"INT(11)"`
	Realonlinepeople int       `xorm:"comment('实时在线人数') INT(11)"`
	Showonlinepeople int       `xorm:"comment('显示在线人数') INT(11)"`
	Status           int       `xorm:"INT(11)"`
	Tag              int       `xorm:"comment('标签：1.热门；2.New') INT(11)"`
	Isalms           int       `xorm:"comment('是否为救济金场，0：关闭，1：开启') INT(11)"`
	Remark           string    `xorm:"VARCHAR(256)"`
	Createtime       time.Time `xorm:"DATETIME"`
	Createby         string    `xorm:"VARCHAR(64)"`
	Updatetime       time.Time `xorm:"DATETIME"`
	Updateby         string    `xorm:"VARCHAR(64)"`
}
