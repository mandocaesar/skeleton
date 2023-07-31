package cache

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v9"
)

// cache is an interface which specify the cache methods
type Cache interface {
	// PING redis command is used to check whether the server is running or not.
	PING(ctx context.Context) (err error)
	// SETEX command is used to set string value with a specified timeout in Redis key.
	SETEX(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error)
	// GET command get the value of redis key. If the key does not exist the special value nil is returned.
	GET(ctx context.Context, key string) (data string, err error)
	// DEL removes the specified redis keys. A key is ignored if it does not exist.
	DEL(ctx context.Context, key string) (err error)
}

// cache struct represent the cache struct used by the caller
type cache struct {
	*redis.Client
}

// NewCache instantiate a redis client that used by the caller
func NewCache(redis *redis.Client) Cache {
	return &cache{
		redis,
	}
}

// PING redis command is used to check whether the server is running or not.
func (c cache) PING(ctx context.Context) (err error) {
	_, err = c.Ping(ctx).Result()
	return
}

// SETEX command is used to set string value with a specified timeout in Redis key.
func (c cache) SETEX(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	_, err = c.SetEx(ctx, key, value, expiration).Result()
	return
}

// GET the value of redis key.
// If the key does not exist the special value nil is returned. An error is returned if
// the value stored at key is not a string, because GET only handles string values.
func (c cache) GET(ctx context.Context, key string) (data string, err error) {
	return c.Get(ctx, key).Result()
}

// DEL removes the specified redis keys. A key is ignored if it does not exist.
func (c cache) DEL(ctx context.Context, key string) (err error) {
	_, err = c.Del(ctx, key).Result()
	return
}
