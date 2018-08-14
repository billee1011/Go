package logic

import (
	"github.com/Sirupsen/logrus"
	"steve/propserver/data"
	"steve/propserver/define"
	"sync"
	"steve/external/configclient"
)

/*
  功能： 道具管理： 加减玩家道具，获取玩家道具,交易序列号去重. 支持redis，db同步存储。交易流水日志对账等.
  作者： SkyWang
  日期： 2018-8-13
*/

var myLogic PropsMgr

func GetMyLogic() *PropsMgr {
	return &myLogic
}

type PropsMgr struct {
	userList sync.Map                 // 用户列表
	muLock   map[uint64]*sync.RWMutex // 用户锁，一个用户一个锁

	propsList map[uint64]*propsInfo
}

func (gm *PropsMgr) Init() error {
	//goldMgr.userList = make(map[uint64]*userGold)
	gm.muLock = make(map[uint64]*sync.RWMutex)
	gm.propsList = make(map[uint64]*propsInfo, 10)

	return gm.getPropsListFromDB()
}

func (gm *PropsMgr) getPropsListFromDB() error {
	strJson, err := configclient.GetConfig("prop", "interactive")
	if err != nil {
		logrus.Errorf("GetPropsListFromDB from config err:", err)
		return err
	}

	return gm.parseJsonPropsList(strJson)

}

func (gm *PropsMgr) parseJsonPropsList(strJson string) error {

	return nil
}

func (gm *PropsMgr) GetMutex(uid uint64) *sync.RWMutex {
	if mu, ok := gm.muLock[uid]; ok {
		return mu
	}
	n := new(sync.RWMutex)
	gm.muLock[uid] = n
	return n
}

// 加玩家道具
func (gm *PropsMgr) AddUserProps(uid uint64, propList map[uint64]int64, seq string, funcId int32, channel int64, createTm int64, gameId, level int32) error {
	// 1. 先获取玩家当前金币值, GetGold()
	// 2. 在内存中对玩家金币进行加减
	// 3. 将变化后的值写到redis和DB
	before := int64(0)
	after := int64(0)

	entry := logrus.WithFields(logrus.Fields{
		"opr":        "add_props",
		"gameId":     gameId,
		"level":      level,
		"uid":        uid,
		"funcId":     funcId,
		"propList":     propList,
		"channel":    channel,
		"seq":        seq,
		"createTime": createTm,
	})

	for propId := range  propList {
		if !gm.checkPropId(propId) {
			entry.Errorln("propId error")
			return define.ErrPropId
		}
	}


	// 按用户ID进行加锁,一个用户一个锁
	mu := gm.GetMutex(uid)
	mu.Lock()
	defer mu.Unlock()

	u, err := gm.getUser(uid)
	if u == nil {
		entry.Errorln("get user error")
		_ = err
		return  define.ErrNoUser
	}

	// 判断交易流水号是否有冲突?
	if !u.CheckSeq(seq) {
		entry.Errorf("seq is same: uid=%d, seq=%s", uid, seq)
		return  define.ErrSeqNo
	}

	// 加道具前，玩家当前道具数量
	for propId, num := range  propList {
		if num >= 0 {
			continue
		}
		before, _ = u.Get(propId)

		if before+num < 0 {
			entry.Errorf("prop num < value: uid=%d, before=%d, add=%d", uid, before, num)
			return define.ErrNoProp
		}
	}

	// 加道具后，玩家当前道具数量
	for propId, num := range  propList {
		before, _ = u.Get(propId)
		if num != 0 {
			after, _ = u.Add(propId, num)
		}


		propList[propId] = after

		entry = logrus.WithFields(logrus.Fields{
			"opr":        "add_props",
			"gameId":     gameId,
			"level":      level,
			"uid":        uid,
			"funcId":     funcId,
			"propId":     propId,
			"before":     before,
			"changed":    num,
			"after":      after,
			"channel":    channel,
			"seq":        seq,
			"createTime": createTm,
		})
		// 交易记录写到日志
		entry.Infoln("add succeed")
	}

	// 交易记录写到redis
	// 交易记录写到DB
	err = gm.saveUserToCacheAndDB(entry, uid, propList)
	if err != nil {
		entry.Errorln("saveUserToCacheAndDB error: ", err)
	}

	return nil
}

// 获取玩家道具
func (gm *PropsMgr) GetUserProps(uid uint64, propId uint64) (map[uint64]int64, error) {
	// 1.先在内存中查找玩家是否存在。
	// 2.不存在，从Redis获取玩家道具.
	// 3.不存在，从DB获取玩家道具.

	if !gm.checkPropId(propId) {
		logrus.Errorf("for={prop id error},uid=%d,goldType=%d", uid, propId)
		return nil, define.ErrNoProp
	}

	// 按用户ID进行加锁,一个用户一个锁
	mu := gm.GetMutex(uid)
	mu.Lock()
	defer mu.Unlock()

	u, _ := gm.getUser(uid)
	if u == nil {
		return nil, define.ErrNoUser
	}
	// 获取玩家指定道具
	g, err := u.GetList(propId)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// 保存玩家变化到Redis和DB
func (gm *PropsMgr) saveUserToCacheAndDB(entry *logrus.Entry, uid uint64, list map[uint64]int64) error {

	// 暂时先保存到Redis
	err := data.SavePropsToRedis(uid, list)
	if err != nil {
		// 记录redis写入失败
		entry.Errorln("SavePropsToRedis error", err)
	}

	// 后续再保存到DB
	for propId, num := range  list {
		err = data.SavePropsToDB(uid, propId, num)
		if err != nil {
			// 记录DB写入失败
			entry.Errorln("SavePropsToDB error:", err)
		}
	}

	return nil
}

// 获取用户
func (gm *PropsMgr) getUser(uid uint64) (*userProps, error) {
	if uid == 0 {
		return nil, nil
	}
	u, ok := gm.userList.Load(uid)
	if !ok {
		return gm.getUserFromCacheOrDB(uid)
	}
	return u.(*userProps), nil
}

// 新建用户
func (gm *PropsMgr) newUser(uid uint64, m map[uint64]int64) *userProps {
	n := newUserProps(uid, m)
	gm.userList.Store(uid, n)
	return n
}

// 从Redis或者DB获取用户
func (gm *PropsMgr) getUserFromCacheOrDB(uid uint64) (*userProps, error) {
	m, err := data.LoadPropsFromRedis(uid)
	if err == nil {
		return gm.newUser(uid, m), nil
	}

	m, err = data.LoadPropsFromDB(uid)
	if err != nil {
		return nil, define.ErrLoadDB
	}
	// 从DB获取到后，马上缓存到Redis
	err = data.SavePropsToRedis(uid, m)
	if err != nil {
		// 记录redis写入失败
		logrus.Errorln("save redis error")
	}
	return gm.newUser(uid, m), nil
}

// 检测道具ID是否有效
func (gm *PropsMgr) checkPropId(propId uint64) bool {

	if _, ok := gm.propsList[propId]; ok {
		return true
	}
	// 先不判断道具是否存在
	return true
	return false
}
