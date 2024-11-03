package models

import (
	"delivery-backend/internal/setting"
	"errors"

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

const (
	ApplicationDisapproved = 1
	ApplicationToBeViewed  = 2
	ApplicationApproved    = 3
)

// 管理员申请表，创建失败的时候会返回error
func CreateMerchantApplication(a *MerchantApplication) error {
	err := db.Model(&MerchantApplication{}).Create(a).Error
	return err
}

// 找到关联申请表的商家账号
func GetRelatedMerchant(application_id int) (*Merchant, error) {
	res := Merchant{}
	err := db.Model(&Merchant{}).
		Where("application_id = ?", application_id).
		First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &res, err
}

// 结果找不到时返回空指针
func GetMerchantApplication(id int) (*MerchantApplication, error) {
	res := MerchantApplication{}
	err := db.Model(&MerchantApplication{}).
		Where("id = ?", id).
		First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &res, err
}

// 从未通过到通过，后续需要创建商家账号
func ApproveApplication(id int) error {
	err := db.Model(&MerchantApplication{}).
		Where("id = ?", id).
		Update("status", ApplicationApproved).Error
	return err
}

// 从已通过到未通过，那么后续需要冻结对应的账号；
// 从未审核到未通过，后续不需要操作
func DisapproveApplication(id int) error {
	err := db.Model(&MerchantApplication{}).
		Where("id = ?", id).
		Update("status", ApplicationDisapproved).Error
	return err
}

// 获得所有的商家申请,注意需要分页查询
// 注意：page从1开始
func GetMerchantApplications(page_cnt int) ([]MerchantApplication, error) {
	page_size := setting.AppSetting.LicensePageSize
	offset := (page_cnt - 1) * page_size
	applications := []MerchantApplication{}
	err := db.Limit(page_size).Offset(offset).Find(&applications).Error
	return applications, err
}
