package majong

import "errors"

var errMessageTypeNotMatch = errors.New("消息体类型不匹配")
var errUnmarshalEvent = errors.New("反序列化事件失败")
