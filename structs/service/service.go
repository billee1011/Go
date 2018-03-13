package service

import "steve/structs"

// Service is an interface must be implemented and exposed by blugin who want to use serviceloader
// services also need to expose "GetService" function by which return the "Service"
type Service interface {
	Start(e *structs.Exposer, param ...string) error
}

type GetServiceFuncType func() Service
