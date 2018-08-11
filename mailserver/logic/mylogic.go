package logic

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"steve/client_pb/mailserver"
	"steve/external/hallclient"
	"steve/mailserver/data"
	"steve/mailserver/define"
	"time"
	"steve/entity/goods"
	"steve/structs"
)

/*
  功能： 邮件管理:
		1.获取未读邮件总数
		2.获取邮件列表
		3.获取指定邮件详情.
		4.删除邮件
		5.领取附件.
  作者： SkyWang
  日期： 2018-8-7
*/

// 是否是主节点
var isMasterNode = false
// 清理过期邮件开始点数
var clearBeginHour = 4
// 清理过期邮件结束点数
var clearEndHour = 6

// 邮件列表
var mailList map[uint64]*define.MailInfo
// 省包节点列表
var provSendList map[int64][]*define.MailInfo



func Init() error {
	err := getDataFromDB()
	if err != nil {
		logrus.Errorf("getDataFromDB first err:%v", err)
		//return err
	} else {
		logrus.Debugf("getDataFromDB first win...")
	}
	//testJsonObject()
	// 启动跑马灯是否开始检测协程
	go runCheckMailChange()
	return nil
}

func testJsonObject() {
	dst := &define.SendDest{}
	dst.SendType = 0
	dst.Prov = 1
	dst.Channel = 0
	jsonDest, _ := data.MarshalSendDest(dst)
	logrus.Debugln(jsonDest)

	gs := &goods.Goods{}
	gs.GoodsType = 0
	gs.GoodsId = 1
	gs.GoodsNum = 100

	jsonGs, _ := data.MarshalAttachGoods([]*goods.Goods{gs,gs})
	logrus.Debugln(jsonGs)
}

// 获取未读消息总数
func GetGetUnReadSum(uid uint64) (int32, error) {
	// 获取玩家渠道ID
	channel, prov, _, ok := getUserInfo(uid)
	if !ok {
		return 0, errors.New("获取玩家渠道ID失败")
	}

	// 从DB获取玩家的已读邮件列表
	readList, err := data.GetUserMailFromDB(uid)
	if err != nil {
		return 0, errors.New("从DB获取玩家已读邮件列表失败")
	}
	sum := int32(0)

	// 获取玩家所属省包的邮件
	if prov >= 0 {
		list := provSendList[prov]
		for _, mail := range list {

			// 检测是否符合省包和渠道ID
			isOk := checkMailProvChannel(uid, mail, channel, prov)
			if !isOk {
				continue
			}

			if one, ok := readList[mail.Id]; !ok || !one.IsRead {
				sum++
			}
		}
	}
	return sum, nil
}

// 获取邮件消息列表
func GetMailList(uid uint64) ([]*mailserver.MailTitle, error) {
	// 获取玩家渠道ID
	channel, prov, _, ok := getUserInfo(uid)
	if !ok {
		return nil, errors.New("获取玩家渠道ID失败")
	}
	if prov < 0 {
		return nil, errors.New("获取玩家省包ID < 0")
	}

	// 从DB获取玩家的已读邮件列表
	readList, err := data.GetUserMailFromDB(uid)
	if err != nil {
		return nil, errors.New("从DB获取玩家已读邮件列表失败")
	}

	titleList := make([]*mailserver.MailTitle, 0, 5)
	// 获取玩家所属省包的邮件

	list := provSendList[prov]
	for _, mail := range list {

		// 检测是否符合省包和渠道ID
		isOk := checkMailProvChannel(uid,mail, channel, prov)
		if !isOk {
			continue
		}
		title := new(mailserver.MailTitle)
		title.MailId = &mail.Id
		title.MailTitle = &mail.Title
		title.CreateTime = &mail.StartTime

		isRead := int32(0)
		one, ok := readList[mail.Id]
		if !ok || !one.IsRead {

		} else {
			isRead = 1
		}
		title.IsRead = &isRead

		isHaveAttach := int32(0)
		if len(mail.Attach) > 0 {
			isHaveAttach = 1
		}
		if one != nil && one.IsGetAttach {
			isHaveAttach = 2
		}
		title.IsHaveAttach = &isHaveAttach

		titleList = append(titleList, title)
	}

	return titleList, nil
}

