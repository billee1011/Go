package logic

import "steve/gold/define"

/*
 功能： 用户金币结构
 作者： SkyWang
 日期： 2018-7-24
*/

const SEQ_LIST_FULL  =  10
const SEQ_LIST_HALF  =  SEQ_LIST_FULL/2

type userGold struct {
	uid      uint64          // 玩家ID
	goldList map[int16]int64 // 货币列表
	lastSeqList map[string]bool  // 最近交易序列号
	lastSeqList2 map[string]bool // 最近交易序列号缓存,双环存，先填第一个，第一个满后，填第2个，第2个满后，清空第一个.
	bFirstSeqList bool			 // 是否是第一个消息队列
}

// 新建一个userGold
func newUserGold(uid uint64, m map[int16]int64) *userGold {
	return &userGold{
		uid:      uid,
		goldList: m,
		bFirstSeqList: true,
		lastSeqList :map[string]bool{},
		lastSeqList2 :map[string]bool{},
	}
}

func (ug *userGold) CheckSeq(seq string) bool {
	if ug.lastSeqList[seq] {
		return false
	}
	if ug.lastSeqList2[seq] {
		return false
	}
	if ug.bFirstSeqList {
		if len(ug.lastSeqList) < SEQ_LIST_FULL {
			ug.lastSeqList[seq] = true
		} else {
			ug.lastSeqList2[seq] = true
			ug.bFirstSeqList = false
		}
	} else {
		if len(ug.lastSeqList2) < SEQ_LIST_FULL {
			ug.lastSeqList2[seq] = true
		} else {
			ug.lastSeqList[seq] = true
			ug.bFirstSeqList = true
		}
	}
	// 如果一个队列满，另一个队列达到一半容量，清空队列满的队列.
	if len(ug.lastSeqList) >= SEQ_LIST_FULL && len(ug.lastSeqList2) >= SEQ_LIST_HALF {
		ug.lastSeqList = make(map[string]bool, SEQ_LIST_FULL)
	} else if len(ug.lastSeqList2) >= SEQ_LIST_FULL && len(ug.lastSeqList) >= SEQ_LIST_HALF {
		ug.lastSeqList2 = make(map[string]bool, SEQ_LIST_FULL)
	}
	return true
}

// 对指定货币加金币
func (ug *userGold) Add(goldType int16, value int64) (int64, error) {
	if goldType < 0 || goldType > 1000 {
		return 0, define.ErrGoldType
	}

	// 可能需要判断加减金币后，金币值变成负值！

	ug.goldList[goldType] += value

	return ug.goldList[goldType], nil
}

// 对指定货币加金币
func (ug *userGold) Get(goldType int16) (int64, error) {
	if goldType < 0 || goldType > 1000 {
		return 0, define.ErrGoldType
	}

	return ug.goldList[goldType], nil
}
