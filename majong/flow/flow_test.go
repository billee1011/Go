package flow

import (
	"errors"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	OLDSTATEID = majongpb.StateID(1)
	NEWSTATEID = majongpb.StateID(2)
	GAMEID     = 1
	EVENTID    = majongpb.EventID(1)
)

var initMajongContext = majongpb.MajongContext{
	CurState: OLDSTATEID,
	GameId:   GAMEID,
}

// Test_flow_ProcessEvent 测试收到事件后， 状态发生转换， 并且期间无错误发生
func Test_flow_ProcessEvent(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	// 创建一些 mock 对象
	mStateFactory := interfaces.NewMockMajongStateFactory(mc)
	mTransitionValidator := interfaces.NewMockTransitionValidator(mc)
	mOldState := interfaces.NewMockMajongState(mc)
	mNewState := interfaces.NewMockMajongState(mc)

	// 状态工厂的创建状态函数被调用时返回对应的 mock 状态对象
	mStateFactory.EXPECT().CreateState(GAMEID, OLDSTATEID).Return(mOldState).AnyTimes()
	mStateFactory.EXPECT().CreateState(GAMEID, NEWSTATEID).Return(mNewState).AnyTimes()

	f := flow{
		context:             initMajongContext,
		stateFactory:        mStateFactory,
		transitionValidator: mTransitionValidator,
	}

	// 旧状态的 ProcessEvent 方法会被调用 1 次
	mOldState.EXPECT().ProcessEvent(EVENTID, gomock.Any(), &f).Return(NEWSTATEID, nil).Times(1)

	// 验证新状态和旧状态是否存在转换
	mTransitionValidator.EXPECT().Valid(OLDSTATEID, NEWSTATEID, EVENTID, &f).Return(nil).Times(1)

	// 旧状态的 OnExit 方法会被调用 1 次
	mOldState.EXPECT().OnExit(&f).Times(1)

	// 新状态的 OnEntry 方法会被调用 1 次
	mNewState.EXPECT().OnEntry(&f).Times(1)

	// 处理事件
	assert.Nil(t, f.ProcessEvent(EVENTID, nil))

	// 状态 ID 变成了新的状态
	assert.Equal(t, f.context.CurState, NEWSTATEID)
}

// Test_flow_ProcessEvent_NoStateChange 测试当前状态和新状态相同时， OnExit 和 OnEntry 不会被调用
func Test_flow_ProcessEvent_NoStateChange(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	// 创建一些 mock 对象
	mStateFactory := interfaces.NewMockMajongStateFactory(mc)
	mTransitionValidator := interfaces.NewMockTransitionValidator(mc)
	mOldState := interfaces.NewMockMajongState(mc)
	mNewState := mOldState

	// 状态工厂的创建状态函数被调用时返回对应的 mock 状态对象
	mStateFactory.EXPECT().CreateState(GAMEID, OLDSTATEID).Return(mOldState).MinTimes(1)

	f := flow{
		context:             initMajongContext,
		stateFactory:        mStateFactory,
		transitionValidator: mTransitionValidator,
	}

	// 旧状态的 ProcessEvent 方法会被调用 1 次
	mOldState.EXPECT().ProcessEvent(EVENTID, gomock.Any(), &f).Return(OLDSTATEID, nil).Times(1)

	// 旧状态的 OnExit 方法会被调用 0 次
	mOldState.EXPECT().OnExit(&f).Times(0)

	// 新状态的 OnEntry 方法会被调用 0 次
	mNewState.EXPECT().OnEntry(&f).Times(0)

	// 处理事件
	assert.Nil(t, f.ProcessEvent(EVENTID, nil))

	// 状态 ID 保留为旧的
	assert.Equal(t, f.context.CurState, OLDSTATEID)
}

// Test_flow_ProcessEvent_CreateStateFail 测试创建当前状态对象失败时 ProcessEvent 的行为
func Test_flow_ProcessEvent_CreateStateFail(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	mStateFactory := interfaces.NewMockMajongStateFactory(mc)
	mTransitionValidator := interfaces.NewMockTransitionValidator(mc)

	// 创建当前状态返回 nil
	mStateFactory.EXPECT().CreateState(GAMEID, OLDSTATEID).Return(nil).MinTimes(1)

	// 不会创建新状态对象
	mStateFactory.EXPECT().CreateState(GAMEID, NEWSTATEID).MaxTimes(0)

	f := flow{
		context:             initMajongContext,
		stateFactory:        mStateFactory,
		transitionValidator: mTransitionValidator,
	}

	// 状态对象创建失败
	assert.Equal(t, errCreateState, f.ProcessEvent(EVENTID, nil))
}

// Test_flow_ProcessEvent_StateProcessFail 测试当前状态对象处理事件失败时的行为
func Test_flow_ProcessEvent_StateProcessFail(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	// 创建一些 mock 对象
	mStateFactory := interfaces.NewMockMajongStateFactory(mc)
	mTransitionValidator := interfaces.NewMockTransitionValidator(mc)
	mOldState := interfaces.NewMockMajongState(mc)
	mNewState := interfaces.NewMockMajongState(mc)

	// 状态工厂的创建状态函数被调用时返回对应的 mock 状态对象
	mStateFactory.EXPECT().CreateState(GAMEID, OLDSTATEID).Return(mOldState).MinTimes(1)

	// 不会创建新状态对象
	mStateFactory.EXPECT().CreateState(GAMEID, NEWSTATEID).Return(mNewState).MaxTimes(0)

	f := flow{
		context:             initMajongContext,
		stateFactory:        mStateFactory,
		transitionValidator: mTransitionValidator,
	}

	// 当前状态处理事件返回错误
	mOldState.EXPECT().ProcessEvent(EVENTID, nil, &f).Times(1).Return(NEWSTATEID, errors.New("some error"))

	// 状态对象创建失败
	assert.Equal(t, errStateProcess, f.ProcessEvent(EVENTID, nil))
}

// Test_flow_ProcessEvent_NewStateCreateFail 测试新状态创建失败的情况
func Test_flow_ProcessEvent_NewStateCreateFail(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	// 创建一些 mock 对象
	mStateFactory := interfaces.NewMockMajongStateFactory(mc)
	mTransitionValidator := interfaces.NewMockTransitionValidator(mc)
	mOldState := interfaces.NewMockMajongState(mc)

	// 状态工厂的创建状态函数被调用时返回对应的 mock 状态对象
	mStateFactory.EXPECT().CreateState(GAMEID, OLDSTATEID).Return(mOldState).AnyTimes()
	mStateFactory.EXPECT().CreateState(GAMEID, NEWSTATEID).Return(nil).AnyTimes()

	f := flow{
		context:             initMajongContext,
		stateFactory:        mStateFactory,
		transitionValidator: mTransitionValidator,
	}

	// 旧状态的 ProcessEvent 方法会被调用 1 次
	mOldState.EXPECT().ProcessEvent(EVENTID, gomock.Any(), &f).Return(NEWSTATEID, nil).Times(1)

	// 验证新状态和旧状态是否存在转换
	mTransitionValidator.EXPECT().Valid(OLDSTATEID, NEWSTATEID, EVENTID, &f).Return(nil).Times(1)

	// 处理事件
	assert.Equal(t, errCreateState, f.ProcessEvent(EVENTID, nil))
}
