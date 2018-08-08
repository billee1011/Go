package db

import (
	"time"
)

type TCurrencyRecord struct {
	Tradeid       string    `xorm:"not null pk unique VARCHAR(64)"`
	Playerid      int64     `xorm:"not null BIGINT(20)"`
	Channel       int       `xorm:"INT(11)"`
	Currencytype  int       `xorm:"comment('1.2.3') INT(11)"`
	Amount        int       `xorm:"INT(11)"`
	Beforebalance int       `xorm:"INT(11)"`
	Afterbalance  int       `xorm:"INT(11)"`
	Tradetime     time.Time `xorm:"DATETIME"`
	Status        int       `xorm:"comment('1.2') INT(11)"`
	Remark        string    `xorm:"VARCHAR(256)"`
}
