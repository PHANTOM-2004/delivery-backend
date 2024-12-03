package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
			log.Trace("got address book list", res)
			unset_id := 0

			for i := range res {
				if res[i].Default {
					unset_id = int(res[i].ID)
					break
				}
			}
			if unset_id == int(addr_id) {
				// 不需要设置
				log.Trace("no need to unset", unset_id)
				return nil
			}
			// 如果不存在unset那么就不需要unset
			log.Trace("address book: unset id")
			if unset_id != 0 {
				// 取消默认
				err = ftx.Model(&AddressBook{}).
					Where("id = ?", unset_id).
					UpdateColumn("default", false).Error
				if err != nil {
					return err
				}
				log.Trace("unset default addressbook", unset_id)
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
