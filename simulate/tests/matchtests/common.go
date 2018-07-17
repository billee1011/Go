package matchtests

import (
	"fmt"
	"net/http"
	"steve/simulate/config"
	"time"
)

// modifyRobotJoinTime 修改机器人加入匹配的时间
func modifyRobotJoinTime(d time.Duration) error {
	matchHTTPAddr := config.GetMatchHTTPAddr()
	url := fmt.Sprintf("%s/set_robot_join_time?robot_join_time=%d", matchHTTPAddr, d/time.Millisecond)
	_, err := http.DefaultClient.Get(url)
	if err != nil {
		return fmt.Errorf("请求修改机器人加入匹配时间失败: %v", err)
	}
	return nil
}
