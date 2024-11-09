package gredis

import (
	"context"
	"delivery-backend/internal/setting"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
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

// Zero expiration means the key has no expiration time.
func Set(key string, data any, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx := context.Background()

	err = Rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func Exists(key string) (bool, error) {
	ctx := context.Background()
	res, err := Rdb.Exists(ctx, key).Result()
	if err == redis.Nil {
		// 不存在
		return false, nil
	}
	return res != 0, err
}

// 找不到的时候返回nil nil
func Get(key string) ([]byte, error) {
	ctx := context.Background()
	value, err := Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// 不存在
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return []byte(value), nil
}

// key不存在仍然delete成功, 返回nil
func Delete(key string) error {
	ctx := context.Background()
	_, err := Rdb.Del(ctx, key).Result()
	if err == redis.Nil {
		// 不存在
		log.Warn("redis: delete non exist key:", key)
		return nil
	}
	return err
}
