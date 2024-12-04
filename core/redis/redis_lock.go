package redis

import (
	"context"
	"time"

	"github.com/aarioai/airis/core/ae"
)

// 申请原子性的锁
func (r *RedisClient) ApplyLock(ctx context.Context, expires time.Duration, k string) *ae.Error {
	err := r.rdb.SetNX(ctx, k+":lock", 1, expires).Err()
	return ae.NewRedisError(err)
}
