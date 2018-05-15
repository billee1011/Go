package utils

import (
	"errors"
	"fmt"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/simulate/connect"
	"steve/simulate/interfaces"
	"time"

	"github.com/golang/protobuf/proto"
)

// DeskPlayer 牌桌玩家数据
type DeskPlayer struct {
	Player interfaces.ClientPlayer
	Seat   int
}

// DeskData 牌桌数据
type DeskData struct {
	Players    map[uint64]DeskPlayer // playerid -> deskPlayer
	BankerSeat int
}

// StartGameParams 启动游戏的参数
type StartGameParams struct {
	Cards      [][]*room.Card // 从庄家位置开始算起，每个位置的固定卡牌
	WallCards  []*room.Card   // 发完牌之后剩下的墙牌
	HszDir     room.Direction // 换三张的方向
	BankerSeat int            // 庄家座号
	ServerAddr string         // 服务器地址
	ClientVer  string         // 客户端版本号

	HszCards [][]*room.Card // 从庄家的位置算起，用来换三张的牌
}

// StartGame 启动一局游戏
func StartGame(params StartGameParams) (*DeskData, error) {
	players, err := createAndLoginUsers(params.ServerAddr, params.ClientVer)
	if err != nil {
		return nil, err
	}
	// TODO : game name
	if err := peipai("scxl", params.Cards, params.WallCards, params.HszDir, params.BankerSeat); err != nil {
		return nil, err
	}
	xipaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_XIPAI_NTF)
	fapaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_FAPAI_NTF)
	hszNotifyExpectors := createHSZNotifyExpector(players)
	seatMap, err := joinDesk(players)
	if err != nil {
		return nil, err
	}
	dd := DeskData{
		BankerSeat: params.BankerSeat,
	}
	dd.Players = map[uint64]DeskPlayer{}
	for _, player := range players {
		dd.Players[player.GetID()] = DeskPlayer{
			Player: player,
			Seat:   calcPlayerSeat(seatMap, player.GetID()),
		}
	}
	checkXipaiNtf(xipaiNtfExpectors, &dd, params.Cards, params.WallCards)
	checkFapaiNtf(fapaiNtfExpectors, &dd, params.Cards, params.WallCards)
	// 执行换三张
	if err := executeHSZ(hszNotifyExpectors, &dd, params.HszCards); err != nil {
		return nil, err
	}
	return &dd, nil
}

// GetDeskPlayerBySeat 根据座号获取牌桌玩家
func GetDeskPlayerBySeat(seat int, deskData *DeskData) *DeskPlayer {
	for _, player := range deskData.Players {
		if player.Seat == seat {
			return &player
		}
	}
	return nil
}

// GetSeatOffset 获取从 src 到 dest 的偏移步数
func GetSeatOffset(src int, dest int, count int) int {
	if dest >= src {
		return dest - src
	}
	return dest + count - src
}

func calcPlayerSeat(seatMap map[int]uint64, playerID uint64) int {
	for seat, pID := range seatMap {
		if pID == playerID {
			return seat
		}
	}
	return 0
}

var errCreateClientFailed = errors.New("创建客户端连接失败")

func createAndLoginUsers(ServerAddr string, ClientVer string) ([]interfaces.ClientPlayer, error) {
	players := []interfaces.ClientPlayer{}
	for i := 0; i < 4; i++ {
		client := connect.NewTestClient(ServerAddr, ClientVer)
		if client == nil {
			return nil, errCreateClientFailed
		}
		player, err := LoginUser(client, fmt.Sprintf("user_%d", i))
		if err != nil {
			return nil, fmt.Errorf("登录用户失败：%v", err)
		}
		players = append(players, player)
	}
	return players, nil
}

func joinDesk(players []interfaces.ClientPlayer) (map[int]uint64, error) {
	expectors := []interfaces.MessageExpector{}
	for _, player := range players {
		e, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)
		if err := ApplyJoinDesk(player); err != nil {
			return nil, err
		}
		expectors = append(expectors, e)
	}
	ntf := room.RoomDeskCreatedNtf{}
	if err := expectors[0].Recv(2*time.Second, &ntf); err != nil {
		return nil, err
	}
	seatMap := map[int]uint64{}
	for _, rplayer := range ntf.GetPlayers() {
		seatMap[int(rplayer.GetSeat())] = rplayer.GetPlayerId()
	}
	return seatMap, nil
}

