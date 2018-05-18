package global

import (
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/utils"
)

func makeRoomCards(card ...room.Card) []*room.Card {
	return utils.MakeRoomCards(card...)
}

var (
	// CommonStartGameParams 通用的启动游戏参数
	CommonStartGameParams = utils.StartGameParams{
		Cards: [][]*room.Card{
			makeRoomCards(Card1W, Card1W, Card1W, Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card4W),
			makeRoomCards(Card5W, Card5W, Card5W, Card5W, Card6W, Card6W, Card6W, Card6W, Card7W, Card7W, Card7W, Card7W, Card8W),
			makeRoomCards(Card1T, Card1T, Card1T, Card1T, Card2T, Card2T, Card2T, Card2T, Card3T, Card3T, Card3T, Card3T, Card4T),
			makeRoomCards(Card5T, Card5T, Card5T, Card5T, Card6T, Card6T, Card6T, Card6T, Card7T, Card7T, Card7T, Card7T, Card8T),
		},
		WallCards: []*room.Card{
			&Card1B,
		},
		HszDir:     room.Direction_AntiClockWise,
		BankerSeat: 0,
		ServerAddr: config.ServerAddr,
		ClientVer:  config.ClientVersion,

		HszCards: [][]*room.Card{
			makeRoomCards(Card1W, Card1W, Card1W),
			makeRoomCards(Card5W, Card5W, Card5W),
			makeRoomCards(Card1T, Card1T, Card1T),
			makeRoomCards(Card5T, Card5T, Card5T),
		},
		DingqueColor: []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO},
	}
)
