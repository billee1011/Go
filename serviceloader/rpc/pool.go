package rpc

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/connectivity"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type connectPool struct {
	connects      sync.Map
	connectMu     sync.Mutex
	caFile        string
	tlsServerName string
}

func newConnectPool(caFile string, tlsServerName string) *connectPool {
	return &connectPool{
		caFile:        caFile,
		tlsServerName: tlsServerName,
	}
}

func (cp *connectPool) getConnect(addr string) (*grpc.ClientConn, error) {
	ico, ok := cp.connects.Load(addr)
	if ok {
		co := ico.(*grpc.ClientConn)
		if co.GetState() == connectivity.Ready {
			return co, nil
		}
		cp.connects.Delete(addr)
	}
	return cp.newConnect(addr)
}

func (cp *connectPool) newConnect(addr string) (*grpc.ClientConn, error) {
	cp.connectMu.Lock()
	defer cp.connectMu.Unlock()

	co, err := cp.connect(addr)
	if err != nil {
		cp.connects.Store(addr, co)
	}
	return co, err
}

func (cp *connectPool) connect(addr string) (*grpc.ClientConn, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "connectPool.connect",
		"address":   addr,
		"ca_file":   cp.caFile,
	})
	opts := []grpc.DialOption{}
	if cp.caFile != "" {
		c, err := credentials.NewClientTLSFromFile(cp.caFile, cp.tlsServerName)
		if err != nil {
			return nil, fmt.Errorf("create client tls failed:%v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(c))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	entry.Info("dial rpc server")
	cc, err := grpc.Dial(addr, opts...)
	return cc, err
}
