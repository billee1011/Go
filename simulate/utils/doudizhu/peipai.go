package doudizhu

import (
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/structs"
	"steve/simulate/utils"
	"testing"

	"github.com/Sirupsen/logrus"
)

// NewStartDDZGameParamsTest1 创建棋牌启动参数测试
func NewStartDDZGameParamsTest1() structs.StartPukeGameParams {
	return structs.StartPukeGameParams{
		GameID: room.GameId_GAMEID_DOUDIZHU, // 斗地主
		Cards: [][]uint32{
			// 第一个人的牌(地主)
			{
				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_A), // 方块A
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_A),    // 梅花A
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_A),   // 红桃A
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_A),   // 黑桃A

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_2), // 方块2
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_2),    // 梅花2
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_2),   // 红桃2
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_2),   // 黑桃2

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_3), // 方块3
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_3),    // 梅花3
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_3),   // 红桃3
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_3),   // 黑桃3

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_4), // 方块4
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_4),    // 梅花4
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_4),   // 红桃4
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_4),   // 黑桃4

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_5), // 方块5

			},

			// 第二个人的牌
			{
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_5),  // 梅花5
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_5), // 红桃5
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_5), // 黑桃5

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_6), // 方块6
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_6),    // 梅花6
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_6),   // 红桃6
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_6),   // 黑桃6

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_7), // 方块7
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_7),    // 梅花7
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_7),   // 红桃7
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_7),   // 黑桃7

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_8), // 方块8
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_8),    // 梅花8
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_8),   // 红桃8
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_8),   // 黑桃8

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_9), // 方块9
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_9),    // 梅花9
			},

			// 第三个人的牌
			{
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_9), // 红桃9
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_9), // 黑桃9

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_10), // 方块10
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_10),    // 梅花10
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_10),   // 红桃10
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_10),   // 黑桃10

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_J), // 方块J
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_J),    // 梅花J
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_J),   // 红桃J
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_J),   // 黑桃J

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_Q), // 方块Q
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_Q),    // 梅花Q
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_Q),   // 红桃Q
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_Q),   // 黑桃Q

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_K), // 方块K
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_K),    // 梅花K
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_K),   // 红桃K
			},
			{
				// 三张底牌
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_K),          // 红桃K
				uint32(room.PokerSuit_PS_NONE) + uint32(room.PokerValue_PV_BLACK_JOKER), // 小王
				uint32(room.PokerSuit_PS_NONE) + uint32(room.PokerValue_PV_RED_JOKER),   // 大王
			},
		},

		// 已废弃
		WallCards: []uint32{},

		// 已废弃
		HszDir: room.Direction_AntiClockWise,

		// 地主的座位号
		BankerSeat: 0,

		// 服务器地址
		ServerAddr: config.ServerAddr,

		// 客户端版本号
		ClientVer: config.ClientVersion,

		// 金币数(座位ID 与 金币 的map)
		PlayerSeatGold: map[int]uint64{0: 100000, 1: 100000, 2: 100000, 3: 100000},

		// 配牌游戏的名字
		PeiPaiGame: "doudizhu",
	}
}

// NewStartDDZGameParamsTest2 创建棋牌启动参数测试2
func NewStartDDZGameParamsTest2() structs.StartPukeGameParams {
	return structs.StartPukeGameParams{
		GameID: room.GameId_GAMEID_DOUDIZHU, // 斗地主
		Cards: [][]uint32{
			// 第一个人的牌(地主)
			{
				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_3), // 方块3
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_3),    // 梅花3
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_3),   // 红桃3
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_3),   // 黑桃3

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_4), // 方块4
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_4),    // 梅花4
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_4),   // 红桃4
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_4),   // 黑桃4

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_5), // 方块5
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_5),    // 梅花5
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_5),   // 红桃5
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_5),   // 黑桃5

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_6), // 方块6
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_6),    // 梅花6
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_6),   // 红桃6
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_6),   // 黑桃6

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_7), // 方块7
			},

			// 第二个人的牌
			{
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_7),  // 梅花7
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_7), // 红桃7
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_7), // 黑桃7

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_8), // 方块8
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_8),    // 梅花8
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_8),   // 红桃8
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_8),   // 黑桃8

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_9), // 方块9
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_9),    // 梅花9
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_9),   // 红桃9
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_9),   // 黑桃9

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_10), // 方块10
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_10),    // 梅花10
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_10),   // 红桃10
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_10),   // 黑桃10

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_J), // 方块J
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_J),    // 梅花J
			},

			// 第三个人的牌
			{
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_J), // 红桃J
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_J), // 黑桃J

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_Q), // 方块Q
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_Q),    // 梅花Q
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_Q),   // 红桃Q
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_Q),   // 黑桃Q

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_K), // 方块K
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_K),    // 梅花K
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_K),   // 红桃K
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_K),   // 红桃K

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_A), // 方块A
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_A),    // 梅花A
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_A),   // 红桃A
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_A),   // 黑桃A

				uint32(room.PokerSuit_PS_DIAMOND) + uint32(room.PokerValue_PV_2), // 方块2
				uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_2),    // 梅花2
				uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_2),   // 红桃2
			},
			{
				// 三张底牌
				uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_2),          // 黑桃2
				uint32(room.PokerSuit_PS_NONE) + uint32(room.PokerValue_PV_BLACK_JOKER), // 小王
				uint32(room.PokerSuit_PS_NONE) + uint32(room.PokerValue_PV_RED_JOKER),   // 大王
			},
		},

		// 已废弃
		WallCards: []uint32{},

		// 已废弃
		HszDir: room.Direction_AntiClockWise,

		// 地主的座位号
		BankerSeat: 0,

		// 服务器地址
		ServerAddr: config.ServerAddr,

		// 客户端版本号
		ClientVer: config.ClientVersion,

		// 金币数(座位ID 与 金币 的map)
		PlayerSeatGold: map[int]uint64{0: 100000, 1: 100000, 2: 100000, 3: 100000},

		// 配牌游戏的名字
		PeiPaiGame: "doudizhu",
	}
}

// PeipaiTest2 配牌测试2
func PeipaiTest2(t *testing.T) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "peipai.go:PeipaiTest1()",
	})

	// 启动参数2
	startParams := NewStartDDZGameParamsTest2()

	// 通知服务器：配牌
	if err := utils.Peipai(startParams.PeiPaiGame, startParams.Cards, startParams.WallCards, startParams.HszDir, startParams.BankerSeat); err != nil {
		logEntry.Error(err)
		return err
	}

	logEntry.Info("peipai.go:TestPeipai() 正常结束")
	return nil
}
