package exec

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"io"
	"time"
)

type SQLExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type CacheExecutor interface {
	io.Closer
	redis.Cmdable
	SetStruct(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	GetStruct(ctx context.Context, key string, stk interface{}) error
	HSetStruct(ctx context.Context, name string, key interface{}, value interface{}) (int64, error)
	HGetStruct(ctx context.Context, key interface{}, field interface{}, value interface{}) error
}
