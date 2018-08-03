package db

import (
	"time"
)

type TPlayerGame struct {
	Id               int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid         int64     `xorm:"not null BIGINT(20)"`
	Gameid           int       `xorm:"INT(11)"`
	Gamename         string    `xorm:"VARCHAR(64)"`
	Winningrate      int       `xorm:"comment('百分比表示，50%，只记录 50，精确到个位数') INT(11)"`
	Winningburea     int       `xorm:"comment('获胜局数') INT(11)"`
	Winbureau        int       `xorm:"comment('胜利局数') INT(11)"`
	Totalbureau      int       `xorm:"comment('总局数') INT(11)"`
	Maxwinningstream int       `xorm:"comment('最高连胜') INT(11)"`
	Maxmultiple      int       `xorm:"comment('最大获胜倍数') INT(11)"`
	Createtime       time.Time `xorm:"DATETIME"`
	Createby         string    `xorm:"VARCHAR(64)"`
	Updatetime       time.Time `xorm:"DATETIME"`
	Updateby         string    `xorm:"VARCHAR(64)"`
}
