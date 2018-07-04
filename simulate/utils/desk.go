package utils

import (
	"errors"
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/structs"
	"time"

	"github.com/Sirupsen/logrus"

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

// StartGame 启动一局游戏
// 开始后停留在等待庄家出牌状态
func StartGame(params structs.StartGameParams) (*DeskData, error) {
	players, err := createAndLoginUsers(4, config.ServerAddr, params.ClientVer)
	if err != nil {
		return nil, err
	}
	// TODO : game name
	if err := peipai(params.PeiPaiGame, params.Cards, params.WallCards, params.HszDir, params.BankerSeat); err != nil {
		return nil, err
	}
	// 配置麻将选项
	if err := majongOption(params.PeiPaiGame, params.IsHsz); err != nil {
		return nil, err
	}
	xipaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_XIPAI_NTF)
	fapaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_FAPAI_NTF)
	// hszNotifyExpectors := createHSZNotifyExpector(players)
	gameID := params.GameID // 设置游戏ID
	seatMap, err := joinDesk(players, gameID)
	if err != nil {
		return nil, err
	}
	// 设置玩家Gold
	if err := majongPlayerGold(params.PlayerSeatGold, seatMap); err != nil {
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
	if err := checkFapaiNtf(fapaiNtfExpectors, &dd, params.Cards, params.WallCards); err != nil {
		return nil, err
	}
	// 是否执行换三张
	if params.IsHsz {
		// 执行换三张
		if err := executeHSZ(&dd, params.HszCards); err != nil {
			return nil, err
		}
	}
	if err := executeDingque(&dd, params.DingqueColor); err != nil {
		return nil, err
	}
	return &dd, nil
}

// StartPokeGame 启动扑克游戏
// 开始后停留在等待庄家出牌状态
func StartPokeGame(params structs.StartPukeGameParams) /*(*DeskData, error)*/ error {

	// 创建并登录3个玩家
	_, err := createAndLoginUsers(3, config.ServerAddr, params.ClientVer)
	if err != nil {
		//return nil, err
		return err
	}

	// 通知服务器：配牌
	if err := peipai(params.PeiPaiGame, params.Cards, params.WallCards, params.HszDir, params.BankerSeat); err != nil {
		// return nil, err
		return err
	}

	// 通知服务器：麻将选项（是否开启换三张）
	//if err := majongOption(params.PeiPaiGame, params.IsHsz); err != nil {
	//	return nil, err
	//}

	// 所有玩家的洗牌通知期望
	xipaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_XIPAI_NTF)

	// 所有玩家的发牌通知期望
	fapaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_FAPAI_NTF)

	// hszNotifyExpectors := createHSZNotifyExpector(players)

	gameID := params.GameID // 设置游戏ID

	/* 	// 加入牌桌
	   	// 返回的 seatMap:座位ID 与errlayerID 的map
	   	seatMap, err := joinDesk(plerrers, gameID)
	   	if err != nil {
	   		return nil, err
	   	}

	   	// 设置玩家金币数
	   	if err := majongPlayerGold(errrams.PlayerSeatGold, seatMap); err != nil {
	   		return nil, err
	   	}

	   	// 拍桌数据
	   	dd := DeskData{
	   		BankerSeat: params.BankerSeat,
	   	}

	   	// 建立playerid -> deskPlayer的map
	   	dd.Players = map[uint64]DeskPlayer{}
	   	for _, player := range players {
	   		dd.Players[player.GetID()] = DeskPlayer{
	   			Player:    player,
	   			Seat:      calcPlayerSeat(seatMap, player.GetID()),
	   			Expectors: createPlayerExpectors(player.GetClient()),
	   		}
	   	}

	   	// 检测是否收到洗牌通知
	   	checkXipaiNtf(xipaiNtfExpectors, &dd, params.Cards, params.WallCards)

	   	// 检测是否收到发牌通知
	   	if err := checkFapaiNtf(fapaiNtfExpectors, &dd, params.Cards, params.WallCards); err != nil {
	   		return nil, err
	   	}

	   	// 是否执行换三张
	   	if params.IsHsz {
	   		// 执行换三张
	   		if err := executeHSZ(&dd, params.HszCards); err != nil {
	   			return nil, err
	   		}
	   	}

	   	// 是否执行定缺
	   	if err := executeDingque(&dd, params.DingqueColor); err != nil {
	   		return nil, err
	   	}

		   return &dd, nil */
	return nil
}

