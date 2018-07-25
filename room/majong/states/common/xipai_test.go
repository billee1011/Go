package common

import (
	"fmt"
	"steve/room/majong/interfaces"
	majongpb "steve/entity/majong"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestXipaiState_xipai(t *testing.T) {
	mc := gomock.NewController(t)
	flow := interfaces.NewMockMajongFlow(mc)

	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
		},
	).AnyTimes()

	s := XipaiState{}

	cards := s.xipai(flow)
	results := ""
	for _, card := range cards {
		results = fmt.Sprintf("%s\n%v %v", results, card.GetColor(), card.GetPoint())
	}
	logrus.Info(results)
	assert.NotEqual(t, 0, len(results))
}
