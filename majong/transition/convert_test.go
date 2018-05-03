package transition

import (
	"testing"

	majongpb "steve/server_pb/majong"

	"github.com/stretchr/testify/assert"
)

func Test_originTransitionToMap(t *testing.T) {
	originTransition := transition{
		GameID: 1,
		States: []struct {
			CurState string      `yaml:"state"`
			Trans    []stateTran `yaml:"transition"`
		}{
			{
				CurState: "state_init",
				Trans: []stateTran{
					{
						Events:    []string{"event_start_game"},
						NextState: "state_xipai",
					},
				},
			},
			{
				CurState: "state_chupaiwenxun",
				Trans: []stateTran{
					{
						Events:    []string{"event_peng_request", "event_gang_request"},
						NextState: "state_peng",
					},
					{
						Events:    []string{"event_hu_request", "event_gang_request"},
						NextState: "state_hu",
					},
				},
			},
		},
	}

	transitionMap, err := originTransitionToMap(&originTransition)

	assert.Nil(t, err)
	assert.NotNil(t, transitionMap)

	assert.Equal(t, 1, len(transitionMap[majongpb.StateID_state_init][majongpb.EventID_event_start_game]))
	assert.Equal(t, 1, len(transitionMap[majongpb.StateID_state_chupaiwenxun][majongpb.EventID_event_peng_request]))
	assert.Equal(t, 2, len(transitionMap[majongpb.StateID_state_chupaiwenxun][majongpb.EventID_event_gang_request]))
	assert.Equal(t, 1, len(transitionMap[majongpb.StateID_state_chupaiwenxun][majongpb.EventID_event_hu_request]))

}
