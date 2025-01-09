package model

import "delivery-backend/common"

type Admin struct {
	common.Model
	AdminName string `gorm:"size:50;not null"`
	Account   string `gorm:"size:50;uniqueIndex;not null"`
	Password  string `gorm:"size:100;not null"`
}
