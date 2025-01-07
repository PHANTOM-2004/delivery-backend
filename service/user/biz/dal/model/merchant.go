package model

import "delivery-backend/common"

type Merchant struct {
	common.Model
	MerchantName string `gorm:"size:50;not null" json:"merchant_name"`
	PhoneNumber  string `gorm:"size:30;not null" json:"phone_number"`
	Account      string `gorm:"size:50;index;not null" json:"account"`
	Password     string `gorm:"size:100;not null" json:"-"`
	//  1启用，0禁用;默认启用
	Status uint8 `gorm:"default:1;not null" json:"status"`
	// 每个商家账号对应的申请表，这个申请表是唯一的
	MerchantApplicationID uint `gorm:"uniqueIndex" json:"merchant_application_id"`
	// 用于preload对应reference model
	MerchantApplication *MerchantApplication `json:"-"` // 关联的对应的MerchantApplication
}

const (
	MerchantAccountEnabled  = 1
	MerchantAccountDisabled = 0
)
