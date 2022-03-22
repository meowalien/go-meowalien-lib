package db_connect

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/meowalien/go-meowalien-lib/db/config_modules"
)

func CreateRedisConnection(dbconf config_modules.RedisConfiguration) (*redis.Client, error) {
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
