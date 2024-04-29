package cache

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/htquangg/a-wasm/config"

	goredis_cache "github.com/go-redis/cache/v9"
	goredis "github.com/redis/go-redis/v9"
)

type Cacher interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttp time.Duration) error
	Delete(ctx context.Context, key string) error
}

type cache struct {
	cfg    *config.Redis
	client *goredis.Client
	cache  *goredis_cache.Cache
}

func Key(k string) string {
	return fmt.Sprintf("awasm:%x", md5.Sum([]byte(k)))
}

func New(ctx context.Context, cfg *config.Redis) (Cacher, error) {
	var tlsConfig *tls.Config
	if cfg.RequireTLS {
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	rdb := goredis.NewClient(&goredis.Options{
		Addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		TLSConfig: tlsConfig,
		Password:  cfg.Password,
		DB:        cfg.DB,
		PoolSize:  cfg.PoolSize,
	})

	status := rdb.Ping(ctx)
	if status == nil {
		return nil, errors.New("connecting to redis: no status")
	}

	if status.Err() != nil {
		return nil, fmt.Errorf("connecting to redis: %w", status.Err())
	}

	return &cache{
		cfg:    cfg,
		client: rdb,
		cache: goredis_cache.New(&goredis_cache.Options{
			Redis: rdb,
		}),
	}, nil
}

func (c *cache) Ping(ctx context.Context) error {
	status := c.client.Ping(ctx)
	if status == nil {
		return errors.New("connecting to redis: no status")
	}

	if status.Err() != nil {
		return fmt.Errorf("connecting to redis: %w", status.Err())
	}

	return nil
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	var value []byte
	key = Key(key)
	if err := c.cache.Get(ctx, key, &value); err != nil {
		if errors.Is(err, goredis_cache.ErrCacheMiss) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return value, true, nil
}

func (c *cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	key = Key(key)
	if err := c.cache.Set(&goredis_cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	}); err != nil {
		return err
	}

	return nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	if err := c.cache.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