// createExpectors 创建消息期望
func createExpectors(players []interfaces.ClientPlayer, msgID msgid.MsgID) map[uint64]interfaces.MessageExpector {
	result := map[uint64]interfaces.MessageExpector{}
	for _, player := range players {
		client := player.GetClient()
		e, _ := client.ExpectMessage(msgID)
		result[player.GetID()] = e
	}
	return result
}

// createHSZNotifyExpector 创建换三张通知消息期望
func createHSZNotifyExpector(players []interfaces.ClientPlayer) map[uint64]interfaces.MessageExpector {
	return createExpectors(players, msgid.MsgID_ROOM_HUANSANZHANG_NTF)
}

// executeHSZ 执行换三张
func executeHSZ(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, HszCards [][]*room.Card) error {
	finishNtfExpectors := map[uint64]interfaces.MessageExpector{}
	for playerID, e := range ntfExpectors {
		hszNtf := room.RoomHuansanzhangNtf{}
		if err := e.Recv(time.Second*2, &hszNtf); err != nil {
			return fmt.Errorf("未收到换三张通知: %v", err)
		}
		player := deskData.Players[playerID]
		offset := GetSeatOffset(deskData.BankerSeat, player.Seat, len(deskData.Players))
		cards := HszCards[offset]

		client := player.Player.GetClient()
		fe, _ := client.ExpectMessage(msgid.MsgID_ROOM_HUANSANZHANG_FINISH_NTF)
		finishNtfExpectors[playerID] = fe
		client.SendPackage(createMsgHead(msgid.MsgID_ROOM_HUANSANZHANG_REQ), &room.RoomHuansanzhangReq{
			Cards: gutils.RoomCards2UInt32(cards),
			Sure:  proto.Bool(true),
		})
	}
	for playerID, e := range finishNtfExpectors {
		finishNtf := room.RoomHuansanzhangFinishNtf{}
		if err := e.Recv(time.Second*2, &finishNtf); err != nil {
			return fmt.Errorf("玩家 %v 未收到换三张完成通知:%v", playerID, err)
		}
	}
	return nil
}

// checkXipaiNtf 检查洗牌通知
func checkXipaiNtf(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, seatCards [][]*room.Card, wallCards []*room.Card) error {
	totalCardCount := len(wallCards)
	for _, sc := range seatCards {
		totalCardCount += len(sc)
	}

	for _, e := range ntfExpectors {
		xipaiNtf := room.RoomXipaiNtf{}
		if err := e.Recv(time.Second*2, &xipaiNtf); err != nil {
			return fmt.Errorf("未收到洗牌通知： %v", err)
		}
		if xipaiNtf.GetBankerSeat() != uint32(deskData.BankerSeat) {
			return fmt.Errorf("庄家索引不对应。 expect: %d  actual:%d", deskData.BankerSeat, xipaiNtf.GetBankerSeat())
		}
		if xipaiNtf.GetTotalCard() != uint32(totalCardCount) {
			return fmt.Errorf("总牌数不对应。 expect: %d actual: %d", totalCardCount, xipaiNtf.GetTotalCard())
		}
	}
	return nil
}

func checkPlayerCardCount(playerCardCounts []*room.PlayerCardCount, deskData *DeskData, seatCards [][]*room.Card) error {
	for _, playerCardCount := range playerCardCounts {
		playerID := playerCardCount.GetPlayerId()
		cardCount := int(playerCardCount.GetCardCount())

		seat := deskData.Players[playerID].Seat
		expectedCount := len(seatCards[seat])
		if cardCount != expectedCount {
			return fmt.Errorf("playerCardCount 卡牌数量不对")
		}
	}
	return nil
}

// checkFapaiNtf 检查发牌通知
func checkFapaiNtf(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, seatCards [][]*room.Card, wallCards []*room.Card) error {
	for playerID, e := range ntfExpectors {
		fapaiNtf := room.RoomFapaiNtf{}
		if err := e.Recv(time.Second*2, &fapaiNtf); err != nil {
			return fmt.Errorf("未收到发牌通知： %v", err)
		}
		seat := deskData.Players[playerID].Seat
		expectCards := gutils.RoomCards2UInt32(seatCards[seat])
		ntfCards := fapaiNtf.GetCards()
		for index, c := range expectCards {
			if c != ntfCards[index] {
				return fmt.Errorf("收到的发牌通知，牌不对。 期望：%v 实际：%v 玩家：%v 座号:%d", expectCards, ntfCards, playerID, seat)
			}
		}
		if err := checkPlayerCardCount(fapaiNtf.GetPlayerCardCounts(), deskData, seatCards); err != nil {
			return err
		}
	}
	return nil
}
