package initialize

import (
	"github.com/go-redis/redis"
)

func InitRedisClient(addr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// PING
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
