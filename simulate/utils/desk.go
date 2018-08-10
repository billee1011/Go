package utils

import (
	"errors"
	"fmt"
	"steve/client_pb/common"
	"steve/client_pb/match"
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/structs"
	"time"

	"github.com/Sirupsen/logrus"

	"steve/room/flows/ddzflow/ddz/states"

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
	BankerSeat int                   // 庄家/地主 的座位号
	DDZData    DDZData               // 斗地主信息
	DiFen      uint64                // 底分
}

// DDZData 斗地主信息
type DDZData struct {
	CurState     int                         // 当前状态
	NextState    *room.NextStage             // 下一状态的数据
	AssignLordID uint64                      // 服务器指定的叫地主玩家的playerID
	ResultLordID uint64                      // 最终的地主的playerID
	Params       structs.StartPukeGameParams // 启动参数
}

// StartGame 启动一局游戏
// 开始后停留在等待庄家出牌状态
func StartGame(params structs.StartGameParams) (*DeskData, error) {
	players, err := CreateAndLoginUsers(params.PlayerNum)
	if err != nil {
		return nil, err
	}
	// TODO : game name
	if err := Peipai(params.PeiPaiGame, params.Cards, params.WallCards, params.HszDir, params.BankerSeat); err != nil {
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
		DiFen:      params.DiFen,
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
	if params.IsDq {
		if err := executeDingque(&dd, params.DingqueColor); err != nil {
			return nil, err
		}
	}
	playerIDs := make([]uint64, 0, len(dd.Players))
	for playerID := range dd.Players {
		playerIDs = append(playerIDs, playerID)
	}
	logrus.WithFields(logrus.Fields{
		"players": playerIDs,
		"params":  params,
	}).Infoln("游戏开始完成")

	return &dd, nil
}

// StartDDZGame 启动斗地主游戏
// 开始后停留在等待庄家出牌状态
func StartDDZGame(params structs.StartPukeGameParams) (*DeskData, error) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::StartDDZGame",
	})

	logEntry.Info("")

	// 创建并登录3个玩家
	players, err := CreateAndLoginUsers(3)
	if err != nil {
		return nil, err
	}

	// 通知服务器：配牌
	if err := Peipai(params.PeiPaiGame, params.Cards, params.WallCards, params.HszDir, params.BankerSeat); err != nil {
		return nil, err
	}

	// 所有玩家的斗地主发牌通知期望
	fapaiNtfExpectors := createExpectors(players, msgid.MsgID_ROOM_DDZ_DEAL_NTF)

	// hszNotifyExpectors := createHSZNotifyExpector(players)

	gameID := params.GameID // 设置游戏ID

	// 牌桌数据
	deskData := DeskData{
		BankerSeat: params.BankerSeat, // 地主的座位号
	}

	// 三个玩家的游戏开始消息的期望
	expectorsStart := []interfaces.MessageExpector{}

	for _, player := range players {
		es, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_START_GAME_NTF)
		expectorsStart = append(expectorsStart, es)
	}

	// 加入牌桌
	// 返回的 seatMap:座位ID 与playerID 的map
	seatMap, err := DDZjoinDesk(players, gameID)
	if err != nil {
		return nil, err
	}

	// 设置玩家金币数
	if err := majongPlayerGold(params.PlayerSeatGold, seatMap); err != nil {
		return nil, err
	}

	// 等待接收游戏开始消息
	for i, player := range players {
		stntf := room.DDZStartGameNtf{}
		if err := expectorsStart[i].Recv(global.DefaultWaitMessageTime, &stntf); err != nil {
			logEntry.Info("接收斗地主游戏开始的通知超时,playerID = ", player.GetID())
			return nil, err
		}

		// 服务器指定的叫地主玩家
		deskData.DDZData.AssignLordID = stntf.GetPlayerId()

		// 下一状态信息
		deskData.DDZData.NextState = stntf.GetNextStage()

		logEntry.Infof("玩家%d收到了斗地主游戏开始的通知,叫地主玩家 = %d，下一状态 = %v, 进入下一状态等待时间 = %d",
			player.GetID(), stntf.GetPlayerId(), deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())
	}

	// 建立playerid -> deskPlayer的map
	deskData.Players = map[uint64]DeskPlayer{}
	for _, player := range players {
		deskData.Players[player.GetID()] = DeskPlayer{
			Player:    player,                                       // clientPlayer
			Seat:      calcPlayerSeat(seatMap, player.GetID()),      // 座位号
			Expectors: createDDZPlayerExpectors(player.GetClient()), // 斗地主所有的消息期望
		}
	}

	// 检测收到的发牌消息是否符合预期
	if err := checkDDZFapaiNtf(fapaiNtfExpectors, &deskData, params.Cards); err != nil {
		return nil, err
	}

	// 保存参数
	deskData.DDZData.Params = params

	return &deskData, nil
}

