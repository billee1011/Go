package scxlai

import (
	"steve/entity/majong"
	"steve/room2/ai"
)

// 注册 AI
func init() {
	// 血流
	ai.GetAtEvent().RegisterAI(scxlGameID, majong.StateID_state_dingque, &dingqueStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, majong.StateID_state_huansanzhang, &huansanzhangStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, majong.StateID_state_chupaiwenxun, &chupaiWenxunStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, majong.StateID_state_zixun, &zixunStateAI{})
	ai.GetAtEvent().RegisterAI(scxlGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
	// 血战
	ai.GetAtEvent().RegisterAI(scxzGameID, majong.StateID_state_dingque, &dingqueStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, majong.StateID_state_huansanzhang, &huansanzhangStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, majong.StateID_state_chupaiwenxun, &chupaiWenxunStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, majong.StateID_state_zixun, &zixunStateAI{})
	ai.GetAtEvent().RegisterAI(scxzGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
	// 二人
	ai.GetAtEvent().RegisterAI(ermjGameID, majong.StateID_state_chupaiwenxun, &chupaiWenxunStateAI{})
	ai.GetAtEvent().RegisterAI(ermjGameID, majong.StateID_state_zixun, &zixunStateAI{})
	ai.GetAtEvent().RegisterAI(ermjGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
}
