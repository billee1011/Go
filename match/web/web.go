package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// web配置信息
var configs = struct {
	robotJoinTime         time.Duration // 机器人加入匹配的时间
	continueDismissTime   time.Duration // 续局牌桌解散时间
	continueRobotTime     time.Duration // 续局牌桌机器人决策时间
	robotContinueRateWin  float32       // 机器人胜利时续局概率
	robotContinueRateLoss float32       // 机器人失败时续局概率
	sameDeskLimitTime     int64         // 同桌限制时间，单位：秒，超过这个时间，匹配时不再限制同桌
	defaulWintRate        int32         // 玩家默认胜率，玩家游戏局数低于最低游戏局数时，采用此值
	minGameTimes          uint32        // 最低游戏局数，玩家局数低于此值，采用此值
	winRateCompuBase      float32       // 计算公式的基础胜率(百分比，例如：0.02表2%)
	goldCompuBase         float32       // 计算公式的基础金币(百分比，例如：0.02表2%)
	maxCompuValidTime     uint32        // 计算公式的最大有效时间(单位：秒，超过此值时匹配正无穷)
}{
	robotJoinTime:         20 * time.Second,
	continueDismissTime:   20 * time.Second,
	continueRobotTime:     3 * time.Second,
	robotContinueRateWin:  0.9,
	robotContinueRateLoss: 0.7,
	sameDeskLimitTime:     60,
	defaulWintRate:        50,
	minGameTimes:          50,
	winRateCompuBase:      0.02,
	goldCompuBase:         0.2,
	maxCompuValidTime:     15,
}

// GetRobotJoinTime 获取机器人加入匹配的时间
// 超过这个时间，就要加入机器人
func GetRobotJoinTime() time.Duration {
	return configs.robotJoinTime
}

// GetContinueDismissTime 获取续局牌桌解散时间
// 超过这个时间，等待中的续局牌桌需要解散
func GetContinueDismissTime() time.Duration {
	return configs.continueDismissTime
}

// GetContinueRobotTime 获取续局牌桌机器人决策时间
//
func GetContinueRobotTime() time.Duration {
	return configs.continueRobotTime
}

// GetRobotContinueRate 获取机器人续局概率
func GetRobotContinueRate(winner bool) float32 {
	if winner {
		return configs.robotContinueRateWin
	}
	return configs.robotContinueRateLoss
}

// GetSameDeskLimitTime 获取同桌限制时间
func GetSameDeskLimitTime() int64 {
	return configs.sameDeskLimitTime
}

// GetDefaultWinRate 获取玩家默认胜率
func GetDefaultWinRate() int32 {
	return configs.defaulWintRate
}

// GetMinGameTimes 获取最低游戏局数
func GetMinGameTimes() uint32 {
	return configs.minGameTimes
}

// GetWinRateCompuBase 获取计算公式的基础胜率
func GetWinRateCompuBase() float32 {
	return configs.winRateCompuBase
}

// GetGoldCompuBase 获取计算公式的基础金币
func GetGoldCompuBase() float32 {
	return configs.goldCompuBase
}

// GetMaxCompuValidTime 获取计算公式的最大有效时间
func GetMaxCompuValidTime() uint32 {
	return configs.maxCompuValidTime
}

func handleChangeDurationVal(d *time.Duration, min, max time.Duration, w http.ResponseWriter, r *http.Request, formField string) {
	result := "OK"
	defer func() {
		w.Write([]byte(result))
	}()

	val, err := strconv.Atoi(r.FormValue(formField))
	if err != nil {
		result = fmt.Sprintf("参数[%s]错误", formField)
		return
	}
	requestD := time.Duration(val) * time.Millisecond
	if requestD < min {
		minString := min.String()
		result = fmt.Sprintf("[%s]时间小于[%s]， 更新为[%s]", formField, minString, minString)
		requestD = min
	}
	if requestD > max {
		maxString := max.String()
		result = fmt.Sprintf("[%s]时间大于[%s]， 更新为[%s]", formField, maxString, maxString)
		requestD = max
	}
	*d = requestD
	return
}

// handleChangeRobotJoinTime 修改机器人加入匹配的时间
func handleChangeRobotJoinTime(w http.ResponseWriter, r *http.Request) {
	handleChangeDurationVal(&configs.robotJoinTime, 1*time.Millisecond, 1*time.Hour, w, r, "robot_join_time")
}

// handleChangeContinueDismissTime 修改续局牌桌解散时间
func handleChangeContinueDismissTime(w http.ResponseWriter, r *http.Request) {
	handleChangeDurationVal(&configs.continueDismissTime, 1*time.Millisecond, 1*time.Hour, w, r, "continue_dismiss_time")
}

// handleChangeContinueRobotTime 修改机器人续局决策时间
func handleChangeContinueRobotTime(w http.ResponseWriter, r *http.Request) {
	handleChangeDurationVal(&configs.continueRobotTime, 1*time.Millisecond, 1*time.Hour, w, r, "continue_robot_time")
}

// handleChangeRobotContinueRate 修改机器人续局概率
func handleChangeRobotContinueRate(w http.ResponseWriter, r *http.Request) {
	var result string
	defer func() {
		w.Write([]byte(result))
	}()

	lossRateString := r.FormValue("loss_rate")
	if lossRateString != "" {
		val, err := strconv.ParseFloat(lossRateString, 32)
		if err != nil {
			result = fmt.Sprintf("%sloss_rate 格式错误\r\n", result)
			return
		}
		configs.robotContinueRateLoss = float32(val)
		result = fmt.Sprintf("%s修改机器人失败时续局概率为 %.2f \r\n", result, val)
	}

	winRateString := r.FormValue("win_rate")
	if winRateString != "" {
		val, err := strconv.ParseFloat(winRateString, 32)
		if err != nil {
			result = fmt.Sprintf("%swin_rate 格式错误\r\n", result)
			return
		}
		configs.robotContinueRateWin = float32(val)
		result = fmt.Sprintf("%s修改机器人胜利时续局概率为 %.2f \r\n", result, val)
	}
}

func init() {
	http.HandleFunc("/set_robot_join_time", handleChangeRobotJoinTime)
	http.HandleFunc("/set_continue_dismiss_time", handleChangeContinueDismissTime)
	http.HandleFunc("/set_continue_robot_time", handleChangeContinueRobotTime)
	http.HandleFunc("/set_robot_continue_rate", handleChangeRobotContinueRate)
}
