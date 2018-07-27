package facade

import (
	"fmt"
	"steve/login/interfaces"
	"steve/structs/proto/base"

	"github.com/golang/protobuf/proto"
)

// SendPackage send package to client
func SendPackage(sender interfaces.MessageSender, clientID uint64, header *base.Header, message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("message marshal failed: %v", err)
	}
	return sender.SendMessage(clientID, header, data)
}
