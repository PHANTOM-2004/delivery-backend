package models

import (
	"errors"

	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type Admin struct {
	gorm.Model
	AdminName string
	Account   string
	Password  string
}

// TODO: 添加单元测试(cyt on 2024-10-19)
func ExistAdmin(account string) (bool, error) {
	var a Admin
	err := db.Model(&Admin{}).Where("account = ?", account).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return a.ID > 0, err
}

func GetAdmin(account string) (*Admin, error) {
	a := &Admin{}
	err := db.Model(&Admin{}).Where("account = ?", account).First(&a).Error
	return a, err
}

func EditAdmin(account string, data any) error {
	err := db.Model(&Admin{}).Where("account = ?", account).Updates(data).Error
	return err
}

func CleanAllAdmin() error {
	log.Info("running admin cleaning")
	defer log.Info("Deleted admins have been cleaned")
	res := db.Unscoped().Where("deleted_at IS NOT NULL").Delete(&Admin{})
	err := res.Error
	log.Infof("rows affected: [%d]", res.RowsAffected)
	return err
}

func DeleteAdmin(account string) error {
	err := db.Where("account = ?", account).Delete(&Admin{}).Error
	return err
}


