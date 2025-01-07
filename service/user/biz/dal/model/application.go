package model

import "delivery-backend/common"

type MerchantApplication struct {
	common.Model
	// 在告诉gorm默认值的时候gorm才知道默认值，否则这里会插入一个0
	Status      uint8  `gorm:"default:2;not null" json:"status"` // 1:不通过审核，2:未审核，3:代表通过审核
	Description string `gorm:"size:300;not null" json:"description"`
	License     string `gorm:"size:200;not null" json:"license"`
	Email       string `gorm:"size:50;not null" json:"email"`
	PhoneNumber string `gorm:"size:30;not null" json:"phone_number"`
	Name        string `gorm:"size:20;not null" json:"name"`
	EmailStatus uint8  `gorm:"default:0;not null" json:"email_status"`
}
