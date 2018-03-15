package sgrpc

// BindType used for BindService
type BindType string

const (
	// UserBind used for bind users to someserver
	UserBind = BindType("user")
	// DeskBind used for bind desks to someserver
	DeskBind = BindType("desk")
)

type ServerBinder interface {
	BindServer(serverName string, bindType BindType, bindData string, serverID string) error
	GetBindServer(serverName string, bindType BindType, bindData string) string
	InvalidateBind(serverName string, bindType BindType, bindData string) error
}
