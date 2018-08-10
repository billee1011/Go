package bean

import (
	"steve/datareport/fixed"
	"strconv"
	"time"
)

type LogBean struct {
	LogType  int32
	Province int32
	City     int32
	Channel  int32
	PlayerId uint64
	Value    string
}

func CreateLogBean(logType int32, province int32, city int32, channel int32, playerId uint64, value string) *LogBean {
	return &LogBean{
		LogType:  logType,
		Province: province,
		City:     city,
		Channel:  channel,
		PlayerId: playerId,
		Value:    value,
	}
}

//把日志对象转换到上报字符串
func (log *LogBean) ToReportFormat() string {
	result := ""
	switch int(log.LogType) {
	case int(fixed.LOG_TYPE_GAME_PERSON_NUM):
		result += fixed.LOG_TABLE_NAME_ONLINE + "|" + log.getHead() + "|" + log.Value + "|"
	case int(fixed.LOG_TYPE_REG),
		int(fixed.LOG_TYPE_ACT),
		int(fixed.LOG_TYPE_GAM),
		int(fixed.LOG_TYPE_GOLD_ADD),
		int(fixed.LOG_TYPE_GOLD_REMV),
		int(fixed.LOG_TYPE_YB_ADD),
		int(fixed.LOG_TYPE_YB_REMV),
		int(fixed.LOG_TYPE_CARD_ADD),
		int(fixed.LOG_TYPE_CARD_REMV):
		result = log.convertToLogRealTimeReport()
	}
	result += time.Now().Format("2006-01-02 15:04:05") //最后添加上报时间字段
	result += "\n"
	return result
}

//转换实时简报表
func (log *LogBean) convertToLogRealTimeReport() string {
	v := fixed.LOG_TABLE_NAME_REALTIME + "|" + log.getHead() + "|" + strconv.FormatUint(log.PlayerId, 10) + "|"
	for i := 0; i < fixed.MAX_PARAM-fixed.HEAD_PARAM; i++ {
		if i == int(log.LogType) {
			v += log.Value
		}
		v += "|"
	}
	return v
}

func (log *LogBean) getHead() string {
	return strconv.Itoa(int(log.Province)) + "|" + strconv.Itoa(int(log.City)) + "|" + strconv.Itoa(int(log.Channel))
}
