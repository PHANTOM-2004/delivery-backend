package models

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/wechat"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// TODO:定期清理超时的订单

type Order struct {
	Model
	PickupNo         string            `gorm:"not null;size:8" json:"pickup_number"`
	OrderNo          string            `gorm:"not null;size:32" json:"order_number"`
	Address          string            `gorm:"size:100;not null" json:"address"`
	CustomerName     string            `gorm:"size:20;not null" json:"customer_name"`
	PhoneNumber      string            `gorm:"size:20;not null" json:"phone_number"`
	Status           uint8             `gorm:"not null;default:0" json:"status"`
	PaymentTime      uint64            `gorm:"not null;default:0" json:"payment_time"`
	RestaurantID     uint              `gorm:"not null" json:"restaurant_id"`
	Restaurant       *Restaurant       `json:"-"`
	RestaurantInfoEx *RestaurantInfoEx `gorm:"-" json:"restaurant_info"`
	WechatUserID     uint              `gorm:"index;not null" json:"-"`
	OrderDetails     []*OrderDetail    `json:"details"`
	// TODO:加入接单骑手号
}

type RestaurantInfoEx struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

const (
	// 订单没有支付
	OrderNotPayed = 0
	// 订单已经支付<->等待抢单
	OrderPayed = 1
	// 订单等待配送<->待取餐
	OrderToDeliver = 2
	// 拿到餐品<->等待送达
	OrderFetched = 3
	// 订单已经完成<->已送达
	OrderFinished = 4
	// 订单被取消
	OrderCanceled = 5
)

// NOTE:
// 口味直接存储，没有必要再联合两张表查一次
// Dish暂时不直接存储，因为涉及图片的展示等等
// 实际上还要保证，这里应当是冗余的，因为可能商家某些菜品会修改过，
// 所以说明细表要保留大多数
type OrderDetail struct {
	Model
	DishID     uint   `gorm:"not null" json:"dish_id"`
	DishName   string `gorm:"not null;size:30" json:"dish_name"`
	FlavorID   uint   `json:"flavor_id"`
	FlavorName string `json:"flavor_name" gorm:"size:30"`
	DishCount  uint16 `gorm:"not null;default:0" json:"dish_count"`
	DishPrice  uint   `gorm:"dish_price"`
	OrderID    uint   `gorm:"index;not null" json:"-"`
}

func CancelOrder(order_id uint) (bool, error) {
	success := false
	err := tx.Transaction(
		func(ftx *gorm.DB) error {
			order := Order{}
			err := tx.Find(&order, order_id).Error
			if err != nil {
				return err
			}
			if order.ID == 0 {
				// 没有找到
				log.Warnf("cancel order not found[%d]", order_id)
				return nil
			}
			if order.Status > OrderNotPayed {
				// 当前状态不可能被取消
				log.Warnf("cannot cancel order with status[%v]", order.Status)
				return nil
			}
			// 满足条件，更新状态
			err = tx.Model(&Order{}).Where("id = ?", order_id).Update("status", OrderCanceled).Error
			if err != nil {
				return err
			}
			success = true
			log.Tracef("order canceled[%d]", order_id)
			return nil
		},
	)
	return success, err
}

func GetOrderByUserID(user_id uint) ([]Order, error) {
	orders := []Order{}
	err := tx.Preload("OrderDetails").Find(&orders, Order{WechatUserID: user_id}).Error
	return orders, err
}

