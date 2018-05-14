package peipai

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	// ctrl := gomock.NewController(t)
	addr := "192.168.8.148:8080"
	Run(addr)
}

func TestLogPeiPaiInfos(t *testing.T) {
	LogPeiPaiInfos()
}

func TestGetPeiPai(T *testing.T) {
	value, err := GetPeiPai("scxl")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}
}
