package utils

import (
	"steve/client_pb/room"
	"steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// ZixunTransform 自询消息转换
func ZixunTransform(record *majong.ZixunRecord) *room.RoomZixunNtf {
	zixunNtf := &room.RoomZixunNtf{}
	zixunNtf.EnableAngangCards = record.GetEnableAngangCards()
	zixunNtf.EnableBugangCards = record.GetEnableBugangCards()
	zixunNtf.EnableChupaiCards = record.GetEnableChupaiCards()
	zixunNtf.EnableQi = proto.Bool(record.GetEnableQi())
	zixunNtf.EnableZimo = proto.Bool(record.GetEnableZimo())
	huType := transformHuType(record.GetHuType())
	if huType != nil {
		zixunNtf.HuType = huType
	}
	zixunNtf.CanTingCardInfo = transformCanTingCardInfo(record.GetCanTingCardInfo())

	return zixunNtf
}

func transformHuType(recordHuType majong.HuType) *room.HuType {
	var huType room.HuType
	switch recordHuType {
	case majong.HuType_hu_dianpao:
		{
			huType = room.HuType_HT_DIANPAO
		}
	case majong.HuType_hu_dihu:
		{
			huType = room.HuType_HT_DIHU
		}
	case majong.HuType_hu_ganghoupao:
		{
			huType = room.HuType_HT_GANGHOUPAO
		}
	case majong.HuType_hu_gangkai:
		{
			huType = room.HuType_HT_GANGKAI
		}
	case majong.HuType_hu_gangshanghaidilao:
		{
			huType = room.HuType_HT_GANGSHANGHAIDILAO
		}
	case majong.HuType_hu_haidilao:
		{
			huType = room.HuType_HT_HAIDILAO
		}
	case majong.HuType_hu_qiangganghu:
		{
			huType = room.HuType_HT_QIANGGANGHU
		}
	case majong.HuType_hu_tianhu:
		{
			huType = room.HuType_HT_TIANHU
		}
	case majong.HuType_hu_zimo:
		{
			huType = room.HuType_HT_ZIMO
		}
	default:
		return nil
	}
	return &huType
}

func transformCanTingCardInfo(minfos []*majong.CanTingCardInfo) []*room.CanTingCardInfo {
	rinfos := []*room.CanTingCardInfo{}
	for _, minfo := range minfos {
		rinfo := &room.CanTingCardInfo{}
		rinfo.OutCard = proto.Uint32(minfo.GetOutCard())
		rinfo.TingCardInfo = transformTingCardInfo(minfo.GetTingCardInfo())
		rinfos = append(rinfos, rinfo)
	}
	return rinfos
}

func transformTingCardInfo(minfos []*majong.TingCardInfo) []*room.TingCardInfo {
	rinfos := []*room.TingCardInfo{}
	for _, minfo := range minfos {
		rinfo := &room.TingCardInfo{}
		rinfo.TingCard = proto.Uint32(minfo.GetTingCard())
		rinfo.Times = proto.Uint32(minfo.GetTimes())
		rinfos = append(rinfos, rinfo)
	}
	return rinfos
}
