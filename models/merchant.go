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

func DeleteMerchant(id uint) (error, int64) {
	res := db.Where("id = ?", id).Delete(&Merchant{})
	return res.Error, res.RowsAffected
}
