package infrastructure

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func connectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
func InitRedis() error {
	var err error
	redisClient, err = connectRedis()
	if err != nil {
		return err
	}
	return nil
}
