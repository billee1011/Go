package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/Sirupsen/logrus"
)

// 平胡:顺子*n+刻子*n+将*1

// card counter
type countMap map[Card]int

func (cm countMap) addAll(cards ...Card) {
	for _, card := range cards {
		cm[card]++
	}
}

func (cm countMap) removeAll(cards ...Card) {
	for _, card := range cards {
		cm[card]--
	}
}

func (cm countMap) contains(sub countMap) bool {
	for card, count := range cm {
		if sub[card] > count {
			return false
		}
	}
	return true
}

func (cm countMap) equal(comp countMap) bool {
	if len(cm) != len(comp) {
		return false
	}
	for card, count := range cm {
		if ycount, ok := comp[card]; !ok || count != ycount {
			return false
		}
	}
	return true
}

func (cm countMap) maxCount() int {
	maxCount := 0
	for _, count := range cm {
		if count > maxCount {
			maxCount = count
		}
	}
	return maxCount
}

// 当前存在的最小的牌，牌的大小等于对应的int的大小
// 使用前应该谨慎注意这个实现，这个方法写的比较粗糙，但是够用
// 只会考虑count>0中的最小牌值
// 如果countMap里有Laizi的情况，Laizi可能会是最小的，因为Laizi等同于0
// 如果countMap里没有count>0的情况，会返回 666 的牌值
func (cm countMap) minValueCard() Card {
	var minCard Card = 666
	for card, count := range cm {
		if count > 0 && card < minCard {
			minCard = card
		}
	}
	return minCard
}

// countMap.len() 和 len(countMap) 的区别在与，前者会屏蔽掉count为0的card
// TODO O(n)是否可以接受？
func (cm countMap) len() int {
	l := 0
	for _, count := range cm {
		if count > 0 {
			l++
		}
	}
	return l
}

//Card 内部牌类型定义
type Card int

//FastCheckHuV1 查表法-判断胡牌，局限性比较大,只适用于没有赖子的情况，并且要求手牌14张，不能得出在现有手牌上所有的成胡牌型
func FastCheckHuV1(cards []Card) (bool, error) {
	gHuCardsMapInitLock.Do(initHuCardsMap)
	if len(cards) != gMaxSize {
		return false, fmt.Errorf("size not match %v != %d\n", cards, gMaxSize)
	}
	_, ok := gHuCardsMap[cardsKey(cards)]
	//fmt.Println("玩家是否可以胡牌：", ok)
	return ok, nil
}

func FastCheckTingV1(cards []Card, avalibleCards []Card) ([]Card, error) {
	gHuCardsMapInitLock.Do(initHuCardsMap)
	if len(cards) != (gMaxSize - 1) {
		return nil, fmt.Errorf("size not match %v != %d\n", cards, gMaxSize-1)
	}
	tingCards := make([]Card, 0)
	for _, avalibleCard := range avalibleCards {
		cards = append(cards, avalibleCard)
		hu, err := FastCheckHuV1(cards)
		if err != nil {
			return nil, err
		}
		if hu {
			tingCards = append(tingCards, avalibleCard)
		}
		cards = cards[:len(cards)-1]
	}
	return tingCards, nil
}

//FastCheckTingInfoV1 查表法-对应的听牌提示，打什么牌可以听什么牌
func FastCheckTingInfoV1(cards []Card, avalibleCards []Card) (map[Card][]Card, error) {
	gHuCardsMapInitLock.Do(initHuCardsMap)
	if len(cards) != gMaxSize {
		return nil, fmt.Errorf("size not match %v != %d\n", cards, gMaxSize)
	}
	tingInfo := make(map[Card][]Card)
	var cm countMap = make(map[Card]int)
	cm.addAll(cards...)
	for index, card := range cards {
		if cm[card] > 0 {
			cm[card] = 0
			cards = append(cards[:index], cards[index+1:]...)
			tingCards, err := FastCheckTingV1(cards, avalibleCards)
			if err != nil {
				return nil, err
			}
			tingInfo[card] = tingCards
			cards = append(cards[:index], append([]Card{card}, cards[index:]...)...)
		}
	}
	return tingInfo, nil
}

