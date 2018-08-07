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
	result := SplitShunZi([]majong.Card{global.Card1W, global.Card3W, global.Card4W, global.Card5W, global.Card6W, global.Card7W, global.Card8W, global.Card9W}) // 13456789万
	assert.Equal(t, result, []Split{
		{SHUNZI, []majong.Card{global.Card7W, global.Card8W, global.Card9W}},
		{SHUNZI, []majong.Card{global.Card3W, global.Card4W, global.Card5W}},
	})
}

func Test_GetShunZi1(t *testing.T) {
	result := SplitShunZi([]majong.Card{global.Card1W, global.Card2W, global.Card3W, global.Card4W, global.Card5W, global.Card6W, global.Card7W, global.Card8W, global.Card9W}) // 123456789万
	assert.Equal(t, result, []Split{
		{SHUNZI, []majong.Card{global.Card1W, global.Card2W, global.Card3W}},
		{SHUNZI, []majong.Card{global.Card7W, global.Card8W, global.Card9W}},
		{SHUNZI, []majong.Card{global.Card4W, global.Card5W, global.Card6W}},
	})
}

func Test_GetKeZi(t *testing.T) {
	result := SplitKeZi([]majong.Card{global.Card3W, global.Card3W, global.Card3W, global.Card4W, global.Card5W, global.Card6W, global.Card7W, global.Card8W, global.Card9W}) // 333456789万
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

func Test_SplitKeZiFirst(t *testing.T) {
	cards := []majong.Card{global.Card5W, global.Card5W, global.Card5W, global.Card6W, global.Card6W, global.Card6W, global.Card6W, global.Card7W, global.Card7W, global.Card7W}
	_, keZis, _, _, _, singles := SplitCards(cards, false)

	assert.Equal(t, keZis, []Split{
		{KEZI, []majong.Card{global.Card5W, global.Card5W, global.Card5W}},
		{KEZI, []majong.Card{global.Card7W, global.Card7W, global.Card7W}},
		{KEZI, []majong.Card{global.Card6W, global.Card6W, global.Card6W}},
	})
	assert.Equal(t, singles, []Split{{SINGLE, []majong.Card{global.Card6W}}})
}

func Test_SplitSpaceShunZi(t *testing.T) {
	cards := []majong.Card{global.Card5W, global.Card7W}
	_, _, _, _, singleChas, _ := SplitCards(cards, false)

	assert.Equal(t, singleChas, []Split{
		{SINGLE_CHA, []majong.Card{global.Card5W, global.Card7W}},
	})
}

func Test_SplitColorCards(t *testing.T) {
	cards := []majong.Card{global.Card5W, global.Card5T, global.Card5B, global.Card6W, global.Card6T, global.Card6B, global.Card6W, global.Card7W, global.Card7T, global.Card7B}
	shunZis, _, _, _, _, singles := SplitCards(cards, false)

	assert.Equal(t, shunZis, []Split{
		{SHUNZI, []majong.Card{global.Card5W, global.Card6W, global.Card7W}},
		{SHUNZI, []majong.Card{global.Card5B, global.Card6B, global.Card7B}},
		{SHUNZI, []majong.Card{global.Card5T, global.Card6T, global.Card7T}},
	})
	assert.Equal(t, singles, []Split{{SINGLE, []majong.Card{global.Card6W}}})
}
