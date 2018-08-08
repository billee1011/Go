package charge

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/entity/db"
	"steve/external/goldclient"
	"steve/hall/data"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleGetChargeInfoReq 获取充值信息请求
func HandleGetChargeInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.GetChargeInfoReq) (rspMsg []exchanger.ResponseMsg) {
	entry := logrus.WithFields(logrus.Fields{
		"player_id": playerID,
		"platform":  req.GetPlatform(),
	})
	result := &common.Result{
		ErrCode: common.ErrCode_EC_FAIL.Enum(),
		ErrDesc: proto.String("获取充值数据失败"),
	}
	response := &hall.GetChargeInfoRsp{
		Result:       result,
		Items:        nil,
		TodayCharge:  proto.Uint64(0),
		DayMaxCharge: proto.Uint64(getDayMaxCharge()),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_GET_CHARGE_INFO_REQ),
			Body:  response,
		},
	}
	todayCharge, err := data.GetPlayerTodayCharge(playerID)
	if err != nil {
		entry.WithError(err).Errorln("获取今日充值失败")
		return
	}
	response.TodayCharge = proto.Uint64(todayCharge)
	dbPlayer, err := data.GetPlayerInfo(playerID, "provinceID")
	if err != nil {
		entry.WithError(err).Errorln("获取省 ID 失败")
		return
	}
	city := dbPlayer.Provinceid
	entry = entry.WithField("city", city)

	items, err := getItemList(city, int(req.GetPlatform()))
	if err != nil {
		entry.WithError(err).Errorln("获取商品列表失败")
		return
	}
	fillItems(items, response)
	response.Result.ErrCode = common.ErrCode_EC_SUCCESS.Enum()
	response.Result.ErrDesc = proto.String("")
	return
}

func fillItems(items []Item, chargeInfoRsp *hall.GetChargeInfoRsp) {
	chargeInfoRsp.Items = make([]*hall.ChargeItem, 0, len(items))
	for _, item := range items {
		chargeInfoRsp.Items = append(chargeInfoRsp.Items, &hall.ChargeItem{
			ItemId:      proto.Uint64(item.ID),
			ItemName:    proto.String(item.Name),
			Price:       proto.Uint64(item.Price),
			Tag:         proto.String(item.Tag),
			Coin:        proto.Uint64(item.Coin),
			PresentCoin: proto.Uint64(item.PresentCoin),
		})
	}
}

// HandleChargeReq 处理充值请求
func HandleChargeReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.ChargeReq) (rspMsg []exchanger.ResponseMsg) {
	entry := logrus.WithFields(logrus.Fields{
		"player_id": playerID,
		"item_id":   req.GetItemId(),
		"cost":      req.GetCost(),
	})
	result := &common.Result{
		ErrCode: common.ErrCode_EC_FAIL.Enum(),
		ErrDesc: proto.String("操作失败"),
	}
	response := &hall.ChargeRsp{
		Result:       result,
		ObtainedCoin: proto.Uint64(0),
		NewCoin:      proto.Uint64(0),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_GET_CHARGE_INFO_REQ),
			Body:  response,
		},
	}
	dbPlayer, err := data.GetPlayerInfo(playerID, "provinceID")
	if err != nil {
		entry.WithError(err).Errorln("获取省 ID 失败")
		return
	}
	item, entry := getChargeItem(dbPlayer, &req, response, entry)

	if item.Price != req.GetCost() {
		entry.Errorln("消费数值错误")
		return
	}
	// TODO: verify
	addCoin(item, playerID, response, entry)
	return
}

// getChargeItem 获取充值商品
func getChargeItem(dbPlayer *db.TPlayer, req *hall.ChargeReq, response *hall.ChargeRsp, entry *logrus.Entry) (*Item, *logrus.Entry) {
	entry = entry.WithField("city", dbPlayer.Cityid)
	items, err := getItemList(dbPlayer.Cityid, int(req.GetPlatform()))
	if err != nil {
		entry.WithError(err).Errorln("获取商品列表失败")
		return nil, entry
	}
	var item *Item
	for _, it := range items {
		if item.ID == req.GetItemId() {
			item = &it
			break
		}
	}
	if item == nil {
		entry.Errorln("商品不存在")
		return nil, entry
	}
	entry = entry.WithField("item", item)
	return item, entry
}

// addCoin 添加充值获得的金币
func addCoin(item *Item, playerID uint64, response *hall.ChargeRsp, entry *logrus.Entry) (bool, *logrus.Entry) {
	coin := int64(item.Coin + item.PresentCoin)
	entry = entry.WithField("coin", coin)
	newCoin, err := goldclient.AddGold(playerID, 1, coin, 0, 0, 0, 0)
	if err != nil {
		entry.Errorln("添加金币失败")
		return false, entry
	}
	entry = entry.WithField("new_coin", newCoin)
	response.ObtainedCoin = proto.Uint64(uint64(coin))
	response.NewCoin = proto.Uint64(uint64(newCoin))
	return true, entry
}
