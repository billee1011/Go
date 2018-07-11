package main

import (
	"steve/serviceexample/dynamicAnalysis/proto"
	"github.com/tkrajina/go-reflector/reflector"
	"fmt"
)



func GetLen(u interface{}) {
	obj := reflector.New(u)
	resp, _ := obj.Method("GetId").Call()
	fmt.Println(resp.Result)
}

func main() {
	m := msg.Msgplus{
		Id: 123,
		Len: 456,
	}

	GetLen(&m)
}