// 查表法的cacheMap-所有成胡的牌型按顺序排序后的string做key
var gHuCardsMap map[string]bool

// 查表法-huCardsMap初始化锁
var gHuCardsMapInitLock sync.Once

//InitHuCardsMap 初始化查表法的huCardsMap
// 长耗时操作
func InitHuCardsMap() {
	gHuCardsMapInitLock.Do(initHuCardsMap)
}

//LoadHuCardsMap 加载查表法的huCardsMap
func LoadHuCardsMap(m map[string]bool) {
	gHuCardsMapInitLock.Do(func() {
		loadHuCardsMap(m)
	})
}

// 可以做将的牌
var gJiangCards = []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47}

// 可以做刻的牌
var gKeCards = []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47}

// 可以做顺的牌
var gShunCards = []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39}

// gHuCardsFileName 查胡表存储文件
var gHuCardsFileName = "huCardsMap.gob.data"

// 初始化查表法的huCardsMap具体实现方法
func initHuCardsMap() {
	gHuCardsMap = make(map[string]bool, 10000000)
	if loadHuCardsFromFile(gHuCardsFileName) {
		return
	}
	dpInitHuCardsMap(1, make([]Card, 0, gMaxSize), make(map[Card]int))
	saveHuCardsToFile(gHuCardsFileName)
}

// 保存胡牌表到文件
func saveHuCardsToFile(fileName string) {
	entry := logrus.WithFields(logrus.Fields{
		"name":      "saveHuCardsToFile",
		"file_name": fileName,
		"map_size":  len(gHuCardsMap),
	})

	// syscall.Umask(0)
	w, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil || w == nil {
		entry.WithError(err).Error("打开文件失败")
		return
	}
	enc := gob.NewEncoder(w)
	err = enc.Encode(gHuCardsMap)
	if err != nil {
		entry.WithError(err).Info("Encode 失败")
	} else {
		entry.Info("保存胡牌表到文件成功")
	}
}

// loadHuCardsFromFile 从文件中加载胡牌表
func loadHuCardsFromFile(fileName string) bool {
	entry := logrus.WithFields(logrus.Fields{
		"name":      "loadHuCardsFileFromFile",
		"file_name": fileName,
	})

	r, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil || r == nil {
		entry.Info("打开文件失败")
		return false
	}
	dec := gob.NewDecoder(r)
	dec.Decode(&gHuCardsMap)
	if len(gHuCardsMap) == 0 {
		entry.Info("没有加载到数据")
		return false
	}
	entry.WithField("map_size", len(gHuCardsMap)).Info("从文件加载胡牌成功")
	return true
}

func loadHuCardsMap(m map[string]bool) {
	gHuCardsMap = make(map[string]bool)
	for k, v := range m {
		gHuCardsMap[k] = v
	}
}

// 查表法最大生成牌数
var gMaxSize = 14

// DFS
func dpInitHuCardsMap(target int, cards []Card, cm countMap) {
	//递归结束条件
	if target == 4 {
		cacheToMap(cards)
	}
	if target == 1 {
		for _, j := range gJiangCards {
			cards = append(cards, j, j)
			cm[j] += 2
			dpInitHuCardsMap(target+1, cards, cm)
			cards = cards[:0]
			cm[j] = 0
		}
	}
	if target == 2 {
		l := len(cards)
		kl := len(gKeCards)
		dpInitHuCardsMap(target+1, cards, cm)
		for kSum := 1; kSum <= 4; kSum++ {
			combineFunc(kSum, kl, func(ints []int) {
				for index, ok := range ints {
					if ok == 1 {
						kCard := gKeCards[index]
						cards = append(cards, kCard, kCard, kCard)
						cm[kCard] += 3
					}
				}
				if cm.maxCount() <= 4 {
					dpInitHuCardsMap(target+1, cards, cm)
				}
				cm.removeAll(cards[l:]...)
				cards = cards[:l]
			})
		}
	}
	if target == 3 {
		l := len(cards)
		sSum := (gMaxSize - l) / 3
		dpInitHuCardsMap(target+1, cards, cm)
		for i := 1; i <= sSum; i++ {
			permutationFunc(sSum, len(gShunCards), func(ints []int) {
				match := true
				for _, index := range ints {
					sCard := gShunCards[index]
					if sCard%10 > 7 {
						match = false
						break
					}
					cards = append(cards, sCard, sCard+1, sCard+2)
					cm[sCard]++
					cm[sCard+1]++
					cm[sCard+2]++
				}
				if match && cm.maxCount() <= 4 {
					dpInitHuCardsMap(target+1, cards, cm)
				}
				cm.removeAll(cards[l:]...)
				cards = cards[:l]
			})
		}
	}
}