// createPlayerExpectors 创建玩家的麻将逻辑消息期望
func createPlayerExpectors(client interfaces.Client) map[msgid.MsgID]interfaces.MessageExpector {
	msgs := []msgid.MsgID{msgid.MsgID_ROOM_DINGQUE_FINISH_NTF, msgid.MsgID_ROOM_HUANSANZHANG_FINISH_NTF, msgid.MsgID_ROOM_CHUPAIWENXUN_NTF,
		msgid.MsgID_ROOM_PENG_NTF, msgid.MsgID_ROOM_GANG_NTF, msgid.MsgID_ROOM_HU_NTF, msgid.MsgID_ROOM_TUOGUAN_NTF,
		msgid.MsgID_ROOM_ZIXUN_NTF, msgid.MsgID_ROOM_CHUPAI_NTF,
		msgid.MsgID_ROOM_MOPAI_NTF, msgid.MsgID_ROOM_WAIT_QIANGGANGHU_NTF,
		msgid.MsgID_ROOM_TINGINFO_NTF, msgid.MsgID_ROOM_INSTANT_SETTLE, msgid.MsgID_ROOM_ROUND_SETTLE, msgid.MsgID_ROOM_DESK_DISMISS_NTF,
		msgid.MsgID_ROOM_CHAT_NTF, msgid.MsgID_ROOM_RESUME_GAME_RSP, msgid.MsgID_ROOM_DESK_QUIT_RSP,
		msgid.MsgID_ROOM_GAMEOVER_NTF,
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

func createAndLoginUsers(userNum int, ServerAddr string, ClientVer string) ([]interfaces.ClientPlayer, error) {
	return CreateAndLoginUsersNum(userNum, ServerAddr, ClientVer)
}

// CreateAndLoginUsersNum 指定创建的人数
func CreateAndLoginUsersNum(num int, ServerAddr string, ClientVer string) ([]interfaces.ClientPlayer, error) {
	players := []interfaces.ClientPlayer{}
	for i := 0; i < num; i++ {

		// 建立客户端
		client := connect.NewTestClient(ServerAddr, ClientVer)
		if client == nil {
			return nil, errCreateClientFailed
		}

		// 登陆
		player, err := LoginUser(client, global.AllocUserName())
		if err != nil {
			return nil, fmt.Errorf("登录用户失败：%v", err)
		}

		players = append(players, player)
	}
	return players, nil
}

// 加入牌桌
// 返回：座位ID 与 playerID的map
func joinDesk(players []interfaces.ClientPlayer, gameID room.GameId) (map[int]uint64, error) {

	// 所有期待的消息
	expectors := []interfaces.MessageExpector{}

	for _, player := range players {

		// 期望收到桌子创建的通知
		e, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)

		// 申请加入牌桌
		if _, err := ApplyJoinDesk(player, gameID); err != nil {
			return nil, err
		}

		expectors = append(expectors, e)
	}

	// 等待接收消息
	ntf := room.RoomDeskCreatedNtf{}
	if err := expectors[0].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
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

func sendCartoonFinish(cartoonType room.CartoonType, deskData *DeskData) error {
	for _, player := range deskData.Players {
		client := player.Player.GetClient()
		if _, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_CARTOON_FINISH_REQ), &room.RoomCartoonFinishReq{
			CartoonType: cartoonType.Enum(),
		}); err != nil {
			return err
		}
	}
	return nil
}

