package contexts

import (
	"errors"
	"steve/common/mjoption"
	server_pb "steve/entity/majong"
	"time"

	"steve/room/common"
	"steve/room/fixed"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// MajongDeskContext 麻将牌桌现场
type MajongDeskContext struct {
	MjContext   server_pb.MajongContext // 牌局现场
	StateNumber int                     // 状态序号
	StateTime   time.Time               // 状态时间
}

// SetStateNumber 更新 stateNumber
func (c *MajongDeskContext) SetStateNumber(state int) {
	c.StateNumber = state
}

var errInitMajongContext = errors.New("初始化麻将现场失败")
var errAllocDeskIDFailed = errors.New("分配牌桌 ID 失败")
var errPlayerNotExist = errors.New("玩家不存在")
var errPlayerNeedXingPai = errors.New("玩家需要参与行牌")

// CreateMajongContext 创建麻将现场
func CreateMajongContext(players []uint64, gameID int, zhuang uint32, fixZhuang bool) (*MajongDeskContext, error) {
	if !fixZhuang {
		zhuang = 0
	}
	param := server_pb.InitMajongContextParams{
		GameId:  int32(gameID),
		Players: players,
		Option: &server_pb.MajongCommonOption{
			MaxFapaiCartoonTime:        uint32(viper.GetInt(fixed.MaxFapaiCartoonTime)),
			MaxHuansanzhangCartoonTime: uint32(viper.GetInt(fixed.MaxHuansanzhangCartoonTime)),
			HasHuansanzhang:            common.GetHsz(gameID),                     //设置玩家是否开启换三张
			Cards:                      common.GetPeiPai(gameID),                  //设置是否配置墙牌
			WallcardsLength:            uint32(common.GetLensOfWallCards(gameID)), //设置墙牌长度
			HszFx: &server_pb.Huansanzhangfx{
				NeedDeployFx:   common.GetHSZFangXiang(gameID) != -1,
				HuansanzhangFx: int32(common.GetHSZFangXiang(gameID)),
			}, //设置换三张方向
		},
		MajongOption:   []byte{},
		ZhuangIndex:    zhuang,
		FixZhuangIndex: true,
	}
	var mjContext server_pb.MajongContext
	var err error
	if mjContext, err = initMajongContext(param); err != nil {
		return nil, err
	}
	if err := fillContextOptions(gameID, &mjContext); err != nil {
		return nil, err
	}
	result := &MajongDeskContext{
		MjContext: mjContext,
		//StateNumber: 0,
		StateTime: time.Now(),
	}
	result.SetStateNumber(0)
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
	mjContext.TempData = new(server_pb.TempDatas) //初始化临时数据

	return
}

func initPlayers(players []uint64) []*server_pb.Player {
	result := []*server_pb.Player{}
	for _, playerID := range players {
		result = append(result, &server_pb.Player{
			PlayerId:          playerID,
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
			TingStateInfo:     &server_pb.TingStateInfo{},
		})
	}
	return result
}

var errNoGameOption = errors.New("没有该游戏的游戏选项")

// fillContextOptions 填充麻将现场的 options
func fillContextOptions(gameID int, mjContext *server_pb.MajongContext) error {
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
