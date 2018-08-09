package almsserver

import (
	"steve/alms/data"
	client_alms "steve/client_pb/alms"

	"github.com/golang/protobuf/proto"
)

// dataToClentPbGameLeveIsOpen data.GemeLeveIsOpentAlms To ClentPb.GameLeveIsOpen
func dataToClentPbGameLeveIsOpen(dataGl []*data.GameLeveConfig) []*client_alms.GemeLeveIsOpent {
	cglio := make([]*client_alms.GemeLeveIsOpent, 0, len(dataGl))
	for _, d := range dataGl {
		c := &client_alms.GemeLeveIsOpent{
			GemeId:  proto.Int32(d.GameID),
			LevelId: proto.Int32(d.LevelID),
			IsOpen:  proto.Int32(int32(d.IsOpen)),
		}
		cglio = append(cglio, c)
	}
	return cglio
}
