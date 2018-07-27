package cheater

import (
	"fmt"
	"net/http"
	"steve/simulate/config"
)

// SetPlayerCoin 设置玩家金豆数
func SetPlayerCoin(playerID uint64, coin uint64) error {
	url := fmt.Sprintf("%s/setgold?player_id=%v&gold=%v", config.GetPeipaiURL(), playerID, coin)
	if _, err := http.DefaultClient.Get(url); err != nil {
		return fmt.Errorf("访问设置金币服务失败:%v", err)
	}
	return nil
}
