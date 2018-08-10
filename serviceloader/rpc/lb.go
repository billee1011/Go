package rpc

import (
	"fmt"
	"sync"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
	"steve/thirdpart/github.com/bsm/grpclb/balancer"
	"steve/thirdpart/github.com/bsm/grpclb/discovery/consul"
)

type loadBalancer struct {
	lbs        *balancer.Server
	lbsInit    sync.Once
	consulAddr string
}

func newLoadBalancer(consulAddr string) *loadBalancer {
	return &loadBalancer{
		consulAddr: consulAddr,
	}
}

func (lb *loadBalancer) getServerAddr(serverName string) (string, error) {
	lb.lbsInit.Do(lb.initLbs)

	servers, err := lb.lbs.GetServers(serverName)
	if err != nil {
		return "", err
	}
	if len(servers) == 0 {
		return "", fmt.Errorf("no server")
	}
	return servers[0].GetAddress(), nil
}

// 通过制定服务ID获取服务连接地址
func (lb *loadBalancer) getServerAddrByServerId(serverName string, serverId string) (string, error) {
	lb.lbsInit.Do(lb.initLbs)

	servers, err := lb.lbs.GetServers(serverName)
	if err != nil {
		return "", err
	}
	if len(servers) == 0 {
		return "", fmt.Errorf("no server")
	}

	for _, addr := range servers {
		strScore := fmt.Sprintf("%d", addr.GetScore())
		if strScore == serverId {
			return addr.GetAddress(), nil
		}
	}
	// logrus.Errorf("err={find no server},server=%s,svrId=%s", serverName, serverId)
	return "", fmt.Errorf("no server")
}

// 通过HashID获取服务连接地址
func (lb *loadBalancer) getServerAddrByHashId(serverName string, hashId uint64) (string, error) {
	lb.lbsInit.Do(lb.initLbs)

	servers, err := lb.lbs.GetServers(serverName)
	if err != nil {
		return "", err
	}
	if len(servers) == 0 {
		return "", fmt.Errorf("no server")
	}
	// logrus.Debugf("server=%s,hashId=%d", serverName, hashId)

	svrSum := uint64(len(servers))
	index := hashId % svrSum
	for _, addr := range servers {
		if addr.GetScore() == int64(index) {
			return addr.GetAddress(), nil
		}
	}
	logrus.Errorf("err={find no server},server=%s,hashId=%d,index=%d", serverName, hashId, index)
	return "", fmt.Errorf("no server")
}

func (lb *loadBalancer) initLbs() {

	config := api.DefaultConfig()
	config.Address = lb.consulAddr
	discovery, err := consul.New(config)
	if err != nil {
		panic(err)
	}

	c := &balancer.Config{}
	// 从consul探测服务列表的频率
	if c.Discovery.Interval == 0 {
		c.Discovery.Interval = 5 * time.Second
	}
	// 检测服务进程是否正常的RPC请求频率
	if c.LoadReport.Interval == 0 {
		c.LoadReport.Interval = 3 * time.Second
	}
	// 认为服务不可用的RPC请求失败最大次数
	if c.LoadReport.MaxFailures == 0 {
		c.LoadReport.MaxFailures = 2
	}
	lb.lbs = balancer.New(discovery, c)
}
