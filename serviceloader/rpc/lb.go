package rpc

import (
	"fmt"
	"sync"

	"github.com/bsm/grpclb/balancer"
	"github.com/bsm/grpclb/discovery/consul"
	"github.com/hashicorp/consul/api"
)

type loadBalancer struct {
	lbs     *balancer.Server
	lbsInit sync.Once
}

func newLoadBalancer() *loadBalancer {
	return &loadBalancer{}
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

func (lb *loadBalancer) initLbs() {
	config := api.DefaultConfig()
	discovery, err := consul.New(config)
	if err != nil {
		panic(err)
	}
	lb.lbs = balancer.New(discovery, nil)
}
