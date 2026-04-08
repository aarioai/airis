package cache

import (
	"context"
	"fmt"
)

func (h *Cache) Ping(ctx context.Context) string {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return "[error] connect redis fail"
	}
	result, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Sprintf("[error] connect redis fail, %s, %s", err.Error(), result)
	}
	return result
}
