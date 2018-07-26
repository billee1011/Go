package sgrpc

import (
	"fmt"
	"net"
	"reflect"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewRPCServer(opts ...grpc.ServerOption) *RPCServerImpl {
	return &RPCServerImpl{
		svr: grpc.NewServer(opts...),
	}
}

type RPCServerImpl struct {
	svr *grpc.Server
	isRetire bool			// 是否退休
	score  int64 			// 负载值
}



func isValidRegister(f interface{}, service interface{}) error {
	tf := reflect.TypeOf(f)
	if tf.Kind() != reflect.Func {
		return fmt.Errorf("RPCServerImpl.RegisterService param is not func。it's %v", tf.Kind())
	}
	if tf.NumIn() != 2 {
		return fmt.Errorf("RPCServerImpl.RegisterService %v should have 2 parameters", tf)
	}
	type0 := tf.In(0)
	if type0 != reflect.TypeOf(&grpc.Server{}) {
		return fmt.Errorf("RPCServerImpl.RegisterService %v should have error first parameter type ", tf)
	}
	if !reflect.TypeOf(service).AssignableTo(tf.In(1)) {
		return fmt.Errorf("RPCServerImpl.RegisterService error service type or func type")
	}
	return nil
}

// 设置退休
func (rsi *RPCServerImpl) EnableRetire(bEnable bool) {
	rsi.isRetire = bEnable
}

// 是否退休
func (rsi *RPCServerImpl) IsRetire() bool {
	return rsi.isRetire
}

// 设置负载值
func (rsi *RPCServerImpl) SetScore(score int64) {
	rsi.score = score
}
// 获取负载值
func (rsi *RPCServerImpl) GetScore() int64 {
	return rsi.score
}

func (rsi *RPCServerImpl) RegisterService(f interface{}, service interface{}) error {
	if err := isValidRegister(f, service); err != nil {
		return err
	}
	logEntry := logrus.WithField("func_name", "RPCServerImpl.RegisterService.register")
	logEntry.WithFields(logrus.Fields{
		"register_func": reflect.TypeOf(f),
		"service":       reflect.TypeOf(service),
	}).Infoln("注册服务")
	reflect.ValueOf(f).Call([]reflect.Value{
		reflect.ValueOf(rsi.svr), reflect.ValueOf(service),
	})
	return nil
}

func (rsi *RPCServerImpl) startServer(lis net.Listener) error {
	return rsi.svr.Serve(lis)
}

func (rsi *RPCServerImpl) stopServer(graceful bool) {
	if graceful {
		rsi.svr.GracefulStop()
	} else {
		rsi.svr.Stop()
	}
}

func (rsi *RPCServerImpl) Work(addr string, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	return rsi.startServer(lis)
}
