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

func (e *exchangerImpl) SendPackageByPlayerID(playerID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	return e.sender.send(playerID, head, body)
}
func (e *exchangerImpl) BroadcastPackageByPlayerID(playerIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	return e.sender.broadcast(playerIDs, head, body)
}
func (e *exchangerImpl) SendPackageBareByPlayerID(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	return e.sender.sendBare(playerID, head, bodyData)
}
func (e *exchangerImpl) BroadcastPackageBareByPlayerID(playerIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	return e.sender.broadcastBare(playerIDs, head, bodyData)
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
