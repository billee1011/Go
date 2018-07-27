package logic

import (
	"github.com/Sirupsen/logrus"
	"steve/gold/data"
	"steve/gold/define"
)

/*
  功能： 金币管理： 加减金币，获取金币.支持redis，db同步存储。
  作者： SkyWang
  日期： 2018-7-24
*/

var goldMgr GoldMgr

func GetGoldMgr() *GoldMgr {
	return &goldMgr
}

type GoldMgr struct {
	userList map[uint64]*userGold // 用户列表
}

func init() {
	goldMgr.userList = make(map[uint64]*userGold)
}

// 支持的货币类型
var gTypeList = map[int16]string{
	1: "coins",
	2: "ingots",
	3: "keyCards",
}

func init() {
	data.SetGoldTypeList(gTypeList)
}

// 加金币
func (gm *GoldMgr) AddGold(uid uint64, goldType int16, value int64, seq string, funcId int32, channel int64, createTm int64) (int64, error) {
	// 1. 先获取玩家当前金币值, GetGold()
	// 2. 在内存中对玩家金币进行加减
	// 3. 将变化后的值写到redis和DB
	before := int64(0)
	after := int64(0)

	entry := logrus.WithFields(logrus.Fields{
		"opr":        "add_gold",
		"uid":        uid,
		"funcId":     funcId,
		"goldType":   goldType,
		"before":     before,
		"changed":    value,
		"after":      after,
		"channel":    channel,
		"seq":        seq,
		"createTime": createTm,
	})

	if !gm.checkGoldType(goldType) {
		entry.Errorln("gold type error")
		return 0, define.ErrGoldType
	}

	// 判断交易流水号是否有冲突?

	u, err := gm.getUser(uid)
	if u == nil {
		entry.Errorln("get user error")
		_ = err
		return 0, define.ErrNoUser
	}

	// 加金币前，玩家当前金币值
	before, err = u.Get(goldType)
	if err != nil {
		entry.Errorln("get gold error")
		return 0, err
	}

	// 加金币后，玩家当前金币值
	after, err = u.Add(goldType, value)
	if err != nil {
		entry.Errorln("add opr error")
		return 0, err
	}

	// 交易记录写到日志
	entry.Infoln("add succeed")

	// 交易记录写到redis
	// 交易记录写到DB
	err = gm.saveUserToCacheAndDB(entry, u, goldType)
	if err != nil {
		entry.Errorln("save cacheordb error")
	}

	return after, nil
}

// 获取金币
func (gm *GoldMgr) GetGold(uid uint64, goldType int16) (int64, error) {
	// 1.先在内存中查找玩家是否存在。
	// 2.不存在，从Redis获取玩家金币.
	// 3.不存在，从DB获取玩家金币.

	if !gm.checkGoldType(goldType) {
		logrus.Errorf("for={gold type error},uid=%d,goldType=%d", uid, goldType)
		return 0, define.ErrGoldType
	}

	u, _ := gm.getUser(uid)
	if u == nil {
		return 0, define.ErrNoUser
	}
	// 获取玩家金币
	g, err := u.Get(goldType)
	if err != nil {
		return 0, err
	}

	return g, nil
}

// 保存玩家变化到Redis和DB
func (gm *GoldMgr) saveUserToCacheAndDB(entry *logrus.Entry, u *userGold, goldType int16) error {

	// 暂时先保存到Redis
	list := u.goldList
	if goldType >= 0 {
		list = make(map[int16]int64)
		list[goldType] = u.goldList[goldType]
	}
	err := data.SaveGoldToRedis(u.uid, list)
	if err != nil {
		// 记录redis写入失败
		entry.Errorln("save redis error")
	}

	// 后续再保存到DB
	err = data.SaveGoldToDB(u.uid, list)
	if err != nil {
		// 记录DB写入失败
		entry.Errorln("save db error")
	}
	return nil
}

// 获取用户
func (gm *GoldMgr) getUser(uid uint64) (*userGold, error) {
	if uid == 0 {
		return nil, nil
	}
	u := gm.userList[uid]
	if u == nil {
		return gm.getUserFromCacheOrDB(uid)
	}
	return u, nil
}

// 新建用户
func (gm *GoldMgr) newUser(uid uint64, m map[int16]int64) *userGold {
	n := newUserGold(uid, m)
	gm.userList[uid] = n
	return n
}

// 从Redis或者DB获取用户
func (gm *GoldMgr) getUserFromCacheOrDB(uid uint64) (*userGold, error) {
	m, err := data.LoadGoldFromRedis(uid)
	if err == nil {
		return gm.newUser(uid, m), nil
	}

	m, err = data.LoadGoldFromDB(uid)
	if err != nil {
		return nil, define.ErrLoadDB
	}
	// 从DB获取到后，马上缓存到Redis
	err = data.SaveGoldToRedis(uid, m)
	if err != nil {
		// 记录redis写入失败
		logrus.Errorln("save redis error")
	}
	return gm.newUser(uid, m), nil
}

// 检测货币类型是否有效
func (gm *GoldMgr) checkGoldType(goldType int16) bool {
	if _, ok := gTypeList[goldType]; ok {
		return true
	}
	return false
}
