package db

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

func FetchCache[T any](redisClient *redis.Client, key string, ttl time.Duration, function func() (T, error)) (T, error) {
	var ctx = context.Background()
	value := redisClient.Get(ctx, key)
	var result T
	if value.Err() == nil {
		err := json.Unmarshal([]byte(value.Val()), &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	result, err := function()
	if err != nil {
		return result, err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return result, err
	}
	redisClient.Set(ctx, key, bytes, ttl)
	return result, nil
}

func GetCaches[T any](redisClient *redis.Client, keys []string) ([]T, error) {
	var ctx = context.Background()
	value := redisClient.MGet(ctx, keys...)
	var result = make([]T, 0)
	for _, v := range value.Val() {
		if v == nil {
			continue
		}
		var ret T
		err := json.Unmarshal(([]byte)(v.(string)), &ret)
		if err != nil {
			return nil, err
		}
		result = append(result, ret)
	}
	return result, nil
}

func SetCache(redisClient *redis.Client, key string, value any, ttl time.Duration) error {
	var ctx = context.Background()
	var s string
	if reflect.TypeOf(value).Kind() == reflect.String {
		s = value.(string)
	} else {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		s = string(data)
	}
	result := redisClient.Set(ctx, key, s, ttl)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
