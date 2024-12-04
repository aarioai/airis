package redis

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/atype"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RedisClient) HSet(ctx context.Context, expires time.Duration, k string, values ...any) *ae.Error {
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err1 := pipe.HSet(ctx, k, values...).Err()
		err2 := pipe.Expire(ctx, k, expires).Err()
		return ae.CatchError(err1, err2)
	})

	return ae.NewRedisError(err)
}

func (r *RedisClient) HMSet(ctx context.Context, expires time.Duration, k string, values ...any) *ae.Error {
	keys := make(map[string]struct{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		f := fmt.Sprintf("%v", values[i])
		if f == "" {
			return ae.NewRedisError(fmt.Errorf("hmset %s empty field name", k))
		}
		if _, ok := keys[f]; ok {
			return ae.NewRedisError(fmt.Errorf("hmset %s conflict %s", k, f))
		}
		keys[f] = struct{}{}
	}
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err1 := pipe.HMSet(ctx, k, values...).Err()
		err2 := pipe.Expire(ctx, k, expires).Err()
		return ae.CatchError(err1, err2)
	})

	return ae.NewRedisError(err)
}

func (r *RedisClient) HScan(ctx context.Context, dest any, k string, fields ...string) *ae.Error {
	c := r.rdb.HMGet(ctx, k, fields...)
	v, err := c.Result()
	if err != nil {
		return ae.NewRedisError(err)
	}
	if len(v) != len(fields) {
		return ae.NotFound
	}
	for _, x := range v {
		if atype.IsNil(x) {
			return ae.NotFound
		}
	}
	e := ae.NewRedisError(c.Scan(dest))
	return e
}

func (r *RedisClient) HGetAll(ctx context.Context, k string, dest any) *ae.Error {
	c := r.rdb.HGetAll(ctx, k)
	result, err := c.Result()
	if err != nil {
		return ae.NewRedisError(err)
	}
	if len(result) == 0 {
		return ae.NotFound
	}
	e := ae.NewRedisError(c.Scan(dest))
	return e
}

func (r *RedisClient) HGetAllInt(ctx context.Context, k string, restrict bool) (map[string]int, *ae.Error) {
	c := r.rdb.HGetAll(ctx, k)
	result, err := c.Result()
	if err != nil {
		return nil, ae.NewRedisError(err)
	}
	n := len(result)
	if n == 0 {
		return nil, ae.NotFound
	}
	value := make(map[string]int, n)
	var x int64
	for k, v := range result {
		if x, err = strconv.ParseInt(v, 10, 32); err != nil {
			if restrict {
				return nil, ae.NewRedisError(fmt.Errorf(`HGetAllInt: invalid int string %s`, v))
			}
			continue
		}
		value[k] = int(x)
	}
	return value, nil
}

// 只要存在一个，就不报错；全是nil，返回 ae.NotFound
func (r *RedisClient) TryHMGet(ctx context.Context, k string, fields ...string) ([]any, bool, *ae.Error) {
	v, err := r.rdb.HMGet(ctx, k, fields...).Result()
	if err != nil {
		return nil, false, ae.NewRedisError(err)
	}
	n := len(v)
	if n != len(fields) {
		return nil, false, ae.NotFound
	}
	ok := true
	e := ae.NotFound
	for _, x := range v {
		if !atype.IsNil(x) {
			e = nil // 只要存在一个不是nil，都正确
			if !ok {
				break
			}
		} else {
			ok = false
			if e == nil {
				break
			}
		}
	}
	return v, ok, e
}

// 只要存在一个，就不报错；全是nil，返回 ae.NotFound
func (r *RedisClient) TryHMGetString(ctx context.Context, k string, fields ...string) ([]string, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]string, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = ""
		} else {
			v[i] = atype.New(x).String()
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetUint64(ctx context.Context, k string, fields []string, defaultValue uint64) ([]uint64, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint64, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultUint64(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetUint(ctx context.Context, k string, fields []string, defaultValue uint) ([]uint, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultUint(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetUint32(ctx context.Context, k string, fields []string, defaultValue uint32) ([]uint32, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint32, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultUint32(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetUint24(ctx context.Context, k string, fields []string, defaultValue atype.Uint24) ([]atype.Uint24, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]atype.Uint24, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultUint24(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetUint16(ctx context.Context, k string, fields []string, defaultValue uint16) ([]uint16, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint16, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultUint16(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetUint8(ctx context.Context, k string, fields []string, defaultValue uint8) ([]uint8, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint8, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultUint8(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetInt64(ctx context.Context, k string, fields []string, defaultValue int64) ([]int64, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int64, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultInt64(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetInt(ctx context.Context, k string, fields []string, defaultValue int) ([]int, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultInt(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetInt32(ctx context.Context, k string, fields []string, defaultValue int32) ([]int32, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int32, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultInt32(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetInt16(ctx context.Context, k string, fields []string, defaultValue int16) ([]int16, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int16, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultInt16(0)
		}
	}
	return v, ok, nil
}
func (r *RedisClient) TryHMGetInt8(ctx context.Context, k string, fields []string, defaultValue int8) ([]int8, bool, *ae.Error) {
	iv, ok, e := r.TryHMGet(ctx, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int8, len(fields))
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = atype.New(x).DefaultInt8(0)
		}
	}
	return v, ok, nil
}

// 不能有一个是nil
func (r *RedisClient) MustHMGet(ctx context.Context, k string, fields ...string) ([]any, *ae.Error) {
	v, err := r.rdb.HMGet(ctx, k, fields...).Result()
	if err != nil {
		return nil, ae.NewRedisError(err)
	}
	if len(v) != len(fields) {
		return nil, ae.NotFound
	}
	for _, x := range v {
		if atype.IsNil(x) {
			return v, ae.NotFound
		}
	}
	return v, nil
}
func (r *RedisClient) MustHMGetUint64(ctx context.Context, k string, fields ...string) ([]uint64, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint64, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultUint64(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetUint(ctx context.Context, k string, fields ...string) ([]uint, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultUint(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetUint32(ctx context.Context, k string, fields ...string) ([]uint32, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint32, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultUint32(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetUint24(ctx context.Context, k string, fields ...string) ([]atype.Uint24, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]atype.Uint24, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultUint24(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetUint16(ctx context.Context, k string, fields ...string) ([]uint16, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint16, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultUint16(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetUint8(ctx context.Context, k string, fields ...string) ([]uint8, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint8, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultUint8(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetInt64(ctx context.Context, k string, fields ...string) ([]int64, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int64, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultInt64(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetInt(ctx context.Context, k string, fields ...string) ([]int, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultInt(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetInt32(ctx context.Context, k string, fields ...string) ([]int32, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int32, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultInt32(0)
	}
	return v, nil
}
func (r *RedisClient) MustHMGetInt16(ctx context.Context, k string, fields ...string) ([]int16, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int16, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultInt16(0)
	}
	return v, nil
}

func (r *RedisClient) MustHMGetInt8(ctx context.Context, k string, fields ...string) ([]int8, *ae.Error) {
	iv, e := r.MustHMGet(ctx, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int8, len(fields))
	for i, a := range iv {
		v[i] = atype.New(a).DefaultInt8(0)
	}
	return v, nil
}
