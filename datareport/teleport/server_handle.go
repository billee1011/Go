package teleport

import (
	"github.com/henrylee2cn/teleport"
	"steve/datareport/bean"
	"fmt"
	"steve/datareport/queue"
)

type ServerHandle struct {
	tp.CallCtx
}

func (report *ServerHandle) Report(log *bean.LogBean) (int, *tp.Rerror) {
	fmt.Printf("params:%v\n", log)
	queue.GetLogQueue().Put(log)
	return 0, nil
}