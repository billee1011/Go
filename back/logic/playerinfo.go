package logic

import (
	"fmt"
	"math"
	"steve/back/data"
	"steve/entity/gamelog"
	"sync"

	"github.com/Sirupsen/logrus"
)

func updatePlayerInfo(detailInfo gamelog.TGameDetail) error {
	playerGame, err := data.GetTPlayerGame(detailInfo.Gameid, detailInfo.Playerid)
	if err != nil {
		return err
	}
	//总局数+1
	playerGame.Totalbureau++
	key := fmt.Sprintf("win_stream:%v@%v", playerGame.Playerid, playerGame.Gameid)
	// contiWin := 0
	winStream, _ := data.GetPlayerMaxwinningstream(key)
	// if v, ok := MaxBurea.LoadOrStore(key, 0); ok {
	// 	winStream = v.(int)
	// }
	if detailInfo.Amount > 0 {
		//胜局+1
		playerGame.Winningburea++
		//连胜+1
		winStream++
	} else {
		//输了，连胜终结
		winStream = 0
	}
	//储存最新连胜
	// MaxBurea.Store(key, winStream)
	if err := data.SetPlayerMaxwinningstream(key, winStream); err != nil {
		logrus.Errorf("failed set maxSream to redis,err:%v", err)
	}
	if winStream > playerGame.Maxwinningstream {
		// 更新连胜
		playerGame.Maxwinningstream = winStream
	}
	// 更新胜率
	playerGame.Winningrate = int(math.Floor((float64(playerGame.Winningburea)/float64(playerGame.Totalbureau))*100 + 0.5))

	if err := data.UpdateTPlayerGame(playerGame); err != nil {
		return err
	}
	if err := data.UpdatePlayerGameToredis(playerGame); err != nil {
		return err
	}
	return nil
}

// MaxBurea 暂时本地储存最大连胜
var MaxBurea = new(sync.Map)
