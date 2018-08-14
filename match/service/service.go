package service

import (
	"context"
	"steve/match/matchv3"
	"steve/server_pb/match"
)

// MatchService 实现 server_pb/match/Match
type MatchService struct {
}

// AddContinueDesk 添加续局牌桌
func (ms *MatchService) AddContinueDesk(ctx context.Context, request *match.AddContinueDeskReq) (response *match.AddContinueDeskRsp, err error) {
	return matchv3.AddContinueDesk(request), nil
}

// ClearAllMatch 清空所有的匹配
func (ms *MatchService) ClearAllMatch(ctx context.Context, request *match.ClearAllMatchReq) (response *match.ClearAllMatchRsp, err error) {
	return matchv3.ClearAllMatch(request), nil
}
