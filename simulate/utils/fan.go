package utils

import (
	"fmt"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// erRenHuChiFanMap 番型对应互斥的番型  room.FanType_FT_
var erRenHuChiFanMap = map[room.FanType][]room.FanType{
	room.FanType_FT_DASIXI: []room.FanType{room.FanType_FT_PENGPENGHU, room.FanType_FT_QUANFENGKE, room.FanType_FT_MENFENGKE, room.FanType_FT_DASANFENG, room.FanType_FT_XIAOSANFENG, room.FanType_FT_SIZIKE},
	room.FanType_FT_DASANYUAN: []room.FanType{
		room.FanType_FT_SHUANGJIANKE, room.FanType_FT_JIANKE,
	},
}

// erRenFanMulMap 番型对应的倍数
var erRenFanMulMap = map[room.FanType]int32{
	room.FanType_FT_TIANHU:          88,
	room.FanType_FT_DIHU:            88,
	room.FanType_FT_DASIXI:          88,
	room.FanType_FT_DASANYUAN:       88,
	room.FanType_FT_DAQIXING:        88,
	room.FanType_FT_LIANQIDUI:       88,
	room.FanType_FT_RENHU:           64,
	room.FanType_FT_ZIYISE:          64,
	room.FanType_FT_XIAOSANYUAN:     64,
	room.FanType_FT_SIANKE:          64,
	room.FanType_FT_QINGYISE:        16,
	room.FanType_FT_HUNYISE:         6,
	room.FanType_FT_GANGSHANGKAIHUA: 2,
	room.FanType_FT_ANGANG:          2,
	room.FanType_FT_JIANKE:          2,
	room.FanType_FT_BUQIUREN:        4,
	room.FanType_FT_WUHUA:           4,
	room.FanType_FT_QUANDAIYAO:      4,
	room.FanType_FT_DANDIAOJIANG:    1,
	room.FanType_FT_BIANZHANG:       1,
	room.FanType_FT_KANZHANG:        1,
	room.FanType_FT_ZIMO:            1,
}

//GetHuChiValueByGameID 根据游戏ID获取互斥番型数组
func GetHuChiValueByGameID(gameID room.GameId, currFan room.FanType) ([]room.FanType, bool) {
	switch gameID {
	case room.GameId_GAMEID_ERRENMJ:
		fans, isExist := erRenHuChiFanMap[currFan]
		return fans, isExist
	default:
		return []room.FanType{}, false
	}
}

//GetFanMulByGameID 根据游戏ID获取互斥番型倍数
func GetFanMulByGameID(gameID room.GameId, currFan room.FanType) int32 {
	switch gameID {
	case room.GameId_GAMEID_ERRENMJ:
		return erRenFanMulMap[currFan]
	default:
		return 0
	}
}

// CheckFanSettle 检测番型结算 winSeat赢家座位，winScore 赢家总赢分，currFan 指定确认都番型
func CheckFanSettle(t *testing.T, deskData *DeskData, gameID room.GameId, winSeat int, winScore int64, currFan room.FanType) {
	winPlayer := GetDeskPlayerBySeat(winSeat, deskData)
	expector, _ := winPlayer.Expectors[msgid.MsgID_ROOM_ROUND_SETTLE]
	ntf := &room.RoomBalanceInfoRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf))
	for _, info := range ntf.BillPlayersInfo {
		fmt.Println(info.GetFan())
		if winPlayer.Player.GetID() == info.GetPid() {
			assert.Equal(t, winScore, info.GetScore())
		} else {
			assert.Equal(t, -winScore, info.GetScore())
		}
		assert.True(t, IsExistAssignFan(currFan, info.GetFan()))
		flag, str := IsExistHuChiFan(gameID, currFan, info.GetFan())
		assert.Falsef(t, flag, str)
		assert.True(t, IsAssignFanMulRight(gameID, currFan, info.GetFan()))
	}
}

// IsExistAssignFan 判断指定番型是否存在
func IsExistAssignFan(currFan room.FanType, Fans []*room.Fan) bool {
	for _, fan := range Fans {
		if fan.GetName() == currFan {
			return true
		}
	}
	return false
}

//IsExistHuChiFan 是否存在互斥的牌
func IsExistHuChiFan(gameID room.GameId, currFan room.FanType, Fans []*room.Fan) (bool, string) {
	fanTyps, isExist := GetHuChiValueByGameID(gameID, currFan)
	if isExist {
		for _, fanTyp := range fanTyps {
			if IsExistAssignFan(fanTyp, Fans) {
				return true, fmt.Sprintf("存在互斥番型")
			}
		}
		return false, ""
	}
	return false, fmt.Sprintf("当前番型不存在：%v", currFan)
}

//IsAssignFanMulRight 判断指定番型倍数是否正确
func IsAssignFanMulRight(gameID room.GameId, currFan room.FanType, Fans []*room.Fan) bool {
	for _, fan := range Fans {
		if fan.GetName() == currFan && fan.GetValue() == erRenFanMulMap[currFan] {
			return true
		}
		if fan.GetName() == currFan {
			fmt.Println("fanName()", fan.GetName())
		}
	}
	return false
}
