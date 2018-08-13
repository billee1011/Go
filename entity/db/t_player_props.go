package db

import (
	"time"
)

type TPlayerProps struct {
	Playerid   int64     `xorm:"not null pk comment('玩家ID') BIGINT(20)"`
	Propid     int64     `xorm:"not null pk comment('道具ID') BIGINT(20)"`
	Count      int64     `xorm:"not null comment('道具数量') BIGINT(20)"`
	Createtime time.Time `xorm:"comment('创建时间') DATETIME"`
	Createby   string    `xorm:"comment('创建人') VARCHAR(100)"`
	Updatetime time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby   string    `xorm:"comment('更新人') VARCHAR(100)"`
}
