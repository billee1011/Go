package db

type TPlayerMail struct {
	NId          int64 `xorm:"not null pk autoincr comment('递增ID') BIGINT(20)"`
	NPlayerid    int64 `xorm:"not null comment('玩家ID') unique(t_player_mail_UN) index BIGINT(20)"`
	NMailid      int64 `xorm:"not null comment('邮件ID') unique(t_player_mail_UN) BIGINT(20)"`
	NIsread      int   `xorm:"comment('是否已读: 0=未读, 1=已读 ') INT(11)"`
	NIsgetattach int   `xorm:"comment('是否已领取附件: 0=未领, 1=已领') INT(11)"`
}
