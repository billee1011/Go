package settle

import (
	"steve/majong/settle/fan"
	"steve/majong/utils"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// DianPaoSettle 点炮胡的结算
type DianPaoSettle struct {
}

// SettleDianPaoHu  点炮胡立即结算,生成结算列表 winnersID 赢家id, loserID 输家id, settleType 结算类型， huType 胡牌类型
func (dianPaoSettle *DianPaoSettle) SettleDianPaoHu(context *majong.MajongContext, winnersID []uint64, loserID uint64, settleType majong.SettleType, huType majong.HuType) ([]*majong.SettleInfo, error) {
	entry := logrus.WithFields(logrus.Fields{
		"name":       "SettleDianPaoHu",
		"winnersID":  winnersID,
		"loserID":    loserID,
		"settleType": settleType,
		"huType":     huType,
	})

	settleInfos := make([]*majong.SettleInfo, 0)
	for i := 0; i < len(winnersID); i++ {
		winner := utils.GetPlayerByID(context.Players, winnersID[i])
		fansMap := make(map[string]uint32)
		gen := uint32(0)
		for i := 0; i < len(fan.ScxlFan); i++ {
			if fan.ScxlFan[i].Condition(*context, huType, winner) {
				fansMap[fan.ScxlFan[i].GetFanName()] = fan.ScxlFan[i].GetFanValue()
			}
		}
		fansMap, gen = scxlFanMutex(fansMap, fan.GetGenCount(winner))

		fanTotal := 1
		for _, value := range fansMap {
			if value != 0 {
				fanTotal = fanTotal * int(value)
			}
		}
		//底数
		ante := GetDi()
		total := int64(fanTotal) * ante * (1 << gen)
		// 结算信息
		settleInfo := NewSettleInfo(context, settleType)
		for _, player := range context.Players {
			if winner.PalyerId == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] + total
			} else if loserID == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] - total
			} else {
				settleInfo.Scores[player.PalyerId] = 0
			}
		}
		settleInfos = append(settleInfos, settleInfo)
	}

	entry.Info("点炮结算")
	return settleInfos, nil
}

func scxlFanMutex(fansMap map[string]uint32, gen uint32) (map[string]uint32, uint32) {
	shiBaLuoHan := fan.FanName[majong.Fan_ShiBaLuoHan]
	jingGouDiao := fan.FanName[majong.Fan_JingGouDiao]
	pengPengHu := fan.FanName[majong.Fan_PengPengHu]

	if value, ok := fansMap[shiBaLuoHan]; ok && value > 0 {
		if value, ok := fansMap[jingGouDiao]; ok && value > 0 { //十八罗汉跟金钩钓互斥
			fansMap[fan.FanName[majong.Fan_JingGouDiao]] = 0
		}
		if value, ok := fansMap[pengPengHu]; ok && value > 0 { //十八罗汉跟碰碰胡互斥
			fansMap[fan.FanName[majong.Fan_PengPengHu]] = 0
		}
		gen = 0
	}

	qingYiSe := fan.FanName[majong.Fan_QingYiSe]
	qiDui := fan.FanName[majong.Fan_QiDui]
	longQiDui := fan.FanName[majong.Fan_LongQiDui]
	qinglongQiDui := fan.FanName[majong.Fan_QingLongQiDui]
	if value, ok := fansMap[qingYiSe]; ok && value > 0 {
		flag := false
		if value, ok := fansMap[shiBaLuoHan]; ok && value > 0 { // 添加清十八罗汉,移除十八罗汉
			qingShiBaLuoHan := fan.FanName[majong.Fan_QingShiBaLuoHan]
			fansMap[qingShiBaLuoHan] = fan.FanValue[majong.Fan_QingShiBaLuoHan]
			fansMap[shiBaLuoHan] = 0
			flag = true
		}
		if value, ok := fansMap[jingGouDiao]; ok && value > 0 { // 添加清金钩钓,移除金钩钓
			qingJingGouDiao := fan.FanName[majong.Fan_QingJingGouDiao]
			fansMap[qingJingGouDiao] = fan.FanValue[majong.Fan_QingJingGouDiao]
			fansMap[jingGouDiao] = 0
			flag = true
		}
		if value, ok := fansMap[pengPengHu]; ok && value > 0 { // 添加清碰,移除碰碰胡
			qingPeng := fan.FanName[majong.Fan_QingPeng]
			fansMap[qingPeng] = fan.FanValue[majong.Fan_QingPeng]
			fansMap[pengPengHu] = 0
			flag = true
		}
		if value, ok := fansMap[qiDui]; ok && value > 0 { // 添加清七对,移除七对
			qingQiDui := fan.FanName[majong.Fan_QingQiDui]
			fansMap[qingQiDui] = fan.FanValue[majong.Fan_QingQiDui]
			fansMap[qiDui] = 0
			flag = true
		}
		if value, ok := fansMap[longQiDui]; ok && value > 0 { // 添加清龙七对,移除七对/清七对/龙七对
			if value, ok := fansMap[qiDui]; ok && value > 0 {
				fansMap[qiDui] = 0
			}
			if value, ok := fansMap[fan.FanName[majong.Fan_QingQiDui]]; ok && value > 0 {
				fansMap[qiDui] = 0
			}
			if value, ok := fansMap[longQiDui]; ok && value > 0 {
				fansMap[longQiDui] = 0
			}
			fansMap[qinglongQiDui] = fan.FanValue[majong.Fan_QingLongQiDui]
			if gen >= 1 {
				gen = gen - 1
			}
			flag = true
		}
		if flag { // 存在可以跟清一色可以合组的牌型，移除清一色
			fansMap[qingYiSe] = 0
		}
	}
	return fansMap, gen
}

// GetDi 获取底注
func GetDi() int64 {
	//return r.Option.(*pb.Option_SiChuangXueLiu).Di
	return 1
}

// NewSettleInfo 初始化生成一条新的结算信息
func NewSettleInfo(context *majong.MajongContext, settleType majong.SettleType) *majong.SettleInfo {
	id := uint64(1)
	len := len(context.SettleInfos)
	scores := make(map[uint64]int64)
	if len != 0 {
		id = context.SettleInfos[len-1].Id + 1
	}
	for _, player := range context.Players {
		scores[player.PalyerId] = 0
	}
	return &majong.SettleInfo{
		Id:         id,
		Scores:     scores,
		SettleType: settleType,
	}
}
