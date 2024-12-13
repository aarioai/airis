package redishelper

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/atype"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func HSet(ctx context.Context, rdb *redis.Client, expires time.Duration, k string, values ...any) *ae.Error {
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err1 := pipe.HSet(ctx, k, values...).Err()
		err2 := pipe.Expire(ctx, k, expires).Err()
		return ae.FirstError(err1, err2)
	})

	return ae.NewRedisError(err)
}

func HMSet(ctx context.Context, rdb *redis.Client, expires time.Duration, k string, values ...any) *ae.Error {
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
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err1 := pipe.HMSet(ctx, k, values...).Err()
		err2 := pipe.Expire(ctx, k, expires).Err()
		return ae.FirstError(err1, err2)
	})

	return ae.NewRedisError(err)
}

func HScan(ctx context.Context, rdb *redis.Client, dest any, k string, fields ...string) *ae.Error {
	c := rdb.HMGet(ctx, k, fields...)
	v, err := c.Result()
	if err != nil {
		return ae.NewRedisError(err)
	}
	if len(v) != len(fields) {
		return ae.ErrorNotFound
	}
	for _, x := range v {
		if atype.IsNil(x) {
			return ae.ErrorNotFound
		}
	}
	e := ae.NewRedisError(c.Scan(dest))
	return e
}

func HGetAll(ctx context.Context, rdb *redis.Client, k string, dest any) *ae.Error {
	c := rdb.HGetAll(ctx, k)
	result, err := c.Result()
	if err != nil {
		return ae.NewRedisError(err)
	}
	if len(result) == 0 {
		return ae.ErrorNotFound
	}
	e := ae.NewRedisError(c.Scan(dest))
	return e
}

func HGetAllInt(ctx context.Context, rdb *redis.Client, k string, restrict bool) (map[string]int, *ae.Error) {
	c := rdb.HGetAll(ctx, k)
	result, err := c.Result()
	if err != nil {
		return nil, ae.NewRedisError(err)
	}
	n := len(result)
	if n == 0 {
		return nil, ae.ErrorNotFound
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

// 只要存在一个，就不报错；全是nil，返回 ae.ErrorNotFound
func TryHMGet(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]any, bool, *ae.Error) {
	v, err := rdb.HMGet(ctx, k, fields...).Result()
	if err != nil {
		return nil, false, ae.NewRedisError(err)
	}
	n := len(v)
	if n != len(fields) {
		return nil, false, ae.ErrorNotFound
	}
	ok := true
	e := ae.ErrorNotFound
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

// 只要存在一个，就不报错；全是nil，返回 ae.ErrorNotFound
func TryHMGetString(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]string, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]string, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = ""
		} else {
			v[i] = newV.Reload(x).String()
		}
	}
	return v, ok, nil
}
func TryHMGetUint64(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue uint64) ([]uint64, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint64, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultUint64(0)
		}
	}
	return v, ok, nil
}
func TryHMGetUint(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue uint) ([]uint, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultUint(0)
		}
	}
	return v, ok, nil
}
func TryHMGetUint32(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue uint32) ([]uint32, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint32, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultUint32(0)
		}
	}
	return v, ok, nil
}
func TryHMGetUint24(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue atype.Uint24) ([]atype.Uint24, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]atype.Uint24, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultUint24(0)
		}
	}
	return v, ok, nil
}
func TryHMGetUint16(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue uint16) ([]uint16, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint16, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultUint16(0)
		}
	}
	return v, ok, nil
}
func TryHMGetUint8(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue uint8) ([]uint8, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, ok, e
	}
	v := make([]uint8, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultUint8(0)
		}
	}
	return v, ok, nil
}
func TryHMGetInt64(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue int64) ([]int64, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int64, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultInt64(0)
		}
	}
	return v, ok, nil
}
func TryHMGetInt(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue int) ([]int, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultInt(0)
		}
	}
	return v, ok, nil
}
func TryHMGetInt32(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue int32) ([]int32, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int32, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultInt32(0)
		}
	}
	return v, ok, nil
}
func TryHMGetInt16(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue int16) ([]int16, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int16, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultInt16(0)
		}
	}
	return v, ok, nil
}
func TryHMGetInt8(ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue int8) ([]int8, bool, *ae.Error) {
	iv, ok, e := TryHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, false, e
	}
	v := make([]int8, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, x := range iv {
		if atype.IsNil(x) {
			v[i] = defaultValue
		} else {
			v[i] = newV.Reload(x).DefaultInt8(0)
		}
	}
	return v, ok, nil
}

// 不能有一个是nil
func MustHMGet(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]any, *ae.Error) {
	v, err := rdb.HMGet(ctx, k, fields...).Result()
	if err != nil {
		return nil, ae.NewRedisError(err)
	}
	if len(v) != len(fields) {
		return nil, ae.ErrorNotFound
	}
	for _, x := range v {
		if atype.IsNil(x) {
			return v, ae.ErrorNotFound
		}
	}
	return v, nil
}
func MustHMGetUint64(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]uint64, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint64, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultUint64(0)
	}
	return v, nil
}
func MustHMGetUint(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]uint, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultUint(0)
	}
	return v, nil
}
func MustHMGetUint32(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]uint32, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint32, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultUint32(0)
	}
	return v, nil
}
func MustHMGetUint24(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]atype.Uint24, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]atype.Uint24, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultUint24(0)
	}
	return v, nil
}
func MustHMGetUint16(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]uint16, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint16, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultUint16(0)
	}
	return v, nil
}
func MustHMGetUint8(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]uint8, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]uint8, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultUint8(0)
	}
	return v, nil
}
func MustHMGetInt64(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]int64, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int64, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultInt64(0)
	}
	return v, nil
}
func MustHMGetInt(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]int, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultInt(0)
	}
	return v, nil
}
func MustHMGetInt32(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]int32, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int32, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultInt32(0)
	}
	return v, nil
}
func MustHMGetInt16(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]int16, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int16, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultInt16(0)
	}
	return v, nil
}

func MustHMGetInt8(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]int8, *ae.Error) {
	iv, e := MustHMGet(ctx, rdb, k, fields...)
	if e != nil {
		return nil, e
	}
	v := make([]int8, len(fields))
	newV := atype.New()
	defer newV.Release()
	for i, a := range iv {
		v[i] = newV.Reload(a).DefaultInt8(0)
	}
	return v, nil
}
