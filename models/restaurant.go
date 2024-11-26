package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Restaurant struct {
	Model
	RestaurantName string `gorm:"size:50;index;not null" json:"restaurant_name"`
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

func (r *Restaurant) AfterDelete(tx *gorm.DB) error {
	log.Trace("running restaurant after delete hook")
	err := tx.Where("restaurant_id = ?", r.ID).Delete(&Category{}).Error
	if err != nil {
		return err
	}
	log.Tracef("categories related to restaurant[%v] are deleted", r.ID)
	err = tx.Where("restaurant_id = ?", r.ID).Delete(&Flavor{}).Error
	if err != nil {
		return err
	}
	log.Tracef("flavors related to restaurant[%v] are deleted", r.ID)

	err = tx.Where("restaurant_id = ?", r.ID).Delete(&Dish{}).Error
	if err != nil {
		return err
	}
	log.Tracef("dishes related to restaurant[%v] are deleted", r.ID)
	return nil
}

func DeleteRestaurant(restaurant_id uint) error {
	err := tx.Delete(&Restaurant{Model: Model{ID: restaurant_id}}).Error
	return err
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

// 如果重名就会返回失败
func CreateRestaurant(data *Restaurant) (bool, error) {
	var err error
	success := false

	err = tx.Transaction(
		func(ftx *gorm.DB) error {
			var err error
			name := data.RestaurantName
			if name != "" {
				// 检测是否存在重名
				r := Restaurant{}
				err = ftx.Find(&r,
					Restaurant{RestaurantName: name}).Error
				if err != nil {
					return err
				}
				if r.ID != 0 {
					success = false
					return nil
				}
			}
			// 不存在重名那么就可以创建
			err = ftx.Create(data).Error
			success = true
			return err
		},
	)

	return success, err
}

// 如果重名就会返回失败
func UpdateRestaurant(restaurant_id uint, data *Restaurant) (bool, error) {
	var err error
	success := false

	err = tx.Transaction(
		func(ftx *gorm.DB) error {
			var err error
			name := data.RestaurantName
			if name != "" {
				// 检测是否存在重名
				r := Restaurant{}
				err = ftx.Find(&r,
					Restaurant{RestaurantName: name}).Error
				if err != nil {
					return err
				}
				// 注意重名的时候, 其实可以是自己
				if r.ID != 0 && r.ID != restaurant_id {
					success = false
					return nil
				}
			}
			// 不存在重名那么就可以插入
			err = ftx.Model(Restaurant{Model: Model{ID: restaurant_id}}).
				Updates(*data).Error
			success = true
			return err
		},
	)

	return success, err
}
