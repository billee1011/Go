package core

import (
	"steve/structs/proto/base"
	"testing"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func Test_Any(t *testing.T) {
	header := steve_proto_base.Header{
		MsgId: proto.Uint32(111),
	}
	buf1, err := proto.Marshal(&header)
	assert.Nil(t, err)

	anyHeader, err := ptypes.MarshalAny(&header)
	assert.Nil(t, err)

	buf2, err := proto.Marshal(anyHeader)
	assert.Nil(t, err)

	// 用 Any 序列化出来的内容和直接序列化出来的内容不同
	assert.NotEqual(t, buf1, buf2)

	header2 := steve_proto_base.Header{}
	assert.Nil(t, proto.Unmarshal(buf2, &header2))

	assert.Equal(t, header2, header)

}
