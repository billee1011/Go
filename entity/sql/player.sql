/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.7.108
 Source Server Type    : MySQL
 Source Server Version : 50722
 Source Host           : 192.168.7.108:3306
 Source Schema         : player

 Target Server Type    : MySQL
 Target Server Version : 50722
 File Encoding         : 65001

 Date: 09/08/2018 11:36:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_currency_record
-- ----------------------------
DROP TABLE IF EXISTS `t_currency_record`;
CREATE TABLE `t_currency_record` ( 
  `tradeID` varchar(64) NOT NULL COMMENT '流水ID',
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `channel` int(11) DEFAULT NULL COMMENT '渠道ID',
  `currencyType` int(11) DEFAULT NULL COMMENT '货币类型: 1=金币, 2=元宝（钻石）， 3=房卡',
  `amount` int(11) DEFAULT NULL COMMENT '加减值',
  `beforeBalance` int(11) DEFAULT NULL COMMENT '操作前金币值',
  `afterBalance` int(11) DEFAULT NULL COMMENT '操作后金币值',
  `tradeTime` datetime DEFAULT NULL COMMENT '创建时间',
  `status` tinyint(1) DEFAULT NULL COMMENT '操作结果： 1=成功，0=失败',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `gameId` bigint(20) DEFAULT NULL COMMENT '游戏ID',
  `level` int(11) DEFAULT NULL COMMENT '场次ID',
  `funcId` int(11) DEFAULT NULL COMMENT '行为ID或功能ID',
  PRIMARY KEY (`tradeID`),
  UNIQUE KEY `t_currency_record_tradeID_uindex` (`tradeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='金币流水表';

-- ----------------------------
-- Table structure for t_hall_info
-- ----------------------------
DROP TABLE IF EXISTS `t_hall_info`;
CREATE TABLE `t_hall_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `recharge` int(11) DEFAULT NULL COMMENT '总充值金额',
  `bust` int(11) DEFAULT NULL COMMENT '累计破产次数',
  `lastGame` int(11) DEFAULT NULL COMMENT '上次金币场玩法',
  `lastLevel` int(11) DEFAULT NULL COMMENT '上次金币场场次',
  `lastFriendsBureauNum` int(11) DEFAULT NULL COMMENT '上次朋友局房号',
  `lastFriendsBureauGame` int(11) DEFAULT NULL COMMENT '上次朋友局玩法',
  `lastGameStartTime` datetime DEFAULT NULL COMMENT '最后游戏开始时间',
  `winningRate` int(11) DEFAULT NULL COMMENT '胜率',
  `backpackID` bigint(20) DEFAULT NULL COMMENT '背包ID',
  `almsGotTimes` int(11) DEFAULT NULL COMMENT '救济已领取次数',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=752 DEFAULT CHARSET=utf8 COMMENT='大厅信息表';

-- ----------------------------
-- Table structure for t_mail
-- ----------------------------
DROP TABLE IF EXISTS `t_mail`;
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统消息表，邮件表';

-- ----------------------------
-- Table structure for t_player
-- ----------------------------
DROP TABLE IF EXISTS `t_player`;
CREATE TABLE `t_player` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `accountID` bigint(20) NOT NULL COMMENT '账户ID',
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `showUID` bigint(20) NOT NULL COMMENT '显示ID',
  `type` int(11) NOT NULL DEFAULT '1' COMMENT '玩家类型1.普通玩家，2.机器人，3.QA\n2.\n3.',
  `channelID` int(11) DEFAULT NULL COMMENT '渠道ID',
  `nickname` varchar(64) DEFAULT NULL COMMENT '昵称',
  `gender` int(11) DEFAULT '1' COMMENT '性别',
  `avatar` varchar(256) DEFAULT NULL COMMENT '头像地址',
  `provinceID` int(11) DEFAULT NULL COMMENT '省ID',
  `cityID` int(11) DEFAULT NULL COMMENT '市ID',
  `name` varchar(64) DEFAULT NULL COMMENT '真实姓名',
  `phone` varchar(11) DEFAULT NULL COMMENT '手机号码',
  `idCard` varchar(20) DEFAULT NULL COMMENT '身份证',
  `isWhiteList` tinyint(1) DEFAULT '0' COMMENT '是否QA，默认否',
  `zipCode` int(11) DEFAULT NULL COMMENT '邮编',
  `shippingAddr` varchar(256) DEFAULT NULL COMMENT '收获地址',
  `status` int(11) DEFAULT '1' COMMENT '1可登录，2冻结，默认为1',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间，通常也是注册时间',
  `createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='玩家表';

-- ----------------------------
-- Table structure for t_player_currency
-- ----------------------------
DROP TABLE IF EXISTS `t_player_currency`;
CREATE TABLE `t_player_currency` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `coins` int(11) DEFAULT NULL COMMENT '当前金币数',
  `ingots` int(11) DEFAULT NULL COMMENT '当前元宝数',
  `keyCards` int(11) DEFAULT NULL COMMENT '当前房卡数',
  `obtainIngots` int(11) DEFAULT NULL COMMENT '总获得元宝数',
  `obtainKeyCards` int(11) DEFAULT NULL COMMENT '总获得房卡数',
  `costIngots` int(11) DEFAULT NULL COMMENT '累计消耗元宝数',
  `costKeyCards` int(11) DEFAULT NULL COMMENT '累计消耗房卡数',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='玩家货币表';

-- ----------------------------
-- Table structure for t_player_game
-- ----------------------------
DROP TABLE IF EXISTS `t_player_game`;
CREATE TABLE `t_player_game` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `gameID` int(11) DEFAULT NULL COMMENT '游戏ID',
  `gameName` varchar(64) DEFAULT NULL COMMENT '游戏名称',
  `winningRate` double DEFAULT NULL COMMENT '胜率，百分比表示，50%，只记录 50，精确到个位数',
  `winningBurea` int(11) DEFAULT NULL COMMENT '胜利局数',
  `totalBureau` int(11) DEFAULT NULL COMMENT '总局数',
  `maxWinningStream` int(11) DEFAULT NULL COMMENT '最高连胜',
  `maxMultiple` int(11) DEFAULT NULL COMMENT '最大倍数',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='玩家游戏信息表';

-- ----------------------------
-- Table structure for t_player_mail
-- ----------------------------
DROP TABLE IF EXISTS `t_player_mail`;
CREATE TABLE `t_player_mail` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '递增ID',
  `n_playerid` bigint(20) NOT NULL COMMENT '玩家ID',
  `n_mailID` bigint(20) NOT NULL COMMENT '邮件ID',
  `n_isRead` tinyint(1) DEFAULT NULL COMMENT '是否已读: 0=未读, 1=已读 ',
  `n_isGetAttach` tinyint(1) DEFAULT NULL COMMENT '是否已领取附件: 0=未领, 1=已领',
  `n_isDel` tinyint(1) DEFAULT '0' COMMENT '是否被用户删除: 0=未删除, 1=删除',
  `n_deleteTime` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`n_id`),
  UNIQUE KEY `t_player_mail_UN` (`n_playerid`,`n_mailID`),
  KEY `t_player_mail_n_playerid_IDX` (`n_playerid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='玩家邮件表'

-- ----------------------------
-- Table structure for t_player_props
-- ----------------------------
DROP TABLE IF EXISTS `t_player_props`;
CREATE TABLE `t_player_props` (
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `propID` bigint(20) NOT NULL COMMENT '道具ID',
  `count` bigint(20) NOT NULL COMMENT '道具数量',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(100) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(100) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`playerID`,`propID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='玩家道具表';

SET FOREIGN_KEY_CHECKS = 1;
