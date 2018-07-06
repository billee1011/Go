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

func (cm countMap) clone() countMap {
	tmp := make(countMap, len(cm))
	for c, v := range cm {
		tmp[c] = v
	}
	return tmp
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

// GroupType 牌组类型
type GroupType int

const (
	TypeJiang GroupType = 0
	TypeKe    GroupType = 1
	TypeShun  GroupType = 2
)

// CardGroup 牌组，将刻顺等
type CardGroup struct {
	GroupType  GroupType
	Cards      []Card
	Replaces   []Card // 癞子所替换的牌
	ReplaceAll bool   // 癞子是否可以替换所有的牌，此时不用理会 Replaces
}

// Combine 牌组列表
type Combine []CardGroup

// Combines Combine 列表
type Combines []Combine

// Laizi 赖子牌标识，一手牌的赖子可能有多个，而且牌值都不一样，
// 为了避免由赖子造成没必要的排列组合，计算过程中赖子都采用这个常量，
var Laizi Card // = 0

// makeJiangGroups 构建能组成将的组合
func makeJiangGroups(cm countMap, laiziCount int) []CardGroup {
	jiangGroups := make([]CardGroup, 0, gMaxSize)
	for card, count := range cm {
		if count >= 2 {
			jiangGroups = append(jiangGroups, CardGroup{GroupType: TypeJiang, Cards: []Card{card, card}})
		}
		if count >= 1 && laiziCount > 0 {
			jiangGroups = append(jiangGroups, CardGroup{GroupType: TypeJiang, Cards: []Card{card, Laizi}, Replaces: []Card{card}})
		}
	}
	// 癞子数大于2，可以单独组成将
	if laiziCount >= 2 {
		jiangGroups = append(jiangGroups, CardGroup{GroupType: TypeJiang, Cards: []Card{Laizi, Laizi}, ReplaceAll: true})
	}
	return jiangGroups
}

// makeCountMap 所有的牌构建乘 countMap
func makeCountMap(cards []Card, laizis map[Card]bool) (cm countMap, laiziCount int) {
	cm = make(countMap, len(cards))
	if laizis == nil {
		laizis = make(map[Card]bool)
	}
	laiziCount = 0
	for _, card := range cards {
		if isLaizi := laizis[card]; isLaizi || card == Laizi {
			laiziCount++
		} else {
			cm[card]++
		}
	}
	return
}

func min(v ...int) int {
	tmp := v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < tmp {
			tmp = v[i]
		}
	}
	return tmp
}

// appendMultiGroup 同一种 group 添加多次
func appendMultiGroup(groups []CardGroup, group CardGroup, count int) []CardGroup {
	for j := 0; j < count; j++ {
		groups = append(groups, group)
	}
	return groups
}

// makeSpecialShunGroup 针对于某个顺建立 groups
func makeSpecialShunGroup(groups []CardGroup, s, m, l Card, sc, mc, lc int, laiziCount int) []CardGroup {
	groups = appendMultiGroup(groups, CardGroup{
		GroupType: TypeShun,
		Cards:     []Card{s, m, l},
	}, min(sc, mc, lc))

	if laiziCount >= 1 {
		groups = appendMultiGroup(groups, CardGroup{
			GroupType: TypeShun,
			Cards:     []Card{s, m, Laizi},
			Replaces:  []Card{l},
		}, min(sc, mc, laiziCount))
		groups = appendMultiGroup(groups, CardGroup{
			GroupType: TypeShun,
			Cards:     []Card{s, Laizi, l},
			Replaces:  []Card{m},
		}, min(sc, laiziCount, lc))
		groups = appendMultiGroup(groups, CardGroup{
			GroupType: TypeShun,
			Cards:     []Card{Laizi, m, l},
			Replaces:  []Card{s},
		}, min(laiziCount, mc, lc))
	}
	if laiziCount >= 2 {
		groups = appendMultiGroup(groups, CardGroup{
			GroupType: TypeShun,
			Cards:     []Card{s, Laizi, Laizi},
			Replaces:  []Card{m, l},
		}, min(laiziCount/2, sc))
		groups = appendMultiGroup(groups, CardGroup{
			GroupType: TypeShun,
			Cards:     []Card{Laizi, m, Laizi},
			Replaces:  []Card{s, l},
		}, min(laiziCount/2, mc))
		groups = appendMultiGroup(groups, CardGroup{
			GroupType: TypeShun,
			Cards:     []Card{Laizi, Laizi, l},
			Replaces:  []Card{s, m},
		}, min(laiziCount/2, lc))
	}
	return groups
}

