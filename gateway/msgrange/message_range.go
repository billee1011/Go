package msgrange

import (
	"fmt"
	"steve/structs/common"
)

type messageRange struct {
	minMsgID uint32
	maxMsgID uint32
}

var gServerMessageRange = map[string]messageRange{
	common.RoomServiceName: messageRange{
		minMsgID: 0x10000,
		maxMsgID: 0x1ffff,
	},
	common.GateServiceName: messageRange{
		minMsgID: 0x1001,
		maxMsgID: 0x1fff,
	},
	common.MatchServiceName: messageRange{
		minMsgID: 0x2001,
		maxMsgID: 0x2fff,
	},
	common.LoginServiceName: messageRange{
		minMsgID: 0x0001,
		maxMsgID: 0x0FFF,
	},
}

// GetMessageServer 获取消息处理服务名字
// 返回值为空表示无对应的服务
func GetMessageServer(msgID uint32) string {
	for serverName, serverRange := range gServerMessageRange {
		if msgID >= serverRange.minMsgID && msgID <= serverRange.maxMsgID {
			return serverName
		}
	}
	return ""
}

// 校验消息 ID 段是否有重复
func checkOverlap() {
	serverRanges := map[string]messageRange{}
	for serverName, serverRange := range gServerMessageRange {
		for checkServerName, checkServerRange := range serverRanges {
			if checkServerRange.minMsgID >= serverRange.minMsgID &&
				checkServerRange.minMsgID <= serverRange.maxMsgID {
				panic(fmt.Sprintf("%s 与 %s 的消息区段有重复", checkServerName, serverName))
			} else if checkServerRange.maxMsgID >= serverRange.minMsgID &&
				checkServerRange.maxMsgID <= serverRange.maxMsgID {
				panic(fmt.Sprintf("%s 与 %s 的消息区段有重复", checkServerName, serverName))
			} else {
				serverRanges[serverName] = serverRange
			}
		}
	}
}

func init() {
	checkOverlap()
}
