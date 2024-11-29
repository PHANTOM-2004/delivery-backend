package wechat_service

import (
	"delivery-backend/internal/gredis"
	"fmt"
)

// TODO:采用定时任务，向数据库中刷新PickUpNumber
// 以及每天定时清零
type PickupNo struct {
	restaurant_id uint
}

func NewPickupNo(restaurant_id uint) *PickupNo {
	return &PickupNo{restaurant_id: restaurant_id}
}

func (p *PickupNo) key() string {
	key := fmt.Sprintf("restaurant[%d]order_counter", p.restaurant_id)
	return key
}

func (p *PickupNo) Get(restaurant_id uint) (int, error) {
	num, err := gredis.Incre(p.key())
	return int(num), err
}
