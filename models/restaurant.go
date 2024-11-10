package models

type Restaurant struct {
	Model
	RestaurantName string `gorm:"size:50;uniqueIndex;not null" json:"restaurant_name"`
	// 店铺的地址
	Address string `gorm:"size:50;not null" json:"address"`
	// 商铺简介
	Description string `gorm:"size:300;not null" json:"description"`
	// 最小起送金额,使用整数存储,默认存储到分
	MinimumDeliveryAmount uint `gorm:"default:0;not null" json:"minimum_delivery_amount"`
	// 商铺评分
	Rating float32 `gorm:"default:0;not null" json:"rating"`
	// 所属的商家ID
	MerchantID uint     `gorm:"index" json:"merchant_id"`
	Merchant   Merchant `json:"-"` // 该部分不参与json化
}

type RestaurantTime struct {
	Model
	// 存放分钟数
	OpenTime uint16 `gorm:"not null" json:"open_time"`
	// 存放分钟数
	CloseTime uint16 `gorm:"not null" json:"close_time"`
	// 时间段，采用格式:x-x 代表周几到周几,只需要3个字符即可
	OpenDays string `gorm:"type:char(3);not null" json:"open_days"`
	// 所属的餐馆ID, 因为一个商家可能存在多个营业时间
	RestaurantID uint `gorm:"index;not null" json:"restaurant_id"`
}

func GetRestaurantByMerchant(merchant_id uint) ([]Restaurant, error) {
	r := []Restaurant{}
	err := tx.
		Where("merchant_id = ?", merchant_id).
		Find(&r).Error
	return r, err
}

func GetRestaurantByID(restaurant_id uint) (*Restaurant, error) {
	r := Restaurant{}
	err := tx.Find(&r,
		Restaurant{Model: Model{ID: restaurant_id}}).
		Error
	return &r, err
}

func CreateRestaurant(data *Restaurant) error {
	err := tx.Create(data).Error
	return err
}

func UpdateRestaurant(restaurant_id uint, data *Restaurant) error {
	err := tx.
		Model(Restaurant{Model: Model{ID: restaurant_id}}).
		Updates(*data).Error
	return err
}

func ExistRestaurant(name string) (bool, error) {
	r := Restaurant{}
	err := tx.Find(&r,
		Restaurant{RestaurantName: name}).Error
	return r.ID != 0, err
}
