package redis

import (
	"context"
	"strings"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

const (
	bucketCap    = 60 // 突发60次请求
	refillCycel  = 10 // 10s重新填充一次
	refillAmount = 10 // 每次添加10个
	bucketKey    = "tokenBucket"
)

const rateLimiterScript = `
local updateTime = ARGV[1]
local capacity = ARGV[2]
local refillCycle = ARGV[3]
local refillAmount = ARGV[4]
local key = KEYS[1]

local remainToken = capacity

-- 查看上一次访问的, 以及当前token余量
local result = redis.call("HMGET", key, "lastUpd", "tokens")

if result[1] then -- 如果上一次有访问记录
	-- 计算填充数量
	local lastUpdate = tonumber(result[1])
	local amount = math.floor((updateTime - lastUpdate) / refillCycle * refillAmount)
	-- 填充令牌桶
	remainToken = math.min(capacity, result[2] + amount)
	if remainToken > 0 then -- 消耗令牌
		remainToken = remainToken - 1
	end
else
	remainToken = capacity - 1
end
-- 更新redis key
redis.call("HMSET", key, "lastUpd", updateTime, "tokens", remainToken)
-- 设置key过期时间, 比如60的桶大小, 60/10 = 6s就填满, 6s内才起作用
redis.call("EXPIRE", key, math.floor(capacity / refillAmount))
-- 返回剩余token
return remainToken

`

// 某一个key的bucket
func AcquireBucket(key string) (int64, error) {
	now := time.Now().Unix()
	cacheKey := strings.Join(
		[]string{bucketKey, key},
		":",
	)

	remain, err := runScript(
		[]string{cacheKey},
		[]any{
			now,
			bucketCap,
			refillCycel,
			refillAmount,
		},
	)
	if err != nil {
		return 0, err
	}
	if remain == nil {
		return 0, nil
	}
	return remain.(int64), nil
}

func runScript(keys []string, args []any) (any, error) {
	val, err := redis.
		NewScript(rateLimiterScript).
		Run(context.Background(), RedisClient, keys, args...).Result()
	if err != nil && err != redis.Nil {
		klog.Error(err)
		return nil, nil
	}
	return val, nil
}