func cardsKey(cards []Card) string {
	copyCards := make([]int, 0, len(cards))
	for _, card := range cards {
		copyCards = append(copyCards, int(card))
	}
	sort.Ints(copyCards)
	var buf bytes.Buffer
	for _, card := range copyCards {
		fmt.Fprintf(&buf, "%d", card)
	}
	return buf.String()
}

func cacheToMap(cards []Card) {
	gHuCardsMap[cardsKey(cards)] = true
}

// 组合算法
// 输入: m:2 n:3
// 输出: 110,101,011
// func([]int)的参数[]int既为输出，每一种组合都会作为参数调用一次f,int数组每一位存的是对应的下标是否可用，1为可用，0为不可用
func combineFunc(m, n int, f func([]int)) {
	if n < m || m == 0 {
		return
	}
	a := make([]int, n)
	for i := 0; i < m; i++ {
		a[i] = 1
	}
	f(a)
	for i := 0; i < n; i++ {
		if (i+1) < n && a[i] == 1 && a[i+1] == 0 {
			a[i], a[i+1] = a[i+1], a[i]
			l, h := 0, i
			for l < h {
				if a[l] == 0 && a[h] == 1 {
					a[l], a[h] = a[h], a[l]
				}
				if a[l] == 1 {
					l++
				}
				if a[h] == 0 {
					h--
				}
			}
			f(a)
			i = -1
		}
	}
}

// 排列算法
// 输入: m:2 n:3
// 输出: 00,01,02,11,12,22
// func([]int)的参数[]int既为输出，每一种排列都会作为参数调用一次f,int数组每一位存的下标
func permutationFunc(m, n int, f func([]int)) {
	if n < m || m == 0 {
		return
	}
	init := make([]int, m)
out:
	for {
		f(init)
		init[m-1]++
		for i := m - 1; i >= 0; i-- {
			if init[i] >= n {
				if i == 0 {
					break out
				}
				init[i-1]++
				init[i] = 0
			} else {
				break
			}
		}
		for i := 1; i < m; i++ {
			if init[i] < init[i-1] {
				init[i] = init[i-1]
			}
		}
	}
}

// 牌组类型
type GroupType int

const (
	TypeJiang GroupType = 0
	TypeKe    GroupType = 1
	TypeShun  GroupType = 2
)

// 牌组，将刻顺等
type CardGroup struct {
	GroupType GroupType
	Cards     []Card
}

// 赖子牌标识，一手牌的赖子可能有多个，而且牌值都不一样，
// 为了避免由赖子造成没必要的排列组合，计算过程中赖子都采用这个常量，
var Laizi Card = 0

