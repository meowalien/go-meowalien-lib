package redis

import (
	"context"
	"crypto/tls"
	"github.com/meowalien/go-meowalien-lib/errs"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type PingSettings struct {
	PingInterval    *time.Duration
	PingRetryTimes  *int
	PingFailRecover chan error
}

type RedisConfiguration struct {
	Network            string
	Addrs              []string
	DB                 int
	Limiter            redis.Limiter
	MaxRedirects       int
	ReadOnly           bool
	RouteByLatency     bool
	RouteRandomly      bool
	Username           string
	Password           string
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolFIFO           bool
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	TLSConfig          *tls.Config
	PingSettings       *PingSettings
}

func CreateRedisConnection(ctx context.Context, config RedisConfiguration) (client redis.Cmdable, err error) {
	switch len(config.Addrs) {
	case 0:
		err = errs.New("redis configuration is empty")
		return
	case 1:
		opt := redis.Options{
			Network:            config.Network,
			Addr:               config.Addrs[0],
			Username:           config.Username,
			Password:           config.Password,
			DB:                 config.DB,
			MaxRetries:         config.MaxRetries,
			MinRetryBackoff:    config.MinRetryBackoff,
			MaxRetryBackoff:    config.MaxRetryBackoff,
			DialTimeout:        config.DialTimeout,
			ReadTimeout:        config.ReadTimeout,
			WriteTimeout:       config.WriteTimeout,
			PoolFIFO:           config.PoolFIFO,
			PoolSize:           config.PoolSize,
			MinIdleConns:       config.MinIdleConns,
			MaxConnAge:         config.MaxConnAge,
			PoolTimeout:        config.PoolTimeout,
			IdleTimeout:        config.IdleTimeout,
			IdleCheckFrequency: config.IdleCheckFrequency,
			TLSConfig:          config.TLSConfig,
			Limiter:            config.Limiter,
		}
		client = redis.NewClient(&opt)
	default:
		opt := redis.ClusterOptions{
			Addrs:              config.Addrs,
			MaxRedirects:       config.MaxRedirects,
			ReadOnly:           config.ReadOnly,
			RouteByLatency:     config.RouteByLatency,
			RouteRandomly:      config.RouteRandomly,
			Username:           config.Username,
			Password:           config.Password,
			MaxRetries:         config.MaxRetries,
			MinRetryBackoff:    config.MinRetryBackoff,
			MaxRetryBackoff:    config.MaxRetryBackoff,
			DialTimeout:        config.DialTimeout,
			ReadTimeout:        config.ReadTimeout,
			WriteTimeout:       config.WriteTimeout,
			PoolFIFO:           config.PoolFIFO,
			PoolSize:           config.PoolSize,
			MinIdleConns:       config.MinIdleConns,
			MaxConnAge:         config.MaxConnAge,
			PoolTimeout:        config.PoolTimeout,
			IdleTimeout:        config.IdleTimeout,
			IdleCheckFrequency: config.IdleCheckFrequency,
			TLSConfig:          config.TLSConfig,
		}
		client = redis.NewClusterClient(&opt)
	}

	err = ping(ctx, client)
	if err != nil {
		err = errs.New(err)
		return
	}

	if config.PingSettings != nil {
		go pingLoop(ctx, client, *config.PingSettings)
	}

	return
}

func ping(ctx context.Context, client redis.Cmdable) (err error) {
	if err = client.Ping(ctx).Err(); err != nil {
		err = errs.New("ping failed: %w", err)
		return
	}
	return
}

const DefaultPingRetryTimes = 20

func pingLoop(ctx context.Context, client redis.Cmdable, config PingSettings) {
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
			if err := ping(ctx, client); err != nil {
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
