
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


drop table if exists `t_game_config`;
create table if not exists t_game_config
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;