package scxlai

import (
	"steve/room2/ai"
	"steve/server_pb/majong"
)

// 注册 AI
func init() {
	// 血流
	ai.GetAtEvent().RegisterAI(scxlGameID, int32(majong.StateID_state_dingque), &dingqueStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, int32(majong.StateID_state_huansanzhang), &huansanzhangStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
	// 血战
	ai.GetAtEvent().RegisterAI(scxzGameID, int32(majong.StateID_state_dingque), &dingqueStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, int32(majong.StateID_state_huansanzhang), &huansanzhangStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
	// 二人
	ai.GetAtEvent().RegisterAI(ermjGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	ai.GetAtEvent().RegisterAI(ermjGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	ai.GetAtEvent().RegisterAI(ermjGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
}