// sendDoubleReq 发送加倍请求
// double : 是否加倍
func sendDoubleReq(player *DeskPlayer, double bool) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "sendJiabeiReq()",
	}).Info("发出加倍请求，玩家 = ", player.Player.GetID())

	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_DDZ_DOUBLE_REQ), &room.DDZDoubleReq{
		IsDouble: &double, // 加倍为true，不加倍为false
	})

	return err
}

// createPlayerExpectors 创建玩家的麻将逻辑消息期望
func createPlayerExpectors(client interfaces.Client) map[msgid.MsgID]interfaces.MessageExpector {
	msgs := []msgid.MsgID{msgid.MsgID_ROOM_DINGQUE_FINISH_NTF, msgid.MsgID_ROOM_HUANSANZHANG_FINISH_NTF, msgid.MsgID_ROOM_CHUPAIWENXUN_NTF,
		msgid.MsgID_ROOM_BUHUA_NTF, msgid.MsgID_ROOM_CHI_NTF, msgid.MsgID_ROOM_PENG_NTF, msgid.MsgID_ROOM_GANG_NTF, msgid.MsgID_ROOM_HU_NTF, msgid.MsgID_ROOM_TUOGUAN_NTF,
		msgid.MsgID_ROOM_ZIXUN_NTF, msgid.MsgID_ROOM_CHUPAI_NTF,
		msgid.MsgID_ROOM_MOPAI_NTF, msgid.MsgID_ROOM_WAIT_QIANGGANGHU_NTF,
		msgid.MsgID_ROOM_TINGINFO_NTF, msgid.MsgID_ROOM_INSTANT_SETTLE, msgid.MsgID_ROOM_ROUND_SETTLE, msgid.MsgID_ROOM_DESK_DISMISS_NTF,
		msgid.MsgID_ROOM_CHAT_NTF, msgid.MsgID_ROOM_RESUME_GAME_RSP, msgid.MsgID_ROOM_DESK_QUIT_RSP,
		msgid.MsgID_ROOM_GAMEOVER_NTF, msgid.MsgID_ROOM_CHANGE_PLAYERS_RSP, msgid.MsgID_MATCH_SUC_CREATE_DESK_NTF, msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF,
		msgid.MsgID_ROOM_DESK_NEED_RESUME_RSP,
		msgid.MsgID_ROOM_GAMEOVER_NTF, msgid.MsgID_ROOM_CHANGE_PLAYERS_RSP, msgid.MsgID_MATCH_SUC_CREATE_DESK_NTF, msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF,
		msgid.MsgID_ROOM_HUANSANZHANG_NTF,
		msgid.MsgID_ROOM_DINGQUE_NTF, msgid.MsgID_HALL_GET_PLAYER_GAME_INFO_RSP, msgid.MsgID_ROOM_USE_PROP_NTF,
	}
	result := map[msgid.MsgID]interfaces.MessageExpector{}
	for _, msg := range msgs {
		result[msg], _ = client.ExpectMessage(msg)
	}
	return result
}

