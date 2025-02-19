package cache

import (
	"context"
	"github.com/aarioai/airis/core"
	"github.com/redis/go-redis/v9"
	"{{APP_BASE}}/conf"
	"{{DRIVER_BASE}}"
	"sync"
	"time"
)

type Cache struct {
	app *core.App
	loc *time.Location
}

var (
	cacheOnce sync.Once
	cacheObj  *Cache
)

func New(app *core.App) *Cache {

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
	cli, e := driver.NewRedis(h.app, conf.RedisConfigSection)
	if e != nil {
		h.app.Log.Error(ctx, e.Text())
		return nil, false
	}
	return cli, true
}
func (h *Cache) persistentRdb(ctx context.Context) (*redis.Client, bool) {
	cli, e := driver.NewRedis(h.app, conf.PersistentRedisConfigSection)
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
