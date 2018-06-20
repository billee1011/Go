package exchanger

import (
	iexchanger "steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"steve/structs/sgrpc"

	"github.com/golang/protobuf/proto"
)

type exchangerImpl struct {
	handlerMgr
	sender   sender
	receiver receiver
}

var _ iexchanger.Exchanger = new(exchangerImpl)

func (e *exchangerImpl) SendPackage(clientID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	return e.sender.send(clientID, head, body)
}

func (e *exchangerImpl) BroadcastPackage(clientIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	return e.sender.broadcast(clientIDs, head, body)
}

func (e *exchangerImpl) SendPackageBare(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	return e.sender.sendBare(clientID, head, bodyData)
}

func (e *exchangerImpl) BroadcastPackageBare(clientIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	return e.sender.broadcastBare(clientIDs, head, bodyData)
}

// NewExchanger 创建 Exchanger
func NewExchanger(rpcServer sgrpc.RPCServer) (iexchanger.Exchanger, error) {
	e := exchangerImpl{}
	e.receiver.handlerMgr = &e.handlerMgr
	if err := rpcServer.RegisterService(steve_proto_gaterpc.RegisterMessageHandlerServer, &e.receiver); err != nil {
		return nil, err
	}
	return &e, nil
}
