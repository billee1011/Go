package game

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/structs"
	"steve/simulate/utils"
)

// DDZGame 斗地主游戏数据
type DDZGame struct {
	players map[int]*DDZPlayer // 座位ID 与 *DDZPlayer 的映射
}

// DDZPlayer 斗地主玩家
type DDZPlayer struct {
	interfaces.ClientPlayer
	Seat int // 座号
}

// StartGame 启动游戏
func StartGame(params structs.StartPukeGameParams) (*DDZGame, error) {
	game := &DDZGame{}

	// 创建并登陆3个用户
	players, err := utils.CreateAndLoginUsersNum(3, config.ServerAddr, config.ClientVersion)
	if err != nil {
		return nil, fmt.Errorf("登录用户失败： %v", err)
	}

	// 创建出DDZPlayer,请求进入桌子
	game.players, err = apply(players)
	if err != nil {
		return nil, err
	}

	// 检测是否收到游戏开始消息
	//if err := checkStartGameNtf(game); err != nil {
	//	return nil, err
	//}
	return game, nil
}

func checkStartGameNtf(game *DDZGame) error {

	// 每一个玩家需收到游戏开启的通知消息
	for _, player := range game.players {
		expector := player.GetExpector(msgid.MsgID_ROOM_DDZ_START_GAME_NTF)
		ntf := room.DDZStartGameNtf{}
		err := expector.Recv(global.DefaultWaitMessageTime, &ntf)
		if err != nil {
			return fmt.Errorf("没有收到游戏开始通知 %v", err)
		}
		fmt.Printf("玩家%d收到游戏开始的通知\n", player.GetID())
	}
	return nil
}

// 根据ClientPlayer 创建DDZPlayer
func createDDZPlayer(players []interfaces.ClientPlayer, playerID uint64, seat int) *DDZPlayer {
	ddzPlayer := &DDZPlayer{
		Seat: seat,
	}
	setClientPlayer(players, playerID, ddzPlayer)

	// 两个消息期待
	ddzPlayer.AddExpectors(msgid.MsgID_ROOM_DDZ_DEAL_NTF, msgid.MsgID_ROOM_DDZ_START_GAME_NTF)

	return ddzPlayer
}

func apply(players []interfaces.ClientPlayer) (map[int]*DDZPlayer, error) {

	// 座位ID 与 *DDZPlayer 的映射
	ddzPlayers := make(map[int]*DDZPlayer, len(players))

	expectors := []interfaces.MessageExpector{}

	// 游戏开启通知的期待
	startExpectors := []interfaces.MessageExpector{}

	for _, player := range players {

		// 创建桌子通知的期望
		e, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)

		// 游戏开启通知的期望
		se, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_START_GAME_NTF)

		// 申请加入桌子
		if _, err := utils.ApplyJoinDesk(player, room.GameId_GAMEID_DOUDIZHU); err != nil {
			return nil, fmt.Errorf("申请加入房间失败 %v", err)
		}

		expectors = append(expectors, e)
		startExpectors = append(startExpectors, se)
	}

	// 等待桌子创建通知的消息
	ntf := room.RoomDeskCreatedNtf{}
	if err := expectors[0].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return nil, fmt.Errorf("没有收到创建房间通知 %v", err)
	}

	// 等待游戏开始通知的消息
	startNtf := room.DDZStartGameNtf{}
	if err := startExpectors[0].Recv(global.DefaultWaitMessageTime, &startNtf); err != nil {
		return nil, fmt.Errorf("没有收到游戏开始通知 %v", err)
	}

	// 根据消息返回值里面的 RoomPlayerInfo 创建所有的 DDZPlayer
	for _, rplayer := range ntf.GetPlayers() {
		playerID := rplayer.GetPlayerId()
		ddzPlayer := createDDZPlayer(players, playerID, int(rplayer.GetSeat()))
		ddzPlayers[int(rplayer.GetSeat())] = ddzPlayer
	}
	return ddzPlayers, nil
}

// 从ClientPlayer数组中找到playerID，并设置DDZPlayer
func setClientPlayer(players []interfaces.ClientPlayer, playerID uint64, ddzPlayer *DDZPlayer) {
	for _, player := range players {
		if player.GetID() == playerID {
			ddzPlayer.ClientPlayer = player
		}
	}
}
