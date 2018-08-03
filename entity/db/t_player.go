package db

import (
	"time"
)

type TPlayer struct {
	Id        int64 `xorm:"pk autoincr BIGINT(20)"`
	Accountid int64 `xorm:"not null BIGINT(20)"`
	Playerid  int64 `xorm:"not null BIGINT(20)"`
	Showuid   int64 `xorm:"comment('展示的用户ID,10位数') BIGINT(20)"`
	Type      int   `xorm:"not null default 1 comment('1.普通玩家
2.机器人
3.管理员') INT(11)"`
	Channelid    int       `xorm:"comment('渠道ID') INT(11)"`
	Nickname     string    `xorm:"comment('昵称') VARCHAR(64)"`
	Gender       int       `xorm:"default 1 comment('性别：1.女，2.男') INT(11)"`
	Avatar       string    `xorm:"comment('头像') VARCHAR(256)"`
	Provinceid   int       `xorm:"comment('省ID') INT(11)"`
	Cityid       int       `xorm:"comment('市ID') INT(11)"`
	Name         string    `xorm:"VARCHAR(64)"`
	Phone        string    `xorm:"VARCHAR(11)"`
	Idcard       string    `xorm:"VARCHAR(20)"`
	Iswhitelist  int       `xorm:"default 0 comment('是否白名单，默认为否，白名单通常是QA') TINYINT(1)"`
	Zipcode      int       `xorm:"INT(11)"`
	Shippingaddr string    `xorm:"VARCHAR(256)"`
	Status       int       `xorm:"default 1 comment('账号状态：1.可登陆，2.冻结，默认1') INT(11)"`
	Remark       string    `xorm:"VARCHAR(256)"`
	Createtime   time.Time `xorm:"DATETIME"`
	Createby     string    `xorm:"VARCHAR(64)"`
	Updatetime   time.Time `xorm:"DATETIME"`
	Updateby     string    `xorm:"VARCHAR(64)"`
}
