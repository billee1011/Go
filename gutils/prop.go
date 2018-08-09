package gutils

import (
	"steve/client_pb/common"
	"steve/entity/prop"
)

// PropTypeServer2Client 道具类型转换 server -> client
func PropTypeServer2Client(propID int32) common.PropType {
	switch propID {
	case prop.Rose:
		return common.PropType_ROSE
	case prop.Beer:
		return common.PropType_BEER
	case prop.Bomb:
		return common.PropType_BOMB
	case prop.GrabChicken:
		return common.PropType_GRAB_CHICKEN
	case prop.EggGun:
		return common.PropType_EGG_GUN
	default:
		return common.PropType_INVALID_PROP
	}
}

// PropTypeClient2Server 道具类型转换 client -> server
func PropTypeClient2Server(propID common.PropType) int32 {
	switch propID {
	case common.PropType_ROSE:
		return prop.Rose
	case common.PropType_BEER:
		return prop.Beer
	case common.PropType_BOMB:
		return prop.Bomb
	case common.PropType_GRAB_CHICKEN:
		return prop.GrabChicken
	case common.PropType_EGG_GUN:
		return prop.EggGun
	default:
		return prop.InvalidProp
	}
}
