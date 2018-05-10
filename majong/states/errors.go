package states

import "errors"

var errInvalidEvent = errors.New("无效的事件")
var errUnmarshalEvent = errors.New("反序列化事件现场失败")
var errInvalidRequestPlayer = errors.New("无效的请求玩家")
