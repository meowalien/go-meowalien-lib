package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClientWrapper redis.Client

func (c *RedisClientWrapper) SetStruct(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	p, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return c.Set(ctx, key, p, expiration).Result()
}

func (c *RedisClientWrapper) GetStruct(ctx context.Context, key string, stk interface{}) error {
	p, err := c.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if p != "" {
		e := json.Unmarshal([]byte(p), stk)
		if e != nil {
			return fmt.Errorf("error when json.Unmarshal: %w", e)
		}
		return nil
	} else {
		return nil
	}
}

func (c *RedisClientWrapper) HSetStruct(ctx context.Context, name string, key interface{}, value interface{}) (int64, error) {
	switch key.(type) {
	case string:
	default:
		p, err := json.Marshal(key)
		if err != nil {
			return 0, fmt.Errorf("error when key: %w", err)
		}
		key = p
	}
	switch value.(type) {
	case string:
	default:
		p, err := json.Marshal(value)
		if err != nil {
			return 0, fmt.Errorf("error when key: %w", err)
		}
		value = p
	}
	return c.HSet(ctx, name, key, value).Result()
}

func (c *RedisClientWrapper) HGetStruct(ctx context.Context, key interface{}, field interface{}, value interface{}) error {
	sKey, sField := "", ""

	switch k := key.(type) {
	case string:
		sKey = k
	default:
		p, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("error when key: %w", err)
		}
		sKey = string(p)
	}

	switch f := field.(type) {
	case string:
		sField = f
	default:
		p, err := json.Marshal(field)
		if err != nil {
			return fmt.Errorf("error when key: %w", err)
		}
		sField = string(p)
	}

	res, err := c.HGet(ctx, sKey, sField).Result()
	if err != nil {
		return fmt.Errorf("error when HGet")
	}

	err = json.Unmarshal([]byte(res), value)
	if err != nil {
		return fmt.Errorf("error when Unmarshal: %w", err)
	}
	return nil
}
