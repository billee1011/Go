package scxlai

import (
	"steve/room/interfaces/global"
	"steve/server_pb/majong"
)

// 注册 AI
func init() {
	g := global.GetDeskAutoEventGenerator()
	// 血流
	g.RegisterAI(scxlGameID, majong.StateID_state_dingque, &dingqueStateAI{})
	g.RegisterAI(scxlGameID, majong.StateID_state_huansanzhang, &huansanzhangStateAI{})
	g.RegisterAI(scxlGameID, majong.StateID_state_chupaiwenxun, &chupaiWenxunStateAI{})
	g.RegisterAI(scxlGameID, majong.StateID_state_zixun, &zixunStateAI{})
	g.RegisterAI(scxlGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
	// 血战
	g.RegisterAI(scxzGameID, majong.StateID_state_dingque, &dingqueStateAI{})
	g.RegisterAI(scxzGameID, majong.StateID_state_huansanzhang, &huansanzhangStateAI{})
	g.RegisterAI(scxzGameID, majong.StateID_state_chupaiwenxun, &chupaiWenxunStateAI{})
	g.RegisterAI(scxzGameID, majong.StateID_state_zixun, &zixunStateAI{})
	g.RegisterAI(scxzGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
	// 二人
	g.RegisterAI(ermjGameID, majong.StateID_state_chupaiwenxun, &chupaiWenxunStateAI{})
	g.RegisterAI(ermjGameID, majong.StateID_state_zixun, &zixunStateAI{})
	g.RegisterAI(ermjGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
}
