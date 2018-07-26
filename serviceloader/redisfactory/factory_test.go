package redisfactory

import (
	"testing"
	"unsafe"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	f := NewFactory("localhost:6379", "")
	client, err := f.NewClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func Test_NewClientConcurrcy(t *testing.T) {
	f := NewFactory("localhost:6379", "")

	clientChannel := make(chan *redis.Client, 100)

	go func() {
		firstclient := <-clientChannel
		for {
			client, ok := <-clientChannel
			if !ok {
				return
			}
			assert.Equal(t, unsafe.Pointer(firstclient), unsafe.Pointer(client))
		}
	}()

	for i := 0; i < 100; i++ {
		go func() {
			client, err := f.NewClient()
			assert.Nil(t, err)
			clientChannel <- client
		}()
	}
	close(clientChannel)
}

func Test_GetRedisClient(t *testing.T) {
	conf := map[string]string{
		"addr":   "localhost:6379",
		"passwd": "",
	}

	viper.SetDefault("redis_list", map[string]interface{}{
		"test": conf,
	})

	f := NewFactory("localhost:6379", "")
	client, err := f.GetRedisClient("test", 0)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
