package models

type Dish struct {
	Model
	Name        string   `gorm:"size:30;not null" json:"name"`
	Price       uint     `gorm:"default:0;not null" json:"price"`
	Image       string   `gorm:"size:256;not null" json:"image"`
	Description string   `gorm:"size:50" json:"description"`
	Sort        uint16   `gorm:"default:0;not null" json:"sort"`
	CategoryID  uint     `gorm:"not null" json:"category_id"`
	Category    Category `json:"-"`
	Flavors     []Flavor `gorm:"many2many:dish_flavor" json:"flavors"`
}

type Flavor struct {
	Model
	Name string `gorm:"size:30;not null" json:"name"`
	// 关联的商店
	RestaurantID uint `gorm:"index"`
}

// 创建一个新菜品
func CreateDish(d *Dish) error {
	err := tx.Create(d).Error
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
	err := tx.Model(&Dish{}).Where("id = ?", flavor_id).Update("name", name).Error
	return err
}

// 获得一个店铺所有的flavor
func GetFlavors(restaurant_id uint) ([]Flavor, error) {
	flavors := []Flavor{}
	err := tx.Find(&flavors, Flavor{RestaurantID: restaurant_id}).Error
	return flavors, err
}

// 一个菜品可以添加多种口味
func AddDishFlavor(dish_id uint, flavors []Flavor) error {
	err := tx.Model(&Dish{Model: Model{ID: dish_id}}).Association("Flavors").Append(flavors)
	return err
}