// dishes id以及对应的口味
// 记得保证stores参数不为空，也就是购物车为空的时候无法下单
func CreateOrder(order *Order, stores []wechat.WXSessionCartStore) error {
	err := tx.Transaction(
		func(ftx *gorm.DB) error {
			// 首先下单
			var err error
			err = ftx.Create(order).Error
			if err != nil {
				return err
			}
			order_id := order.ID
			if order_id == 0 {
				log.Error("order id不可能为0,因为创建成功")
			}

			// 0. 把store按照dish id进行排序
			sort.Slice(stores, func(i, j int) bool {
				return stores[i].DishID < stores[j].DishID
			})

			// 1. 查询出需要使用的所有dishes, 不包含口味信息,按照dishid进行排序
			// 准备好每一个dish的price
			dishes_id := make([]uint, len(stores))
			for i := range dishes_id {
				dishes_id[i] = stores[i].DishID
			}
			dishes := []Dish{}
			err = ftx.Order("id").Find(&dishes, dishes_id).Error
			if err != nil {
				return err
			}
			dishes_map := map[uint]*Dish{}
			for i := range dishes {
				dishes_map[dishes[i].ID] = &dishes[i]
			}
			log.Trace("prepared dishes:\n", dishes_map)

			// 2. 准备好所有口味信息，为了后续保留口味从中查找。
			// 由于口味数量比较小，所以使用线性查找即可
			flavors_id := []uint{}
			for i := range stores {
				if stores[i].FlavorID != 0 {
					flavors_id = append(flavors_id, stores[i].FlavorID)
				}
			}
			flavors := []Flavor{}
			err = ftx.Find(&flavors, flavors_id).Error
			if err != nil {
				return err
			}
			// 建立一个id -> string的哈希表
			flavors_map := map[uint]*Flavor{}
			for i := range flavors {
				flavors_map[flavors[i].ID] = &flavors[i]
			}
			log.Trace("prepared flavors:\n", flavors_map)

			// 3. 此时store中dish id是升序，对应到的dishes中的id也是升序，可以对应上
			order_details := make([]OrderDetail, len(stores))
			for i := range stores {
				// order id
				order_details[i].OrderID = order_id
				// dish info
				store_dish_id := stores[i].DishID
				order_details[i].DishID = dishes_map[store_dish_id].ID
				order_details[i].DishName = dishes_map[store_dish_id].Name
				order_details[i].DishPrice = dishes_map[store_dish_id].Price
				order_details[i].DishCount = uint16(stores[i].Cnt)

				// flavor info
				order_details[i].FlavorID = stores[i].FlavorID // flavor id 可能是0
				if flavor_id := stores[i].FlavorID; flavor_id != 0 {
					order_details[i].FlavorName = flavors_map[flavor_id].Name
				}
			}

			log.Trace("prepared order details", order_details)

			// 4. 创建订单明细
			err = ftx.Create(order_details).Error
			return err
		},
	)
	return err
}

func PayOrder(order_id uint) (bool, error) {
	success := false
	err := tx.Transaction(func(ftx *gorm.DB) error {
		var err error
		order := Order{}
		err = ftx.Find(&order, order_id).Error
		if err != nil {
			return err
		}
		if order.Status != OrderNotPayed {
			// 状态不符合支付
			return nil
		}
		now := time.Now().Unix()
		duration := now - int64(order.CreatedAt)
		if duration > 60*15 {
			// 超过15分钟
			return nil
		}

		err = ftx.Model(&Order{}).Where("id = ?", order_id).UpdateColumn("status", OrderPayed).Error
		if err != nil {
			return err
		}

		// 成功支付
		success = true
		return err
	})
	return success, err
}

func GetOrderByStatus(status uint8) ([]Order, error) {
	orders := []Order{}
	err := tx.Where("status = ?", status).
		Preload("OrderDetails").Preload("Restaurant").
		Find(&orders).Error
	return orders, err
}

func SetOrderStatus(order_id uint, status uint8) (bool, error) {
	res := tx.Model(&Order{}).Where("id = ?", order_id).UpdateColumn("status", status)
	return res.RowsAffected > 0, res.Error
}

// 1...n
func GetOrderByPage(page_id int) ([]Order, error) {
	page_size := setting.AppSetting.ApplicationPageSize
	offset := max(page_id-1, 0) * page_size
	orders := []Order{}
	err := tx.Limit(page_size).Offset(offset).
		Preload("OrderDetails").Preload("Restaurant").
		Find(&orders).Error
	return orders, err
}
