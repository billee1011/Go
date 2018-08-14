package logic

import (
	"steve/external/goldclient"
	"github.com/Sirupsen/logrus"
	"time"
	"sync"
)

func startTestGoldServer() {
	testAddGold()
}


func testAddGold() {

	begin := time.Now().UnixNano()
	uid := uint64(18)
	v , err := goldclient.GetGold(uid, 1)
	logrus.Infof("getGold=%d, err=%v++++++++++", v, err)

	last, err := goldclient.AddGold(uid, 1, 1, 1, 100, 0, 0)
	logrus.Debugf("uid=%d, AddGold=%d, curGold=%v, err=%v",uid, uid,  last, err)
	v , err = goldclient.GetGold(uid, 1)
	logrus.Infof("getGold=%d, err=%v++++++++++", v, err)
	end := time.Now().UnixNano()
	logrus.Infof("AddGold=%d, curGold=%v, useTime=%d(ms), err=%v------------------------",uid,  last, (end - begin)/1000000, err )
}

func testAddGoldBig() {

	wg := sync.WaitGroup{}
	begin := time.Now().UnixNano()
	uid := uint64(1300)
	v , err := goldclient.GetGold(uid, 1)
	logrus.Infof("getGold=%d, err=%v++++++++++", v, err)

	uid = 1300
	last := v
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			v , err := goldclient.GetGold(uid+uint64(idx), 1)
			logrus.Infof("getGold=%d, err=%v++++++++++", v, err)
			//last, err := goldclient.AddGold(uid+uint64(idx), 1, 1, 1, 100)
			//logrus.Debugf("uid=%d, AddGold=%d, curGold=%v, err=%v",uid+uint64(idx), uid,  last, err)
		}(i)

	}
	wg.Wait()

	last, err = goldclient.AddGold(uid, 1, 1, 1, 100, 0, 0)
	logrus.Debugf("uid=%d, AddGold=%d, curGold=%v, err=%v",uid, uid,  last, err)
	v , err = goldclient.GetGold(uid, 1)
	logrus.Infof("getGold=%d, err=%v++++++++++", v, err)
	end := time.Now().UnixNano()
	logrus.Infof("AddGold=%d, curGold=%v, useTime=%d(ms), err=%v------------------------",uid,  last, (end - begin)/1000000, err )
}