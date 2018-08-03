package db

import (
	"time"
)

type TGameLevelConfig struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)"`
	Gameid     int       `xorm:"INT(11)"`
	Levelid    int       `xorm:"INT(11)"`
	Name       string    `xorm:"VARCHAR(256)"`
	Fee        int       `xorm:"comment('桌费') INT(11)"`
	Basescores int       `xorm:"INT(11)"`
	Lowscores  int       `xorm:"INT(11)"`
	Highscores int       `xorm:"INT(11)"`
	Minpeople  int       `xorm:"INT(11)"`
	Maxpeople  int       `xorm:"INT(11)"`
	Showpeople int       `xorm:"comment('显示实时人数') INT(11)"`
	Status     int       `xorm:"INT(11)"`
	Tag        int       `xorm:"comment('标签，1.热门；2.New;') INT(11)"`
	Remark     string    `xorm:"VARCHAR(256)"`
	Createtime time.Time `xorm:"DATETIME"`
	Createby   string    `xorm:"VARCHAR(64)"`
	Updatetime time.Time `xorm:"DATETIME"`
	Updateby   string    `xorm:"VARCHAR(64)"`
}
