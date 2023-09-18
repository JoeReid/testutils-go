package containers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/cache/v9"
	"github.com/go-redis/redis_rate/v10"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	resource *dockertest.Resource
}

func (r *Redis) Port() string {
	return r.resource.GetPort("6379/tcp")
}

func (r *Redis) GoRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", r.Port()),
	})
}

func (r *Redis) GoRedisCache() *cache.Cache {
	return cache.New(&cache.Options{
		Redis: r.GoRedisClient(),
	})
}

func (r *Redis) RedisLock() *redislock.Client {
	return redislock.New(r.GoRedisClient())
}

func (r *Redis) RedisRate() *redis_rate.Limiter {
	return redis_rate.NewLimiter(r.GoRedisClient())
}

func NewRedis(t *testing.T) *Redis {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to construct docker pool: %v", err)
	}

	if err := pool.Client.Ping(); err != nil {
		t.Fatalf("failed to connect to docker: %v", err)
	}

	resource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	r := &Redis{resource: resource}

	if err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		return r.GoRedisClient().Ping(ctx).Err()
	}); err != nil {
		t.Fatalf("failed to connect to redis container: %v", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("failed to purge redis container: %v", err)
		}
	})

	return r
}