func FastCheckHuV2(cards []Card, laizis map[Card]bool) bool {
	if laizis == nil {
		laizis = make(map[Card]bool)
	}
	//构造计算所需参数
	l := len(cards)
	if l%3 != 2 {
		panic(fmt.Errorf("FastCheckHuV2 wrong cards size %v,%v", l, cards))
	}
	var cm countMap = make(map[Card]int, l)
	laiziCount := 0
	for _, card := range cards {
		if isLaizi := laizis[card]; isLaizi || card == Laizi {
			laiziCount++
		} else {
			cm[card]++
		}
	}

	//将
	jiangGroups := make([]CardGroup, 0, gMaxSize)
	for card, count := range cm {
		if count >= 2 {
			jiangGroups = append(jiangGroups, CardGroup{GroupType: TypeJiang, Cards: []Card{card, card}})
		}
		if count >= 1 && laiziCount > 0 {
			jiangGroups = append(jiangGroups, CardGroup{GroupType: TypeJiang, Cards: []Card{card, Laizi}})
		}
	}
	if len(jiangGroups) == 0 && laiziCount >= 2 {
		jiangGroups = append(jiangGroups, CardGroup{GroupType: TypeJiang, Cards: []Card{Laizi, Laizi}})
	}

	//刻
	keGroups := make([]CardGroup, 0, gMaxSize/3)
	for card, count := range cm {
		if count >= 3 {
			keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{card, card, card}})
		}
		if count >= 2 && laiziCount > 0 {
			keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{card, card, Laizi}})
		}
		if count >= 1 && laiziCount > 1 {
			keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{card, Laizi, Laizi}})
		}
	}
	if len(keGroups) == 0 && laiziCount > 2 {
		keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{Laizi, Laizi, Laizi}})
	}

	groupCount := (l + 1) / 3
	stop := false
	//遍历将
	for _, jGroup := range jiangGroups {
		//遍历刻
		kl := len(keGroups)
		for kSum := groupCount - 1; kSum > 0; kSum-- { //从大到小开始，在较好的情况下会更快
			//刻的排列组合
			combineFuncWithStop(kSum, kl, func(ints []int) {
				var usedCM countMap = make(map[Card]int, len(cards))
				usedCM.addAll(jGroup.Cards...)
				for index, ok := range ints {
					if ok == 1 {
						usedCM.addAll(keGroups[index].Cards...)
					}
				}
				for card, count := range usedCM {
					if card == Laizi {
						if count > laiziCount {
							return
						}
						continue
					}
					if cm[card] < count {
						return
					}
				}
				var unusedCM countMap = make(map[Card]int, len(cards))
				for card, count := range cm {
					unusedCM[card] = count - usedCM[card]
				}
				lastLaiziCount := laiziCount - usedCM[Laizi]
				//delete(usedCM,Laizi)
				if matchShun(unusedCM, lastLaiziCount) {
					stop = true
				}
			}, &stop)
			if stop {
				return true
			}
		}
		//0刻的情况
		var usedCM countMap = make(map[Card]int, len(cards))
		usedCM.addAll(jGroup.Cards...)
		var unusedCM countMap = make(map[Card]int, len(cards))
		for card, count := range cm {
			unusedCM[card] = count - usedCM[card]
		}
		if matchShun(unusedCM, laiziCount) {
			return true
		}
	}
	return false
}

func matchShun(unusedCM countMap, laiziCount int) bool {
	for unusedCM.len() > 0 {
		s := unusedCM.minValueCard() //TODO 优先队列?
		m := s + 1
		l := s + 2
		mCount := unusedCM[m]
		lCount := unusedCM[l]
		if mCount > 0 && lCount > 0 {
			unusedCM[s]--
			unusedCM[m]--
			unusedCM[l]--
		} else if mCount > 0 && laiziCount > 0 {
			unusedCM[s]--
			unusedCM[m]--
			laiziCount--
		} else if lCount > 0 && laiziCount > 0 {
			unusedCM[s]--
			unusedCM[l]--
			laiziCount--
		} else if laiziCount > 1 {
			unusedCM[s]--
			laiziCount -= 2
		} else {
			return false
		}
	}
	return true
}

// 带stop变量的 combineFunc ，避免没必要的排列组合的搜索
// 这个方法的设计是stop变量的值只会在 f 里改变，所以每次在调用 f 后才会判断一次stop
// TODO 这样做可以降低平均时间界，但是最坏情况等同于普通的 combineFunc
func combineFuncWithStop(m, n int, f func([]int), stop *bool) {
	if n < m || m == 0 {
		return
	}
	a := make([]int, n)
	for i := 0; i < m; i++ {
		a[i] = 1
	}
	f(a)
	if *stop {
		return
	}
	for i := 0; i < n; i++ {
		if (i+1) < n && a[i] == 1 && a[i+1] == 0 {
			a[i], a[i+1] = a[i+1], a[i]
			l, h := 0, i
			for l < h {
				if a[l] == 0 && a[h] == 1 {
					a[l], a[h] = a[h], a[l]
				}
				if a[l] == 1 {
					l++
				}
				if a[h] == 0 {
					h--
				}
			}
			f(a)
			if *stop {
				return
			}
			i = -1
		}
	}
}

