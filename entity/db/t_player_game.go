package db

import (
	"time"
)

type TPlayerGame struct {
	Id               int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid         int64     `xorm:"not null BIGINT(20)"`
	Gameid           int       `xorm:"INT(11)"`
	Gamename         string    `xorm:"VARCHAR(64)"`
	Winningrate      int       `xorm:"comment('50% 50') INT(11)"`
	Winningburea     int       `xorm:"INT(11)"`
	Totalbureau      int       `xorm:"INT(11)"`
	Maxwinningstream int       `xorm:"INT(11)"`
	Maxmultiple      int       `xorm:"INT(11)"`
	Createtime       time.Time `xorm:"DATETIME"`
	Createby         string    `xorm:"VARCHAR(64)"`
	Updatetime       time.Time `xorm:"DATETIME"`
	Updateby         string    `xorm:"VARCHAR(64)"`
}
