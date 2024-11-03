package models

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	MerchantName string
	PhoneNumber  string
	Account      string
	Password     string
	// 在告诉gorm默认值的时候gorm才知道默认值，否则这里会插入一个0
	Status                int8 `gorm:"default:1"`
	MerchantApplicationID int
	MerchantApplication   MerchantApplication
}

// 优先判断其他错误， 找不到时id返回为0,
func GetMerchantID(account string) (uint, error) {
	var m Merchant
	err := db.Model(&Merchant{}).Where("account = ?", account).First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return m.ID, err
}

func GetMerchant(account string) (*Merchant, error) {
	var m Merchant
	err := db.Model(&Merchant{}).Where("account = ?", account).First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func CreateMerchant(m *Merchant) error {
	err := db.Model(&Merchant{}).Create(m).Error
	return err
}

func GetMerchantByID(id uint) (*Merchant, error) {
	m := &Merchant{}
	err := db.Model(&Merchant{}).Where("id = ?", id).First(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return m, err
}

func EnableMerchant(id uint) error {
	err := db.Model(&Merchant{}).Where("id = ?", id).Update("status", 1).Error
	return err
}

// 禁用merchant账号
func DisableMerchant(id uint) error {
	err := db.Model(&Merchant{}).Where("id = ?", id).Update("status", 0).Error
	return err
}

func EditMerchant(id uint, data any) error {
	err := db.Model(&Merchant{}).Where("id = ?", id).Updates(data).Error
	return err
}

func CleanAllMerchants() error {
	log.Info("running merchant cleaning")
	defer log.Info("Deleted merchant have been cleaned")
	res := db.Unscoped().Where("deleted_at IS NOT NULL").Delete(&Merchant{})
	err := res.Error
	log.Infof("rows affected: [%d]", res.RowsAffected)
	return err
}

func DeleteMerchant(account string) (error, int64) {
	res := db.Where("account = ?", account).Delete(&Merchant{})
	return res.Error, res.RowsAffected
}
