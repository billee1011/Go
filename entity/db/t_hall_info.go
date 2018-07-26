package db

import (
	"time"
)

type THallInfo struct {
	Id                    int64     `xorm:"pk autoincr BIGINT(20)"`
	Playerid              int64     `xorm:"not null BIGINT(20)"`
	Recharge              int       `xorm:"comment('总充值金额') INT(11)"`
	Coins                 int       `xorm:"comment('当前金币数') INT(11)"`
	Ingots                int       `xorm:"comment('当前面元宝数') INT(11)"`
	Keycards              int       `xorm:"comment('当前房卡') INT(11)"`
	Obtainingots          int       `xorm:"comment('总获得房卡') INT(11)"`
	Obtainkeycards        int       `xorm:"comment('总获得房卡') INT(11)"`
	Costingots            int       `xorm:"comment('累计消耗元宝数') INT(11)"`
	Costkeycards          int       `xorm:"comment('累计消耗房卡数') INT(11)"`
	Bust                  int       `xorm:"comment('总破产次数：单次金豆减少触发破产的次数') INT(11)"`
	Lastgame              int       `xorm:"comment('上次金币场场次') INT(11)"`
	Lastlevel             int       `xorm:"comment('上次金币场场次') INT(11)"`
	Lastfriendsbureaunum  int       `xorm:"comment('上次朋友局房号') INT(11)"`
	Lastfriendsbureaugame int       `xorm:"comment('上次朋友局玩法') INT(11)"`
	Lastgamestarttime     time.Time `xorm:"comment('最后游戏时间的开始时间') DATETIME"`
	Winningrate           int       `xorm:"comment('胜率') INT(11)"`
	Backpackid            int64     `xorm:"comment('背包ID') BIGINT(20)"`
	Remark                string    `xorm:"VARCHAR(256)"`
	Createtime            time.Time `xorm:"DATETIME"`
	Createby              string    `xorm:"VARCHAR(64)"`
	Updatetime            time.Time `xorm:"DATETIME"`
	Updateby              string    `xorm:"VARCHAR(64)"`
}
