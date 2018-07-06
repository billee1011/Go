package utils

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_initHuCardsMap(t *testing.T) {
	initHuCardsMap()
	// fmt.Println(len(gHuCardsMap))
	// pwd, _ := os.Getwd()
	// filename := fmt.Sprintf("%s%s%s", pwd, string(os.PathSeparator), "huCardsMap.data")
	// bytes, _ := json.Marshal(gHuCardsMap)
	// ioutil.WriteFile(filename, bytes, os.ModePerm)
}

// func Test_initHuCardsMap_gob(t *testing.T) {
// 	initHuCardsMap()
// 	fmt.Println(len(gHuCardsMap))
// 	// pwd, _ := os.Getwd()
// 	// filename := fmt.Sprintf("%s%s%s", pwd, string(os.PathSeparator), "huCardsMap.data")
// 	// bytes, _ := json.Marshal(gHuCardsMap)
// 	syscall.Umask(0)
// 	w, err := os.OpenFile("huCardsMap.gob.data", os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
// 	assert.Nil(t, err)

// 	enc := gob.NewEncoder(w)
// 	assert.Nil(t, enc.Encode(gHuCardsMap))
// }

// func Test_loadHuCardsMap_gob(t *testing.T) {
// 	r, err := os.OpenFile("huCardsMap.gob.data", os.O_RDONLY, 0644)
// 	assert.Nil(t, err)
// 	dec := gob.NewDecoder(r)
// 	gHuCardsMap := make(map[string]bool, 10000000)
// 	dec.Decode(&gHuCardsMap)
// 	assert.NotEqual(t, len(gHuCardsMap), 0)
// }

// func Test_loadHuCardsMap(t *testing.T) {
// 	tLoadHuCardsMap()
// 	fmt.Println(len(gHuCardsMap))
// }

// func tLoadHuCardsMap() {
// 	pwd, _ := os.Getwd()
// 	filename := fmt.Sprintf("%s%s%s", pwd, string(os.PathSeparator), "huCardsMap.data")
// 	bytes, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	m := make(map[string]bool)
// 	err = json.Unmarshal(bytes, &m)
// 	if err != nil {
// 		panic(err)
// 	}
// 	LoadHuCardsMap(m)
// }

func Test_combineFunc(t *testing.T) {
	type args struct {
		m int
		n int
		f func([]int)
	}
	count := 0
	tests := []struct {
		name string
		args args
	}{
		{name: "m3n5", args: args{m: 4, n: 18, f: func(ints []int) {
			count++
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			combineFunc(tt.args.m, tt.args.n, tt.args.f)
		})
	}
	fmt.Println(count)
}

func Test_permutationFunc(t *testing.T) {
	type args struct {
		m int
		n int
		f func([]int)
	}
	//count := 0
	tests := []struct {
		name string
		args args
	}{
		//{name: "m3n5", args: args{m: 3, n: 5, f: func(ints []int) {
		//	count++
		//	fmt.Println(ints)
		//}}},
		{name: "1516171718192324253233343838", args: args{m: 4, n: len(gShunCards), f: func(ints []int) {
			if gShunCards[ints[0]] == 15 && gShunCards[ints[1]] == 17 && gShunCards[ints[2]] == 23 && gShunCards[ints[3]] == 32 {
				fmt.Println(ints)
			}
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permutationFunc(tt.args.m, tt.args.n, tt.args.f)
		})
	}
	//fmt.Println(count)
}

func BenchmarkPermutationFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		permutationFunc(4, len(gShunCards), func(ints []int) {

		})
	}
}

func BenchmarkCombineFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		combineFunc(3, 27, func(ints []int) {

		})
	}
}

//func init() {
//	InitHuCardsMap()
//}

func TestFastCheckHuV1(t *testing.T) {
	//cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	cards := []Card{24, 15, 33, 16, 38, 19, 17, 25, 23, 17, 34, 18, 32, 38}
	fmt.Println(cardsKey(cards))
	fmt.Println(FastCheckHuV1(cards))
}

func TestFastCheckTingV1(t *testing.T) {
	cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19}
	fmt.Println(FastCheckTingV1(cards, []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47}))
}

func TestFastCheckTingInfoV1(t *testing.T) {
	cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	fmt.Println(FastCheckTingInfoV1(cards, []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47}))
}

func getRandomTestCards(num int) []Card {
	base := []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47}
	allCards := make([]Card, 0, 4*len(base))
	allCards = append(allCards, base...)
	allCards = append(allCards, base...)
	allCards = append(allCards, base...)
	allCards = append(allCards, base...)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(allCards)
	for i := 0; i < l; i++ {
		swapIndex := r.Intn(l-i) + i
		allCards[i], allCards[swapIndex] = allCards[swapIndex], allCards[i]
	}
	for i := l - 1; i >= 0; i-- {
		swapIndex := r.Intn(i + 1)
		allCards[i], allCards[swapIndex] = allCards[swapIndex], allCards[i]
	}
	return allCards[:num]
}

