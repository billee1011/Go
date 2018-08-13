package db

import (
	"time"
)

type TGameSumary struct {
	Sumaryid      int64     `xorm:"not null pk comment('汇总信息ID') BIGINT(20)"`
	Deskid        int64     `xorm:"comment('桌子ID') BIGINT(20)"`
	Gameid        int       `xorm:"not null comment('游戏ID') INT(11)"`
	Levelid       int       `xorm:"not null comment('场次ID') INT(11)"`
	Playerids     string    `xorm:"comment('当前桌的所有玩家ID用","分割') VARCHAR(256)"`
	Scoreinfo     string    `xorm:"comment('玩家得分情况') VARCHAR(256)"`
	Winnerids     string    `xorm:"comment('赢家IDs') VARCHAR(256)"`
	Roundcurrency string    `xorm:"comment('牌局日志信息') TEXT"`
	Gamestarttime time.Time `xorm:"comment('游戏开始时间') DATETIME"`
	Gameovertime  time.Time `xorm:"comment('游戏结束时间') DATETIME"`
	Createtime    time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby      string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime    time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby      string    `xorm:"comment('更新人') VARCHAR(64)"`
}
