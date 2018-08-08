package logic

import (
	"steve/client_pb/mailserver"
	"steve/mailserver/data"
	"github.com/Sirupsen/logrus"
	"steve/mailserver/define"
	"steve/external/hallclient"
	"errors"
	"time"
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

var mailList map[int64]*define.MailInfo

var channelSendList map[int64][]*define.MailInfo
var provSendList map[int64][]*define.MailInfo

func Init() error {
	err := getDataFromDB()
	if err != nil {
		return err
	}
	return nil
}

// 获取未读消息总数
func GetGetUnReadSum(uid uint64) (int32, error) {
	// 获取玩家渠道ID
	channel, prov, _, ok := getUserInfo(uid)
	if !ok {
		return 0, errors.New("获取玩家渠道ID失败")
	}

	// 从DB获取玩家的已读邮件列表
	readList, err:= data.GetUserMailFromDB(uid)
	if err != nil {
		return 0,  errors.New("从DB获取玩家已读邮件列表失败")
	}

	sum := int32(0)
	// 获取玩家所属渠道的邮件
	if channel > 0 {
		list  := channelSendList[channel]
		for _, mail := range list {
			if one, ok := readList[mail.Id]; !ok  || !one.IsRead{
				sum++
			}
		}

	}
	// 获取玩家所属省包的邮件
	if prov > 0 {
		list  := provSendList[prov]
		for _, mail := range list {
			if one, ok := readList[mail.Id]; !ok  || !one.IsRead{
				sum++
			}
		}
	}

	return sum, nil
}

// 获取邮件消息列表
func GetMailList(uid uint64) ([]*mailserver.MailTitle, error) {

	return nil, nil
}

// 获取指定邮件详情
func GetMailDetail(uid uint64, mail uint64) (*mailserver.MailDetail, error) {

	return nil, nil
}

// 删除邮件
func DelMail(uid uint64, mail uint64) error {

	return nil
}

// 领取附件奖励请求
func AwardAttach(uid uint64, mail uint64) (string, error) {

	return "", nil
}

// 从DB获取邮件列表
func getDataFromDB() error {
	mailList , err := data.LoadMailListFromDB()
	if err != nil {
		logrus.Errorln("load email list from db err:", mailList)
		return err
	}
	logrus.Debugln("email list:" , mailList)
	// 检测邮件状态
	checkMailStatus()
	return err
}

// 检测邮件状态是否变化
func checkMailStatus() error {

	curDate := time.Now().Format("2006-01-02 00:00:00")

	bUpdate := false

	for _, mail := range  mailList {
		if mail.State == define.StateChecked {
			// 检测是否开始
			if curDate >= mail.StartTime {
				mail.State = define.StateSending
				bUpdate = true
			}
		} else if mail.State == define.StateSending {
			// 检测是否结束
			if mail.IsUseEndTime &&  curDate >= mail.EndTime {
				mail.State = define.StateSended
				bUpdate = true
			}
		} else if  mail.State == define.StateSended {
			// 检测是否达到删除时间
			if mail.IsUseDelTime &&  curDate >= mail.DelTime {
				mail.State = define.StateDelete
				bUpdate = true
			}
		}
	}

	// 更新发送列表provSendList
	if bUpdate {
		// 将发送中和发送截至的加入到指定列表中
		myList := make(map[int64][]*define.MailInfo)
		for _, mail := range  mailList {
			if mail.State == define.StateSending || mail.State == define.StateSended {
				//myList[mail.]
				_ = myList
			}
		}
	}

	return nil
}

// 调用hall接口获取用户信息
// 返回:渠道ID，省ID，城市ID
func getUserInfo(uid uint64) (int64, int64, int64, bool) {
	//return 1, 0, 0
	info, err := hallclient.GetPlayerInfo(uid)
	if err != nil {
		return 0, 0, 0, false
	}
	if info == nil {
		return 0, 0, 0, false
	}

	return int64(info.ChannelId), int64(info.ProvinceId), int64(info.CityId), true
}



