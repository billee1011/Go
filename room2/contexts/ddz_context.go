package contexts

import (
	"steve/server_pb/ddz"
)

// DDZDeskContext 斗地主游戏现场
type DDZDeskContext struct {
	DDZContext ddz.DDZContext // 牌局现场
}
