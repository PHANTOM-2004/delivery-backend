package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Dish struct {
	Model
	Name         string    `gorm:"size:30;not null" json:"name"`
	Price        uint      `gorm:"default:0;not null" json:"price"`
	Image        string    `gorm:"size:256;not null" json:"image"`
	Description  string    `gorm:"size:50" json:"description"`
	Sort         uint16    `gorm:"default:0;not null" json:"sort"`
	RestaurantID uint      `gorm:"not null" json:"-"`
	Flavors      []*Flavor `gorm:"many2many:dish_flavor" json:"flavors"`
}

type Flavor struct {
	Model
	Name string `gorm:"size:30;not null" json:"name"`
	// 关联的商店
	RestaurantID uint `gorm:"index" json:"-"`
}

// 从口味得到商家
func GetMerchantIDByFlavor(flavor_id uint) (uint, error) {
	var ID uint

	err := tx.Transaction(func(ftx *gorm.DB) error {
		f := Flavor{Model: Model{ID: flavor_id}}
		err := ftx.First(&f, Flavor{}).Error
		if err != nil {
			return err
		}
		log.Trace(f)
		log.Tracef("flavor[%v] belongs to restaurant[%v]", flavor_id, f.RestaurantID)

		r := Restaurant{}
		err = ftx.Find(&r, f.RestaurantID).Error
		if err != nil {
			return err
		}
		log.Trace(r)
		log.Tracef("restaurant[%v] belongs to merchant[%v]", f.RestaurantID, r.MerchantID)

		ID = r.MerchantID
		return nil
	})

	if err == nil {
		defer log.Tracef("transaction get merchant[%v] by flavor[%v]", ID, flavor_id)
	}

	return ID, err
}

// 删除一个口味
func DeleteFlavor(flavor_id uint) error {
	err := tx.Delete(&Flavor{Model: Model{ID: flavor_id}}).Error
	return err
}

// 创建一个新菜品
func CreateDish(d *Dish) error {
	err := tx.Create(d).Error
	return err
}

// 删除一个菜品
func DeleteDish(dish_id uint) error {
	err := tx.Delete(&Dish{Model: Model{ID: dish_id}}).Error
	return err
}

// get dish的时候同时获得所有flavor
func GetDishFlavors(dish_id uint) (*Dish, error) {
	d := Dish{}
	err := tx.Preload("Flavors").Find(&d, dish_id).Error
	return &d, err
}

// 注意更新不存在的dish的情况
func UpdateDish(id uint, d *Dish) error {
	err := tx.Model(&Dish{}).Where("id = ?", id).Updates(*d).Error
	return err
}

func GetFlavor(flavor_id uint) (*Flavor, error) {
	f := Flavor{}
	err := tx.Find(&f, flavor_id).Error
	return &f, err
}

// 创建一个新口味
func CreateFlavor(f *Flavor) error {
	err := tx.Create(f).Error
	return err
}

func UpdateFlavor(flavor_id uint, name string) error {
	err := tx.Model(&Flavor{}).Where("id = ?", flavor_id).Update("name", name).Error
	return err
}

// 获得一个店铺所有的flavor
func GetFlavors(restaurant_id uint) ([]Flavor, error) {
	flavors := []Flavor{}
	err := tx.Find(&flavors, Flavor{RestaurantID: restaurant_id}).Error
	return flavors, err
}

func GetDishes(restaurant_id uint) ([]Dish, error) {
	dishes := []Dish{}
	err := tx.Preload("Flavors").Find(&dishes).Error
	return dishes, err
}

// 只取几个dishes作为demo
func GetTopDishes(restaurant_id uint) ([]Dish, error) {
	dishes := []Dish{}
	page_size := 4
	err := tx.Limit(page_size).Find(&dishes).Error
	return dishes, err
}

// 一个菜品可以添加多种口味
func AddDishFlavor(dish_id uint, flavors_id []uint) error {
	flavors := make([]Flavor, len(flavors_id))
	for i := range flavors {
		id := flavors_id[i]
		flavors[i] = Flavor{Model: Model{ID: id}}
	}
	err := tx.Model(&Dish{Model: Model{ID: dish_id}}).Association("Flavors").Append(flavors)
	return err
}

func DeleteDishFlavor(dish_id uint, flavors_id []uint) error {
	flavors := make([]Flavor, len(flavors_id))
	for i := range flavors {
		id := flavors_id[i]
		flavors[i] = Flavor{Model: Model{ID: id}}
	}
	err := tx.Model(&Dish{Model: Model{ID: dish_id}}).Association("Flavors").Delete(flavors)
	return err
}
