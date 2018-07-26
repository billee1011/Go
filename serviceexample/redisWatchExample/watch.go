package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
)

var incr func(string) error

func notWatch(client *redis.Client) error {
	n, err := client.Get("counter3").Int64()
	if err != nil && err != redis.Nil {
		return err
	}
	client.Set("counter3", strconv.FormatInt(n+1, 10), 0)
	return nil
}

func useWatch(client *redis.Client) {
	incr = func(key string) error {
		err := client.Watch(func(tx *redis.Tx) error {
			n, err := tx.Get(key).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Set(key, strconv.FormatInt(n+1, 10), 0)
				return nil
			})
			return err
		}, key)
		if err == redis.TxFailedErr {
			return incr(key)
		}
		return err
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	useWatch(client)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := incr("counter3")
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()

	n, err := client.Get("counter3").Int64()
	fmt.Println(n, err)
}
