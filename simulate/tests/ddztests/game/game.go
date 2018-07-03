package game

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
)

// StartGameParams 启动游戏的参数
type StartGameParams struct{}

// DDZGame 斗地主游戏数据
type DDZGame struct {
	players map[int]*DDZPlayer
}

// DDZPlayer 斗地主玩家
type DDZPlayer struct {
	interfaces.ClientPlayer
	Seat int // 座号
}

// StartGame 启动游戏
func StartGame(params StartGameParams) (*DDZGame, error) {
	game := &DDZGame{}

	players, err := utils.CreateAndLoginUsersNum(3, config.ServerAddr, config.ClientVersion)
	if err != nil {
		return nil, fmt.Errorf("登录用户失败： %v", err)
	}
	game.players, err = apply(players)
	if err != nil {
		return nil, err
	}
	if err := checkStartGameNtf(game); err != nil {
		return nil, err
	}
	return game, nil
}

func checkStartGameNtf(game *DDZGame) error {
	for _, player := range game.players {
		expector := player.GetExpector(msgid.MsgID_ROOM_DDZ_START_GAME_NTF)
		ntf := room.DDZStartGameNtf{}
		err := expector.Recv(global.DefaultWaitMessageTime, &ntf)
		if err != nil {
			return fmt.Errorf("没有收到游戏开始通知 %v", err)
		}
	}
	return nil
}

func createDDZPlayer(players []interfaces.ClientPlayer, playerID uint64, seat int) *DDZPlayer {
	ddzPlayer := &DDZPlayer{
		Seat: seat,
	}
	setClientPlayer(players, playerID, ddzPlayer)
	ddzPlayer.AddExpectors(msgid.MsgID_ROOM_DDZ_DEAL_NTF, msgid.MsgID_ROOM_DDZ_START_GAME_NTF)
	return ddzPlayer
}

func apply(players []interfaces.ClientPlayer) (map[int]*DDZPlayer, error) {
	ddzPlayers := make(map[int]*DDZPlayer, len(players))

	expectors := []interfaces.MessageExpector{}
	for _, player := range players {
		e, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)
		if _, err := utils.ApplyJoinDesk(player, room.GameId_GAMEID_DDZ); err != nil {
			return nil, fmt.Errorf("申请加入房间失败 %v", err)
		}
		expectors = append(expectors, e)
	}
	ntf := room.RoomDeskCreatedNtf{}
	if err := expectors[0].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return nil, fmt.Errorf("没有收到创建房间通知 %v", err)
	}
	for _, rplayer := range ntf.GetPlayers() {
		playerID := rplayer.GetPlayerId()
		ddzPlayer := createDDZPlayer(players, playerID, int(rplayer.GetSeat()))
		ddzPlayers[int(rplayer.GetSeat())] = ddzPlayer
	}
	return ddzPlayers, nil
}

func setClientPlayer(players []interfaces.ClientPlayer, playerID uint64, ddzPlayer *DDZPlayer) {
	for _, player := range players {
		if player.GetID() == playerID {
			ddzPlayer.ClientPlayer = player
		}
	}
}
