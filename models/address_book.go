package models

import "gorm.io/gorm"

type AddressBook struct {
	Model
	Default         bool   `gorm:"default:false;not null" json:"default"`
	Address         string `gorm:"size:80;not null" json:"address"`
	DetailedAddress string `gorm:"size:80" json:"detailed_address"`
	Name            string `gorm:"size:20;not null" json:"name"`
	Gender          uint8  `gorm:"not null" json:"gender"`
	PhoneNumber     string `gorm:"size:20;not null" json:"phone_number"`
	WechatUserID    uint   `gorm:"index;not null" json:"-"`
}

const (
	GenderMale   = 1
	GenderFemale = 2
)

func CreateAddressBook(a *AddressBook) error {
	err := tx.Create(a).Error
	return err
}

func DeleteAddressBook(addr_id uint) error {
	err := tx.Delete(&AddressBook{}, addr_id).Error
	return err
}

func GetAddressBooks(user_id uint) ([]AddressBook, error) {
	res := []AddressBook{}
	err := tx.Find(&res, AddressBook{WechatUserID: user_id}).Error
	return res, err
}

func SetDefaultAddressBook(user_id uint, addr_id uint) error {
	err := tx.Transaction(
		func(ftx *gorm.DB) error {
			res := []AddressBook{}
			err := ftx.Find(&res, AddressBook{WechatUserID: user_id}).Error
			if err != nil {
				return err
			}
			unset_id := 0

			for i := range res {
				if !res[i].Default {
					continue
				}
				unset_id = i
			}
			if unset_id == int(addr_id) {
				// 不需要设置
				return nil
			}
			// 取消默认
			err = ftx.Model(&AddressBook{}).
				Where("id = ?", unset_id).
				UpdateColumn("default", false).Error
			if err != nil {
				return err
			}
			// 设置默认
			err = ftx.Model(&AddressBook{}).
				Where("id = ?", addr_id).
				UpdateColumn("default", true).Error
			return err
		},
	)
	return err
}

// 不应当包含Default字段
func UpdateAddressBook(id uint, a *AddressBook) error {
	err := tx.Model(&AddressBook{}).Where("id = ?", id).Updates(*a).Error
	return err
}
