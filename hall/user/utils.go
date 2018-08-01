package user

import (
	"steve/client_pb/common"
	"steve/server_pb/user"

	"github.com/golang/protobuf/proto"
)

// ServerGameConfig2Client server端GameConfig转化为client端
func ServerGameConfig2Client(gameInfos []*user.GameConfig) []*common.GameConfig {
	cGameConfigs := make([]*common.GameConfig, 0)
	for _, gameInfo := range gameInfos {
		cGameConfig := &common.GameConfig{
			GameId:   proto.Uint32(gameInfo.GameId),
			GameName: proto.String(gameInfo.GameName),
			GameType: proto.Uint32(gameInfo.GameType),
		}
		cGameConfigs = append(cGameConfigs, cGameConfig)
	}
	return cGameConfigs
}

// ServerGameLevelConfig2Client server端GameLvelConfig转化为client端
func ServerGameLevelConfig2Client(gameInfos []*user.GameLevelConfig) []*common.GameLevelConfig {
	cGameLevelConfigs := make([]*common.GameLevelConfig, 0)
	for _, gameInfo := range gameInfos {
		cGameLevelConfig := &common.GameLevelConfig{
			GameId:     proto.Uint32(gameInfo.GameId),
			LevelId:    proto.Uint32(gameInfo.LevelId),
			LevelName:  proto.String(gameInfo.LevelName),
			BaseScores: proto.Uint32(gameInfo.BaseScores),
			LowScores:  proto.Uint32(gameInfo.LowScores),
			HighScors:  proto.Uint32(gameInfo.HighScores),
			MinPeople:  proto.Uint32(gameInfo.MinPeople),
			MaxPeople:  proto.Uint32(gameInfo.MaxPeople),
		}
		cGameLevelConfigs = append(cGameLevelConfigs, cGameLevelConfig)
	}
	return cGameLevelConfigs
}