// createDDZPlayerExpectors 创建斗地主玩家的消息期望
func createDDZPlayerExpectors(client interfaces.Client) map[msgid.MsgID]interfaces.MessageExpector {

	logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::createDDZPlayerExpectors",
	}).Info("")

	// 所有期望的消息
	msgs := []msgid.MsgID{
		msgid.MsgID_ROOM_DDZ_START_GAME_NTF, // 斗地主 开始游戏通知
		msgid.MsgID_ROOM_DDZ_DEAL_NTF,       // 斗地主 发牌通知
		msgid.MsgID_ROOM_DDZ_GRAB_LORD_RSP,  // 斗地主 叫/抢地主响应
		msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF,  // 斗地主 叫/抢地主广播
		msgid.MsgID_ROOM_DDZ_LORD_NTF,       // 斗地主 叫/抢地主通知
		msgid.MsgID_ROOM_DDZ_DOUBLE_RSP,     // 斗地主 加倍响应
		msgid.MsgID_ROOM_DDZ_DOUBLE_NTF,     // 斗地主 加倍通知
		//msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP,  // 斗地主 出牌响应
		//msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF,  // 斗地主 出牌通知
		msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF, // 斗地主 结束通知
		msgid.MsgID_ROOM_DDZ_RESUME_RSP,    // 斗地主 回复对局响应
		//msgid.MsgID_ROOM_DDZ_RESUME_RSP,    // 斗地主 回复对局响应
		msgid.MsgID_HALL_GET_PLAYER_GAME_INFO_RSP,
		msgid.MsgID_ROOM_USE_PROP_NTF,
	}

	result := map[msgid.MsgID]interfaces.MessageExpector{}

	// 为每一个消息建立期待
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

// calcPlayerSeat 根据playerID获取座位号
func calcPlayerSeat(seatMap map[int]uint64, playerID uint64) int {
	for seat, pID := range seatMap {
		if pID == playerID {
			return seat
		}
	}
	return 0
}

var errCreateClientFailed = errors.New("创建客户端连接失败")

// CreateAndLoginUsers 创建玩家并登录
func CreateAndLoginUsers(num int) ([]interfaces.ClientPlayer, error) {
	return CreateAndLoginUsersNum(num)
}

// CreateAndLoginUsersNum 指定创建的人数
func CreateAndLoginUsersNum(num int) ([]interfaces.ClientPlayer, error) {
	players := []interfaces.ClientPlayer{}
	for i := 0; i < num; i++ {
		player, err := LoginNewPlayer()
		if err != nil {
			return nil, fmt.Errorf("登录用户失败： %v", err)
		}

		players = append(players, player)
	}
	return players, nil
}

// 加入牌桌
// 返回：座位ID 与 playerID的map
func joinDesk(players []interfaces.ClientPlayer, gameID common.GameId) (map[int]uint64, error) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "joinDesk",
		"GameID":    gameID,
	})

	logEntry.Info("申请加入牌桌")

	// 所有期待的消息
	expectors := []interfaces.MessageExpector{}

	for _, player := range players {

		// 期望收到桌子创建的通知
		e, _ := player.GetClient().ExpectMessage(msgid.MsgID_MATCH_SUC_CREATE_DESK_NTF)

		// 申请加入牌桌
		if _, err := ApplyJoinDesk(player, gameID); err != nil {
			return nil, fmt.Errorf("请求加入房间失败: %v", err)
		}

		expectors = append(expectors, e)
	}

	// 等待接收桌子创建消息
	ntf := match.MatchSucCreateDeskNtf{}
	if err := expectors[0].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return nil, err
	}

	logEntry.Info("收到了桌子创建的通知")

	seatMap := map[int]uint64{}
	for _, rplayer := range ntf.GetPlayers() {
		seatMap[int(rplayer.GetSeat())] = rplayer.GetPlayerId()
	}
	return seatMap, nil
}

// DDZjoinDesk 斗地主加入牌桌
// 返回：座位ID 与 playerID的map
func DDZjoinDesk(players []interfaces.ClientPlayer, gameID common.GameId) (map[int]uint64, error) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "DDZjoinDesk",
		"GameID":    gameID,
	})

	logEntry.Info("申请加入牌桌")

	// 所有的期待桌子创建消息
	expectors := []interfaces.MessageExpector{}

	for _, player := range players {

		// 期望收到桌子创建的通知
		e, _ := player.GetClient().ExpectMessage(msgid.MsgID_MATCH_SUC_CREATE_DESK_NTF)
		expectors = append(expectors, e)

		// 申请加入牌桌
		if _, err := ApplyJoinDesk(player, gameID); err != nil {
			return nil, err
		}
	}

	// 等待接收桌子创建消息
	ntf := match.MatchSucCreateDeskNtf{}
	if err := expectors[0].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return nil, fmt.Errorf("没有收到创建房间通知: %v", err)
	}

	logEntry.Info("收到了桌子创建的通知")

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

// 通知服务器：每个玩家的动画播放完成
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