// func BenchmarkFastCheckHuV1(b *testing.B) {
// 	b.StopTimer()
// 	InitHuCardsMap()
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()
// 		cards := getRandomTestCards(14)
// 		b.StartTimer()
// 		FastCheckHuV1(cards)
// 	}
// }

// func BenchmarkFastCheckTingV1(b *testing.B) {
// 	b.StopTimer()
// 	InitHuCardsMap()
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()
// 		cards := getRandomTestCards(14)
// 		b.StartTimer()
// 		FastCheckTingV1(cards[:len(cards)-1], []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47})
// 	}
// }

// func BenchmarkFastCheckTingInfoV1(b *testing.B) {
// 	b.StopTimer()
// 	InitHuCardsMap()
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()
// 		cards := getRandomTestCards(14)
// 		b.StartTimer()
// 		FastCheckTingInfoV1(cards, []Card{11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47})
// 	}
// }

//用随机的牌来比较正确性
func TestCheckHuCorrect1(t *testing.T) {
	laizis := map[Card]bool{}
	for i := 0; i < 10000000; i++ {
		if i%1000000 == 0 {
			fmt.Println(i)
		}
		cards := getRandomTestCards(14)
		hu1, _ := FastCheckHuV1(cards)
		hu2, _ := FastCheckHuV2(cards, laizis, false)
		if hu1 != hu2 {
			panic(fmt.Errorf("cards:[%v] not correct hu1[%v] hu2[%v]", cards, hu1, hu2))
		}
	}
}

//用查表法所有成胡的key去判断正确性
func TestCheckHuCorrectV2(t *testing.T) {
	InitHuCardsMap()
	laizis := map[Card]bool{}
	for k := range gHuCardsMap {
		cards := make([]Card, 0, len(k)/2)
		for i := 2; i <= len(k); i += 2 {
			cardStr := k[i-2 : i]
			cardInt, err := strconv.Atoi(cardStr)
			if err != nil {
				assert.Errorf(t, err, err.Error())
			}
			cards = append(cards, Card(cardInt))
		}
		if ok, _ := FastCheckHuV2(cards, laizis, false); !ok {
			err := fmt.Errorf("cards:[%v] can not hu k[%s]", cards, k)
			assert.Errorf(t, err, err.Error())
		}
	}
}

func TestFastCheckHuV2(t *testing.T) {
	//cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	//cards := []Card{14, 15, 16, 25, 26, 26, 27, 27, 28, 29, 29, 35, 36, 37,}
	//cards := []Card{14, 15, 15, 15, 16, 35, 36, 37, 37, 37, 38, 38, 39, 39,}
	//cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19, 21}
	cards := []Card{11, 12, 13, 17, 17, 22, 23, 24, 33, 33, 34, 34, 35, 35}
	laizis := map[Card]bool{21: true}
	ok, _ := FastCheckHuV2(cards, laizis, false)
	assert.True(t, ok)
}

func TestFastCheckHuV2ProfileOnce(t *testing.T) {
	//cards := []Card{17, 46, 32, 22, 28, 13, 37, 34, 14, 32, 35, 34, 33, 0,}
	//cards := []Card{35, 21, 38, 14, 17, 46, 37, 16, 35, 38, 14, 41, 24, 11,}
	cards := []Card{45, 43, 21, 42, 15, 34, 46, 12, 38, 33, 0, 0, 0, 0}
	laizis := map[Card]bool{Laizi: true}
	now := time.Now()
	FastCheckHuV2(cards, laizis, false)
	subCost := time.Now().Sub(now).Nanoseconds()
	fmt.Println(subCost)
}

func TestFastCheckHuV2Profile(t *testing.T) {
	var cost int64 = 0
	var count int64 = 100
	var maxCost int64 = 0
	var minCost int64 = math.MaxInt64
	testMap := make(map[int64][]Card)
	for i := int64(0); i < count; i++ {
		cards := getRandomTestCards(13)
		cards = append(cards, Laizi)
		laizis := map[Card]bool{Laizi: true}
		now := time.Now()
		FastCheckHuV2(cards, laizis, false)
		subCost := time.Now().Sub(now).Nanoseconds()
		if subCost > maxCost {
			maxCost = subCost
		}
		if subCost < minCost {
			minCost = subCost
		}
		testMap[subCost] = cards
		cost += subCost
	}
	fmt.Println("avg: ", cost/count)
	fmt.Println("max: ", maxCost)
	fmt.Println("min: ", minCost)
	for k, v := range testMap {
		fmt.Printf("%v: %v", v, k)
		fmt.Println()
	}
}

