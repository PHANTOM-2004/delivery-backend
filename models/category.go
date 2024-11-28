package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Category struct {
	Model
	// 分类名称，比如素菜，荤菜，小吃，套餐
	Name string `gorm:"size:30;not null" json:"name"`
	// 分类类型，1代表菜品，2代表套餐;默认是菜品
	Type uint8 `gorm:"default:1;not null" json:"type"`
	// 排序值，用于决定顺序；
	Sort uint16 `gorm:"default:0;not null" json:"sort"`
	// 1代表禁用，2代表启用，默认禁用
	Status uint8 `gorm:"default:1;not null" json:"status"`
	// 考虑到每一家店铺大概率有一个独立的分类，因此每一个category对应一家店铺
	Restaurant   *Restaurant `json:"-"`
	RestaurantID uint        `gorm:"index;not null" json:"-"`
	Dishes       []*Dish     `gorm:"many2many:category_dish" json:"dishes"`
}

func (c *Category) AfterDelete(tx *gorm.DB) error {
	log.Trace("running category after delete hook")
	// 目前category和dishes是分立的，因此不必主动删除
	return nil
}

// NOTE:注意为什么得到空结构体
// https://github.com/go-gorm/gorm/issues/3686
func DeleteCategory(category_id uint) error {
	err := tx.Delete(&Category{Model: Model{ID: category_id}}).Error
	return err
}

// 同时会返回对应的dishes
func GetCategoryByRestaurant(restaurant_id uint) ([]Category, error) {
	c := []Category{}
	err := tx.Preload("Dishes").
		Find(&c, Category{RestaurantID: restaurant_id}).Error
	return c, err
}

// 同时会返回对应的dishes, 以及flavors
func GetCategoryDishFlavor(restaurant_id uint) ([]Category, error) {
	c := []Category{}
	err := tx.Preload("Dishes.Flavors").
		Find(&c, Category{RestaurantID: restaurant_id}).Error
	return c, err
}

// 同时返回对应的dishes
func GetCategoryDish(category_id uint) (*Category, error) {
	res := Category{}
	err := tx.Preload("Dishes").Find(&res, category_id).Error
	return &res, err
}

func CreateCategory(data *Category) error {
	err := tx.Create(data).Error
	return err
}

// 注意更新不存在的category的情况
func UpdateCategory(category_id uint, data *Category) error {
	err := tx.Model(&Category{}).Where("id = ?", category_id).Updates(data).Error
	return err
}

func AddCategoryDish(category_id uint, dishes_id []uint) error {
	dishes := make([]Dish, len(dishes_id))
	for i := range dishes {
		id := dishes_id[i]
		dishes[i] = Dish{Model: Model{ID: id}}
	}
	err := tx.Model(&Category{Model: Model{ID: category_id}}).Association("Dishes").Append(dishes)
	return err
}

func DeleteCategoryDish(category_id uint, dishes_id []uint) error {
	dishes := make([]Dish, len(dishes_id))
	for i := range dishes {
		id := dishes_id[i]
		dishes[i] = Dish{Model: Model{ID: id}}
	}
	err := tx.Model(&Category{Model: Model{ID: category_id}}).Association("Dishes").Delete(dishes)
	return err
}
