package redis

import (
	"context"
	"github.com/meowalien/go-meowalien-lib/errs"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfiguration struct {
	*redis.Options
	*redis.ClusterOptions
	PingInterval    *time.Duration
	PingRetryTimes  *int
	PingFailRecover chan error
}

func CreateRedisConnection(ctx context.Context, config RedisConfiguration) (client redis.Cmdable) {
	switch {
	case config.ClusterOptions != nil:
		client = redis.NewClusterClient(config.ClusterOptions)
	case config.Options != nil:
		client = redis.NewClient(config.Options)
	default:
		panic(errs.New("redis configuration is empty"))
	}

	go pingLoop(ctx, client, config)

	return
}

const DefaultPingRetryTimes = 20

func pingLoop(ctx context.Context, client redis.Cmdable, config RedisConfiguration) {
	if config.PingInterval == nil {
		return
	}
	var pingRetryLimit = DefaultPingRetryTimes
	if config.PingRetryTimes != nil {
		pingRetryLimit = *config.PingRetryTimes
	}
	pingRetry := pingRetryLimit
	timer := time.NewTimer(*config.PingInterval)
	for {
		select {
		case <-ctx.Done():
			if config.PingFailRecover != nil {
				select {
				case config.PingFailRecover <- errs.New("ping loop is closed by context"):
				default:
				}
			}
			return
		case <-timer.C:
			if pingRetry <= 0 {
				err := errs.New("ping redis failed after %d times retry", pingRetryLimit)
				select {
				case config.PingFailRecover <- err:
				default:
					log.Println(err.Error())
				}

				return
			}
			if err := client.Ping(ctx).Err(); err != nil {
				select {
				case config.PingFailRecover <- err:
				default:
					log.Println(err.Error())
				}
				pingRetry--
			} else {
				pingRetry = pingRetryLimit
			}
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(*config.PingInterval)
		}
	}
}
