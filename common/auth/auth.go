package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// GenerateAuthToken 生成认证 token
func GenerateAuthToken(playerID uint64, gateIP string, gatePort int, expire int64, key string) string {
	data := fmt.Sprintf("%v%s%d%v%s", playerID, gateIP, gatePort, expire, key)
	result := md5.Sum([]byte(data))
	return hex.EncodeToString(result[:])
}
