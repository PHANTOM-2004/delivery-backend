package redis

import (
	"context"
	"delivery-backend/service/api/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init() {
	// RedisClient = redis.NewClient(&redis.Options{
	// 	Addr:     conf.GetConf().Redis.Address,
	// 	Username: conf.GetConf().Redis.Username,
	// 	Password: conf.GetConf().Redis.Password,
	// 	DB:       conf.GetConf().Redis.DB,
	// })
	//
	defer klog.Info("api service: redis init done")
	RedisClient = redis.NewFailoverClient(
		&redis.FailoverOptions{
			MasterName:    "delivery-master",
			SentinelAddrs: []string{conf.GetConf().Redis.Address},
			Username:      conf.GetConf().Redis.Username,
			DB:            conf.GetConf().Redis.DB,
			Password:      conf.GetConf().Redis.Password,
		},
	)

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
