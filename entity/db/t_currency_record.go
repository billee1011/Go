package db

import (
	"time"
)

type TCurrencyRecord struct {
	Tradeid       string    `xorm:"not null pk unique VARCHAR(64)"`
	Playerid      int64     `xorm:"not null BIGINT(20)"`
	Channel       int       `xorm:"INT(11)"`
	Currencytype  int       `xorm:"comment('货币类型：1.金币，2.元宝，3，房卡') INT(11)"`
	Amount        int       `xorm:"comment('变化金额') INT(11)"`
	Beforebalance int       `xorm:"comment('变化前余额') INT(11)"`
	Afterbalance  int       `xorm:"comment('变化后余额') INT(11)"`
	Tradetime     time.Time `xorm:"comment('交易时间') DATETIME"`
	Status        int       `xorm:"comment('1.成功，2失败') INT(11)"`
	Remark        string    `xorm:"VARCHAR(256)"`
	Gameid        int       `xorm:"comment('游戏ID') INT(11)"`
	Level         int       `xorm:"comment('场次ID') INT(11)"`
	Funcid        int       `xorm:"comment('行为ID或功能ID') INT(11)"`
}
