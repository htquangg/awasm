package cache

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"time"

	goredis_cache "github.com/go-redis/cache/v9"
	goredis "github.com/redis/go-redis/v9"

	"github.com/htquangg/a-wasm/config"
)

var ErrMissingFetchFunc = fmt.Errorf("missing fetch function")

type Cacher interface {
	Ping(context.Context) error
	Get(context.Context, *Key) ([]byte, bool, error)
	Fetch(context.Context, *Key, time.Duration, FetchFunc) ([]byte, bool, error)
	Set(context.Context, *Key, []byte, time.Duration) error
	Delete(context.Context, *Key) error
}

type cache struct {
	cfg     *config.Config
	client  *goredis.Client
	cache   *goredis_cache.Cache
	keyFunc KeyFunc
}

type FetchFunc func() (interface{}, error)

type Key struct {
	Namespace string
	Key       string
}

func (k *Key) Compute(f KeyFunc) (string, error) {
	key := k.Key

	if f != nil {
		var err error
		key, err = f(key)
		if err != nil {
			return "", err
		}
	}

	if k.Namespace != "" {
		return k.Namespace + ":" + key, nil
	}
	return key, nil
}

type KeyFunc func(string) (string, error)

func HashKeyFunc(hasher func() hash.Hash) KeyFunc {
	return func(in string) (string, error) {
		h := hasher()
		n, err := h.Write([]byte(in))
		if err != nil {
			return "", err
		}
		if got, want := n, len(in); got < want {
			return "", fmt.Errorf("only hashed %d of %d bytes", got, want)
		}
		dig := h.Sum(nil)
		return fmt.Sprintf("%x", dig), nil
	}
}

func HMACKeyFunc(hasher func() hash.Hash, key []byte) KeyFunc {
	return HashKeyFunc(func() hash.Hash {
		return hmac.New(hasher, key)
	})
}

func New(ctx context.Context, cfg *config.Config) (Cacher, error) {
	var tlsConfig *tls.Config
	if cfg.Redis.RequireTLS {
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	rdb := goredis.NewClient(&goredis.Options{
		Addr:      fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		TLSConfig: tlsConfig,
		Password:  cfg.Redis.Password,
		DB:        cfg.Redis.DB,
		PoolSize:  cfg.Redis.PoolSize,
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
		keyFunc: HMACKeyFunc(sha1.New, cfg.Key.CacheKeyHMACBytes),
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

func (c *cache) Get(ctx context.Context, k *Key) ([]byte, bool, error) {
	key, err := k.Compute(c.keyFunc)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compute key: %w", err)
	}

	var value []byte

	if err := c.cache.Get(ctx, key, &value); err != nil {
		if errors.Is(err, goredis_cache.ErrCacheMiss) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return value, true, nil
}

func (c *cache) Fetch(ctx context.Context, k *Key, ttl time.Duration, f FetchFunc) ([]byte, bool, error) {
	key, err := k.Compute(c.keyFunc)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compute key: %w", err)
	}

	var value []byte

	fn := func(tx *goredis.Tx) error {
		// Get current value or zero.
		err := c.cache.Get(ctx, key, &value)
		if err != nil && !errors.Is(err, goredis_cache.ErrCacheMiss) {
			return fmt.Errorf("failed to GET key: %w", err)
		}

		if len(value) != 0 {
			return nil
		}

		// No value found
		if f == nil {
			return ErrMissingFetchFunc
		}
		val, err := f()
		if err != nil {
			return err
		}

		value, err = json.Marshal(val)
		if err != nil {
			return err
		}

		// Operation is committed only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
			return pipe.Set(ctx, key, value, ttl).Err()
		})

		return err
	}

	// This is a CAS operation, so retry
	for i := 0; i < 5; i++ {
		err = c.client.Watch(ctx, fn, key)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, false, err
	}

	return value, len(value) != 0, nil
}

func (c *cache) Set(ctx context.Context, k *Key, value []byte, ttl time.Duration) error {
	key, err := k.Compute(c.keyFunc)
	if err != nil {
		return fmt.Errorf("failed to compute key: %w", err)
	}

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

func (c *cache) Delete(ctx context.Context, k *Key) error {
	key, err := k.Compute(c.keyFunc)
	if err != nil {
		return fmt.Errorf("failed to compute key: %w", err)
	}

	if err := c.cache.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
