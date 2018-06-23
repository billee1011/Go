package desks

import (
	"steve/gutils"
	"steve/room/interfaces/global"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

type optionxuezhan struct{}
type optionxueliu struct{}
type optiondefaul struct{}

// OptionHandle 选项处理接口
type OptionHandle interface {
	handleQuitByPlayerState(desk *desk, playerID uint64)
}

func (xl *optionxueliu) handleQuitByPlayerState(desk *desk, playerID uint64) {

}

func (xz *optionxuezhan) handleQuitByPlayerState(desk *desk, playerID uint64) {
	mjContext := desk.dContext.mjContext
	player := gutils.GetMajongPlayer(playerID, &mjContext)
	if player.GetXpState() != majong.XingPaiState_normal {
		deskMgr := global.GetDeskMgr()
		deskMgr.RemoveDeskPlayerByPlayerID(playerID)
	}
	logrus.WithFields(logrus.Fields{
		"funcName":    "handleQuitByPlayerState",
		"gameID":      mjContext.GetGameId(),
		"playerState": player.GetXpState(),
	}).Infof("玩家:%v退出后的相关处理", playerID)
}

func (xl *optiondefaul) handleQuitByPlayerState(desk *desk, playerID uint64) {

}

// GetOptionByFactory 根据游戏id从选项工厂拿具体
func GetOptionByFactory(gameID int) OptionHandle {
	switch gameID {
	case gutils.SCXLGameID:
		return &optionxueliu{}
	case gutils.SCXZGameID:
		return &optionxuezhan{}
	default:
		return &optiondefaul{}
	}
}