func BenchmarkFastCheckHuV2Laizi0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(14)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(13)
		cards = append(cards, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(12)
		cards = append(cards, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(11)
		cards = append(cards, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(10)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(9)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(8)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func BenchmarkFastCheckHuV2Laizi7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(7)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckHuV2(cards, laizis, false)
	}
}

func TestFastCheckTingV2(t *testing.T) {
	cards := []Card{11, 11, 11, 11}
	//cards := []Card{11, 11, 11, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	laizis := map[Card]bool{21: true}
	result := FastCheckTingV2(cards, laizis)
	assert.Nil(t, result)
}

// TestJiulianbaodengTing 测试九莲宝灯听牌
func TestJiulianbaodengTing(t *testing.T) {
	cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19}
	//cards := []Card{11, 11, 11, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	laizis := map[Card]bool{21: true}
	tingCards := FastCheckTingV2(cards, laizis)
	assert.Len(t, tingCards, 9)
	assert.Contains(t, tingCards, []Card{11, 12, 13, 14, 15, 16, 17, 18, 19})
}

func TestFastCheckTingInfoV2(t *testing.T) {
	//map[11:[11 14 17 21] 12:[12 13 16 21] 14:[14 17 21] 15:[15 16 21] 16:[12 13 15 16 21] 19:[11 12 13 14 15 16 17 18 19 21] 13:[12 13 21] 17:[17 21] 18:[18 21]]
	//map[14:[17 14 21] 17:[17 21] 18:[18 21] 19:[12 16 13 15 17 11 14 18 21] 11:[17 11 14 21] 12:[12 16 13 21] 13:[12 13 21] 15:[16 15 21] 16:[12 16 13 15 21]]
	cards := []Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 29, 29, 29, 21}
	laizis := map[Card]bool{}
	tingInfos := FastCheckTingInfoV2(cards, laizis)
	canTingCards := []Card{}
	for card := range tingInfos {
		canTingCards = append(canTingCards, card)
	}
	assert.Nil(t, canTingCards)
}

func TestFastCheckTingInfoV2ProfileOnce(t *testing.T) {
	//cards := []Card{17, 46, 32, 22, 28, 13, 37, 34, 14, 32, 35, 34, 33, 0,}
	//cards := []Card{35, 21, 38, 14, 17, 46, 37, 16, 35, 38, 14, 41, 24, 11,}
	//cards := []Card{45, 43, 21, 42, 15, 34, 46, 12, 38, 33, 0, 0, 0, 0,} //checkHu:6311618 checkTingInfo:1040296114
	//cards := []Card{38, 31, 33, 13, 37, 11, 15, 24, 39, 16, 0, 0, 0, 0}
	cards := []Card{36, 33, 37, 24, 13, 42, 33, 23, 35, 15, 0, 0, 0, 0}
	laizis := map[Card]bool{Laizi: true}
	now := time.Now()
	fmt.Println(FastCheckTingInfoV2(cards, laizis))
	subCost := time.Now().Sub(now).Nanoseconds()
	fmt.Println(subCost)
}

func TestFastCheckTingInfoV2Profile(t *testing.T) {
	var cost int64 = 0
	var count int64 = 100
	var maxCost int64 = 0
	var minCost int64 = math.MaxInt64
	testMap := make(map[int64][]Card)
	for i := int64(0); i < count; i++ {
		cards := getRandomTestCards(10)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		now := time.Now()
		FastCheckTingInfoV2(cards, laizis)
		subCost := time.Now().Sub(now).Nanoseconds()
		if subCost > maxCost {
			maxCost = subCost
		}
		if subCost < minCost {
			minCost = subCost
		}
		testMap[subCost] = cards
		cost += subCost
	}
	fmt.Println("avg: ", cost/count)
	fmt.Println("max: ", maxCost)
	fmt.Println("min: ", minCost)

	var keys sort.IntSlice = make([]int, 0, len(testMap))
	for k := range testMap {
		keys = append(keys, int(k))
	}
	sort.Sort(sort.Reverse(keys))
	for _, key := range keys {
		fmt.Printf("%v: %v", testMap[int64(key)], key)
		fmt.Println()
	}
}

func BenchmarkFastCheckTingInfoV2Laizi0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(14)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(13)
		cards = append(cards, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(12)
		cards = append(cards, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(11)
		cards = append(cards, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(10)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(9)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(8)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}

func BenchmarkFastCheckTingInfoV2Laizi7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cards := getRandomTestCards(7)
		cards = append(cards, Laizi, Laizi, Laizi, Laizi, Laizi, Laizi, Laizi)
		laizis := map[Card]bool{Laizi: true}
		b.StartTimer()
		FastCheckTingInfoV2(cards, laizis)
	}
}
