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

// 只允许通过已经没通过的以及未审核的, 也就是说除了已经通过的其他的都可以通过
func ApproveApplication(id int) (bool, error) {
	res := tx.Model(&MerchantApplication{}).
		Where("id = ? AND status <> ?", id, ApplicationApproved).
		Update("status", ApplicationApproved)
	return res.RowsAffected == 1, res.Error
}

// 只能不通过没有审核的
func DisapproveApplication(id int) (bool, error) {
	res := tx.Model(&MerchantApplication{}).
		Where("id = ? AND status = ?", id, ApplicationToBeViewed).
		Update("status", ApplicationDisapproved)
	return res.RowsAffected == 1, res.Error
}

// 获得所有的商家申请,注意需要分页查询
// 注意：page从1开始
func GetMerchantApplications(page_cnt int) ([]MerchantApplication, error) {
	page_size := setting.AppSetting.LicensePageSize
	offset := max(page_cnt-1, 0) * page_size
	applications := []MerchantApplication{}
	err := tx.Limit(page_size).Offset(offset).Find(&applications).Error
	return applications, err
}
