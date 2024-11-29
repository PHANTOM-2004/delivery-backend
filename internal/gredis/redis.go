package gredis

import (
	"context"
	"delivery-backend/internal/setting"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var rdb *redis.Client

func Setup() {
	rdb = redis.NewClient(&redis.Options{
		Addr:           setting.RedisSetting.Host,
		Password:       setting.RedisSetting.Password,
		MaxActiveConns: setting.RedisSetting.MaxActive,
		DB:             0,
	})
}

// 设置string value, 不经过序列化
func SetStr(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	err := rdb.Set(ctx, key, value, expiration).Err()
	return err
}

// 作为json存进去
// Zero expiration means the key has no expiration time.
func Set(key string, data any, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = rdb.Set(ctx, key, value, expiration).Err()
	return err
}

func HSet(key string, pairs []string) error {
	if len(pairs)%2 != 0 {
		log.Panic("invalid usage")
	}
	ctx := context.Background()
	err := rdb.HSet(ctx, key, pairs).Err()
	return err
}

func HGet(key string, field string) (string, error) {
	ctx := context.Background()
	res, err := rdb.HGet(ctx, key, field).Result()
	return res, err
}

func Expire(key string, expires time.Duration) error {
	ctx := context.Background()
	err := rdb.Expire(ctx, key, expires).Err()
	return err
}

func Exists(key string) (bool, error) {
	ctx := context.Background()
	res, err := rdb.Exists(ctx, key).Result()
	if err == redis.Nil {
		// 不存在
		return false, nil
	}
	return res != 0, err
}

// 找不到的时候返回nil nil
func Get(key string) ([]byte, error) {
	ctx := context.Background()
	value, err := rdb.Get(ctx, key).Result()
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
	_, err := rdb.Del(ctx, key).Result()
	if err == redis.Nil {
		// 不存在
		log.Warn("redis: delete non exist key:", key)
		return nil
	}
	return err
}

func Incre(key string) (int64, error) {
	ctx := context.Background()
	num, err := rdb.Incr(ctx, key).Result()
	return num, err
}
