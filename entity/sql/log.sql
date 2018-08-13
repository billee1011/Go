/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.7.108
 Source Server Type    : MySQL
 Source Server Version : 50722
 Source Host           : 192.168.7.108:3306
 Source Schema         : log

 Target Server Type    : MySQL
 Target Server Version : 50722
 File Encoding         : 65001

 Date: 09/08/2018 11:36:36
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
-- Table structure for t_game_detail
-- ----------------------------
DROP TABLE IF EXISTS `t_game_detail`;
CREATE TABLE `t_game_detail` (
  `detailID` bigint(20) NOT NULL COMMENT '明细ID',
  `sumaryID` bigint(20) NOT NULL COMMENT '汇总ID',
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `deskID` bigint(20) DEFAULT NULL COMMENT '桌子ID',
  `gameID` int(11) DEFAULT NULL COMMENT '游戏ID',
  `levelID` int(11) DEFAULT NULL COMMENT '场次ID',
  `amount` bigint(20) DEFAULT NULL COMMENT '输赢金额',
  `isWinner` tinyint(1) DEFAULT NULL COMMENT '是否赢家',
  `brokerCount` int(11) DEFAULT NULL COMMENT '破产次数',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`detailID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='游戏明细信息';

-- ----------------------------
-- Table structure for t_game_sumary
-- ----------------------------
DROP TABLE IF EXISTS `t_game_sumary`;
CREATE TABLE `t_game_sumary` (
  `sumaryID` bigint(20) NOT NULL COMMENT '汇总信息ID',
  `deskID` bigint(20) DEFAULT NULL COMMENT '桌子ID',
  `gameID` int(11) NOT NULL COMMENT '游戏ID',
  `levelID` int(11) NOT NULL COMMENT '场次ID',
  `playerIDs` varchar(256) CHARACTER SET latin1 DEFAULT NULL COMMENT '当前桌的所有玩家ID用","分割',
  `scoreInfo` varchar(256) CHARACTER SET latin1 DEFAULT NULL COMMENT '玩家得分情况',
  `winnerIDs` varchar(256) CHARACTER SET latin1 DEFAULT NULL COMMENT '赢家IDs',
  `roundCurrency` text CHARACTER SET latin1 COMMENT '牌局日志信息',
  `gamestartTime` datetime DEFAULT NULL COMMENT '游戏开始时间',
  `gameoverTime` datetime DEFAULT NULL COMMENT '游戏结束时间',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) CHARACTER SET latin1 DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`sumaryID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='游戏记录汇总表';


-- ----------------------------
-- Table structure for t_login_record
-- ----------------------------
DROP TABLE IF EXISTS `t_login_record`;
CREATE TABLE `t_login_record` (
  `recordID` bigint(20) NOT NULL COMMENT '记录id',
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `onlineDuration` int(11) DEFAULT '0' COMMENT '在线时长（分钟）',
  `gamingDuration` int(11) DEFAULT '0' COMMENT '游戏时长（分钟）',
  `area` varchar(64) DEFAULT NULL COMMENT '所选地区',
  `loginChannel` int(11) DEFAULT NULL COMMENT '登录渠道：省ID + 渠道ID',
  `loginType` int(11) DEFAULT NULL COMMENT '登录方式',
  `loginTime` datetime DEFAULT NULL COMMENT '登录时间',
  `logoutTime` datetime DEFAULT NULL COMMENT '登出时间',
  `ip` varchar(16) DEFAULT NULL COMMENT '登录IP',
  `loginDevice` varchar(32) DEFAULT NULL COMMENT '登录设备',
  `deviceCode` varchar(128) DEFAULT NULL COMMENT '设备IMEI（唯一识别码）',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(64) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`recordID`),
  UNIQUE KEY `t_login_record_recordID_uindex` (`recordID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='登录记录表';

SET FOREIGN_KEY_CHECKS = 1;