// 13张牌查听，检索出听哪些牌
// 赖子牌也会参与考虑
func FastCheckTingV2(cards []Card, laizis map[Card]bool) []Card {
	if laizis == nil {
		laizis = make(map[Card]bool, 0)
	}
	l := len(cards)
	if l%3 != 1 {
		panic(fmt.Errorf("FastCheckTingV2 wrong cards size %v", l))
	}
	var cm countMap = make(map[Card]int, l)
	for _, card := range cards {
		cm[card]++
	}
	tingCards := make([]Card, 0)
	availableCards := availableCards(cards, laizis)
	for _, availableCard := range availableCards {
		if cm[availableCard] == 4 {
			continue
		}
		cards = append(cards, availableCard)
		if FastCheckHuV2(cards, laizis) {
			tingCards = append(tingCards, availableCard)
		}
		cards = cards[:len(cards)-1]
	}
	return tingCards
}

// 赖子也是可用牌
func availableCards(cards []Card, laizis map[Card]bool) []Card {
	l := len(cards)
	var cm countMap = make(map[Card]int, l)
	for _, card := range cards {
		cm[card]++
		if !laizis[card] {
			y := card % 10
			if y > 2 {
				cm[card-2]++
			}
			if y > 1 {
				cm[card-1]++
			}
			if y < 8 {
				cm[card+2]++
			}
			if y < 9 {
				cm[card+1]++
			}
		}
	}
	ret := make([]Card, 0, len(cm))
	for card := range cm {
		ret = append(ret, card)
	}
	for card := range laizis {
		ret = append(ret, card)
	}
	return ret
}

// 14张牌查听,检索出分别打掉哪张牌可以听哪些牌
func FastCheckTingInfoV2(cards []Card, laizis map[Card]bool) map[Card][]Card {
	if laizis == nil {
		laizis = make(map[Card]bool, 0)
	}
	l := len(cards)
	if l%3 != 2 {
		panic(fmt.Errorf("FastCheckTingInfoV2 wrong cards size %v", l))
	}
	checkCards := make([]Card, 0, 14)
	tingInfos := make(map[Card][]Card)
	var cm countMap = make(map[Card]int, l)
	for _, card := range cards {
		cm[card]++
	}
	var cache countMap = make(map[Card]int, l)
	var laiziTingCards = make([]Card, 0)
	availableCards := availableCards(cards, laizis)
	for i, card := range cards {
		isLaizi := laizis[card]
		//已经遍历过的牌不重复
		if cache[card] > 0 {
			continue
		} else if isLaizi && cache[Laizi] > 0 {
			//赖子的牌值不一样，但是结果可以复用
			if len(laiziTingCards) > 0 {
				tingInfos[card] = laiziTingCards
			}
			continue
		} else {
			if isLaizi {
				cache[Laizi]++
			} else {
				cache[card]++
			}
		}
		//开始检查
		checkCards = append(checkCards, cards[:i]...)
		checkCards = append(checkCards, cards[i+1:]...)
		//剪枝
		checkCards = append(checkCards, Laizi)
		if FastCheckHuV2(checkCards, laizis) {
			checkCards = checkCards[:len(checkCards)-1]
			//细化搜索 TODO 时间级过大
			tingCards := make([]Card, 0)
			for _, availableCard := range availableCards {
				// if cm[availableCard] == 4 && card != availableCard {
				// 	continue
				// }
				checkCards = append(checkCards, availableCard)
				if FastCheckHuV2(checkCards, laizis) {
					tingCards = append(tingCards, availableCard)
				}
				checkCards = checkCards[:len(checkCards)-1]
			}
			if len(tingCards) > 0 {
				tingInfos[card] = tingCards
				if isLaizi {
					laiziTingCards = tingCards
				}
			}
		}
		checkCards = checkCards[:0]
	}
	return tingInfos
}
