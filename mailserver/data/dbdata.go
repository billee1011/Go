package data

import (
	"encoding/json"
	"fmt"
	"steve/mailserver/define"
	"steve/structs"
	"strconv"
	"steve/entity/goods"
	"github.com/Sirupsen/logrus"
)

/*
	功能： 服务数据保存到Mysql.
	作者： SkyWang
	日期： 2018-8-7

CREATE TABLE `t_mail` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '递增ID',
  `n_title` varchar(150) DEFAULT NULL COMMENT '邮件标题',
  `n_detail` text COMMENT '邮件内容',
  `n_attach` varchar(256) DEFAULT NULL COMMENT '邮件附件：json格式 ',
  `n_dest` text COMMENT '发送对象:json格式',
  `n_state` int(11) DEFAULT NULL COMMENT '邮件状态',
  `n_starttime` datetime DEFAULT NULL COMMENT '发送开始时间： ',
  `n_endtime` datetime DEFAULT NULL COMMENT '发送截至时间',
  `n_deltime` datetime DEFAULT NULL COMMENT '邮件删除时间',
  `n_createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `n_createBy` varchar(64) DEFAULT NULL COMMENT '创建人',
  `n_updateTime` datetime DEFAULT NULL COMMENT '最后更新时间',
  `n_updateBy` varchar(64) DEFAULT NULL COMMENT '最后更新人',
  PRIMARY KEY (`n_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统消息表，邮件表'

CREATE TABLE `t_player_mail` (
  `n_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '递增ID',
  `n_playerid` bigint(20) NOT NULL COMMENT '玩家ID',
  `n_mailID` bigint(20) NOT NULL COMMENT '邮件ID',
  `n_isRead` int(11) DEFAULT NULL COMMENT '是否已读',
  `n_isGetAttach` int(11) DEFAULT NULL COMMENT '是否已领取附件',
  PRIMARY KEY (`n_id`),
  UNIQUE KEY `t_player_mail_UN` (`n_playerid`,`n_mailID`),
  KEY `t_player_mail_n_playerid_IDX` (`n_playerid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='玩家邮件表'

*/

const dbName = "player"


// 设置邮件为已读
func SetEmailReadTagFromDB(uid uint64, mailId uint64)  error {

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return fmt.Errorf("connect db error")
	}

	strCol := "n_playerid, n_mailID, n_isRead, n_isGetAttach, n_isDel "

	sql := fmt.Sprintf("insert into t_player_mail (%s) values('%d','%d','%d','%d','%d');",
		strCol, uid, mailId, 1, 0, 0)
	res, err := engine.Exec(sql)
	if err != nil {
		return err
	}
	if aff, err := res.RowsAffected(); aff == 0 {
		return err
	}

	return nil
}

// 从DB获取指定玩家的邮件列表
func GetUserMailFromDB(uid uint64) (map[uint64]*define.PlayerMail, error) {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return nil, fmt.Errorf("connect db error")
	}

	sql := fmt.Sprintf("select n_id, n_mailID, n_isRead, n_isGetAttach,n_isDel from t_player_mail where n_playerid='%d' ;", uid)
	res, err := engine.QueryString(sql)
	if err != nil {

		return nil, err
	}
	list := make(map[uint64]*define.PlayerMail)
	for _, row := range res {

		id, _ := strconv.ParseInt(row["n_id"], 10, 64)
		if id == 0 {
			continue
		}

		info := new(define.PlayerMail)
		info.Id = id

		info.PlayerId = uid
		mailId, _ := strconv.ParseInt(row["n_mailID"], 10, 64)
		info.MailId = uint64(mailId)

		isRead, _ := strconv.ParseInt(row["n_isRead"], 10, 64)
		if isRead !=  0 {
			info.IsRead = true
		} else {
			info.IsRead = false
		}
		isGet, _ := strconv.ParseInt(row["n_isGetAttach"], 10, 64)
		if isGet !=  0 {
			info.IsGetAttach = true
		} else {
			info.IsGetAttach = false
		}
		isDel, _ := strconv.ParseInt(row["n_isDel"], 10, 64)
		if isDel !=  0 {
			info.IsDel = true
		} else {
			info.IsDel = false
		}

		list[info.MailId] = info
	}

	return list, nil
}

// 从DB加载邮件列表
func LoadMailListFromDB() (map[uint64]*define.MailInfo, error) {

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return nil, fmt.Errorf("connect db error")
	}

	sql := fmt.Sprintf("select n_id, n_title, n_detail, n_attach, n_dest, n_state, n_starttime , n_endtime, n_deltime, n_updateTime from t_mail ;")
	res, err := engine.QueryString(sql)
	if err != nil {

		return nil, err
	}
	list := make(map[uint64]*define.MailInfo)
	for _, row := range res {

		id, _ := strconv.ParseInt(row["n_id"], 10, 64)
		if id == 0 {
			continue
		}

		info := new(define.MailInfo)
		info.Id = uint64(id)
		info.Title = row["n_title"]
		info.Detail = row["n_detail"]
		info.Attach = row["n_attach"]
		info.Dest = row["n_dest"]
		state, _ := strconv.ParseInt(row["n_state"], 10, 16)
		info.State = int8(state)
		info.StartTime = row["n_starttime"]
		info.EndTime = row["n_endtime"]
		info.DelTime = row["n_deltime"]
		info.UpdateTime = row["n_updateTime"]

		// 解析发送目标
		info.DestList = parseSendDest(info.Dest)
		if info.DestList == nil {
			logrus.Errorf("parseSendDest error: mailid=%d, tilte=%s", id, info.Title)
		}
		// 解析附件物品列表
		info.AttachGoods = parseAttachGoods(info.Attach)
		if info.AttachGoods == nil {
			logrus.Errorf("parseAttachGoods error: mailid=%d, tilte=%s", id, info.Title)
		}
		list[info.Id] = info
	}

	return list, nil
}

// 解析发送目标json
func parseSendDest( strJson string) []*define.SendDest {

	jsonObject := make([]*define.SendDest,0, 2)
	err := json.Unmarshal([]byte(strJson), jsonObject)
	if err != nil {
		return nil
	}
	return jsonObject
}

func MarshalSendDest(dest *define.SendDest) (string, error) {
	data, err := json.Marshal(dest)
	return string(data), err
}

// 解析附件物品json
func parseAttachGoods( strJson string) []*goods.Goods {
	jsonObject := make([]*goods.Goods,0, 2)
	err := json.Unmarshal([]byte(strJson), jsonObject)
	if err != nil {
		return nil
	}

	return jsonObject
}

func MarshalAttachGoods(g []*goods.Goods) (string, error) {
	data, err := json.Marshal(g)
	return string(data), err
}