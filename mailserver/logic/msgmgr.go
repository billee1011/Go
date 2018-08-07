package logic

import (
	"steve/mailserver/define"
	"time"
	"fmt"
	"steve/mailserver/data"
	"errors"
	"github.com/Sirupsen/logrus"
	"steve/external/hallclient"
	"steve/external/gateclient"
	"steve/client_pb/msgid"
	"steve/client_pb/msgserver"
)

/*
  功能： 跑马灯管理:
		1.获取玩家跑马灯列表
		2.下发跑马灯变化通知给所有玩家
		3.每隔1分钟从DB获取跑马灯列表.
		4.每隔1分钟检测未开始的跑马灯是否开始，并且检测已开始的跑马灯是否结束.
  作者： SkyWang
  日期： 2018-7-24
*/

var myMgr MsgMgr

func GetMsgMgr() *MsgMgr {
	return &myMgr
}

type MsgMgr struct {
	horseList   map[int64]*define.HorseRace // 跑马灯列表
	cityList    map[int64]*define.HorseRace // 城市级别的跑马灯列表
	provList    map[int64]*define.HorseRace // 省份级别的跑马灯列表
	channelList map[int64]*define.HorseRace // 渠道级别的跑马灯列表
}

func (gm *MsgMgr) Init() error {
	// 初始化跑马灯数据
	err := gm.getHorseRaceFromDB()
	if err != nil {
		return err
	}
	// 启动跑马灯是否开始检测协程
	go gm.runCheckHorseChange()

	return nil
}

// 获取跑马灯列表
func (gm *MsgMgr) GetHorseRace(uid uint64) ([]string, int32, int32, error) {
	// 1.先调用hall的接口获取玩家的渠道ID，省份ID和市ID
	// 2.从redis获取目前开启的跑马灯(达到开始时间和结束时间之间)。
	// 3.按照从市ID-》省ID-》渠道ID的顺序获取跑马灯，如果获取到马上返回。
	// 4.如果开启不使用上级配置的开关，则不向上检索.
	channel, prov, city, ok := gm.getUserInfo(uid)
	logrus.Debugf("getUserInfo: uid=%d,channel=%d,prov=%d,city=%d", uid, channel, prov, city)

	if !ok {
		return nil, 0, 0, errors.New("获取玩家渠道ID失败")
	}

	if city != 0 {
		// 先读取城市级别的跑马灯
		str, tick, sleep, ok := gm.getLevelHorseRace(city, gm.cityList)
		if ok {
			return str, tick, sleep, nil
		}
	}

	if prov != 0 {
		// 再读取省份级别的跑马灯
		str, tick, sleep, ok := gm.getLevelHorseRace(prov, gm.provList)
		if ok {
			return str, tick, sleep, nil
		}
	}

	if channel >= 0 {
		// 最后读取渠道级别的跑马灯
		str, tick, sleep, ok := gm.getLevelHorseRace(channel, gm.channelList)
		if ok {
			return str, tick, sleep, nil
		}
	}
	return nil, 0, 0, nil
}


// 从DB获取最新的跑马灯列表
func (gm *MsgMgr) getHorseRaceFromDB() error {

	gm.horseList, _ = data.LoadHorseFromDB()

	if gm.horseList == nil {
		myMgr.horseList = make(map[int64]*define.HorseRace)
		return  errors.New("获取跑马灯配置失败")
	}

	gm.cityList = make(map[int64]*define.HorseRace)
	gm.provList = make(map[int64]*define.HorseRace)
	gm.channelList = make(map[int64]*define.HorseRace)
	for _, horse := range  gm.horseList {
		if horse.Channel < 0 {
			continue
		}
		if horse.Prov == 0 && horse.City == 0 {
			gm.channelList[horse.Channel] = horse
		} else if horse.Prov > 0 && horse.City == 0 {
			gm.provList[horse.Prov] = horse
		} else if horse.City > 0 {
			gm.cityList[horse.City] = horse
		}
	}

	//gm.testHorseJson()

	return nil
}
func (gm *MsgMgr) testHorseJson() {
	testJson := &define.HorseRaceJson{}
	testJson.SleepTime = 300
	testJson.TickTime = 5
	for i := 0; i < 4; i++ {
		hc := new(define.HorseContentJson)
		hc.PlayType = 1
		hc.WeekDate = []int8{int8(i+1),6}
		hc.BeginTime = "12:00"
		hc.EndTime = "18:00"
		hc.Content = fmt.Sprintf("循环播放:跑马灯%d", i+1)
		testJson.Horse = append(testJson.Horse, hc)
	}
	hc := new(define.HorseContentJson)
	hc.PlayType = 2
	hc.WeekDate = []int8{}
	hc.BeginDate = "2018-07-30"
	hc.EndDate = "2018-08-15"
	hc.BeginTime = "12:00"
	hc.EndTime = "18:00"
	hc.Content = "指定时间播放"

	testJson.Horse = append(testJson.Horse, hc)

	str, _ := data.MarshalHorseJson(testJson)

	logrus.Debugln(str)
}


