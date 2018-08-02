package fixed

import "errors"

const (
	// ListenClientAddr 代表监听客户端的IP地址，默认值为 127.0.0.1
	ListenClientAddr = "lis_client_addr"

	// ListenClientPort 代表监听客户端的端口， 默认值为 36001
	ListenClientPort = "lis_client_port"

	// ListenPeipaiAddr 配牌監聽地址
	ListenPeipaiAddr = "peipai_addr"

	// XingPaiTimeOut 行牌超时时间，单位为second，默认值为 10
	XingPaiTimeOut = "xp_timeout"

	// HuStateTimeOut 胡牌状态下的超时时间，单位为second，默认值为 1
	HuStateTimeOut = "hs_timeout"

	// TingStateTimeOut 听牌状态下的超时时间，单位为second，默认值为 1
	TingStateTimeOut = "ts_timeout"

	// MaxFapaiCartoonTime 发牌的动画时间
	MaxFapaiCartoonTime = "fp_cartoontime"

	// MaxHuansanzhangCartoonTime 换三张动画时间
	MaxHuansanzhangCartoonTime = "hsz_cartoontime"

	/*event type*/

	// NormalEvent 普通事件
	NormalEvent int = iota
	// OverTimeEvent 超时事件
	OverTimeEvent
	// TuoGuanEvent 托管事件
	TuoGuanEvent
	// RobotEvent 机器人事件
	RobotEvent
	// HuStateEvent 胡状态事件
	HuStateEvent
	// TingStateEvent 听状态事件
	TingStateEvent
	// SpecialOverTimeEvent 胡听状态下的超时事件
	SpecialOverTimeEvent

	/* model Name */
	EventModelName   = "EventModel"
	RequestModelName = "RequestModel"

	PlayerModelName   = "PlayerModel"
	MessageModelName  = "MessageModel"
	ChatModelName     = "ChatModel"
	ContinueModelName = "ContinueModel"
)

var (
	// ErrInvalidEvent invalid event
	ErrInvalidEvent = errors.New("invalid error")
	// ErrInvalidRequestPlayer invalid request player
	ErrInvalidRequestPlayer = errors.New("invalid request player")
)
