package initial

import (
	"errors"
	server_pb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

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

var errCreateEmptyContextFailed = errors.New("创建空的麻将现场失败")
var errInvalidParam = errors.New("参数错误")

// InitMajongContext 初始化麻将现场
func InitMajongContext(param server_pb.InitMajongContextParams) (mjContext server_pb.MajongContext, err error) {
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

	if err = proto.Unmarshal([]byte{}, &mjContext); err != nil {
		logEntry.WithError(err).Errorln(errCreateEmptyContextFailed)
		err = errCreateEmptyContextFailed
		return
	}
	mjContext.GameId = param.GetGameId()
	mjContext.CurState = server_pb.StateID_state_init
	mjContext.Players = initPlayers(param.GetPlayers())
	mjContext.ActivePlayer = param.GetPlayers()[param.GetZhuangIndex()]
	mjContext.ZhuangjiaIndex = param.GetZhuangIndex()
	mjContext.FixZhuangjiaIndex = param.GetFixZhuangIndex()

	mjContext.Option = param.GetOption()
	mjContext.MajongOption = param.GetMajongOption()

	return
}
