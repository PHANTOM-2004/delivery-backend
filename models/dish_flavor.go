package models

type Dish struct {
	Model
	Name        string `gorm:"size:30;not null" form:"name" validate:"max=30" json:"name"`
	Price       uint   `gorm:"default:0;not null" form:"price" json:"price"`
	Image       string `gorm:"size:256;not null" json:"image"`
	Description string `gorm:"size:50" form:"description" validate:"max=50" json:"description"`
	Sort        uint16 `gorm:"default:0;not null" form:"sort" json:"sort"`
	CategoryID  uint   `gorm:"not null" json:"category_id"`
}

func CreateDish(d *Dish) error {
	err := tx.Create(d).Error
	return err
}

func GetDish(dish_id uint) (*Dish, error) {
	d := Dish{}
	err := tx.Find(&d, dish_id).Error
	return &d, err
}

// 注意更新不存在的dish的情况
func UpdateDish(id uint, d Dish) error {
	err := tx.Model(&Dish{}).Where("id = ?", id).Updates(d).Error
	return err
}
