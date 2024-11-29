package models

type Order struct {
	Model
	PickupNo     string         `gorm:"not null;size:8" json:"pickup_number"`
	OrderNo      string         `gorm:"not null;size:32" json:"order_number"`
	Address      string         `gorm:"size:100;not null" json:"address"`
	CustomerName string         `gorm:"size:20;not null" json:"customer_name"`
	PhoneNumber  string         `gorm:"size:20;not null" json:"phone_number"`
	Status       uint8          `gorm:"not null;default:0" json:"status"`
	PaymentTime  uint64         `gorm:"not null;default:0" json:"payment_time"`
	WechatUserID uint           `gorm:"index;not null" json:"-"`
	OrderDetails []*OrderDetail `json:"details"`
	// TODO:加入接单骑手号
}

const (
	// 订单没有支付
	OrderNotPayed = 0
	// 订单已经支付, 等待抢单
	OrderPayed = 1
	// 订单等待配送
	OrderToDeliver = 2
	// 订单已经完成
	OrderFinished = 3
)

// NOTE:
// 口味直接存储，没有必要再联合两张表查一次
// Dish暂时不直接存储，因为涉及图片的展示等等
type OrderDetail struct {
	Model
	DishID    string `gorm:"not null" json:"dish_id"`
	Flavor    string `gorm:"size:20" json:"dish_flavor"`
	DishCount uint16 `gorm:"not null;default 0" json:"dish_count"`
	DishPrice uint   `gorm:"dish_price"`
	OrderID   uint   `gorm:"index;not null" json:"-"`
}
