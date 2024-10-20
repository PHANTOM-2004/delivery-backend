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
// 建议在service层添加单元测试
func ExistAdmin(account string) (bool, error) {
	var a Admin
	err := db.Model(&Admin{}).Where("account = ?", account).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	// 此时必定是存在的， 但可能发生其他错误
	// 因此对于调用者来说需要优先判断error, 如果error不是空
	// 那么是不可信的
	return true, err
}

func CreateAdmin(a *Admin) error {
	err := db.Model(&Admin{}).Create(a).Error
	return err
}

func GetAdmin(account string) (*Admin, error) {
	a := &Admin{}
	err := db.Model(&Admin{}).Where("account = ?", account).First(a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return a, err
}

// 默认情况下GORM 只会更新非零值的字段
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

func DeleteAdmin(account string) (error, int64) {
	res := db.Where("account = ?", account).Delete(&Admin{})
	return res.Error, res.RowsAffected
}