// 检测所有玩家的牌数量
// playerCardCounts ：每个玩家手牌信息的数组
// deskData			: 桌子数据
// seatCards		: 从地主开始，每个玩家的牌，客户端配置的
func checkPlayerCardCount(playerCardCounts []*room.PlayerCardCount, deskData *DeskData, seatCards [][]uint32) error {
	for _, playerCardCount := range playerCardCounts {
		playerID := playerCardCount.GetPlayerId()

		// 该玩家的手牌数量，服务器通知的
		cardCount := int(playerCardCount.GetCardCount())

		seat := deskData.Players[playerID].Seat

		// 客户端配置的手牌数量，所以也是期待的数量
		expectedCount := len(seatCards[seat])

		// 两者不等，则报错
		if cardCount != expectedCount {
			return fmt.Errorf("playerCardCount 卡牌数量不对，玩家手牌数量:%v，期待数量:%v", cardCount, expectedCount)
		}
	}
	return nil
}

// checkFapaiNtf 检查发牌通知
func checkFapaiNtf(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, seatCards [][]uint32, wallCards []uint32) error {
	for playerID, e := range ntfExpectors {
		fapaiNtf := room.RoomFapaiNtf{}
		if err := e.Recv(global.DefaultWaitMessageTime, &fapaiNtf); err != nil {
			return fmt.Errorf("%v未收到发牌通知： %v", playerID, err)
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

// 	checkDDZFapaiNtf 检测斗地主的发牌消息是否符合预期
//	ntfExpectors 	: 所有的playerID与消息期待的map
//  deskData		: 牌桌数据
//  seatCards		: 从地主开始，每个玩家的牌，客户端配置的
//
func checkDDZFapaiNtf(ntfExpectors map[uint64]interfaces.MessageExpector, deskData *DeskData, seatCards [][]uint32) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::checkDDZFapaiNtf",
	})

	for playerID, e := range ntfExpectors {

		// 发牌通知消息
		fapaiNtf := room.DDZDealNtf{}

		// 接收该消息
		if err := e.Recv(global.DefaultWaitMessageTime, &fapaiNtf); err != nil {
			return fmt.Errorf("未收到发牌通知： %v", err)
		}

		// 座位号
		seat := deskData.Players[playerID].Seat

		// 期待的牌（由于seatCards里面的牌是客户端配置的，所以服务器发下来时应该一致）
		expectCards := states.DDZSortDescend(seatCards[seat])

		// 服务器下发的牌
		ntfCards := fapaiNtf.GetCards()

		// 逐个比较，不一致则报错
		for index, c := range expectCards {
			if c != ntfCards[index] {
				return fmt.Errorf("收到的发牌通知，牌不对。 期望：%v 实际：%v 玩家playerID：%v 座号:%d", expectCards, ntfCards, playerID, seat)
			}
		}

		// 下一状态信息
		deskData.DDZData.NextState = fapaiNtf.GetNextStage()

		// 检测每个玩家的手牌数量是否和期待的相同
		//if err := checkPlayerCardCount(fapaiNtf.GetPlayerCardCounts(), deskData, seatCards); err != nil {
		//	return err
		//}
	}

	// 暂停2秒（因为2秒之后才进入叫地主状态）
	//time.Sleep(time.Second * 2)

	logEntry.Infof("发牌处理中，下一状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

	logEntry.Infof("发送发牌动画播放完成的请求")

	// 通知服务器，每个玩家的发牌动画播放完成
	return sendCartoonFinish(room.CartoonType_CTNT_DDZ_FAPAI, deskData)
}

// executeHSZ 执行换三张
func executeHSZ(deskData *DeskData, HszCards [][]uint32) error {
	if HszCards == nil {
		logrus.Infoln("换三张牌没配置，不执行换三张")
		return nil
	}
	for playerID, player := range deskData.Players {
		hszNtf := room.RoomHuansanzhangNtf{}
		e := player.Expectors[msgid.MsgID_ROOM_HUANSANZHANG_NTF]
		if err := e.Recv(global.DefaultWaitMessageTime, &hszNtf); err != nil {
			return fmt.Errorf("玩家 %v 未收到换三张推荐通知:%v", playerID, err)
		}
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
	for playerID, player := range deskData.Players {
		hszNtf := room.RoomDingqueNtf{}
		e := player.Expectors[msgid.MsgID_ROOM_DINGQUE_NTF]
		if err := e.Recv(global.DefaultWaitMessageTime, &hszNtf); err != nil {
			return fmt.Errorf("玩家 %v 未收到定缺推荐通知:%v", playerID, err)
		}
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
