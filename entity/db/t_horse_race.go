package db

type THorseRace struct {
	NId         int64  `xorm:"not null pk autoincr comment('ID') BIGINT(20)"`
	NChannel    int64  `xorm:"not null comment('ID') index BIGINT(20)"`
	NProv       int64  `xorm:"comment('ID') BIGINT(20)"`
	NCity       int64  `xorm:"comment('ID') BIGINT(20)"`
	NBuse       int    `xorm:"default 1 TINYINT(1)"`
	NBuseparent int    `xorm:"default 1 TINYINT(1)"`
	NHorsedata  string `xorm:"comment('json') TEXT"`
}
