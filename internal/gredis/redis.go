package gredis

import (
	"context"
	"delivery-backend/internal/setting"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Setup() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:           setting.RedisSetting.Host,
		Password:       setting.RedisSetting.Password,
		MaxActiveConns: setting.RedisSetting.MaxActive,
		DB:             0,
	})
}

func Set(key string, data any, expiresTime int) (bool, error) {
	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	err = Rdb.Set(ctx, key, value, time.Second*time.Duration(expiresTime)).Err()
	if err != nil {
		return false, err
	}

	return true, err
}

func Exists(key string) bool {
	ctx := context.Background()
	res, err := Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return res != 0
}

func Get(key string) ([]byte, error) {
	ctx := context.Background()
	value, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}

func Delete(key string) (bool, error) {
	ctx := context.Background()
	res, err := Rdb.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res != 0, nil
}