// makeShunGroups 构建顺的组合
func makeShunGroups(cm countMap, laiziCount int) []CardGroup {
	result := make([]CardGroup, 0, len(cm)/3)
	maked := map[Card]struct{}{} // 已经构建了的顺, Card为该顺中的最小的值
	for card := range cm {
		if card/10 == 4 || card == Laizi { // 风牌不能作顺
			continue
		}
		for i := Card(-2); i <= 0; i++ {
			s, m, l := card+i, card+i+1, card+i+2       // 最小的牌， 中间的牌，最大的牌
			if s/10 != l/10 || s%10 == 0 || l%10 == 0 { // 花色不同，不能组成顺
				continue
			}
			if _, ok := maked[s]; ok { // 已经构建过了，不重复构建
				continue
			}
			osc, omc, olc := cm[s], cm[m], cm[l] // 牌的数量
			result = makeSpecialShunGroup(result, s, m, l, osc, omc, olc, laiziCount)
			maked[s] = struct{}{}
		}
	}
	if laiziCount >= 3 {
		result = appendMultiGroup(result, CardGroup{
			GroupType:  TypeShun,
			Cards:      []Card{Laizi, Laizi, Laizi},
			ReplaceAll: true,
		}, laiziCount/3)
	}
	return result
}

// makeKeGroups 构建刻的组合
func makeKeGroups(cm countMap, laiziCount int) []CardGroup {
	keGroups := make([]CardGroup, 0, gMaxSize/3)
	for card, count := range cm {
		if count >= 3 {
			keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{card, card, card}})
		}
		if count >= 2 && laiziCount >= 1 {
			keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{card, card, Laizi}, Replaces: []Card{card}})
		}
		if count >= 1 && laiziCount >= 2 {
			keGroups = append(keGroups, CardGroup{GroupType: TypeKe, Cards: []Card{card, Laizi, Laizi}, Replaces: []Card{card}})
		}
	}

	if laiziCount >= 3 {
		keGroups = appendMultiGroup(keGroups, CardGroup{
			GroupType:  TypeKe,
			Cards:      []Card{Laizi, Laizi, Laizi},
			ReplaceAll: true,
		}, laiziCount/3)
	}
	return keGroups
}

func getLeft(cm countMap, laiziCount int, groups ...CardGroup) (bool, int, countMap) {
	used := make(countMap, len(cm))
	for _, group := range groups {
		for _, card := range group.Cards {
			if card == Laizi {
				if laiziCount > 0 {
					laiziCount--
				} else {
					return false, 0, nil
				}
			} else {
				if used[card] < cm[card] {
					used[card]++
				} else {
					return false, 0, nil
				}
			}
		}
	}
	unused := make(countMap, len(cm))
	for card, count := range cm {
		if count > used[card] {
			unused[card] = count - used[card]
		}
	}
	return true, laiziCount, unused
}

// FastCheckHuV2 查胡，
// cards : 所有的牌
// laizis: 哪些牌是癞子
// needAll: 是否需要查出所有的胡牌组合
// TODO : 优化算法，按照不同花色分组查询
// TODO : 查胡方法另外建一个包
func FastCheckHuV2(cards []Card, laizis map[Card]bool, needAll bool) (bool, []Combine) {
	combines := []Combine{}
	cardsCount := len(cards)
	if cardsCount%3 != 2 {
		return false, combines
	}
	groupCount := cardsCount/3 + 1
	cm, laiziCount := makeCountMap(cards, laizis)
	jiangGroups := makeJiangGroups(cm, laiziCount)

	for _, jGroup := range jiangGroups {
		stop := false
		// 算出剩余的
		_, leftLaizi, unused := getLeft(cm, laiziCount, jGroup)
		keGroups := makeKeGroups(unused, leftLaizi)

		for kecount := groupCount - 1; kecount >= 0; kecount-- {
			combineFuncWithStop(kecount, len(keGroups), func(indexs []int) {
				useKeGroups := make([]CardGroup, 0, kecount)
				for index, ok := range indexs {
					if ok == 1 {
						useKeGroups = append(useKeGroups, keGroups[index])
					}
				}
				ok1, leftLaizi1, unused1 := getLeft(unused, leftLaizi, useKeGroups...)
				if !ok1 {
					return
				}
				shunGroups := makeShunGroups(unused1, leftLaizi1)
				useShunCount := groupCount - kecount - 1
				combineFuncWithStop(useShunCount, len(shunGroups), func(indexs []int) {
					useShunGroups := make([]CardGroup, 0, useShunCount)
					for index, ok := range indexs {
						if ok == 1 {
							useShunGroups = append(useShunGroups, shunGroups[index])
						}
					}
					if ok2, _, _ := getLeft(unused1, leftLaizi1, useShunGroups...); !ok2 {
						return
					}
					groups := append([]CardGroup{jGroup}, useKeGroups...)
					groups = append(groups, useShunGroups...)
					combines = append(combines, groups)
					if !needAll {
						stop = true
					}
				}, &stop)
			}, &stop)
		}
		if !needAll && len(combines) > 0 {
			return true, combines
		}
	}
	return len(combines) > 0, combines
}

