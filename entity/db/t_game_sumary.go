package db

import (
	"time"
)

type TGameSumary struct {
	Sumaryid      int64     `xorm:"not null pk BIGINT(20)"`
	Deskid        int64     `xorm:"BIGINT(20)"`
	Gameid        int       `xorm:"not null INT(11)"`
	Levelid       int       `xorm:"not null comment('ID') INT(11)"`
	Playerids     string    `xorm:"comment('|') VARCHAR(256)"`
	Scoreinfo     string    `xorm:"comment(',ID') VARCHAR(256)"`
	Winnerids     string    `xorm:"comment('ID|') VARCHAR(256)"`
	Roundcurrency string    `xorm:"TEXT"`
	Gameovertime  time.Time `xorm:"DATETIME"`
	Createtime    time.Time `xorm:"DATETIME"`
	Createby      string    `xorm:"VARCHAR(64)"`
	Updatetime    time.Time `xorm:"DATETIME"`
	Updateby      string    `xorm:"VARCHAR(64)"`
}
