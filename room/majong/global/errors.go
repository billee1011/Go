package global

// 定义全局的错误变量

import "errors"

// ErrInvalidEvent 无效事件
var ErrInvalidEvent = errors.New("无效的事件")

// ErrUnmarshalEvent 反序列化事件现场失败
var ErrUnmarshalEvent = errors.New("反序列化事件现场失败")

// ErrInvalidRequestPlayer 无效的请求玩家
var ErrInvalidRequestPlayer = errors.New("无效的请求玩家")
