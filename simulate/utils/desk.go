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
	Player    interfaces.ClientPlayer
	Seat      int
	Expectors map[msgid.MsgID]interfaces.MessageExpector // 消息期望， 消息 ID -> expector
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

	HszCards     [][]*room.Card   // 从庄家的位置算起，用来换三张的牌
	DingqueColor []room.CardColor // 定缺花色。 从庄家位置算起

}

func (p *StartGameParams) copyCards(cards []*room.Card) []*room.Card {
	return append([]*room.Card{}, cards...)
}

func (p *StartGameParams) copy2dCards(cards [][]*room.Card) [][]*room.Card {
	result := [][]*room.Card{}
	for _, c := range cards {
		result = append(result, p.copyCards(c))
	}
	return result
}
func (p *StartGameParams) copyColors(colors []room.CardColor) []room.CardColor {
	return append([]room.CardColor{}, colors...)
}

// Clone 创建一个新的副本
func (p *StartGameParams) Clone() StartGameParams {
	return StartGameParams{
		Cards:        p.copy2dCards(p.Cards),
		WallCards:    p.copyCards(p.WallCards),
		HszDir:       p.HszDir,
		BankerSeat:   p.BankerSeat,
		ServerAddr:   p.ServerAddr,
		ClientVer:    p.ClientVer,
		HszCards:     p.copy2dCards(p.HszCards),
		DingqueColor: p.copyColors(p.DingqueColor),
	}
}

// StartGame 启动一局游戏
// 开始后停留在等待庄家出牌状态
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
	// hszNotifyExpectors := createHSZNotifyExpector(players)
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
			Player:    player,
			Seat:      calcPlayerSeat(seatMap, player.GetID()),
			Expectors: createPlayerExpectors(player.GetClient()),
		}
	}
	checkXipaiNtf(xipaiNtfExpectors, &dd, params.Cards, params.WallCards)
	checkFapaiNtf(fapaiNtfExpectors, &dd, params.Cards, params.WallCards)
	// 执行换三张
	if err := executeHSZ(&dd, params.HszCards); err != nil {
		return nil, err
	}
	if err := executeDingque(&dd, params.DingqueColor); err != nil {
		return nil, err
	}
	return &dd, nil
}

// createPlayerExpectors 创建玩家的麻将逻辑消息期望
func createPlayerExpectors(client interfaces.Client) map[msgid.MsgID]interfaces.MessageExpector {
	msgs := []msgid.MsgID{msgid.MsgID_ROOM_CHUPAIWENXUN_NTF, msgid.MsgID_ROOM_PENG_NTF, msgid.MsgID_ROOM_GANG_NTF, msgid.MsgID_ROOM_HU_NTF,
		msgid.MsgID_ROOM_ZIXUN_NTF, msgid.MsgID_ROOM_CHUPAI_NTF,
		msgid.MsgID_ROOM_MOPAI_NTF, msgid.MsgID_ROOM_WAIT_QIANGGANGHU_NTF,
		msgid.MsgID_ROOM_TINGINFO_NTF, msgid.MsgID_ROOM_INSTANT_SETTLE, msgid.MsgID_ROOM_ROUND_SETTLE, msgid.MsgID_ROOM_DESK_DISMISS_NTF,
	}
	result := map[msgid.MsgID]interfaces.MessageExpector{}
	for _, msg := range msgs {
		result[msg], _ = client.ExpectMessage(msg)
	}
	return result
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

// executeHSZ 执行换三张
func executeHSZ(deskData *DeskData, HszCards [][]*room.Card) error {
	finishNtfExpectors := map[uint64]interfaces.MessageExpector{}
	for playerID, player := range deskData.Players {
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

// checkDingqueColor 检查定缺花色是否正确
func checkDingqueColor(deskData *DeskData, playerColors []*room.PlayerDingqueColor, expectedColors []room.CardColor) bool {
	if len(playerColors) != len(expectedColors) {
		return false
	}
	for _, playerColor := range playerColors {
		playerID := playerColor.GetPlayerId()
		seat := deskData.Players[playerID].Seat
		offset := GetSeatOffset(deskData.BankerSeat, seat, len(deskData.Players))
		expected := expectedColors[offset]
		actual := playerColor.GetColor()
		if expected != actual {
			return false
		}
	}
	return true
}

// executeDingque 执行定缺
func executeDingque(deskData *DeskData, colors []room.CardColor) error {
	finishExpectors := map[uint64]interfaces.MessageExpector{}
	for playerID, player := range deskData.Players {
		seat := player.Seat
		offset := GetSeatOffset(deskData.BankerSeat, seat, len(deskData.Players))

		client := player.Player.GetClient()
		finishExpector, _ := client.ExpectMessage(msgid.MsgID_ROOM_DINGQUE_FINISH_NTF)
		finishExpectors[playerID] = finishExpector

		client.SendPackage(createMsgHead(msgid.MsgID_ROOM_DINGQUE_REQ), &room.RoomDingqueReq{
			Color: colors[offset].Enum(),
		})
	}
	for playerID, e := range finishExpectors {
		ntf := room.RoomDingqueFinishNtf{}
		if err := e.Recv(time.Second*2, &ntf); err != nil {
			return fmt.Errorf("玩家 %d 未收到定缺完成通知: %v", playerID, err)
		}
		if !checkDingqueColor(deskData, ntf.GetPlayerDingqueColor(), colors) {
			return fmt.Errorf("玩家收到的定缺通知中花色不正确。 expected:%v actual:%v", colors, ntf.GetPlayerDingqueColor())
		}
	}
	return nil
}
