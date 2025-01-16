package redishelper

import (
	"context"
	"errors"
	"fmt"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/types"
	"golang.org/x/exp/constraints"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUnexpectedNilResult = errors.New("unexpected nil result")
)

// HSet
//
//   - HSet("key1", "value1", "key2", "value2")
//
//   - HSet([]string{"key1", "value1", "key2", "value2"})
//
//   - HSet(map[string]interface{}{"key1": "value1", "key2": "value2"})
//     Playing struct With "redis" tag.
//     type MyHash struct { Key1 string `redis:"key1"`; Key2 int `redis:"key2"` }
//
//   - HSet("myhash", MyHash{"value1", "value2"}) Warn: redis-server >= 4.0
//
// HMSet deprecated after redis 3
func HSet(ctx context.Context, rdb *redis.Client, expires time.Duration, key string, data ...any) error {
	if len(data) == 0 {
		return ae.ErrInvalidInput
	}
	dst := data
	if len(data) == 1 {
		arg := data[0]
		switch arg := arg.(type) {
		case []string:
			for _, s := range arg {
				dst = append(dst, s)
			}
		case []interface{}:
			dst = append(dst, arg...)
		case map[string]interface{}:
			for k, v := range arg {
				dst = append(dst, k, v)
			}
		case map[string]string:
			for k, v := range arg {
				dst = append(dst, k, v)
			}
		default:
			// scan struct field
			v := reflect.ValueOf(arg)
			if v.Type().Kind() == reflect.Ptr {
				if v.IsNil() {
					return ae.ErrInvalidInput
				}
				v = v.Elem()
			}

			if v.Type().Kind() == reflect.Struct {
				dst = appendStructField(dst, v)
			}
			return ae.ErrInvalidInput
		}
	}
	ttl, err := rdb.TTL(ctx, key).Result()
	if err == nil && ttl > 0 {
		return rdb.HSet(ctx, key, dst...).Err()
	}

	_, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err1 := pipe.HSet(ctx, key, dst...).Err()
		err2 := pipe.Expire(ctx, key, expires).Err()
		return ae.FirstErr(err1, err2)
	})

	return err
}

func HScan(ctx context.Context, rdb *redis.Client, k string, dest any, fields ...string) error {
	c := rdb.HMGet(ctx, k, fields...)
	v, err := c.Result()
	if err != nil {
		return err
	}
	if len(v) != len(fields) {
		return errors.New("fields mismatch")
	}
	for _, x := range v {
		if types.IsNil(x) {
			return ErrUnexpectedNilResult
		}
	}
	return c.Scan(dest)
}

func HGetAll(ctx context.Context, rdb *redis.Client, k string, dest any) error {
	c := rdb.HGetAll(ctx, k)
	result, err := c.Result()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return redis.Nil
	}
	return c.Scan(dest)
}

func HGetAllInt(ctx context.Context, rdb *redis.Client, k string, restrict bool) (map[string]int, error) {
	c := rdb.HGetAll(ctx, k)
	result, err := c.Result()
	if err != nil {
		return nil, err
	}
	n := len(result)
	if n == 0 {
		return nil, redis.Nil
	}
	value := make(map[string]int, n)
	var x int64
	for k, v := range result {
		if x, err = types.ParseInt64(v); err != nil {
			if restrict {
				return nil, fmt.Errorf(`invalid int string %s`, v)
			}
			continue
		}
		value[k] = int(x)
	}
	return value, nil
}

// 只要存在一个，就不报错；全是nil，返回 ae.ErrorNotFound
func TryHMGet(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]any, bool, error) {
	v, err := rdb.HMGet(ctx, k, fields...).Result()
	if err != nil {
		return nil, false, err
	}
	n := len(v)
	if n != len(fields) {
		return nil, false, redis.Nil
	}
	ok := true
	err = redis.Nil
	for _, x := range v {
		if !types.IsNil(x) {
			err = nil // 只要存在一个不是nil，都正确
			if !ok {
				break
			}
		} else {
			ok = false
			if err == nil {
				break
			}
		}
	}
	return v, ok, err
}

// 只要存在一个，就不报错；全是nil，返回 ae.ErrorNotFound
func TryHMGetString(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]string, bool, error) {
	iv, ok, err := TryHMGet(ctx, rdb, k, fields...)
	if err != nil {
		return nil, ok, err
	}
	v := make([]string, len(fields))
	for i, x := range iv {
		if types.IsNil(x) {
			v[i] = ""
		} else {
			v[i] = atype.String(x)
		}
	}
	return v, ok, nil
}
func TryHMGetN[T constraints.Integer](ctx context.Context, rdb *redis.Client, k string, fields []string, defaultValue T, panicOnNil bool) ([]T, bool, error) {
	iv, ok, err := TryHMGet(ctx, rdb, k, fields...)
	if err != nil {
		return nil, ok, err
	}

	result := make([]T, 0, len(fields))
	for _, v := range iv {
		if types.IsNil(v) {
			if panicOnNil {
				return nil, ok, ErrUnexpectedNilResult
			}
			result = append(result, defaultValue)
		} else {
			switch x := v.(type) {
			case uint8:
				result = append(result, T(x))
			case uint16:
				result = append(result, T(x))
			case uint32:
				result = append(result, T(x))
			case uint64:
				result = append(result, T(x))
			case int64:
				result = append(result, T(x))
			case int32:
				result = append(result, T(x))
			case int:
				result = append(result, T(x))
			case int16:
				result = append(result, T(x))
			case int8:
				result = append(result, T(x))
			case float64:
				result = append(result, T(x))
			case float32:
				result = append(result, T(x))
			case string:
				if x == "" {
					if panicOnNil {
						return nil, ok, ErrUnexpectedNilResult
					}
					result = append(result, defaultValue)
				} else {
					isSigned := x[0] == '-'
					if isSigned {
						var y int64
						if y, err = types.ParseInt64(x, 10); err != nil {
							return nil, ok, fmt.Errorf("invalid number value:%s", x)
						}
						result = append(result, T(y))
					} else {
						var y uint64
						if y, err = types.ParseUint64(x, 10); err != nil {
							return nil, ok, fmt.Errorf("invalid number value:%s", x)
						}
						result = append(result, T(y))
					}
				}
			default:
				return nil, ok, fmt.Errorf("invalid number value:%v", x)
			}
		}
	}
	return result, ok, nil
}

// 不能有一个是nil
func MustHMGet(ctx context.Context, rdb *redis.Client, k string, fields ...string) ([]any, error) {
	v, err := rdb.HMGet(ctx, k, fields...).Result()
	if err != nil {
		return nil, err
	}
	if len(v) != len(fields) {
		return nil, redis.Nil
	}
	for _, x := range v {
		if types.IsNil(x) {
			return v, ErrUnexpectedNilResult
		}
	}
	return v, nil
}
