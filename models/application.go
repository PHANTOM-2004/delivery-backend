package models

import (
	"delivery-backend/internal/setting"

	log "github.com/sirupsen/logrus"
)

type MerchantApplication struct {
	Model
	// 在告诉gorm默认值的时候gorm才知道默认值，否则这里会插入一个0
	Status      uint8  `gorm:"default:2;not null" json:"status"` // 1:不通过审核，2:未审核，3:代表通过审核
	Description string `gorm:"size:300;not null" json:"description"`
	License     string `gorm:"size:200;not null" json:"license"`
	Email       string `gorm:"size:50;not null" json:"email"`
	PhoneNumber string `gorm:"size:30;not null" json:"phone_number"`
	Name        string `gorm:"size:20;not null" json:"name"`
	EmailStatus uint8  `gorm:"default:0;not null" json:"email_status"`
}

const (
	EmailNotSent   = 0
	EmailSent      = 1
	EmailSentError = 2
)

const (
	ApplicationDisapproved = 1
	ApplicationToBeViewed  = 2
	ApplicationApproved    = 3
)

func UpdateEmailStatus(application_id uint, s uint8) error {
	err := tx.Model(&MerchantApplication{}).
		Where("id", application_id).
		Update("email_status", s).Error
	return err
}

// 管理员申请表，创建失败的时候会返回error
func CreateMerchantApplication(a *MerchantApplication) error {
	// NOTE:这里如果加上model就会出错，非常奇怪。
	log.Tracef("creating merchant application [%v]", *a)
	err := tx.Create(a).Error
	return err
}

// 找到关联申请表的商家账号
func GetRelatedMerchant(application_id uint) (*Merchant, error) {
	res := Merchant{}
	err := tx.Find(&res, Merchant{MerchantApplicationID: application_id}).Error
	return &res, err
}

// 结果找不到时ID为0
func GetMerchantApplication(id uint) (*MerchantApplication, error) {
	res := MerchantApplication{}
	err := tx.Find(&res, MerchantApplication{Model: Model{ID: id}}).Error
	return &res, err
}

// 从未通过到通过，后续需要创建商家账号
func ApproveApplication(id int) error {
	err := tx.Model(&MerchantApplication{}).
		Where("id = ?", id).
		Update("status", ApplicationApproved).Error
	return err
}

// 从已通过到未通过，那么后续需要冻结对应的账号；
// 从未审核到未通过，后续不需要操作
func DisapproveApplication(id int) error {
	err := tx.Model(&MerchantApplication{}).
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
	err := tx.Limit(page_size).Offset(offset).Find(&applications).Error
	return applications, err
}
