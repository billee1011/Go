package scxlai

import (
	"github.com/magiconair/properties/assert"
	"steve/entity/majong"
	"steve/room/majong/global"
	"testing"
)

func Test_DivideCard(t *testing.T) {
	colors := divideByColor([]*majong.Card{&global.Card1W, &global.Card2W, &global.Card1T, &global.Card1B, &global.Card2T, &global.Card1Z})
	assert.Equal(t, colors[majong.CardColor_ColorWan], []majong.Card{global.Card1W, global.Card2W})
	assert.Equal(t, colors[majong.CardColor_ColorTiao], []majong.Card{global.Card1T, global.Card2T})
	assert.Equal(t, colors[majong.CardColor_ColorTong], []majong.Card{global.Card1B})
	assert.Equal(t, colors[majong.CardColor_ColorZi], []majong.Card{global.Card1Z})
}

func Test_GetShunZi(t *testing.T) {
	result := getShunZi([]majong.Card{global.Card1W, global.Card3W, global.Card4W, global.Card5W, global.Card6W, global.Card7W, global.Card8W, global.Card9W}) // 13456789万
	assert.Equal(t, result, []Split{
		{SHUNZI, []majong.Card{global.Card7W, global.Card8W, global.Card9W}},
		{SHUNZI, []majong.Card{global.Card3W, global.Card4W, global.Card5W}},
	})
}

func Test_GetShunZi1(t *testing.T) {
	result := getShunZi([]majong.Card{global.Card1W, global.Card2W, global.Card3W, global.Card4W, global.Card5W, global.Card6W, global.Card7W, global.Card8W, global.Card9W}) // 123456789万
	assert.Equal(t, result, []Split{
		{SHUNZI, []majong.Card{global.Card1W, global.Card2W, global.Card3W}},
		{SHUNZI, []majong.Card{global.Card7W, global.Card8W, global.Card9W}},
		{SHUNZI, []majong.Card{global.Card4W, global.Card5W, global.Card6W}},
	})
}

func Test_GetKeZi(t *testing.T) {
	result := getKeZi([]majong.Card{global.Card3W, global.Card3W, global.Card3W, global.Card4W, global.Card5W, global.Card6W, global.Card7W, global.Card8W, global.Card9W}) // 333456789万
	assert.Equal(t, result, []Split{
		{KEZI, []majong.Card{global.Card3W, global.Card3W, global.Card3W}},
	})
}

func Test_Remove(t *testing.T) {
	cards := []majong.Card{global.Card3W, global.Card3W, global.Card3W}
	cards = Remove(cards, global.Card3W)
	assert.Equal(t, cards, []majong.Card{global.Card3W, global.Card3W})
}

func Test_SplitShunZiFirst(t *testing.T) {
	cards := []majong.Card{global.Card5W, global.Card5W, global.Card5W, global.Card6W, global.Card6W, global.Card6W, global.Card6W, global.Card7W, global.Card7W, global.Card7W}
	shunZis, _, _, _, _, singles := SplitCards(cards, true)

	assert.Equal(t, shunZis, []Split{
		{SHUNZI, []majong.Card{global.Card5W, global.Card6W, global.Card7W}},
		{SHUNZI, []majong.Card{global.Card5W, global.Card6W, global.Card7W}},
		{SHUNZI, []majong.Card{global.Card5W, global.Card6W, global.Card7W}},
	})
	assert.Equal(t, singles, []Split{{SINGLE, []majong.Card{global.Card6W}}})
}
