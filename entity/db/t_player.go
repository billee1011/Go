package db

import (
	"time"
)

type TPlayer struct {
	Id        int64 `xorm:"pk autoincr BIGINT(20)"`
	Accountid int64 `xorm:"not null BIGINT(20)"`
	Playerid  int64 `xorm:"not null BIGINT(20)"`
	Showuid   int64 `xorm:"not null BIGINT(20)"`
	Type      int   `xorm:"not null default 1 comment('1.
2.
3.') INT(11)"`
	Channelid    int       `xorm:"comment('ID') INT(11)"`
	Nickname     string    `xorm:"VARCHAR(64)"`
	Gender       int       `xorm:"default 1 comment('1.2.') INT(11)"`
	Avatar       string    `xorm:"VARCHAR(256)"`
	Provinceid   int       `xorm:"comment('ID') INT(11)"`
	Cityid       int       `xorm:"comment('ID') INT(11)"`
	Name         string    `xorm:"VARCHAR(64)"`
	Phone        string    `xorm:"VARCHAR(11)"`
	Idcard       string    `xorm:"VARCHAR(20)"`
	Iswhitelist  int       `xorm:"default 0 comment('QA') TINYINT(1)"`
	Zipcode      int       `xorm:"INT(11)"`
	Shippingaddr string    `xorm:"VARCHAR(256)"`
	Status       int       `xorm:"default 1 comment('1.2.1') INT(11)"`
	Remark       string    `xorm:"VARCHAR(256)"`
	Createtime   time.Time `xorm:"DATETIME"`
	Createby     string    `xorm:"VARCHAR(64)"`
	Updatetime   time.Time `xorm:"DATETIME"`
	Updateby     string    `xorm:"VARCHAR(64)"`
}
