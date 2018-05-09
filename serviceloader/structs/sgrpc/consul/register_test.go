package consul

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redis"
)

func Test_allocServiceID(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	stringID, err := allocServiceID(redisCli)
	assert.Nil(t, err)
	ID, err := strconv.Atoi(stringID)
	assert.Nil(t, err)

	stringID2, err := allocServiceID(redisCli)
	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(ID+1), stringID2)
}

func Test_registerService(t *testing.T) {
	initConsulClient()
	redisCli := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	assert.Nil(t, registerService("testservice", "localhost", 8080, redisCli))
	assert.Nil(t, registerService("testservice", "localhost", 8080, redisCli))
}
