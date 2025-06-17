package cache

import (
	"context"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa"
	"github.com/redis/go-redis/v9"
	"project/simple/app/app_aajs/conf"
	"sync"
	"time"
)

type Cache struct {
	app *aa.App
	loc *time.Location
}

var (
	cacheOnce sync.Once
	cacheObj  *Cache
)

func New(app *aa.App) *Cache {
	cacheOnce.Do(func() {
		cacheObj = &Cache{
			app: app,
			loc: app.Config.TimeLocation,
		}
	})
	return cacheObj
}


// go redis will create connection pool automatically. so its no need to close a connection.
// @doc https://pkg.go.dev/github.com/redis/go-redis/v9#Client.Close
// It is rare to Close a Client, as the Client is meant to be long-lived and shared between many goroutines.
func (h *Cache) rdb(ctx context.Context) (*redis.Client, bool) {
	cli, e := driver.NewRedisPool(h.app, conf.RedisCfgSection)
	if e != nil {
		h.app.Log.Error(ctx, e.Text())
		return nil, false
	}
	return cli, true
}

func (h *Cache) check(ctx context.Context, err error)  bool  {
	return h.app.CheckErrors(ctx, err)
}

func (h *Cache) rdbNoExpire(ctx context.Context) (*redis.Client, bool) {
	cli, e := driver.NewRedisPool(h.app, conf.RedisNoExpireCfgSection)
	if !h.app.Check(ctx, e) {
		return nil, false
	}
	return cli, true
}

func (h *Cache) del(ctx context.Context, keys ...string) bool {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return false
	}
	// @doc https://pkg.go.dev/github.com/redis/go-redis/v9#Client.Close
    // It is rare to Close a Client, as the Client is meant to be long-lived and shared between many goroutines.
	// defer rdb.Close()
	_, err := rdb.Del(ctx, keys...).Result()
	return h.app.Check(ctx, driver.NewRedisError(err))
}
