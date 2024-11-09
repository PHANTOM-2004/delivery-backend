package models

import (
	log "github.com/sirupsen/logrus"
)

type Admin struct {
	Model
	AdminName string `gorm:"size:50;not null"`
	Account   string `gorm:"size:50;uniqueIndex;not null"`
	Password  string `gorm:"size:100;not null"`
}

// TODO: 添加单元测试(cyt on 2024-10-19)
// 建议在service层添加单元测试
func ExistAdmin(account string) (bool, error) {
	a := Admin{}
	err := tx.Find(&a, Admin{Account: account}).Error
	// 此时必定是存在的， 但可能发生其他错误
	// 因此对于调用者来说需要优先判断error, 如果error不是空
	// 那么是不可信的
	defer log.Tracef("called exist admin[%s]", account)
	return a.ID != 0, err
}

func CreateAdmin(a *Admin) error {
	err := tx.Create(a).Error
	log.Tracef("called create admin[%v]", *a)
	return err
}

func GetAdmin(account string) (*Admin, error) {
	a := Admin{}
	err := tx.Find(&a, Admin{Account: account}).Error
	defer log.Tracef("called get admin[%v]", a)
	return &a, err
}

// 默认情况下GORM 只会更新非零值的字段
func EditAdmin(id uint, data any) error {
	err := tx.Find(&Admin{Model: Model{ID: id}}).Updates(data).Error
	return err
}

func CleanAllAdmin() error {
	log.Info("running admin cleaning")
	defer log.Info("Deleted admins have been cleaned")
	res := tx.Unscoped().Where("deleted_at IS NOT NULL").Delete(&Admin{})
	err := res.Error
	log.Infof("rows affected: [%d]", res.RowsAffected)
	return err
}

func DeleteAdmin(account string) (error, int64) {
	res := tx.Where("account = ?", account).Delete(&Admin{})
	return res.Error, res.RowsAffected
}
