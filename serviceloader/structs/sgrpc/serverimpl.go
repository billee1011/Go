package sgrpc

import (
	"fmt"
	"net"
	"reflect"

	"google.golang.org/grpc"
)

func NewRPCServer(opts ...grpc.ServerOption) *RPCServerImpl {
	return &RPCServerImpl{
		svr:       grpc.NewServer(opts...),
		registers: []interface{}{},
	}
}

type RPCServerImpl struct {
	svr       *grpc.Server
	registers []interface{} // registers
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

func (rsi *RPCServerImpl) RegisterService(f interface{}, service interface{}) error {
	if err := isValidRegister(f, service); err != nil {
		return err
	}
	rsi.registers = append(rsi.registers, func() {
		vf := reflect.ValueOf(f)
		vf.Call([]reflect.Value{reflect.ValueOf(rsi.svr), reflect.ValueOf(service)})
	})
	return nil
}

func (rsi *RPCServerImpl) startServer(lis net.Listener) error {
	for _, register := range rsi.registers {
		reflect.ValueOf(register).Call([]reflect.Value{})
	}
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
	if 0 == len(rsi.registers) {
		return nil
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	return rsi.startServer(lis)
}
