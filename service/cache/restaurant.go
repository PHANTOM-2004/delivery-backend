package cache

import (
	"delivery-backend/internal/gredis"
	"strconv"
)

// NOTE:餐馆营业时缓存，暂时支持营业时间的缓存;营业状态只依赖于商家的手动设定。
//
// 对于商家来说，暂时只需要考虑商家是否设置了打烊状态；
// 目前的营业状态仅仅根据商家手动设定决定。后续再考虑是否加入自动设置营业状态

// 用于读取redis的缓存
type RedisRestaurantStatus struct {
	key string
}

func NewRedisRestaurantStatus(restaurant_id uint) *RedisRestaurantStatus {
	return &RedisRestaurantStatus{key: "RESTAURANT_" + strconv.Itoa(int(restaurant_id))}
}

func (r *RedisRestaurantStatus) GetKeyName() string {
	return r.key
}

// 设置不会过期的缓存
func (r *RedisRestaurantStatus) Set(status string) error {
	err := gredis.Set(r.key, status, 0)
	return err
}

// 如果商家没有设置状态，默认返回关店状态0.
func (r *RedisRestaurantStatus) Get() (uint8, error) {
	res, err := gredis.Get(r.key)
	if err != nil {
		return 0, err
	}
	if res == nil {
		// 如果找不到返回关门状态
		return 0, nil
	}
	status, _ := strconv.Atoi(string(res))
	return uint8(status), nil
}
