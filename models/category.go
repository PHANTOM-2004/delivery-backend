package models

type Category struct {
	Model
	// 分类名称，比如素菜，荤菜，小吃，套餐
	Name string `gorm:"size:30;not null" form:"name" validate:"max=30"json:"name"`
	// 分类类型，0代表菜品，1代表套餐;默认是菜品
	Type uint8 `gorm:"default:0;not null" form:"type" validate:"gte=0,lte=1" json:"type"`
	// 排序值，用于决定顺序；
	Sort uint `gorm:"default:0;not null" form:"sort" validate:"gte=0" json:"sort"`
	// 0代表禁用，1代表启用，默认禁用
	Status uint8 `gorm:"default:0;not null" form:"status" validate:"gte=0,lte=1" json:"status"`
	// 考虑到每一家店铺大概率有一个独立的分类，因此每一个category对应一家店铺
	RestaurantID uint `gorm:"index;not null" validate:"gte=0,lte=0" json:"restaurant_id"`
}

func GetCategoryByRestaurant(restaurant_id uint) ([]Category, error) {
	c := []Category{}
	err := tx.Find(&c, Category{RestaurantID: restaurant_id}).Error
	return c, err
}

func GetCategory(category_id uint) (*Category, error) {
	res := Category{Model: Model{ID: category_id}}
	err := tx.Find(&res).Error
	return &res, err
}

func CreateCategory(data *Category) error {
	err := tx.Create(data).Error
	return err
}

func UpdateCategory(category_id uint, data any) error {
	err := tx.Model(&Category{}).Where("id = ?", category_id).Updates(data).Error
	return err
}
