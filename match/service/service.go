package service

import (
	"context"
	"steve/match/matchv2"
	"steve/server_pb/match"
)

// MatchService 实现 server_pb/match/Match
type MatchService struct {
}

// AddContinueDesk 添加续局牌桌
func (ms *MatchService) AddContinueDesk(ctx context.Context, request *match.AddContinueDeskReq) (response *match.AddContinueDeskRsp, err error) {
	return matchv2.AddContinueDesk(request), nil
}
