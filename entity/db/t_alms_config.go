package db

import (
	"time"
)

type TAlmsConfig struct {
	Id               int64     `xorm:"pk autoincr BIGINT(20)"`
	Almscountdonw    int       `xorm:"comment('救济倒计时，时间是秒') INT(11)"`
	Depositcountdonw int       `xorm:"comment('快充倒计时，时间是秒') INT(11)"`
	Getnorm          int       `xorm:"comment('救济线') INT(11)"`
	Gettimes         int       `xorm:"comment('救济领取次数') INT(11)"`
	Getnumber        int       `xorm:"comment('领取数量') INT(11)"`
	Version          int       `xorm:"comment('配置版本号，每次改变增加1,初始1') INT(11)"`
	Createtime       time.Time `xorm:"DATETIME"`
	Createby         string    `xorm:"VARCHAR(64)"`
	Updatetime       time.Time `xorm:"DATETIME"`
	Updateby         string    `xorm:"VARCHAR(64)"`
}
