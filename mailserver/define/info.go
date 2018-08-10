package define

import (
	"steve/entity/goods"
)

/*
	功能: 邮件结构定义
	作者: Skywang
	日期: 2018-8-7

CREATE TABLE `t_mail` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '递增ID',
  `n_title` varchar(150) DEFAULT NULL COMMENT '邮件标题',
  `n_detail` text COMMENT '邮件内容',
  `n_attach` varchar(256) DEFAULT NULL COMMENT '邮件附件：json格式 ',
  `n_dest` text COMMENT '发送对象:json格式',
  `n_state` int(11) DEFAULT NULL COMMENT '邮件状态：未发送=0＞审核中=1＞已审核=2＞发送中=3＞发送结束=4＞已拒绝=5＞已撤回=6＞已失效=7 ',
  `n_starttime` datetime DEFAULT NULL COMMENT '发送开始时间: 2018-08-08 12:00:00',
  `n_endtime` datetime DEFAULT NULL COMMENT '发送截至时间: 2018-08-18 12:00:00',
  `n_deltime` datetime DEFAULT NULL COMMENT '邮件删除时间: 2018-09-18 12:00:00',
  `n_createTime` datetime DEFAULT NULL COMMENT '创建时间: 2018-08-08 12:00:00',
  `n_createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `n_updateTime` datetime DEFAULT NULL COMMENT '最后更新时间: 2018-08-08 12:00:00',
  `n_updateBy` varchar(64) DEFAULT NULL COMMENT '最后更新人',
  `n_isUseEndTime` tinyint(1) DEFAULT '1' COMMENT '是否启用截至时间',
  `n_isUseDelTime` tinyint(1) DEFAULT '1' COMMENT '是否启用删除时间',
  PRIMARY KEY (`n_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='系统消息表，邮件表'

CREATE TABLE `t_player_mail` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '递增ID',
  `n_playerid` bigint(20) NOT NULL COMMENT '玩家ID',
  `n_mailID` bigint(20) NOT NULL COMMENT '邮件ID',
  `n_isRead` int(11) DEFAULT NULL COMMENT '是否已读',
  `n_isGetAttach` int(11) DEFAULT NULL COMMENT '是否已领取附件',

*/

// 玩家邮件信息
type PlayerMail struct {
	Id          int64  // 唯一编号
	PlayerId    uint64 // 玩家ID
	MailId      uint64 // 邮件ID
	IsRead      bool   // 是否已读
	IsGetAttach bool   // 是否已领取附件
	IsDel		bool   // 是否被玩家删除
}

type MailInfo struct {
	Id          uint64         // 唯一编号
	Title       string         // 邮件标题
	Detail      string         // 邮件内容
	Attach      string         // 邮件附件：json格式
	Dest        string         // 发送对象:json格式
	State       int8           // 邮件状态: 未发送=0＞审核中=1＞已审核=2＞发送中=3＞发送结束=4＞已拒绝=5＞已撤回=6＞已 失效=7
	StartTime   string         // 发送开始时间
	EndTime     string         // 发送截至时间
	DelTime     string         // 邮件删除时间
	UpdateTime  string         // 最后更新时间
	AttachGoods []*goods.Goods // 附件奖励列表
	DestList    []*SendDest    // 发送目标

	IsUseEndTime bool 			// 是否使用截至时间
	IsUseDelTime bool			//  是否使用删除时间
}

// 发送目标
type SendDest struct {
	SendType    int8     	`json:"sendType"`    		// 发送类型: 0=全部玩家, 1=指定玩家
	Channel 	int64  		`json:"channel"` 			// 渠道ID
	Prov    	int64  		`json:"prov"`    			// 省包ID
	PlayerList  []uint64 	`json:"playerList"`  	// 玩家列表
}
