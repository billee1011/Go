package db

import (
	"time"
)

type TMail struct {
	NId         int64     `xorm:"not null pk autoincr comment('递增ID') BIGINT(20)"`
	NTitle      string    `xorm:"comment('邮件标题') VARCHAR(150)"`
	NDetail     string    `xorm:"comment('邮件内容') TEXT"`
	NAttach     string    `xorm:"comment('邮件附件：json格式 ') VARCHAR(256)"`
	NDest       string    `xorm:"comment('发送对象:json格式') TEXT"`
	NState      int       `xorm:"comment('邮件状态') INT(11)"`
	NStarttime  time.Time `xorm:"comment('发送开始时间：未发送=0＞审核中=1＞已审核=2＞发送中=3＞发送结束=4＞已拒绝=5＞已撤回=6＞已 失效=7 ') DATETIME"`
	NEndtime    time.Time `xorm:"comment('发送截至时间') DATETIME"`
	NDeltime    time.Time `xorm:"comment('邮件删除时间') DATETIME"`
	NCreatetime time.Time `xorm:"comment('创建时间') DATETIME"`
	NCreateby   string    `xorm:"comment('创建人') VARCHAR(64)"`
	NUpdatetime time.Time `xorm:"comment('最后更新时间') DATETIME"`
	NUpdateby   string    `xorm:"comment('最后更新人') VARCHAR(64)"`
}
