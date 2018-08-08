package db

import (
	"time"
)

type THallInfo struct {
	Id                    int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid              int64     `xorm:"not null unique BIGINT(20)"`
	Recharge              int       `xorm:"comment('总充值金额') INT(11)"`
	Bust                  int       `xorm:"comment('总破产次数：单次金豆减少触发破产的次数') INT(11)"`
	Lastgame              int       `xorm:"comment('上次金币场场次') INT(11)"`
	Lastlevel             int       `xorm:"comment('上次金币场场次') INT(11)"`
	Lastfriendsbureaunum  int       `xorm:"comment('上次朋友局房号') INT(11)"`
	Lastfriendsbureaugame int       `xorm:"comment('上次朋友局玩法') INT(11)"`
	Lastgamestarttime     time.Time `xorm:"comment('最后游戏时间的开始时间') DATETIME"`
	Backpackid            int64     `xorm:"comment('背包ID') BIGINT(20)"`
	Almsgottimes          int       `xorm:"comment('救济已领取次数') INT(11)"`
	Remark                string    `xorm:"VARCHAR(256)"`
	Createtime            time.Time `xorm:"DATETIME"`
	Createby              string    `xorm:"VARCHAR(64)"`
	Updatetime            time.Time `xorm:"DATETIME"`
	Updateby              string    `xorm:"VARCHAR(64)"`
}
