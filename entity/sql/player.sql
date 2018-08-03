

create table t_hall_info
(
  id                    bigint auto_increment    
    primary key,
  playerID              bigint       not null,
  recharge              int          null
  comment '总充值金额',
  bust                  int          null
  comment '总破产次数：单次金豆减少触发破产的次数',
  lastGame              int          null
  comment '上次金币场场次',
  lastLevel             int          null
  comment '上次金币场场次',
  lastFriendsBureauNum  int          null
  comment '上次朋友局房号',
  lastFriendsBureauGame int          null
  comment '上次朋友局玩法',
  lastGameStartTime     datetime     null
  comment '最后游戏时间的开始时间',
  winningRate           int          null
  comment '胜率',
  backpackID            bigint       null
  comment '背包ID',
  remark                varchar(256) null,
  createTime            datetime     null,
  createBy              varchar(64)  null,
  updateTime            datetime     null,
  updateBy              varchar(64)  null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '玩家大厅信息表'; 
  

create table t_player
(
  id           bigint auto_increment
    primary key,
  accountID    bigint                 not null,
  playerID     bigint                 not null,
  showUID     int                 not null,
  type         int default '1'        not null
  comment '1.普通玩家
2.机器人
3.管理员',
  channelID    int                    null
  comment '渠道ID',
  nickname     varchar(64)            null
  comment '昵称',
  gender       int default '1'        null
  comment '性别：1.女，2.男',
  avatar       varchar(256)           null
  comment '头像',
  provinceID   int                    null 
  comment '省ID',
  cityID       int                    null
  comment '市ID',
  name         varchar(64)            null,  
  phone        varchar(11)            null,
  idCard       varchar(20)            null,
  isWhiteList  tinyint(1) default '0' null
  comment '是否白名单，默认为否，白名单通常是QA',
  zipCode      int                    null,
  shippingAddr varchar(256)           null,
  status       int default '1'        null
  comment '账号状态：1.可登陆，2.冻结，默认1',
  remark       varchar(256)           null,
  createTime   datetime               null,
  createBy     varchar(64)            null,
  updateTime   datetime               null,
  updateBy     varchar(64)            null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8  comment '玩家信息表';    

create table t_player_currency
(
  id             bigint auto_increment
    primary key,
  playerID       bigint       not null,
  coins          int          null
  comment '当前金币数',
  ingots         int          null
  comment '当前面元宝数',
  keyCards       int          null
  comment '当前房卡',
  obtainIngots   int          null
  comment '总获得元宝',
  obtainKeyCards int          null
  comment '总获得房卡',
  costIngots     int          null
  comment '累计消耗元宝数',
  costKeyCards   int          null
  comment '累计消耗房卡数',
  remark         varchar(256) null,
  createTime     datetime     null,
  createBy       varchar(64)  null,
  updateTime     datetime     null,
  updateBy       varchar(64)  null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '玩家虚拟货币表';  

CREATE TABLE `t_player_game` (
	`id` bigint(20) NOT NULL AUTO_INCREMENT,
	`playerID` bigint(20) NOT NULL,
	`gameID` int(11) DEFAULT NULL,
	`gameName` varchar(64) DEFAULT NULL,
	`winningRate` int(11) DEFAULT NULL COMMENT '百分比表示，50%，只记录 50，精确到个位数',
	`winningBurea` int(11) DEFAULT NULL COMMENT '获胜局数',
	`totalBureau` int(11) DEFAULT NULL COMMENT '总局数',
	`maxWinningStream` int(11) DEFAULT NULL COMMENT '最高连胜',
	`maxMultiple` int(11) DEFAULT NULL COMMENT '最大获胜倍数',
	`createTime` datetime DEFAULT NULL,
	`createBy` varchar(64) DEFAULT NULL,
	`updateTime` datetime DEFAULT NULL,
	`updateBy` varchar(64) DEFAULT NULL,
	PRIMARY KEY (`id`)
  ) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=latin1 COMMENT='玩家游戏汇总信息';

