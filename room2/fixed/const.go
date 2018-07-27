package fixed

const (
	// ListenClientAddr 代表监听客户端的IP地址，默认值为 127.0.0.1
	ListenClientAddr = "lis_client_addr"

	// ListenClientPort 代表监听客户端的端口， 默认值为 36001
	ListenClientPort = "lis_client_port"

	// ListenPeipaiAddr 配牌監聽地址
	ListenPeipaiAddr = "peipai_addr"

	// XingPaiTimeOut 行牌超时时间，单位为second，默认值为 10
	XingPaiTimeOut = "xp_timeout"

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


	/* Event Name */
	Event       = "EventModel"
	Player      = "PlayerModel"
	Request     = "RequestModel"
	Message     = "MessageModel"
	Trusteeship = "TrusteeshipModel"
	Chat = "ChatModel"
)