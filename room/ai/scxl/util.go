package scxlai

import (
	"sort"
	"steve/entity/majong"
	"steve/gutils"
)

func RemoveSplits(cards []majong.Card, splits []Split) []majong.Card {
	var clone = Clone(cards) //克隆，避免修改原slice的底层数组，产生副作用
	for _, split := range splits {
		for _, card := range split.cards {
			clone = Remove(clone, card)
		}
	}
	return clone
}

func Remove(cards []majong.Card, removeCard majong.Card) []majong.Card {
	for i, card := range cards {
		if card == removeCard {
			cards = append(cards[:i], cards[i+1:]...)
			break
		}
	}
	return cards
}

func Clone(cards []majong.Card) []majong.Card {
	var clone []majong.Card
	for _, card := range cards {
		clone = append(clone, card)
	}
	return clone
}

func ContainsEdge(cards []majong.Card) bool {
	for _, card := range cards {
		if card.Point == 1 || card.Point == 9 {
			return true
		}
	}
	return false
}

func Contains(splits []Split, inCard majong.Card) bool {
	for _, split := range splits {
		if ContainsCard(split.cards, inCard) {
			return true
		}
	}
	return false
}

func ContainsCard(cards []majong.Card, inCard majong.Card) bool {
	for _, card := range cards {
		if card == inCard {
			return true
		}
	}
	return false
}

func NonPointer(cards []*majong.Card) []majong.Card {
	var result []majong.Card
	for _, card := range cards {
		result = append(result, *card)
	}
	return result
}

type MJCardSlice []majong.Card

func (cs MJCardSlice) Len() int      { return len(cs) }
func (cs MJCardSlice) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs MJCardSlice) Less(i, j int) bool {
	return gutils.ServerCard2Number(&cs[i]) < gutils.ServerCard2Number(&cs[j])
}

func MJCardSort(cards []majong.Card) {
	cs := MJCardSlice(cards)
	sort.Sort(cs)
}
