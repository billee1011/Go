package fantests

// //shibaluohan 共同步骤
// // 玩家换三张后的牌
// //庄家0手牌 11, 11, 11, 11, 22, 22, 22, 22, 13, 13, 13, 13, 14, 14
// //1玩家手牌 18, 18, 18, 17, 15, 19, 15, 16, 16, 17, 17, 18, 19
// //2玩家手牌 31, 31, 31, 32, 32, 32, 33, 33, 33, 34, 34, 34, 35
// //3玩家手牌 26, 26, 26, 27, 27, 27, 28, 28, 28, 29, 29, 36, 36
// // 庄家暗杠11,摸14，暗杠22,摸14,暗杠13,摸15,暗杠14,摸36,出36
// // 结果牌型：
// // 庄家：gang{11,22,13,14},handcard{15}
// func shibaluohan(t *testing.T) *utils.DeskData {
// 	params := global.NewCommonStartGameParams()
// 	params.BankerSeat = 0
// 	params.Cards = [][]uint32{
// 		{31, 31, 31, 11, 22, 22, 22, 22, 13, 13, 13, 13, 14, 14},
// 		{26, 26, 26, 17, 15, 19, 15, 16, 16, 17, 17, 18, 19},
// 		{11, 11, 11, 32, 32, 32, 33, 33, 33, 34, 34, 34, 35},
// 		{18, 18, 18, 27, 27, 27, 28, 28, 28, 29, 29, 36, 36},
// 	}
// 	params.WallCards = []uint32{14, 14, 15, 36, 15, 36}
// 	// 对家换牌
// 	params.HszDir = room.Direction_Opposite
// 	params.HszCards = [][]uint32{
// 		{31, 31, 31},
// 		{26, 26, 26},
// 		{11, 11, 11},
// 		{18, 18, 18},
// 	}
// 	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TIAO}
// 	deskData, err2 := utils.StartGame(params)
// 	assert.NotNil(t, deskData)
// 	assert.Nil(t, err2)

// 	//开局 0 自询能暗杠
// 	utils.CheckZixunNtf(t, deskData, 0, false, true, true)
// 	// 0 暗杠11
// 	utils.SendGangReq(deskData, 0, 11, room.GangType_AnGang)
// 	// 暗杠立即结算6分
// 	utils.CheckInstantSettleScoreNotify(t, deskData, 0, 6)
// 	// 0 自询能暗杠
// 	utils.CheckZixunNtf(t, deskData, 0, false, true, false)
// 	// 0 暗杠22
// 	utils.SendGangReq(deskData, 0, 22, room.GangType_AnGang)
// 	// 暗杠立即结算6分
// 	utils.CheckInstantSettleScoreNotify(t, deskData, 0, 6)
// 	// 0 自询能暗杠
// 	utils.CheckZixunNtf(t, deskData, 0, false, true, false)
// 	// 0 暗杠13
// 	utils.SendGangReq(deskData, 0, 13, room.GangType_AnGang)
// 	// 暗杠立即结算6分
// 	utils.CheckInstantSettleScoreNotify(t, deskData, 0, 6)
// 	// 0 自询能暗杠
// 	utils.CheckZixunNtf(t, deskData, 0, false, true, false)
// 	// 0 暗杠14
// 	utils.SendGangReq(deskData, 0, 14, room.GangType_AnGang)
// 	// 暗杠立即结算6分
// 	utils.CheckInstantSettleScoreNotify(t, deskData, 0, 6)
// 	// 0 自询
// 	utils.CheckZixunNtf(t, deskData, 0, false, false, false)
// 	// 0 出牌 36
// 	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 36))
// 	// 3玩家能碰36
// 	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 3, true, false, false))
// 	return deskData
// }

// //TestFan_Shibaluohan_Zimo 十八罗汉立即结算自摸测试
// // 3碰36,出29,庄摸牌15，自摸15
// // 期望总赢分 384 = 64 *2 *3
// func TestFan_Shibaluohan_Zimo(t *testing.T) {
// 	deskData := shibaluohan(t)
// 	// 3发送碰
// 	assert.Nil(t, utils.SendPengReq(deskData, 3))
// 	// 3玩家碰牌成功进入自询
// 	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
// 	// 3 出牌 29
// 	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 29))
// 	//0 号玩家摸牌15后 检测自询,能自摸
// 	utils.CheckZixunNtf(t, deskData, 0, false, false, true)
// 	// 0玩家发送胡,自摸胡15
// 	utils.SendHuReq(deskData, 0)
// 	// 0 号玩家发送胡请求
// 	assert.Nil(t, utils.SendHuReq(deskData, 0))

// 	// 检测所有玩家收到自摸通知
// 	utils.CheckHuNotify(t, deskData, []int{0}, 0, 15, room.HuType_HT_ZIMO)

// 	// 检测十八罗汉自摸分数,十八罗汉64倍*自摸2倍
// 	winScro := 64 * 2 * (len(deskData.Players) - 1)
// 	utils.CheckInstantSettleScoreNotify(t, deskData, 0, int64(winScro))
// }

// //TestFan_Shibaluohan_Dianpao 十八罗汉立即点炮自摸测试
// // 3弃碰36,1摸牌15，出15,庄点炮15
// // 期望总赢分 64
// func TestFan_Shibaluohan_Dianpao(t *testing.T) {
// 	deskData := shibaluohan(t)
// 	// 3玩家发送弃碰,胡36
// 	utils.SendQiReq(deskData, 3)
// 	// 1 号玩家摸牌15后 检测自询
// 	utils.CheckZixunNtf(t, deskData, 1, false, true, false)
// 	// 1 出 15
// 	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 15))
// 	//0玩家能点炮胡15
// 	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, false, true, false))
// 	// 0 号玩家发送胡请求
// 	assert.Nil(t, utils.SendHuReq(deskData, 0))

// 	// 检测所有玩家收到点炮通知x
// 	utils.CheckHuNotify(t, deskData, []int{0}, 1, 15, room.HuType_HT_DIANPAO)
// 	// 检测十八罗汉点炮分数
// 	utils.CheckInstantSettleScoreNotify(t, deskData, 0, 64)
// }
