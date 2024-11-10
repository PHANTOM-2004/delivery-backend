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