// 获取指定邮件详情
func GetMailDetail(uid uint64, mailId uint64) (*mailserver.MailDetail, error) {

	mail, ok := mailList[mailId]
	if !ok {
		return nil, errors.New("指定邮件不存在")
	}
	if mail.State != define.StateSended && mail.State != define.StateSending {
		return nil, errors.New("指定邮件状态错误")
	}
	// 从DB获取玩家的已读邮件列表
	one, _ := data.GetTheMailFromDB(uid, mailId)

	if one != nil && one.IsDel {
		return nil, errors.New("邮件已被用户删除")
	}

	if one == nil {
		// 设置邮件=已读
		data.SetEmailReadTagFromDB(uid, mailId, true, mail.DelTime)
	} else {
		if !one.IsRead {
			// 设置邮件=已读
			data.SetEmailReadTagFromDB(uid, mailId, false, mail.DelTime)
		}
	}

	detail := new(mailserver.MailDetail)
	detail.MailId = &mail.Id
	detail.MailTitle = &mail.Title
	detail.Content = &mail.Detail

	t := mailserver.GoodsType_GOODSTYPE_PROPS

	for _, ach := range  mail.AttachGoods {
		newGoods := new(mailserver.Goods)
		newGoods.GoodsType = &t
		newGoods.GoodsId = &ach.GoodsId
		newGoods.GoodsNum = &ach.GoodsNum
	}

	isRead := int32(0)
	detail.IsRead = &isRead

	isHaveAttach := int32(0)
	if len(mail.Attach) > 0 {
		isHaveAttach = 1
	}
	if one != nil && one.IsGetAttach {
		isHaveAttach = 2
	}
	detail.IsHaveAttach = &isHaveAttach

	return detail, nil
}

// 标记邮件为已读请求
func SetReadTag(uid uint64, mailId uint64) error {
	mail, ok := mailList[mailId]
	if !ok {
		return errors.New("指定邮件不存在")
	}

	if mail.State != define.StateSended && mail.State != define.StateSending {
		return errors.New("指定邮件状态错误")
	}

	// 从DB获取玩家的已读邮件列表
	one, _ := data.GetTheMailFromDB(uid, mailId)

	if one != nil && one.IsDel {
		return errors.New("邮件已被用户删除")
	}

	if one == nil {
		// 设置邮件=已读
		data.SetEmailReadTagFromDB(uid, mailId, true, mail.DelTime)
	} else {
		if !one.IsRead {
			// 设置邮件=已读
			data.SetEmailReadTagFromDB(uid, mailId, false, mail.DelTime)
		}
	}

	return nil
}

// 删除邮件
func DelMail(uid uint64, mailId uint64) error {

	mail, ok := mailList[mailId]
	if !ok {
		return errors.New("指定邮件不存在")
	}
	// 从DB获取玩家的已读邮件列表
	one, _ := data.GetTheMailFromDB(uid, mailId)
	if one != nil && one.IsDel {
		return errors.New("邮件已被用户删除")
	}

	if one == nil {
		// 设置邮件为删除状态
		return data.DelEmailFromDB(uid, mailId, true,mail.DelTime)
	}
	return data.DelEmailFromDB(uid, mailId, false,mail.DelTime)

}

// 领取附件奖励请求
func AwardAttach(uid uint64, mailId uint64) ([]*mailserver.Goods, error) {

	// 从DB获取玩家的已读邮件列表
	one, _ := data.GetTheMailFromDB(uid, mailId)
	if one != nil && one.IsDel {
		return nil, errors.New("邮件已被用户删除")
	}

	if one == nil {
		return nil, errors.New("邮件不存在")
	}

	if !one.IsRead {
		return nil, errors.New("必须先读邮件, 才能领取附件")
	}

	// 如果已领取，直接返回
	if one.IsGetAttach  {
		return nil, errors.New("邮件已领取")
	}

	// 发放附件奖励

	// 标记为已领取
	data.SetAttachGettedDB(uid, mailId)

	return nil, nil
}


