package gutils

import (
	"steve/client_pb/room"
	"steve/entity/prop"
)

// PropTypeServer2Client 道具类型转换 server -> client
func PropTypeServer2Client(propID int32) room.PropType {
	switch propID {
	case prop.Rose:
		return room.PropType_ROSE
	case prop.Beer:
		return room.PropType_BEER
	case prop.Bomb:
		return room.PropType_BOMB
	case prop.GrabChicken:
		return room.PropType_GRAB_CHICKEN
	case prop.EggGun:
		return room.PropType_EGG_GUN
	default:
		return room.PropType_INVALID_PROP
	}
}

// PropTypeClient2Server 道具类型转换 client -> server
func PropTypeClient2Server(propID room.PropType) int32 {
	switch propID {
	case room.PropType_ROSE:
		return prop.Rose
	case room.PropType_BEER:
		return prop.Beer
	case room.PropType_BOMB:
		return prop.Bomb
	case room.PropType_GRAB_CHICKEN:
		return prop.GrabChicken
	case room.PropType_EGG_GUN:
		return prop.EggGun
	default:
		return prop.InvalidProp
	}
}