// checkXipaiNtf 检查洗牌通知
func checkXipaiNtf(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, seatCards [][]uint32, wallCards []uint32) error {
	totalCardCount := len(wallCards)
	for _, sc := range seatCards {
		totalCardCount += len(sc)
	}

	for _, e := range ntfExpectors {
		xipaiNtf := room.RoomXipaiNtf{}
		if err := e.Recv(global.DefaultWaitMessageTime, &xipaiNtf); err != nil {
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

func checkPlayerCardCount(playerCardCounts []*room.PlayerCardCount, deskData *DeskData, seatCards [][]uint32) error {
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
func checkFapaiNtf(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, seatCards [][]uint32, wallCards []uint32) error {
	for playerID, e := range ntfExpectors {
		fapaiNtf := room.RoomFapaiNtf{}
		if err := e.Recv(global.DefaultWaitMessageTime, &fapaiNtf); err != nil {
			return fmt.Errorf("未收到发牌通知： %v", err)
		}
		seat := deskData.Players[playerID].Seat
		expectCards := seatCards[seat]
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
	return sendCartoonFinish(room.CartoonType_CTNT_FAPAI, deskData)
}

// executeHSZ 执行换三张
func executeHSZ(deskData *DeskData, HszCards [][]uint32) error {
	if HszCards == nil {
		logrus.Infoln("换三张牌没配置，不执行换三张")
		return nil
	}
	finishNtfExpectors := map[uint64]interfaces.MessageExpector{}
	for playerID, player := range deskData.Players {
		offset := GetSeatOffset(deskData.BankerSeat, player.Seat, len(deskData.Players))
		cards := HszCards[offset]

		client := player.Player.GetClient()
		fe := player.Expectors[msgid.MsgID_ROOM_HUANSANZHANG_FINISH_NTF]
		finishNtfExpectors[playerID] = fe
		client.SendPackage(createMsgHead(msgid.MsgID_ROOM_HUANSANZHANG_REQ), &room.RoomHuansanzhangReq{
			Cards: cards,
			Sure:  proto.Bool(true),
		})
	}
	for playerID, e := range finishNtfExpectors {
		finishNtf := room.RoomHuansanzhangFinishNtf{}
		if err := e.Recv(global.DefaultWaitMessageTime, &finishNtf); err != nil {
			return fmt.Errorf("玩家 %v 未收到换三张完成通知:%v", playerID, err)
		}
	}
	return sendCartoonFinish(room.CartoonType_CTNT_HUANSANZHANG, deskData)
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
	if colors == nil {
		logrus.Infoln("定缺花色没配置，不执行定缺")
		return nil
	}
	for _, player := range deskData.Players {
		seat := player.Seat
		offset := GetSeatOffset(deskData.BankerSeat, seat, len(deskData.Players))
		SendDingqueReq(seat, deskData, colors[offset])
	}
	WaitDingqueFinish(deskData, global.DefaultWaitMessageTime, colors, []int{0, 1, 2, 3})
	return nil
}

// SendDingqueReq 发送定缺请求
func SendDingqueReq(seat int, deskData *DeskData, color room.CardColor) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(createMsgHead(msgid.MsgID_ROOM_DINGQUE_REQ), &room.RoomDingqueReq{
		Color: color.Enum(),
	})
	return err
}

// SendHuansanzhangReq 发送换三张请求
func SendHuansanzhangReq(seat int, deskData *DeskData, hszCards []uint32, sure bool) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(createMsgHead(msgid.MsgID_ROOM_HUANSANZHANG_REQ), &room.RoomHuansanzhangReq{
		Cards: hszCards,
		Sure:  proto.Bool(sure),
	})
	return err
}

// WaitDingqueFinish 等定缺完成通知
func WaitDingqueFinish(deskData *DeskData, duration time.Duration, expectedColors []room.CardColor, waitSeats []int) error {
	for _, seat := range waitSeats {
		ntf := room.RoomDingqueFinishNtf{}
		player := GetDeskPlayerBySeat(seat, deskData)
		expector := player.Expectors[msgid.MsgID_ROOM_DINGQUE_FINISH_NTF]
		if err := expector.Recv(duration, &ntf); err != nil {
			return err
		}
		if expectedColors != nil && !checkDingqueColor(deskData, ntf.GetPlayerDingqueColor(), expectedColors) {
			return fmt.Errorf("玩家收到的定缺通知中花色不正确。 expected:%v actual:%v", expectedColors, ntf.GetPlayerDingqueColor())
		}
	}
	return nil
}

// WaitHuansanzhangFinish 等待换三张完成
func WaitHuansanzhangFinish(deskData *DeskData, duration time.Duration, waitSeats []int, expectInCards []uint32, expectSeat int) error {
	for _, seat := range waitSeats {
		ntf := room.RoomHuansanzhangFinishNtf{}
		player := GetDeskPlayerBySeat(seat, deskData)
		expector := player.Expectors[msgid.MsgID_ROOM_HUANSANZHANG_FINISH_NTF]
		if err := expector.Recv(duration, &ntf); err != nil {
			return err
		}
		if expectInCards != nil && seat == expectSeat {
			for _, expectInCard := range expectInCards {
				result := func(targetCard uint32, cards []uint32) (rs bool) {
					rs = false
					for _, card := range cards {
						if card == targetCard {
							rs = true
							break
						}
					}
					return
				}(expectInCard, ntf.GetInCards())
				if !result {
					return fmt.Errorf("玩家收到的换三张通知中牌不正确。 expected:%v actual:%v", expectInCards, ntf.GetInCards())
				}
			}
		}
	}
	return nil
}

// WaitMoPaiNtf 等待摸牌通知
func WaitMoPaiNtf(deskData *DeskData, duration time.Duration, waitSeats []int, moCard uint32, moSeat int) error {
	for _, seat := range waitSeats {
		ntf := room.RoomMopaiNtf{}
		player := GetDeskPlayerBySeat(seat, deskData)
		expector := player.Expectors[msgid.MsgID_ROOM_MOPAI_NTF]
		if err := expector.Recv(duration, &ntf); err != nil {
			return err
		}
		if seat == moSeat && moSeat != -1 {
			if moCard != 0 && ntf.GetCard() != moCard {
				return fmt.Errorf("玩家摸到的牌不正确。 expected:%v actual:%v", moCard, ntf.GetCard())
			}
		}
	}
	return nil
}

// WaitTuoGuanNtf 等待托管通知
func WaitTuoGuanNtf(deskData *DeskData, duration time.Duration, waitSeats []int) error {
	for _, seat := range waitSeats {
		ntf := room.RoomTuoGuanNtf{}
		player := GetDeskPlayerBySeat(seat, deskData)
		expector := player.Expectors[msgid.MsgID_ROOM_TUOGUAN_NTF]
		if err := expector.Recv(duration, &ntf); err != nil {
			return err
		}
	}
	return nil
}
