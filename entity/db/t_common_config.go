package db

type TCommonConfig struct {
	Id     int64  `xorm:"pk autoincr BIGINT(20)"`
	Key    string `xorm:"not null comment('config key') unique(key) VARCHAR(128)"`
	Subkey string `xorm:"not null comment('config sub key') unique(key) VARCHAR(128)"`
	Value  string `xorm:"comment('config context, json format') TEXT"`
}
