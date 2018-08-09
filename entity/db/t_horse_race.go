package db

type THorseRace struct {
	NId         int64  `xorm:"not null pk autoincr comment('数据递增ID') BIGINT(20)"`
	NChannel    int64  `xorm:"not null comment('渠道ID') index BIGINT(20)"`
	NProv       int64  `xorm:"comment('省包ID') BIGINT(20)"`
	NCity       int64  `xorm:"comment('城市ID') BIGINT(20)"`
	NBuse       int    `xorm:"default 1 comment('是否启用') TINYINT(1)"`
	NBuseparent int    `xorm:"default 1 comment('是否启用上级配置') TINYINT(1)"`
	NHorsedata  string `xorm:"comment('json格式的跑马灯配置，具体格式参考相关说明文件') TEXT"`
}
