package rpc

import (
	"fmt"
	"sync"

	"github.com/bsm/grpclb/balancer"
	"github.com/bsm/grpclb/discovery/consul"
	"github.com/hashicorp/consul/api"
	"github.com/Sirupsen/logrus"
	"time"
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
	logrus.Debugln("getServerAddr>>", serverName,":", servers)
	return servers[0].GetAddress(), nil
}

func (lb *loadBalancer) initLbs() {

	config := api.DefaultConfig()
	config.Address = lb.consulAddr
	discovery, err := consul.New(config)
	if err != nil {
		panic(err)
	}

	c := &balancer.Config{ }
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
