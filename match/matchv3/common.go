package matchv3

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
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
