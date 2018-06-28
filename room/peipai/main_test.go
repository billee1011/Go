package peipai

import (
	"fmt"
	"steve/room/peipai/handle"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// ctrl := gomock.NewController(t)
	addr := "192.168.8.148:8080"
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		Run(addr)
		defer wg.Done()
	}()
	go func() {
		defer wg.Done()
		for {
			TestLogPeiPaiInfos(t)
			TestLogOptionInfos(t)
			time.Sleep(time.Second * 5)
		}
	}()
	wg.Wait()
}

func TestLogPeiPaiInfos(t *testing.T) {
	handle.LogPeiPaiInfos()
}

func TestLogOptionInfos(t *testing.T) {
	handle.LogOptionInfos()
}

func TestGetPeiPai(T *testing.T) {
	value, err := handle.GetPeiPai("scxl")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}
}
