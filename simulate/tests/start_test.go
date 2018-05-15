package tests

import (
	"steve/client_pb/room"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRoomCards(card ...room.Card) []*room.Card {
	result := []*room.Card{}
	for i := range card {
		result = append(result, &card[i])
	}
	return result
}

// Test_StartGame 测试游戏开始
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌，
func Test_StartGame(t *testing.T) {
	deskData, err := utils.StartGame(utils.StartGameParams{
		Cards: [][]*room.Card{
			makeRoomCards(Card1W, Card1W, Card1W, Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card4W),
			makeRoomCards(Card5W, Card5W, Card5W, Card5W, Card6W, Card6W, Card6W, Card6W, Card7W, Card7W, Card7W, Card7W, Card8W),
			makeRoomCards(Card1T, Card1T, Card1T, Card1T, Card2T, Card2T, Card2T, Card2T, Card3T, Card3T, Card3T, Card3T, Card4T),
			makeRoomCards(Card5T, Card5T, Card5T, Card5T, Card6T, Card6T, Card6T, Card6T, Card7T, Card7T, Card7T, Card7T, Card8T),
		},
		WallCards: []*room.Card{
			&Card1B,
		},
		HszDir:     0, // TODO
		BankerSeat: 0,
		ServerAddr: ServerAddr,
		ClientVer:  ClientVersion,

		HszCards: [][]*room.Card{
			makeRoomCards(Card1W, Card1W, Card1W),
			makeRoomCards(Card5W, Card5W, Card5W),
			makeRoomCards(Card1T, Card1T, Card1T),
			makeRoomCards(Card5T, Card5T, Card5T),
		},
		DingqueColor: []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO},
	})
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
}
