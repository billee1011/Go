package scxlai

import (
	"steve/room/interfaces/global"
	"steve/server_pb/majong"
)

// 注册 AI
func init() {
	g := global.GetDeskAutoEventGenerator()
	// 血流
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_dingque), &dingqueStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_huansanzhang), &huansanzhangStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	// 血战
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_dingque), &dingqueStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_huansanzhang), &huansanzhangStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
	// 二人
	g.RegisterAI(ermjGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	g.RegisterAI(ermjGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	g.RegisterAI(ermjGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
}
