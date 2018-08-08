package user

import (
	"steve/client_pb/common"
	"steve/entity/db"
	"steve/server_pb/user"

	"github.com/golang/protobuf/proto"
)

// DBGameConfig2Client db转client端使用
func DBGameConfig2Client(dbGameConfigs []*db.TGameConfig) []*common.GameConfig {
	cGameConfigs := make([]*common.GameConfig, 0)
	for _, dbGameConfig := range dbGameConfigs {
		cGameConfig := &common.GameConfig{
			GameId:    proto.Uint32(uint32(dbGameConfig.Gameid)),
			GameName:  proto.String(dbGameConfig.Name),
			GameType:  proto.Uint32(uint32(dbGameConfig.Type)),
			MinPeople: proto.Uint32(uint32(dbGameConfig.Minpeople)),
			MaxPeople: proto.Uint32(uint32(dbGameConfig.Maxpeople)),
		}

		cGameConfigs = append(cGameConfigs, cGameConfig)
	}
	return cGameConfigs

}

// DBGameConfig2Server db转server端使用
func DBGameConfig2Server(dbGameConfigs []*db.TGameConfig) (gameInfos []*user.GameConfig) {
	gameInfos = make([]*user.GameConfig, 0)
	for _, dbGameConfig := range dbGameConfigs {
		gameInfo := &user.GameConfig{
			GameId:    uint32(dbGameConfig.Gameid),
			GameName:  dbGameConfig.Name,
			GameType:  uint32(dbGameConfig.Type),
			MinPeople: uint32(dbGameConfig.Minpeople),
			MaxPeople: uint32(dbGameConfig.Maxpeople),
		}

		gameInfos = append(gameInfos, gameInfo)
	}
	return
}

// DBGamelevelConfig2Sercer db转server端使用
func DBGamelevelConfig2Sercer(dbGameConfigs []*db.TGameLevelConfig) (gamelevelConfigs []*user.GameLevelConfig) {
	gamelevelConfigs = make([]*user.GameLevelConfig, 0)
	for _, dbGameConfig := range dbGameConfigs {
		gamelevelConfig := &user.GameLevelConfig{
			GameId:     uint32(dbGameConfig.Gameid),
			LevelId:    uint32(dbGameConfig.Levelid),
			LevelName:  dbGameConfig.Name,
			BaseScores: uint32(dbGameConfig.Basescores),
			LowScores:  uint32(dbGameConfig.Lowscores),
			HighScores: uint32(dbGameConfig.Highscores),
		}

		gamelevelConfigs = append(gamelevelConfigs, gamelevelConfig)
	}
	return
}

// DBGamelevelConfig2Client db转client端使用
func DBGamelevelConfig2Client(dbGameConfigs []*db.TGameLevelConfig) []*common.GameLevelConfig {
	cGameLevelConfigs := make([]*common.GameLevelConfig, 0)
	for _, dbGameConfig := range dbGameConfigs {
		cGameLevelConfig := &common.GameLevelConfig{
			GameId:     proto.Uint32(uint32(dbGameConfig.Gameid)),
			LevelId:    proto.Uint32(uint32(dbGameConfig.Levelid)),
			LevelName:  proto.String(dbGameConfig.Name),
			BaseScores: proto.Uint32(uint32(dbGameConfig.Basescores)),
			LowScores:  proto.Uint32(uint32(dbGameConfig.Lowscores)),
			HighScors:  proto.Uint32(uint32(dbGameConfig.Highscores)),
			ShowPeople: proto.Uint32(uint32(dbGameConfig.Showonlinepeople)),
			RealPeople: proto.Uint32(uint32(dbGameConfig.Realonlinepeople)),
		}
		cGameLevelConfigs = append(cGameLevelConfigs, cGameLevelConfig)
	}
	return cGameLevelConfigs
}
