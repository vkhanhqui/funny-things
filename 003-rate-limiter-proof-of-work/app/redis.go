package app

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(endpoint string, port int, password string) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", endpoint, port),
		Password: password,
		DB:       0,
	})

	_, err = client.Ping(context.Background()).Result()
	return
}
