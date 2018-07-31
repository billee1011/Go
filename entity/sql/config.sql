

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
)
  comment '游戏场次配置表';