// 启动邮件列表变化检测协程
func  runCheckMailChange() error{

	// 1分钟更新一次邮件列表
	for {
		time.Sleep(time.Minute)

		// 判断当时是否是主节点
		isMasterNode = structs.GetGlobalExposer().ConsulReq.IsMasterNode()
		getDataFromDB()
	}
	return nil
}

// 从DB获取邮件列表
func getDataFromDB() error {

	list, err := data.LoadMailListFromDB()
	if err != nil {
		logrus.Errorln("load email list from db err:", list)
		return err
	}
	logrus.Debugln("email list:", list)
	mailList = list
	// 检测邮件状态
	checkMailStatus(mailList)

	// 主节点每日半夜4-6点，清理过期邮件
	clearExpiredEmail()
	return err
}

// 主节点每日半夜4-6点，清理过期邮件
var thisDay = 0
func clearExpiredEmail() {
	if !isMasterNode {
		// 非主节点，直接返回
		return
	}

	now := time.Now()
	// 每日只执行1次
	if thisDay == now.YearDay() {
		return
	}
	if now.Hour() >= clearBeginHour && now.Hour() < clearEndHour {
		data.ClearExpiredEmailFromDB()
		data.ClearExpiredUserEmailFromDB()
		thisDay = now.YearDay()
	}
}

// 检测邮件状态是否变化
func checkMailStatus(mailList map[uint64]*define.MailInfo) error {

	curDate := time.Now().Format("2006-01-02 15:04:05")

	bUpdate := false

	for _, mail := range mailList {
		if mail.State == define.StateChecked {
			// 检测是否开始
			if curDate >= mail.StartTime {
				mail.State = define.StateSending
				bUpdate = true

				// 主节点负责保存邮件状态到DB
				if isMasterNode {
					// 保存邮件状态变化到DB
					data.SetEmailStateToDB(mail.Id, mail.State)
				}

			}
		} else if mail.State == define.StateSending {
			// 检测是否结束
			if mail.IsUseEndTime && curDate >= mail.EndTime {
				mail.State = define.StateSended
				bUpdate = true
				// 主节点负责保存邮件状态到DB
				if isMasterNode {
					// 保存邮件状态变化到DB
					data.SetEmailStateToDB(mail.Id, mail.State)
				}
			}
		} else if mail.State == define.StateSended {
			// 检测是否达到删除时间
			if mail.IsUseDelTime && curDate >= mail.DelTime {
				mail.State = define.StateDelete
				bUpdate = true
				// 主节点负责保存邮件状态到DB
				if isMasterNode {
					// 保存邮件状态变化到DB
					data.SetEmailStateToDB(mail.Id, mail.State)
				}
			}
		}
	}

	bUpdate = true
	// 更新发送列表provSendList
	if bUpdate {
		// 将发送中和发送截至的加入到指定列表中
		myList := make(map[int64][]*define.MailInfo)
		for _, mail := range mailList {
			if mail.State == define.StateSending || mail.State == define.StateSended {

				for _, dest := range mail.DestList {

					myList[dest.Prov] = append(myList[dest.Prov], mail)
				}
			}
		}
		provSendList = myList
	}

	return nil
}

// 调用hall接口获取用户信息
// 返回:渠道ID，省ID，城市ID
func getUserInfo(uid uint64) (int64, int64, int64, bool) {
	return 0, 1, 0, true
	info, err := hallclient.GetPlayerInfo(uid)
	if err != nil {
		return 0, 0, 0, false
	}
	if info == nil {
		return 0, 0, 0, false
	}

	return int64(info.ChannelId), int64(info.ProvinceId), int64(info.CityId), true
}

// 检测是否符合省包和渠道ID
func checkMailProvChannel(uid uint64,mail *define.MailInfo, channel int64, prov int64) bool {
	isOk := false
	for _, dest := range mail.DestList {
		if dest.Prov != 0 && prov != dest.Prov {
			continue
		}

		if dest.Channel != 0 && channel != dest.Channel {
			continue
		}
		if dest.SendType == define.SendAll {
			// 发送给所有人
			isOk = true
		} else {
			// 发送给指定玩家列表
			for _, id := range dest.PlayerList {
				if id == uid {
					isOk = true
					break
				}
			}
		}
		if isOk {
			break
		}

	}
	return isOk
}