// 带stop变量的 combineFunc ，避免没必要的排列组合的搜索
// 这个方法的设计是stop变量的值只会在 f 里改变，所以每次在调用 f 后才会判断一次stop
// TODO 这样做可以降低平均时间界，但是最坏情况等同于普通的 combineFunc
func combineFuncWithStop(m, n int, f func([]int), stop *bool) {
	if n < m {
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

// FastCheckTingV2Old 检查当前可以听哪些牌
// 返回所有正在听的牌与对应的组合
func FastCheckTingV2Old(cards []Card, laizis map[Card]bool) CardCombines {
	checkCards := make([]Card, 0, len(cards)+1)
	checkCards = append(checkCards, cards...)
	checkCards = append(checkCards, Laizi) // 加一张癞子牌查胡，然后从组合中取出癞子可以替换的牌

	result := CardCombines{}

	ok, combines := FastCheckHuV2(checkCards, laizis, true)
	if !ok {
		return result
	}
	for _, combine := range combines {
		tingCards := getLaiziCanReplaceCards(combine)
		for _, card := range tingCards {
			if result[card] == nil {
				result[card] = Combines{combine}
			} else {
				result[card] = append(result[card], combine)
			}
		}
	}
	return result
}

// FastCheckTingV2 检查当前可以听哪些牌
// 返回所有正在听的牌与对应的组合
// TODO: 需要优化
func FastCheckTingV2(cards []Card, laizis map[Card]bool) CardCombines {

	laizicount := 0
	for _, card := range cards {
		if card == Laizi || (laizis != nil && laizis[card]) {
			laizicount++
		}
	}
	if laizicount >= 2 {
		return FastCheckTingV2Old(cards, laizis)
	}

	checkCards := make([]Card, 0, len(cards)+1)
	checkCards = append(checkCards, cards...)
	checkCards = append(checkCards, Laizi) // 加一张癞子牌查胡，然后从组合中取出癞子可以替换的牌

	result := CardCombines{}

	ok, combines := FastCheckHuV2(checkCards, laizis, true)
	if !ok {
		return result
	}
	for _, combine := range combines {
		tingCards := getLaiziCanReplaceCards(combine)
		for _, card := range tingCards {
			if result[card] != nil {
				continue
			}
			checkCards := append(checkCards[:len(checkCards)-1], card)
			_, newcombines := FastCheckHuV2(checkCards, laizis, true)
			result[card] = newcombines
		}
	}
	return result
}

// getLaiziCanReplaceCards 获取癞子可以替换的牌列表
// 以及替换每张牌后的 Combine
func getLaiziCanReplaceCards(combine Combine) []Card {
	result := []Card{}
	for _, group := range combine {
		if group.ReplaceAll {
			return allCards
		}
		if group.Replaces == nil {
			continue
		}
		result = append(result, group.Replaces...)
	}
	return result
}

var allCards = []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47}

// CardCombines Card->Combines 的映射
type CardCombines map[Card]Combines

func isLaizi(card Card, laizis map[Card]bool) bool {
	if laizis == nil {
		return card == Laizi
	}
	_, ok := laizis[card]
	return ok
}

// FastCheckTingInfoV2 14张牌查听,检索出分别打掉哪张牌可以听哪些牌以及对应的组合
func FastCheckTingInfoV2(cards []Card, laizis map[Card]bool) map[Card]CardCombines {
	result := make(map[Card]CardCombines)
	// 打出癞子的听牌结果
	var laiziCache CardCombines
	checked := map[Card]bool{} // 已经查过的牌

	for index, card := range cards {
		laizi := isLaizi(card, laizis)
		// 癞子结果复用
		if laizi && laiziCache != nil {
			result[card] = laiziCache
			continue
		}
		if _, ok := checked[card]; ok { // 已经查过了不再查
			continue
		}
		checkCards := make([]Card, 0, len(cards)-1)
		checkCards = append(checkCards, cards[:index]...)
		checkCards = append(checkCards, cards[index+1:]...)

		ting := FastCheckTingV2(checkCards, laizis)
		if ting != nil && len(ting) > 0 {
			result[card] = ting
		}
		checked[card] = true
		if laizi {
			laiziCache = ting
		}
	}
	return result
}
