package handle

import (
	"net/http"
	"steve/room/interfaces/global"
	"strconv"
)

// SetGoldHandle set player gold
func SetGoldHandle(resp http.ResponseWriter, req *http.Request) {
	playerID, err := strconv.ParseUint(req.FormValue("player_id"), 10, 64)
	response := "OK"
	defer resp.Write([]byte(response))

	if err != nil {
		response = "player_id 数据错误"
		return
	}
	gold, err := strconv.ParseUint(req.FormValue("gold"), 10, 64)
	if err != nil {
		response = "gold 数据错误"
		return
	}

	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	player.SetCoin(gold)
}
