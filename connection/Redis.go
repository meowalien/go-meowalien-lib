package connection

import (
	"context"
	"github.com/go-redis/redis/v8"
)
type RedisConfiguration struct {
	Host     string
	Port     string
	Password string
	ID       int
}

func CreateRedisConnection(dbconf RedisConfiguration) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     dbconf.Host + ":" + dbconf.Port,
		Password: dbconf.Password,
		DB:       dbconf.ID,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, err
}
