package states

import (
	"steve/client_pb/room"
)

var StageTime = map[room.DDZStage]uint32 {
	room.DDZStage_DDZ_STAGE_DEAL    :2,
	room.DDZStage_DDZ_STAGE_CALL    :15,
	room.DDZStage_DDZ_STAGE_GRAB    :15,
	room.DDZStage_DDZ_STAGE_DOUBLE  :15,
	room.DDZStage_DDZ_STAGE_PLAYING :15,
	room.DDZStage_DDZ_STAGE_OVER    :4,
}

func genNextStage(stage room.DDZStage) *room.NextStage {
	stageTime := StageTime[stage]
	return &room.NextStage{
		Stage: stage.Enum(),
		Time: &stageTime,
	}
}

func genResult(errCode uint32, errDesc string) *room.Result {
	return &room.Result{ErrCode: &errCode, ErrDesc: &errDesc}
}