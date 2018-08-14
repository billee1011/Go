package matchv3

import (
	"fmt"
	"net"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/external/gateclient"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

// IPUInt32ToString 整形IP地址转为字符串型IP
func IPUInt32ToString(intIP uint32) string {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

// IPStringToUInt32 字符串型IP转为uint32型
func IPStringToUInt32(ipStr string) uint32 {
	bits := strings.Split(ipStr, ".")

	if len(bits) != 4 {
		logrus.Errorln("IPStringToUInt32() 参数错误，ipStr = ", ipStr)
		return 0
	}

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

// GetServerAddr 获取本match服的IP地址
func GetServerAddr() string {
	localIP := viper.GetString("rpc_addr")
	localPort := viper.GetInt("rpc_port")
	localAddr := fmt.Sprintf("%s:%d", localIP, localPort)

	return localAddr
}

// NotifyCancelMatch 通知指定的玩家:取消匹配,目前发送取消匹配成功的回复
func NotifyCancelMatch(playersID []uint64) {

	// 取消匹配成功的回复
	response := match.CancelMatchRsp{
		ErrCode: proto.Int32(int32(match.MatchError_EC_SUCCESS)),
		ErrDesc: proto.String("成功"),
	}

	// 所有的玩家
	for _, playerID := range playersID {
		gateclient.SendPackageByPlayerID(playerID, uint32(msgid.MsgID_CANCEL_MATCH_RSP), &response)
	}
}

// match有关消息
const (
	ClearMatch = iota // 清空所有的匹配
)
