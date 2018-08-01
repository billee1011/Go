package logic

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

}

func init() {

}

// 从DB获取最新的跑马灯列表
func (gm *MsgMgr) getHorseRaceFromDB()  error {

	return nil
}

// 检测跑马灯跑马灯是否达到了开始时间，或结束时间
func (gm *MsgMgr) checkHorseRace()  error {

	return nil
}

// 发送跑马灯变化通知
func (gm *MsgMgr) sendHorseRaceChangedNtf()  error {

	return nil
}

// 获取跑马灯列表
func (gm *MsgMgr) GetHorseRace(uid uint64) ([]string, error) {
	// 1.先调用hall的接口获取玩家的渠道ID，省份ID和市ID
	// 2.从redis获取目前开启的跑马灯(达到开始时间和结束时间之间)。
	// 3.按照从市ID-》省ID-》渠道ID的顺序获取跑马灯，如果获取到马上返回。
	// 4.如果开启不使用上级配置的开关，则不向上检索.


	return nil, nil
}



