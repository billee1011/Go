package teleport

import (
	"github.com/henrylee2cn/teleport"
	"steve/datareport/bean"
	"strconv"
	"steve/datareport/conn"
)

type Client struct {
	config  conn.Config
	per     tp.Peer
	session tp.Session
}

func NewClient(config conn.Config) *Client{
	return &Client{
		config:config,
	}
}
func (client *Client) Connect() {
	//tp.SetLoggerLevel("ERROR")
	client.per = tp.NewPeer(tp.PeerConfig{})
	sess, err := client.per.Dial(client.config.Address + ":" + strconv.Itoa(client.config.Port))
	if err != nil {
		tp.Fatalf("---------->%v", err)
		return
	}
	client.session = sess
}

func (client *Client) Stop() {
	client.per.Close()
}

func (client *Client) Send(log *bean.LogBean) {
	var result int
	client.session.Call(
		"/server_handle/report",
		log,
		&result,
	).Rerror()
}
/*
type push struct {
	tp.PushCtx
}

func (p *push) Status(arg *string) *tp.Rerror {
	tp.Printf("server status: %s", *arg)
	return nil
}
*/