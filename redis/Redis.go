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
	PingInterval   time.Duration `json:"ping_interval,omitempty"`
	PingRetryLimit int           `json:"ping_retry_times,omitempty"`
}

type ConnectionConfig struct {
	Address            []string      `json:"addr"`
	Password           string        `json:"password,omitempty"`
	PoolSize           int           `json:"pool_size,omitempty"`
	MaxRetries         int           `json:"max_retries,omitempty"`
	MinRetryBackoff    time.Duration `json:"min_retry_backoff,omitempty"`
	MaxRetryBackoff    time.Duration `json:"max_retry_backoff,omitempty"`
	Network            string        `json:"network,omitempty"`
	DB                 int           `json:"db,omitempty"`
	Limiter            redis.Limiter `json:"limiter,omitempty"`
	MaxRedirects       int           `json:"max_redirects,omitempty"`
	ReadOnly           bool          `json:"read_only,omitempty"`
	RouteByLatency     bool          `json:"route_by_latency,omitempty"`
	RouteRandomly      bool          `json:"route_randomly,omitempty"`
	Username           string        `json:"username,omitempty"`
	DialTimeout        time.Duration `json:"dial_timeout,omitempty"`
	ReadTimeout        time.Duration `json:"read_timeout,omitempty"`
	WriteTimeout       time.Duration `json:"write_timeout,omitempty"`
	PoolFIFO           bool          `json:"pool_fifo,omitempty"`
	MinIdleConns       int           `json:"min_idle_conns,omitempty"`
	MaxConnAge         time.Duration `json:"max_conn_age,omitempty"`
	PoolTimeout        time.Duration `json:"pool_timeout,omitempty"`
	IdleTimeout        time.Duration `json:"idle_timeout,omitempty"`
	IdleCheckFrequency time.Duration `json:"idle_check_frequency,omitempty"`
	TLSConfig          *tls.Config   `json:"tls_config,omitempty"`
	PingSettings       *PingSettings `json:"ping_settings,omitempty"`
	OnPingError        func(error)   `json:"-"`
}

type Client interface {
	redis.Cmdable
	AddHook(hook redis.Hook)
}

func NewClient(ctx context.Context, config ConnectionConfig) (client Client, err error) { //nolint:gocritic
	if config.OnPingError == nil {
		config.OnPingError = func(err error) {
			log.Println(err.Error())
		}
	}
	switch len(config.Address) {
	case 0:
		err = errs.New("redis configuration is empty")
		return
	case 1:
		opt := redis.Options{
			Network:            config.Network,
			Addr:               config.Address[0],
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
			Addrs:              config.Address,
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
		go pingLoop(ctx, client, *config.PingSettings, config.OnPingError)
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
const DefaultPingInterval = time.Second * 3

func pingLoop(ctx context.Context, client redis.Cmdable, config PingSettings, onError func(error)) {
	if config.PingInterval == 0 {
		config.PingInterval = DefaultPingInterval
		return
	}
	if config.PingRetryLimit == 0 {
		config.PingRetryLimit = DefaultPingRetryTimes
	}
	pingRetry := config.PingRetryLimit
	timer := time.NewTimer(config.PingInterval)
	for {
		select {
		case <-ctx.Done():
			log.Println(errs.New("ping loop is closed by context"))
			return
		case <-timer.C:
			if pingRetry <= 0 {
				err := errs.New("ping redis failed after %d times retry", config.PingRetryLimit)
				onError(err)
				return
			}
			if err := ping(ctx, client); err != nil {
				onError(err)
				pingRetry--
			} else {
				pingRetry = config.PingRetryLimit
			}
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(config.PingInterval)
		}
	}
}
