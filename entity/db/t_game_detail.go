package db

import (
	"time"
)

type TGameDetail struct {
	Detailid    int64     `xorm:"not null pk comment('明细ID') BIGINT(20)"`
	Sumaryid    int64     `xorm:"not null comment('汇总ID') BIGINT(20)"`
	Playerid    int64     `xorm:"not null comment('玩家ID') BIGINT(20)"`
	Deskid      int64     `xorm:"comment('桌子ID') BIGINT(20)"`
	Gameid      int       `xorm:"comment('游戏ID') INT(11)"`
	Levelid     int       `xorm:"comment('场次ID') INT(11)"`
	Amount      int64     `xorm:"comment('输赢金额') BIGINT(20)"`
	Iswinner    int       `xorm:"comment('是否赢家') TINYINT(1)"`
	Brokercount int       `xorm:"comment('破产次数') INT(11)"`
	Createtime  time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby    string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime  time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby    string    `xorm:"comment('更新人') VARCHAR(64)"`
}
