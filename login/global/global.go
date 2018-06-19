package global

import "steve/login/interfaces"

var gMessageSender interfaces.MessageSender

// SetMessageSender set global message sender
func SetMessageSender(s interfaces.MessageSender) {
	gMessageSender = s
}

// GetMessageSender get global message sender
func GetMessageSender() interfaces.MessageSender {
	return gMessageSender
}
