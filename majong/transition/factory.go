package transition

import (
	"errors"
	"steve/majong/interfaces"

	majongpb "steve/entity/majong"

	"github.com/Sirupsen/logrus"
)

const (
	transitionFile = "./transition.yaml"
)

var gOriginTransitionMap map[int]transition

type transitionMap map[majongpb.StateID](map[majongpb.EventID][]majongpb.StateID)

var gTransitionMaps map[int]transitionMap

type factory struct{}

// NewFactory 创建工厂
func NewFactory() interfaces.TransitionValidatorFactory {
	return new(factory)
}

func (f *factory) CreateTransitionValidator(gameID int) interfaces.TransitionValidator {
	return &transitionValidator{
		transitionMap: gTransitionMaps[gameID],
	}
}

type transitionValidator struct {
	transitionMap transitionMap
}

func (t *transitionValidator) Valid(oldState majongpb.StateID, newState majongpb.StateID, event majongpb.EventID, f interfaces.MajongFlow) error {
	eventStateMap, exists := t.transitionMap[oldState]
	if !exists {
		return errors.New("不能存在转换关系")
	}

	stateIDs, exist := eventStateMap[event]
	if !exist {
		return errors.New("不能存在转换关系")
	}

	for _, stateID := range stateIDs {
		if stateID == newState {
			return nil
		}
	}
	return errors.New("不能存在转换关系")
}

func init() {
	gOriginTransitionMap := make(map[int]transition, 20)
	tt, err := loadTransition(transitionFile)
	if err != nil {
		logrus.WithError(err).Panic("加载转换表配置失败")
	}
	for _, t := range tt {
		gOriginTransitionMap[t.GameID] = t
	}
	logrus.WithField("count", len(gOriginTransitionMap)).Info("转换表加载完成")

	gTransitionMaps = make(map[int]transitionMap)
	for gameID, transition := range gOriginTransitionMap {
		gTransitionMaps[gameID], err = originTransitionToMap(&transition)
		if err != nil {
			logrus.WithError(err).Panic("转换转换表错误")
		}
	}
}
