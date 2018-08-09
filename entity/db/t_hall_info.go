package db

import (
	"time"
)

type THallInfo struct {
	Id                    int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid              int64     `xorm:"not null BIGINT(20)"`
	Recharge              int       `xorm:"INT(11)"`
	Bust                  int       `xorm:"INT(11)"`
	Lastgame              int       `xorm:"INT(11)"`
	Lastlevel             int       `xorm:"INT(11)"`
	Lastfriendsbureaunum  int       `xorm:"INT(11)"`
	Lastfriendsbureaugame int       `xorm:"INT(11)"`
	Lastgamestarttime     time.Time `xorm:"DATETIME"`
	Winningrate           int       `xorm:"INT(11)"`
	Backpackid            int64     `xorm:"comment('ID') BIGINT(20)"`
	Remark                string    `xorm:"VARCHAR(256)"`
	Createtime            time.Time `xorm:"DATETIME"`
	Createby              string    `xorm:"VARCHAR(64)"`
	Updatetime            time.Time `xorm:"DATETIME"`
	Updateby              string    `xorm:"VARCHAR(64)"`
}
