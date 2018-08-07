package scxlai

import (
	"sort"
	"steve/entity/majong"
	"steve/gutils"
)

func RemoveSplits(cards []majong.Card, splits []Split) []majong.Card {
	for _, split := range splits {
		for _, card := range split.cards {
			cards = Remove(cards, card)
		}
	}
	return cards
}

func Remove(cards []majong.Card, removeCard majong.Card) []majong.Card {
	for i, card := range cards {
		if card.Equals(removeCard) {
			cards = append(cards[:i], cards[i+1:]...)
			break
		}
	}
	return cards
}

func ContainsEdge(split Split) bool {
	for _, card := range split.cards {
		if card.Point == 1 || card.Point == 9 {
			return true
		}
	}
	return false
}

func Contains(splits []Split, inCard majong.Card) bool {
	for _, split := range splits {
		for _, card := range split.cards {
			if card == inCard {
				return true
			}
		}
	}
	return false
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
