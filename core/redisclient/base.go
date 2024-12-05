package redisclient

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis 操作的上下文类型
type RedisClient struct {
	rdb *redis.Client
}

// TTL 常量定义
const (
	HourlyTTL  = 24 * time.Hour // 24小时 TTL，用于小时级缓存   要求每小时会自动清除之前表；为了避免宕机等影响，ttl设计长一点，24小时内宕机恢复，就能使用
	DailyTTL   = 72 * time.Hour // 3天 TTL，用于天级缓存  要求每天会自动清除之前表；为了避免宕机等影响，ttl设计长一点
	DefaultTTL = time.Hour      // 默认 TTL
)

// 创建新的 Redis 上下文
func NewClient(rdb *redis.Client) *RedisClient {
	return &RedisClient{
		rdb: rdb,
	}
}
