package redis_wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type JsonWrapper struct {
	redis.Cmdable
}

func (j JsonWrapper) Wrap(c redis.Cmdable) redis.Cmdable {
	return &JsonWrapper{
		Cmdable: c,
	}
}

func (j *JsonWrapper) SetStruct(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	p, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return j.Set(ctx, key, p, expiration).Result()
}

func (j *JsonWrapper) GetStruct(ctx context.Context, key string, stk interface{}) error {
	p, err := j.Get(ctx, key).Result()
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

func (j *JsonWrapper) HSetStruct(ctx context.Context, name string, key interface{}, value interface{}) (int64, error) {
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
	return j.HSet(ctx, name, key, value).Result()
}

func (j *JsonWrapper) HGetStruct(ctx context.Context, key interface{}, field interface{}, value interface{}) error {
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

	res, err := j.HGet(ctx, sKey, sField).Result()
	if err != nil {
		return fmt.Errorf("error when HGet")
	}

	err = json.Unmarshal([]byte(res), value)
	if err != nil {
		return fmt.Errorf("error when Unmarshal: %w", err)
	}
	return nil
}
