package redishelper

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/aarioai/airis/core/ae"
)

// 申请原子性的锁
func ApplyLock(ctx context.Context, rdb *redis.Client, expires time.Duration, k string) *ae.Error {
	err := rdb.SetNX(ctx, k+":lock", 1, expires).Err()
	return ae.NewRedisError(err)
}
