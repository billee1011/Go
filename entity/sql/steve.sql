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
)
  comment '虚拟货币流水表';

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
)
  comment '游戏配置表';

create table t_game_detail
(
  detailID   varchar(64) not null
    primary key,
  playerID   bigint      not null,
  deskID     int         null,
  gameID     int         null,
  amount     int         null,
  isWinner   tinyint(1)  null,
  createTime datetime    null,
  createBy   varchar(64) null,
  updateTime datetime    null,
  updateBy   varchar(64) null
)
  comment '游戏记录明细表';

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
  status     int          null,
  remark     varchar(256) null,
  createTime datetime     null,
  createBy   varchar(64)  null,
  updateTime datetime     null,
  updateBy   varchar(64)  null
)
  comment '游戏场次配置表';

create table t_game_sumary
(
  sumaryID   bigint       not null
    primary key,
  deskID     bigint       null,
  gameID     int          null,
  levelID    int          null
  comment '场次ID',
  playerIDs  varchar(256) null
  comment '桌子内玩家，多个玩家用|分割',
  winnerIDs  varchar(256) null
  comment '赢家ID，多个赢家用|分割',
  createTime datetime     null,
  createBy   varchar(64)  null,
  updateTime datetime     null,
  updateBy   varchar(64)  null
)
  comment '游戏记录汇总表';

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
)
  comment '玩家大厅信息表';

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
)
  comment '玩家登录记录表';

create table t_player
(
  id           bigint auto_increment
    primary key,
  accountID    bigint                 not null,
  playerID     bigint                 not null,
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
)
  comment '玩家信息表';

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
)
  comment '玩家虚拟货币表';

