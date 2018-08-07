package teleport

import (
	"github.com/henrylee2cn/teleport"
	"steve/datareport/conn"
)

type Server struct {
	config conn.Config
	per tp.Peer
}

func NewServer(config conn.Config) *Server{
	return &Server{
		config:config,
	}
}

func (server *Server) Start(){
	server.per = tp.NewPeer(tp.PeerConfig{
		LocalIP  : server.config.Address,
		ListenPort:  uint16(server.config.Port),
		CountTime:   true,
		PrintDetail: true,
	})
	server.per.RouteCall(new(ServerHandle))
	server.per.ListenAndServe()
}

func (server *Server) Stop(){
	server.per.Close()
}