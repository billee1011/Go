package db

import (
	"time"
)

type TPlayerGame struct {
	Id               int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid         int64     `xorm:"not null comment('玩家ID') BIGINT(20)"`
	Gameid           int       `xorm:"comment('游戏ID') INT(11)"`
	Gamename         string    `xorm:"comment('游戏名称') VARCHAR(64)"`
	Winningrate      float64   `xorm:"comment('胜率，百分比表示，50%，只记录 50，精确到个位数') DOUBLE"`
	Winningburea     int       `xorm:"comment('胜利局数') INT(11)"`
	Totalbureau      int       `xorm:"comment('总局数') INT(11)"`
	Maxwinningstream int       `xorm:"comment('最高连胜') INT(11)"`
	Maxmultiple      int       `xorm:"comment('最大倍数') INT(11)"`
	Createtime       time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby         string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime       time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby         string    `xorm:"comment('更新人') VARCHAR(64)"`
}
