package models

import (
	"delivery-backend/internal/setting"

	"gorm.io/gorm"
)

type MerchantApplication struct {
	gorm.Model
	Status      int8
	Description string
	License     string
	Email       string
	PhoneNumber string
}

// 管理员申请表，创建失败的时候会返回error
func CreateMerchantApplication(a *MerchantApplication) error {
	err := db.Model(&MerchantApplication{}).Create(a).Error
	return err
}

func ApproveApplication(id int) error {
	// TODO:
	return nil
}

func DisapproveApplication(id int) error {
	// TODO
	return nil
}

// 获得所有的商家申请,注意需要分页查询
// 注意：page从1开始
func GetMerchantApplication(page_cnt int) ([]MerchantApplication, error) {
	page_size := setting.AppSetting.LicensePageSize
	offset := (page_cnt - 1) * page_size
	applications := []MerchantApplication{}
	err := db.Limit(page_size).Offset(offset).Find(&applications).Error
	return applications, err
}
