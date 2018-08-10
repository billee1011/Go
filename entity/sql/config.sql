
/* 通用配置表 */
drop table if exists `t_common_config`;
create table if not exists `t_common_config` (
  `id` bigint auto_increment,
  `key` varchar(128) NOT NULL comment 'config key', 
  `subkey` varchar(128) NOT NULL comment 'config sub key',
  `value` text comment 'config context, json format',
  PRIMARY KEY ( `id` ),
  UNIQUE KEY(`key`, `subkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 comment '通用配置表';


create table t_game_config
(
  id         bigint auto_increment
    primary key,
  gameID     int          null,
  name       varchar(128) null
  comment '游戏名称',
  type       int          null
  comment '游戏类型',
  createTime datetime     null,
  createBy   varchar(64)  null,
  updateTime datetime     null,
  updateBy   varchar(64)  null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '游戏配置表';           



create table t_game_level_config
(
  id         bigint auto_increment
    primary key,
  gameID     int          null,
  levelID    int          null,
  name       varchar(256) null,
  baseScores int          null,
  lowScores  int          null,
  highScores int          null,
  minPeople  int          null,
  maxPeople int           null,
  status     int          null,
  remark     varchar(256) null,
  createTime datetime     null,
  createBy   varchar(64)  null,
  updateTime datetime     null,
  updateBy   varchar(64)  null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '游戏场次配置表';



CREATE TABLE `t_horse_race` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '数据递增ID',
  `n_channel` bigint(20) NOT NULL COMMENT '渠道ID',
  `n_prov` bigint(20) DEFAULT NULL COMMENT '省包ID',
  `n_city` bigint(20) DEFAULT NULL COMMENT '城市ID',
  `n_bUse` tinyint(1) DEFAULT '1' COMMENT '是否启用',
  `n_bUseParent` tinyint(1) DEFAULT '1' COMMENT '是否启用上级配置',
  `n_horseData` text COMMENT 'json格式的跑马灯配置，具体格式参考相关说明文件',
  PRIMARY KEY (`n_id`),
  KEY `t_horse_race_n_channel_IDX` (`n_channel`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COMMENT='跑马灯表'


CREATE TABLE `t_mail` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '递增ID',
  `n_title` varchar(150) DEFAULT NULL COMMENT '邮件标题',
  `n_detail` text COMMENT '邮件内容',
  `n_attach` varchar(256) DEFAULT NULL COMMENT '邮件附件：json格式 ',
  `n_dest` text COMMENT '发送对象:json格式',
  `n_state` int(11) NOT NULL COMMENT '邮件状态：未发送=0＞审核中=1＞已审核=2＞发送中=3＞发送结束=4＞已拒绝=5＞已撤回=6＞已失效=7 ',
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