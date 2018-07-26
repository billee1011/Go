package contexts

import (
	"errors"
	"steve/common/mjoption"
	"steve/entity/majong"
	server_pb "steve/entity/majong"
	"steve/room2/common"
	"steve/room2/fixed"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var errInitMajongContext = errors.New("初始化麻将现场失败")
var errAllocDeskIDFailed = errors.New("分配牌桌 ID 失败")
var errPlayerNotExist = errors.New("玩家不存在")
var errPlayerNeedXingPai = errors.New("玩家需要参与行牌")

func CreateMajongContext(players []uint64, gameId int) (*MjContext, error) {
	param := server_pb.InitMajongContextParams{
		GameId:  int32(gameId),
		Players: players,
		Option: &server_pb.MajongCommonOption{
			MaxFapaiCartoonTime:        uint32(viper.GetInt(fixed.MaxFapaiCartoonTime)),
			MaxHuansanzhangCartoonTime: uint32(viper.GetInt(fixed.MaxHuansanzhangCartoonTime)),
			HasHuansanzhang:            common.GetHsz(gameId),                     //设置玩家是否开启换三张
			Cards:                      common.GetPeiPai(gameId),                  //设置是否配置墙牌
			WallcardsLength:            uint32(common.GetLensOfWallCards(gameId)), //设置墙牌长度
			HszFx: &server_pb.Huansanzhangfx{
				NeedDeployFx:   common.GetHSZFangXiang(gameId) != -1,
				HuansanzhangFx: int32(common.GetHSZFangXiang(gameId)),
			}, //设置换三张方向
			/*Zhuang: &server_pb.Zhuang{
				NeedDeployZhuang: util.GetZhuangIndex(desk.GetGameId()) != -1,
				ZhuangIndex:      int32(util.GetZhuangIndex(desk.GetGameId())),
			},*/
		}, //设置庄家
		MajongOption: []byte{},
	}
	var mjContext server_pb.MajongContext
	var err error
	if mjContext, err = initMajongContext(param); err != nil {
		return nil, err
	}
	if err := fillContextOptions(gameId, &mjContext); err != nil {
		return nil, err
	}
	result := &MjContext{
		MjContext:   mjContext,
		StateNumber: 0,
		StateTime:   time.Now(),
	}
	return result, nil
}

var errCreateEmptyContextFailed = errors.New("创建空的麻将现场失败")
var errInvalidParam = errors.New("参数错误")

// InitMajongContext 初始化麻将现场
func initMajongContext(param server_pb.InitMajongContextParams) (mjContext server_pb.MajongContext, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "InitMajongContext",
		"zhuangjia_index": param.GetZhuangIndex(),
		"players":         param.GetPlayers(),
	})
	if int(param.GetZhuangIndex()) >= len(param.GetPlayers()) {
		logEntry.Errorln(errInvalidParam)
		err = errInvalidParam
		return
	}

	mjContext.GameId = param.GetGameId()
	mjContext.CurState = server_pb.StateID_state_init
	mjContext.Players = initPlayers(param.GetPlayers())
	mjContext.ActivePlayer = param.GetPlayers()[param.GetZhuangIndex()]
	mjContext.ZhuangjiaIndex = param.GetZhuangIndex()
	// mjContext.FixZhuangjiaIndex = param.GetFixZhuangIndex()

	mjContext.Option = param.GetOption()
	mjContext.MajongOption = param.GetMajongOption()

	return
}

func initPlayers(players []uint64) []*server_pb.Player {
	result := []*server_pb.Player{}
	for _, playerID := range players {
		result = append(result, &server_pb.Player{
			PalyerId:          playerID,
			HandCards:         []*server_pb.Card{},
			OutCards:          []*server_pb.Card{},
			ChiCards:          []*server_pb.ChiCard{},
			PengCards:         []*server_pb.PengCard{},
			HuCards:           []*server_pb.HuCard{},
			GangCards:         []*server_pb.GangCard{},
			PossibleActions:   []server_pb.Action{},
			HasSelected:       false,
			HasDingque:        false,
			HuansanzhangSure:  false,
			HuansanzhangCards: []*server_pb.Card{},
			Properties:        make(map[string][]byte, 0),
		})
	}
	return result
}

var errNoGameOption = errors.New("没有该游戏的游戏选项")

// fillContextOptions 填充麻将现场的 options
func fillContextOptions(gameID int, mjContext *majong.MajongContext) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "fillContextOptions",
		"game_id":   gameID,
	})
	gameOption := mjoption.GetGameOptions(gameID)
	if gameOption == nil {
		entry.Errorln(errNoGameOption)
		return errNoGameOption
	}
	mjContext.SettleOptionId = uint32(gameOption.SettleOptionID)
	mjContext.CardtypeOptionId = uint32(gameOption.CardTypeOptionID)
	mjContext.XingpaiOptionId = uint32(gameOption.XingPaiOptionID)
	return nil
}
