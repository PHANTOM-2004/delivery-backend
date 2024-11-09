package models

import "github.com/shopspring/decimal"

type Dish struct {
	Model
	Name        string          `gorm:"size:30;not null" json:"name"`
	Price       decimal.Decimal `gorm:"type:decimal(10,4);not null" json:"price"`
	Image       string          `gorm:"size:256;not null" json:"image"`
	Description string          `gorm:"size:50" json:"decription"` // 描述可以为空
	Sort        uint            `gorm:"default:0;not null" json:"sort"`
	CategoryID  uint            `gorm:"not null" json:"category_id"`
}
