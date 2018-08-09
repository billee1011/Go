create table t_currency_record
(
  tradeID       varchar(64)  not null
    primary key,
  playerID      bigint       not null,
  channel       int          null,
  currencyType  int          null
  comment '货币类型：1.金币，2.元宝，3，房卡',
  amount        int          null
  comment '变化金额',
  beforeBalance int          null
  comment '变化前余额',
  afterBalance  int          null
  comment '变化后余额',
  tradeTime     datetime     null
  comment '交易时间',
  status        int          null
  comment '1.成功，2失败',
  remark        varchar(256) null,
  constraint t_currency_record_tradeID_uindex
  unique (tradeID)
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '虚拟货币流水表';



CREATE TABLE `t_game_detail` (
  `detailID` bigint(20) NOT NULL,
  `sumaryID` bigint(20) NOT NULL,
  `playerID` bigint(20) NOT NULL,
  `deskID` bigint(20) DEFAULT NULL,
  `gameID` int(11) DEFAULT NULL,
  `amount` bigint(20) DEFAULT NULL,
  `isWinner` tinyint(1) DEFAULT NULL,
  `brokerCount` int(11) DEFAULT NULL,
  `createTime` datetime DEFAULT NULL,
  `createBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL,
  `updateTime` datetime DEFAULT NULL,
  `updateBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`detailID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;





CREATE TABLE `t_game_sumary` (
  `sumaryID` bigint(20) NOT NULL,
  `deskID` bigint(20) DEFAULT NULL,
  `gameID` int(11) NOT NULL,
  `levelID` int(11) NOT NULL COMMENT 'ID',
  `playerIDs` varchar(256) CHARACTER SET latin1 DEFAULT NULL COMMENT '|',
  `scoreInfo` varchar(256) CHARACTER SET latin1 DEFAULT NULL COMMENT ',ID',
  `winnerIDs` varchar(256) CHARACTER SET latin1 DEFAULT NULL COMMENT 'ID|',
  `roundCurrency` text CHARACTER SET latin1,
  `gameoverTime` datetime DEFAULT NULL,
  `createTime` datetime DEFAULT NULL,
  `createBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL,
  `updateTime` datetime DEFAULT NULL,
  `updateBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`sumaryID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;





create table t_login_record
(
  recordID       bigint          not null
    primary key,
  playerID       bigint          not null,
  onlineDuration int default '0' null
  comment '在线时长',
  gamingDuration int default '0' null
  comment '游戏时长',
  area           varchar(64)     null,
  loginChannel   int             null
  comment '上一次登录游戏的渠道号：省ID + 渠道ID',
  loginType      int             null
  comment '玩家上一次登陆游戏时，所选方式。',
  loginTime      datetime        null,
  logoutTime     datetime        null,
  ip             varchar(16)     null,
  loginDevice    varchar(32)     null,
  deviceCode     varchar(128)    null,
  createTime     datetime        null,
  createBy       varchar(64)     null,
  updateTime     datetime        null,
  updateBy       varchar(64)     null,
  constraint t_login_record_recordID_uindex
  unique (recordID)
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '玩家登录记录表';
