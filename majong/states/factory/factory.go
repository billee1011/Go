package factory
/*
 功能： 状态机工厂, 可以实现多种不同的状态机系统，比如麻将类状态机，扑克类状态机.common是通用状态实现目录，可以通过定义子游戏的目录来定义自己的状态实现目录.
 作者： SkyWang
 日期： 2018-7-17
 */
import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"steve/majong/bus"
)
// 状态机管理器列表: 在下面列表中，定义自己的特殊状态机。
var mapMgr  =  map[int32] NewStater {
	0 :&mjStateMgr{},		// 默认状态机(麻将状态机)
	//1: &mjStateMgr{},		// 其他游戏自定义状态机
}

// 状态机管理器接口
type NewStater interface {
	// 新建状态
	NewState(stateID majongpb.StateID) interfaces.MajongState
}

func init() {
	f := new(factory)
	bus.SetMajongStateFacotry(f)
}

type factory struct {

}

var _ interfaces.MajongStateFactory = new(factory)

func (f *factory) CreateState(gameID int32, stateID majongpb.StateID) interfaces.MajongState {
	// 使用自己的状态机
	if t, ok := mapMgr[gameID]; ok {
		return t.NewState(stateID)
	}
	// 使用通用的麻将状态机
	return mapMgr[0].NewState(stateID)
}

