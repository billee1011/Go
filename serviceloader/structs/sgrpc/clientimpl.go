package sgrpc

// TODO 健康检查

import (
	"errors"
	"fmt"
	"steve/serviceloader/structs/sgrpc/consul"
	"sync"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ClientImpl implements sgrpc.Client
type ClientImpl struct {
	caFile          string
	tlsServerName   string
	serviceQueryMap map[string]*serviceQueryData
	mutex           sync.RWMutex
}

// NewClientImpl 创建客户端 Value
func NewClientImpl(caFile, tlsServerName string) *ClientImpl {
	return &ClientImpl{
		tlsServerName:   tlsServerName,
		serviceQueryMap: make(map[string]*serviceQueryData),
		caFile:          caFile,
	}
}

type serviceQueryData struct {
	serviceID  string
	cc         *grpc.ClientConn
	pickCnt    uint64
	serverAddr string
}

// GetClientConnByServerName 根据服务名称获取 grpc 连接对象
func (ci *ClientImpl) GetClientConnByServerName(serverName string) (*grpc.ClientConn, error) {
	serverData, addPickCount, updateConn, err := ci.filterBestServer(serverName)
	if err != nil {
		return nil, err
	}
	if serverData.cc == nil {
		serverData.cc, err = ci.dialRPCServer(serverData.serverAddr)
		if err != nil {
			return nil, fmt.Errorf("连接 RPC 服务 [%s@%s] 失败: %v", serverName, serverData.serverAddr, err)
		}
		updateConn(serverData.cc)
	}

	defer addPickCount()
	return serverData.cc, nil
}

var errServiceNotExist = errors.New("service not exist")

// filterBestServer 选取出最优的服务，并且返回 queryData 的拷贝
// 调用者通过 addPickCount 来添加被选中次数， 通过 updateConn 来更新 grpc 客户端端连接信息
// TODO : 外部并发调用 updateConn 只会保留最后一个， 但不是问题
func (ci *ClientImpl) filterBestServer(serverName string) (serverData serviceQueryData, addPickCount func(), updateConn func(cc *grpc.ClientConn), err error) {
	serviceDatas := consul.GetServiceDatasByName(serverName)

	var picked *serviceQueryData
	for _, serviceData := range serviceDatas {
		qd := ci.getServiceQueryData(&serviceData)
		if picked == nil || qd.pickCnt < picked.pickCnt {
			picked = qd
		}
	}
	if picked == nil {
		return serviceQueryData{}, nil, nil, errServiceNotExist
	}
	return *picked, func() {
			ci.mutex.Lock()
			defer ci.mutex.Unlock()
			picked.pickCnt++
		}, func(cc *grpc.ClientConn) {
			ci.mutex.Lock()
			defer ci.mutex.Unlock()
			picked.cc = cc
		}, nil
}

func (ci *ClientImpl) dialRPCServer(addr string) (*grpc.ClientConn, error) {
	entry := logrus.WithFields(logrus.Fields{
		"name":    "ClientImpl.dialRPCServer",
		"address": addr,
		"ca_file": ci.caFile,
	})
	opts := []grpc.DialOption{}
	if ci.caFile != "" {
		c, err := credentials.NewClientTLSFromFile(ci.caFile, ci.tlsServerName)
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

func (ci *ClientImpl) getServiceQueryData(sd *consul.ServiceData) *serviceQueryData {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()
	var qd *serviceQueryData
	var found bool
	if qd, found = ci.serviceQueryMap[sd.ServiceID]; !found {
		qd = &serviceQueryData{
			serviceID:  sd.ServiceID,
			serverAddr: sd.Addr + ":" + sd.Port,
		}
		ci.serviceQueryMap[sd.ServiceID] = qd
	}
	return qd
}

// GetClientConnByServerID 根据服务 ID 获取 grpc 客户端连接
func (ci *ClientImpl) GetClientConnByServerID(serverID string) (*grpc.ClientConn, error) {
	sd := consul.GetServiceDataByID(serverID)
	if sd == nil {
		return nil, nil
	}
	qd := ci.getServiceQueryData(sd)
	if qd.cc == nil {
		var err error
		if qd.cc, err = ci.dialRPCServer(qd.serverAddr); err != nil {
			return qd.cc, fmt.Errorf("dial rpc server [%s] faield.%v", qd.serverAddr, err)
		}
	}
	return qd.cc, nil
}
