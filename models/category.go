package models

type Category struct {
	Model
	// 分类名称，比如素菜，荤菜，小吃，套餐
	Name string `gorm:"size:30;not null" json:"name"`
	// 分类类型，0代表菜品，1代表套餐;默认是菜品
	Type uint8 `gorm:"default:0;not null" json:"type"`
	// 排序值，用于决定顺序；
	Sort uint16 `gorm:"default:0;not null" json:"sort"`
	// 0代表禁用，1代表启用，默认禁用
	Status uint8 `gorm:"default:0;not null" json:"status"`
	// 考虑到每一家店铺大概率有一个独立的分类，因此每一个category对应一家店铺
	Restaurant   Restaurant `json:"-"`
	RestaurantID uint       `gorm:"index;not null" json:"restaurant_id"`
	Dishes       []Dish     `json:"dishes"`
}

// 同时会返回对应的dishes
func GetCategoryByRestaurant(restaurant_id uint) ([]Category, error) {
	c := []Category{}
	err := tx.Preload("Dishes").
		Find(&c, Category{RestaurantID: restaurant_id}).Error
	return c, err
}

// 同时返回对应的dishes
func GetCategory(category_id uint) (*Category, error) {
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
