package game

import (
	mahjonginit "steve/majong/export/initial"
	"steve/room/peipai/handle"
	mahjong "steve/server_pb/majong"
)

func InitGameContext(gameID int) interface{} {
	switch gameID {
	case 1:
		return initMahjongContext
	default:
		return nil
	}
}

func initMahjongContext(playerIDs []uint64, gameID int) (mahjongContext mahjong.MajongContext, err error) {
	param := mahjong.InitMajongContextParams{
		GameId:  int32(gameID),
		Players: playerIDs, // TODO 需要携带位置信息
		Option: &mahjong.MajongCommonOption{
			MaxFapaiCartoonTime:        10 * 1000,
			MaxHuansanzhangCartoonTime: 10 * 1000,
			HasHuansanzhang:            handle.GetHsz(gameID),                     //设置玩家是否开启换三张
			Cards:                      handle.GetPeiPai(gameID),                  //设置是否配置墙牌
			WallcardsLength:            uint32(handle.GetLensOfWallCards(gameID)), //设置墙牌长度
			HszFx: &mahjong.Huansanzhangfx{
				NeedDeployFx:   handle.GetHSZFangXiang(gameID) != -1,
				HuansanzhangFx: int32(handle.GetHSZFangXiang(gameID)),
			}, //设置换三张方向
			Zhuang: &mahjong.Zhuang{
				NeedDeployZhuang: handle.GetZhuangIndex(gameID) != -1,
				ZhuangIndex:      int32(handle.GetZhuangIndex(gameID)),
			},
		}, //设置庄家
		MajongOption: []byte{},
	}

	if mahjongContext, err = mahjonginit.InitMajongContext(param); err != nil {
		return
	}

	//if err := fillContextOptions(gameID, &mahjongContext); err != nil {
	//	return
	//}

	return
}

func HandleEvent(gameID int, eventID mahjong.EventID, eventContext []byte) {
	switch gameID {
	case 1:

	default:

	}
}
