package db

import (
	"time"
)

type TPlayer struct {
	Id        int64 `xorm:"pk autoincr BIGINT(20)"`
	Accountid int64 `xorm:"not null comment('账户ID') BIGINT(20)"`
	Playerid  int64 `xorm:"not null comment('玩家ID') BIGINT(20)"`
	Showuid   int64 `xorm:"not null comment('显示ID') BIGINT(20)"`
	Type      int   `xorm:"not null default 1 comment('玩家类型1.普通玩家，2.机器人，3.QA
2.
3.') INT(11)"`
	Channelid    int       `xorm:"comment('渠道ID') INT(11)"`
	Nickname     string    `xorm:"comment('昵称') VARCHAR(64)"`
	Gender       int       `xorm:"default 1 comment('性别') INT(11)"`
	Avatar       string    `xorm:"comment('头像地址') VARCHAR(256)"`
	Provinceid   int       `xorm:"comment('省ID') INT(11)"`
	Cityid       int       `xorm:"comment('市ID') INT(11)"`
	Name         string    `xorm:"comment('真实姓名') VARCHAR(64)"`
	Phone        string    `xorm:"comment('手机号码') VARCHAR(11)"`
	Idcard       string    `xorm:"comment('身份证') VARCHAR(20)"`
	Iswhitelist  int       `xorm:"default 0 comment('是否QA，默认否') TINYINT(1)"`
	Zipcode      int       `xorm:"comment('邮编') INT(11)"`
	Shippingaddr string    `xorm:"comment('收获地址') VARCHAR(256)"`
	Status       int       `xorm:"default 1 comment('1可登录，2冻结，默认为1') INT(11)"`
	Remark       string    `xorm:"comment('备注') VARCHAR(256)"`
	Createtime   time.Time `xorm:"comment('创建时间，通常也是注册时间') DATETIME"`
	Createby     string    `xorm:"comment('创建人') VARCHAR(64)"`
	Updatetime   time.Time `xorm:"comment('更新时间') DATETIME"`
	Updateby     string    `xorm:"comment('更新人') VARCHAR(64)"`
}
