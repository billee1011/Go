package logic

import (
	"time"
	"github.com/Sirupsen/logrus"
	"steve/external/propsclient"
)

func startTestPropsServer() {
	testPropsGold()
}


func testPropsGold() {

	begin := time.Now().UnixNano()
	uid := uint64(11)
	v , err := propsclient.GetUserProps(uid, 0)
	logrus.Infof("GetUserProps=%d, err=%v++++++++++", v, err)

	m := make(map[uint64]int64)
	m[1001] = 100
	m[1002] = 105
	err = propsclient.AddUserProps(uid, m, 1, 100, 0, 0)
	logrus.Debugf("uid=%d, AddGold=%d, err=%v",uid, uid,   err)
	v , err = propsclient.GetUserProps(uid, 0)
	logrus.Infof("GetUserProps=%v, err=%v++++++++++", v, err)
	end := time.Now().UnixNano()
	logrus.Infof("GetUserProps=%d,  useTime=%d(ms), err=%v------------------------",uid,   (end - begin)/1000000, err )
}
