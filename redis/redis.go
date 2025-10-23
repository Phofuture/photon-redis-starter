package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func Get[T any](ctx context.Context, key string) (T, error) {
	var value T
	if str, err := rdb.Get(ctx, key).Result(); err != nil {
		return value, err
	} else if err := json.Unmarshal([]byte(str), &value); err != nil {
		return value, err
	} else {
		return value, nil
	}
}

func HGet[T any](ctx context.Context, key string, field string) (value T, err error) {
	str, err := rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return
	} else if err = json.Unmarshal([]byte(str), &value); err != nil {
		return
	}
	return
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

func HSet(ctx context.Context, key string, values map[string]interface{}) error {
	return rdb.HSet(ctx, key, values).Err()
}

func Exists(ctx context.Context, key string) bool {
	result, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return rdb.SetNX(ctx, key, value, expiration).Result()
}

func RPush(ctx context.Context, key string, value interface{}) (int64, error) {
	return rdb.RPush(ctx, key, value).Result()
}

func RPop[T any](ctx context.Context, key string) (T, error) {
	var value T
	if str, err := rdb.RPop(ctx, key).Result(); err != nil {
		return value, err
	} else if err := json.Unmarshal([]byte(str), &value); err != nil {
		return value, err
	} else {
		return value, nil
	}
}

func LRange[T any](ctx context.Context, key string, start, stop int64) ([]T, error) {
	var values []T
	if strs, err := rdb.LRange(ctx, key, start, stop).Result(); err != nil {
		return values, err
	} else {
		for _, str := range strs {
			var value T
			if err := json.Unmarshal([]byte(str), &value); err != nil {
				return values, err
			}
			values = append(values, value)
		}
		return values, nil
	}
}

func LRangeAll[T any](ctx context.Context, key string) ([]T, error) {
	return LRange[T](ctx, key, 0, -1)
}

func EvalBool(ctx context.Context, script string, keys []string, args ...interface{}) (ok bool, err error) {
	var result interface{}
	result, err = rdb.Eval(ctx, script, keys, args).Result()
	if err != nil {
		return
	}
	switch v := result.(type) {
	case int64:
		ok = v == 1
	default:
		return false, fmt.Errorf("unexpected type %T for result", v)
	}
	return
}

func Del(ctx context.Context, keys ...string) error {
	return rdb.Del(ctx, keys...).Err()
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rdb.Expire(ctx, key, expiration).Err()
}
