package wechat_service

import (
	"delivery-backend/internal/gredis"
	"fmt"
	"strconv"
	"time"
)

// TODO:采用定时任务，向数据库中刷新PickUpNumber
// 以及每天定时清零
type OrderGen struct {
	restaurant_id uint
}

func NewOrderGen(restaurant_id uint) *OrderGen {
	return &OrderGen{restaurant_id: restaurant_id}
}

func (p *OrderGen) counterKey() string {
	key := fmt.Sprintf("restaurant[%d]order_counter", p.restaurant_id)
	return key
}

func (p *OrderGen) numberKey() string {
	key := fmt.Sprintf("restaurant[%d]order_number", p.restaurant_id)
	return key
}

func (p *OrderGen) GetPickupNo() (string, error) {
	num, err := gredis.Incre(p.counterKey())
	return strconv.Itoa(int(num)), err
}

func (p *OrderGen) GetOrderNo() (string, error) {
	num, err := gredis.Incre(p.counterKey())
	if err != nil {
		return "", err
	}
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102150405")
	res := fmt.Sprintf("%s%d%d",
		formattedTime,
		190514*p.restaurant_id,
		num)
	return res, nil
}
