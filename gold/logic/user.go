package logic

import "steve/gold/define"

/*
 功能： 用户金币结构
 作者： SkyWang
 日期： 2018-7-24
*/

type userGold struct {
	uid      uint64          // 玩家ID
	goldList map[int16]int64 // 货币列表
}

// 新建一个userGold
func newUserGold(uid uint64, m map[int16]int64) *userGold {
	return &userGold{
		uid:      uid,
		goldList: m,
	}
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
