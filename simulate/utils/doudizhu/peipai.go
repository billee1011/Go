package doudizhu

import (
	"steve/client_pb/common"
	"steve/simulate/structs"
	"steve/simulate/utils"
	"testing"

	"github.com/Sirupsen/logrus"
)

// NewStartDDZGameParamsTest1 创建棋牌启动参数测试
func NewStartDDZGameParamsTest1() structs.StartPukeGameParams {
	return structs.StartPukeGameParams{
		GameID: common.GameId_GAMEID_DOUDIZHU, // 斗地主
		Cards: [][]uint32{
			// 第一个人的牌(地主)
			{
				0x11, 0x21, 0x31, 0x41,
				0x12, 0x22, 0x32, 0x42,
				0x13, 0x23, 0x33, 0x43,
				0x14, 0x24, 0x34, 0x44,
				0x15,
			},

			// 第二个人的牌
			{
				0x25, 0x35, 0x45,
				0x16, 0x26, 0x36, 0x46,
				0x17, 0x27, 0x37, 0x47,
				0x18, 0x28, 0x38, 0x48,
				0x19, 0x29,
			},

			// 第三个人的牌
			{
				0x39, 0x49,
				0x1A, 0x2A, 0x3A, 0x4A,
				0x1B, 0x2B, 0x3B, 0x4B,
				0x1C, 0x2C, 0x3C, 0x4C,
				0x1D, 0x2D, 0x3D,
			},
			{
				// 三张底牌
				0x4D, 0x0E, 0x0F,
			},
		},

		// 金币数(座位ID 与 金币 的map)
		PlayerSeatGold: map[int]uint64{0: 10000, 1: 5000, 2: 8000},

		// 配牌游戏的名字
		PeiPaiGame: "doudizhu",
	}
}

// NewStartDDZGameParamsTest2 创建棋牌启动参数测试2
func NewStartDDZGameParamsTest2() structs.StartPukeGameParams {
	return structs.StartPukeGameParams{
		GameID: common.GameId_GAMEID_DOUDIZHU, // 斗地主
		Cards: [][]uint32{
			// 第一个人的牌(地主)
			{
				0x11, 0x21, 0x31, 0x41,
				0x12, 0x22, 0x32, 0x42,
				0x13, 0x23, 0x33, 0x43,
				0x14, 0x24, 0x34, 0x44,
				0x15,
			},

			// 第二个人的牌
			{
				0x25, 0x35, 0x45,
				0x16, 0x26, 0x36, 0x46,
				0x17, 0x27, 0x37, 0x47,
				0x18, 0x28, 0x38, 0x48,
				0x19, 0x29,
			},

			// 第三个人的牌
			{
				0x39, 0x49,
				0x1A, 0x2A, 0x3A, 0x4A,
				0x1B, 0x2B, 0x3B, 0x4B,
				0x1C, 0x2C, 0x3C, 0x4C,
				0x1D, 0x2D, 0x3D,
			},
			{
				// 三张底牌
				0x4D, 0x0E, 0x0F,
			},
		},

		// 金币数(座位ID 与 金币 的map)
		PlayerSeatGold: map[int]uint64{0: 10000, 1: 10000, 2: 10000},

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
	startParams := NewStartDDZGameParamsTest1()

	// 通知服务器：配牌
	if err := utils.Peipai(startParams.PeiPaiGame, startParams.Cards, startParams.WallCards, startParams.HszDir, startParams.BankerSeat); err != nil {
		logEntry.Error(err)
		return err
	}

	logEntry.Info("peipai.go:TestPeipai() 正常结束")
	return nil
}
