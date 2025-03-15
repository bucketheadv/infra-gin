package db

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"reflect"
	"slices"
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

type TableWithID[T cmp.Ordered] interface {
	TableName() string
	GetID() T
}

func GetModelCaches[T TableWithID[R], R cmp.Ordered](client *redis.Client, cacheKeyFormat string, ids []R, expires time.Duration, fallback func(missingIds []R) *gorm.DB) ([]T, error) {
	if len(ids) == 0 {
		return make([]T, 0), nil
	}

	var result = make([]T, 0)
	var missingIds = make([]R, 0)
	var keys = make([]string, 0)
	for _, id := range ids {
		var key = fmt.Sprintf(cacheKeyFormat, id)
		keys = append(keys, key)
	}
	foundModels, err := GetCaches[T](client, keys)
	if err != nil {
		return nil, err
	}
	var foundModelIds = make([]R, 0)
	for _, u := range foundModels {
		foundModelIds = append(foundModelIds, u.GetID())
		result = append(result, u)
	}
	for _, id := range ids {
		if !slices.Contains(foundModelIds, id) {
			missingIds = append(missingIds, id)
		}
	}

	if len(missingIds) > 0 {
		var models []T
		fallback(missingIds).Find(&models)
		for _, model := range models {
			var key = fmt.Sprintf(cacheKeyFormat, model.GetID())
			err := SetCache(client, key, model, expires)
			if err != nil {
				panic(err)
			}
			result = append(result, model)
		}
	}
	return result, nil
}
