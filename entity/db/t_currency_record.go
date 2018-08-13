package db

import (
	"time"
)

type TCurrencyRecord struct {
	Tradeid       string    `xorm:"not null pk comment('流水ID') VARCHAR(64)"`
	Playerid      int64     `xorm:"not null comment('玩家ID') index BIGINT(20)"`
	Channel       int       `xorm:"comment('渠道ID') INT(11)"`
	Currencytype  int       `xorm:"comment('货币类型: 1=金币, 2=元宝（钻石）， 3=房卡') INT(11)"`
	Amount        int       `xorm:"comment('加减值') INT(11)"`
	Beforebalance int       `xorm:"comment('操作前金币值') INT(11)"`
	Afterbalance  int       `xorm:"comment('操作后金币值') INT(11)"`
	Tradetime     time.Time `xorm:"comment('创建时间') DATETIME"`
	Status        int       `xorm:"comment('操作结果： 1=成功，0=失败') TINYINT(1)"`
	Remark        string    `xorm:"comment('备注') VARCHAR(256)"`
	Gameid        int64     `xorm:"comment('游戏ID') BIGINT(20)"`
	Level         int       `xorm:"comment('场次ID') INT(11)"`
	Funcid        int       `xorm:"comment('行为ID或功能ID') INT(11)"`
}