// 启动跑马灯状态变化检测协程
func (gm *MsgMgr) runCheckHorseChange() error{

	// 1分钟检测一次跑马灯状态
	for {
		bUpdate := false
		for _, horse := range gm.cityList {
			if gm.checkHorseChanged(horse) {
				bUpdate = true
			}
		}
		time.Sleep(time.Millisecond * 20)
		for _, horse := range gm.provList {
			if gm.checkHorseChanged(horse) {
				bUpdate = true
			}
		}
		time.Sleep(time.Millisecond * 20)
		for _, horse := range gm.channelList {
			if gm.checkHorseChanged(horse) {
				bUpdate = true
			}
		}

		// 如果有变化发送通知消息
		if bUpdate {
			gm.sendHorseRaceChangedNtf()
		}
		//gm.testGetHorseRace()
		time.Sleep(time.Minute)

	}

	return nil
}

// 检测跑马灯跑马灯是否达到了开始时间，或结束时间
func (gm *MsgMgr) checkHorseChanged(horse *define.HorseRace) bool {

	if horse.IsUse == 0 {
		// 不启用配置
		return false
	}

	isUpdate := false
	for _, hc := range horse.Content {
		if gm.checkHorseBegin(hc) {
			if gm.SetUpdateTag(hc, define.StateStart) {
				isUpdate = true

			}
		} else {
			if gm.SetUpdateTag(hc, define.StateStop) {
				isUpdate = true

			}
		}
	}

	return isUpdate
}

// 发送跑马灯变化通知
func (gm *MsgMgr) sendHorseRaceChangedNtf() error {
	req := &msgserver.MsgSvrHorseRaceChangeNtf{}
	channel := int32(0)
	req.Channel = &channel
	logrus.Debugln("horse race status changed...")
	gateclient.NsqBroadcastAllMsg(uint32(msgid.MsgID_MSGSVR_HORSE_RACE_UPDATE_NTF), req)
	return nil
}

// 调用hall接口获取用户信息
// 返回:渠道ID，省ID，城市ID
func (gm *MsgMgr) getUserInfo(uid uint64) (int64, int64, int64, bool) {

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

func (gm *MsgMgr)SetUpdateTag(hc *define.HorseContent, status int8) bool {
	if hc.CheckStatus == status {
		return false
	}
	hc.CheckStatus = status
	return true
}

// 检测指定跑马灯内容是否达到开始时间
func (gm *MsgMgr) checkHorseBegin(hc *define.HorseContent) bool {
	tm := time.Now()
	wd := tm.Weekday()
	// 1.判断类型
	if hc.PlayType == define.PlayLoop {
		//  重复
		// 周判断
		if !hc.WeekDate[int8(wd)] {
			return false
		}

	} else if  hc.PlayType == define.PlayFix {
		// 指定
		// 日期判断
		curDate := tm.Format("2006-01-02")

		if curDate < hc.BeginDate {
			return false
		}

		if curDate > hc.EndDate {
			return false
		}
	} else {
		return false
	}

	// 判断是否达到时间段
	curTime := fmt.Sprintf("%02d:%02d", tm.Hour(), tm.Minute())

	if curTime < hc.BeginTime  {
		return false
	}
	if curTime >= hc.EndTime {
		return false
	}

	return true
}

// 得到玩家指定层级的跑马灯
// 返回：跑马灯内容, 间隔, 一轮休眠时间
func (gm *MsgMgr) getLevelHorseRace(level int64, m map[int64]*define.HorseRace) ([]string, int32, int32, bool) {

	if m == nil {
		return nil, 0, 0, false
	}
	horse, ok := m[level]
	if !ok {
		return nil, 0, 0, false
	}

	if horse == nil {
		return nil, 0, 0, false
	}


	if horse.IsUse == 0 {
		// 不启用配置
		if horse.IsUseParent == 0 {
			// 不启用上级配置
			return nil, 0, 0, true
		} else {
			return nil, 0, 0, false
		}
	}

	list := make([]string, 0, len(horse.Content))
	for _, hc := range horse.Content {
		if hc.CheckStatus == define.StateStart {
			list = append(list, hc.Content)
		}
	}

	return list, horse.TickTime, horse.SleepTime, true
}


func (gm *MsgMgr) testGetHorseRace() {
	gm.GetHorseRace(1001)
	//gm.GetHorseRace(1002)

}