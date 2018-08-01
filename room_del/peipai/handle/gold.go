package handle

import (
	"fmt"
	"net/http"
	"steve/room/interfaces/global"
	"strconv"
)

// SetGoldHandle 设置玩家金币
func SetGoldHandle(resp http.ResponseWriter, req *http.Request) {

	// 玩家ID
	playerID, err := strconv.ParseUint(req.FormValue(PlayerIDKey), 10, 64)
	response := "OK"
	defer resp.Write([]byte(response))

	if err != nil {
		response = "player_id 数据错误"
		return
	}

	// 金币数
	gold, err := strconv.ParseUint(req.FormValue(GoldKey), 10, 64)
	if err != nil {
		response = "gold 数据错误"
		return
	}

	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	if player == nil {
		response = "player_id 不存在"
		return
	}
	respMSG(resp, fmt.Sprintf("配置玩家金币数成功,当前为:\n玩家ID[%v] -- 金币[%v]\n", playerID, gold), 200)
	player.SetCoin(gold)
}
