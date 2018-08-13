package main

import (
	"testing"
	"runtime"
	"sync"
)

func Test_Routi(t *testing.T) {
	runtime.GOMAXPROCS(1)

	fail := 0

	for m := 0; m < 30000; m++ {
		wg := sync.WaitGroup{}

		uid := 1300
		t.Logf("Test_GoRoutine...beginId=%d", uid)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				uid++
			}(i)

		}
		wg.Wait()
		if uid == 2300 {
			t.Logf("Test_GoRoutine ... endId=%d",uid)
		} else {
			t.Errorf("Test_GoRoutine ... endId=%d",uid)
			fail++
		}

	}

	t.Logf("Test_GoRoutine finished ... failSum=%d",fail)
}