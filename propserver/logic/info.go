package logic

import "fmt"

/*
 功能： 用户金币结构
 作者： SkyWang
 日期： 2018-7-24
*/

const SEQ_LIST_FULL  =  10
const SEQ_LIST_HALF  =  SEQ_LIST_FULL/2


// 道具信息
type propsInfo struct {
	attrType  int32			// 属性类型
	attrId    uint64		// 类型ID
	attrValue int64			// 属性值
	attrLimit int64			// 叠加上限
}


type userProps struct {
	uid      uint64          		// 玩家ID
	propsList map[uint64]int64 		// 道具列表
	lastSeqList map[string]bool  	// 最近交易序列号
	lastSeqList2 map[string]bool 	// 最近交易序列号缓存,双环存，先填第一个，第一个满后，填第2个，第2个满后，清空第一个.
	bFirstSeqList bool			 	// 是否是第一个消息队列
}

// 新建一个userProps
func newUserProps(uid uint64, m map[uint64]int64) *userProps {
	return &userProps{
		uid:      uid,
		propsList: m,
		bFirstSeqList: true,
		lastSeqList :map[string]bool{},
		lastSeqList2 :map[string]bool{},
	}
}

func (ug *userProps) CheckSeq(seq string) bool {
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

// 对指定道具加
func (ug *userProps) Add(propId uint64, value int64) (int64, error) {

	ug.propsList[propId] += value
	return ug.propsList[propId], nil
}

// 对指定道具
func (ug *userProps) Get(propId uint64) (int64, error) {

	return ug.propsList[propId], nil
}

// 对指定道具
func (ug *userProps) GetList(propId uint64) (map[uint64]int64, error) {
	if propId > 0 {
		num, ok := ug.propsList[propId]
		if ok {
			m := make(map[uint64]int64, 1)
			m[propId] = num
			return m, nil
		}
		return nil, fmt.Errorf("no the props")
	}
	return ug.propsList, nil
}