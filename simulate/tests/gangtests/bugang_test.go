package gangtests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuGang(t *testing.T) {
	param := global.NewCommonStartGameParams()
	param.BankerSeat = 0
	param.WallCards = []*room.Card{}
	param.Cards[0][4] = &global.Card6W
	param.Cards[0][5] = &global.Card6W
	param.Cards[0][6] = &global.Card6W
	param.Cards[1][4] = &global.Card2W
	param.Cards[1][5] = &global.Card2W
	param.Cards[1][6] = &global.Card2W
	param.WallCards = []*room.Card{&global.Card8W, &global.Card8W, &global.Card8W, &global.Card9W, &global.Card3B}
	deskData, err := utils.StartGame(param)
	assert.Nil(t, err)
	utils.WaitZixunNtf(deskData, deskData.BankerSeat)
	//庄家出牌,出12
	utils.SendChupaiReq(deskData, deskData.BankerSeat, uint32(12))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(12), deskData.BankerSeat)
	//下家请求碰12
	xjSeat := (deskData.BankerSeat + 1) % len(deskData.Players)
	utils.SendPengReq(deskData, xjSeat)
	//检查碰的通知
	utils.CheckPengNotify(t, deskData, xjSeat, 12)
	//碰成功后收到自询通知
	utils.CheckZixunNotify(t, deskData, xjSeat)
	//下家出牌请求
	utils.SendChupaiReq(deskData, xjSeat, 17)
	//下家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 17, xjSeat)
	//对家摸牌(自寻)响应
	djSeat := (xjSeat + 1) % len(deskData.Players)
	utils.CheckMoPaiNotify(t, deskData, djSeat, 18)
	//对家出牌请求
	utils.SendChupaiReq(deskData, djSeat, 22)
	//对家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 22, djSeat)
	//尾家摸牌(自寻)响应
	wjSeat := (djSeat + 1) % len(deskData.Players)
	utils.CheckMoPaiNotify(t, deskData, wjSeat, 18)
	//尾家出牌请求
	utils.SendChupaiReq(deskData, wjSeat, 27)
	//尾家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 27, wjSeat)
	//庄家摸牌(自寻)响应
	utils.CheckMoPaiNotify(t, deskData, deskData.BankerSeat, 18)
	//庄家出牌请求
	utils.SendChupaiReq(deskData, deskData.BankerSeat, 16)
	//庄家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 16, deskData.BankerSeat)
	//下家摸牌(自寻)响应
	utils.CheckMoPaiNotify(t, deskData, xjSeat, 19)
	//下家补杠
	utils.SendGangReq(deskData, xjSeat, 12, room.GangType_BuGang)
	//下家补杠响应
	player := utils.GetDeskPlayerBySeat(xjSeat, deskData)
	utils.CheckGangNotify(t, deskData, player.Player.GetID(), player.Player.GetID(), 12, room.GangType_BuGang)

}
