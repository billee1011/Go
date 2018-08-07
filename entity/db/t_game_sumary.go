package db

import (
	"time"
)

type TGameSumary struct {
	Sumaryid      int64     `xorm:"not null pk BIGINT(20)"`
	Deskid        int64     `xorm:"BIGINT(20)"`
	Gameid        int       `xorm:"not null INT(11)"`
	Levelid       int       `xorm:"not null comment('场次ID') INT(11)"`
	Playerids     string    `xorm:"comment('桌子内玩家，多个玩家用|分割') VARCHAR(256)"`
	Scoreinfo     string    `xorm:"comment('输赢分,顺序和ID相同') VARCHAR(256)"`
	Winnerids     string    `xorm:"comment('赢家ID，多个赢家用|分割') VARCHAR(256)"`
	Roundcurrency string    `xorm:"TEXT"`
	Gameovertime  time.Time `xorm:"DATETIME"`
	Createtime    time.Time `xorm:"DATETIME"`
	Createby      string    `xorm:"VARCHAR(64)"`
	Updatetime    time.Time `xorm:"DATETIME"`
	Updateby      string    `xorm:"VARCHAR(64)"`
}
