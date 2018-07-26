package db

import (
	"time"
)

type TGameLevelConfig struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)"`
	Gameid     int       `xorm:"INT(11)"`
	Levelid    int       `xorm:"INT(11)"`
	Name       string    `xorm:"VARCHAR(256)"`
	Basescores int       `xorm:"INT(11)"`
	Lowscores  int       `xorm:"INT(11)"`
	Highscores int       `xorm:"INT(11)"`
	Status     int       `xorm:"INT(11)"`
	Remark     string    `xorm:"VARCHAR(256)"`
	Createtime time.Time `xorm:"DATETIME"`
	Createby   string    `xorm:"VARCHAR(64)"`
	Updatetime time.Time `xorm:"DATETIME"`
	Updateby   string    `xorm:"VARCHAR(64)"`
}
