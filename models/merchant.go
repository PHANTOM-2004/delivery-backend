package models

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 在告诉gorm默认值的时候gorm才知道默认值，否则这里会插入一个0

type Merchant struct {
	Model
	MerchantName string `gorm:"size:50;not null" json:"merchant_name"`
	PhoneNumber  string `gorm:"size:30;not null" json:"phone_number"`
	Account      string `gorm:"size:50;index;not null" json:"account"`
	Password     string `gorm:"size:100;not null" json:"-"`
	//  1启用，0禁用;默认启用
	Status uint8 `gorm:"default:1;not null" json:"status"`
	// 每个商家账号对应的申请表，这个申请表是唯一的
	MerchantApplicationID uint `gorm:"uniqueIndex" json:"merchant_application_id"`
	// 用于preload对应reference model
	MerchantApplication *MerchantApplication `json:"-"` // 关联的对应的MerchantApplication
}

const (
	MerchantAccountEnabled  = 1
	MerchantAccountDisabled = 0
)

// reference: https://gorm.io/docs/preload.html#nested_preloading
func GetMerchantByCategory(category_id uint) (*Merchant, error) {
	c := Category{}
	err := tx.Preload("Restaurant.Merchant").Find(&c, category_id).Error
	return c.Restaurant.Merchant, err
}

// func GetMerchantByDish(dish_id uint) (*Merchant, error) {
// 	d := Dish{}
// 	err := tx.Preload("Category.Restaurant.Merchant").Find(&d, dish_id).Error
// 	return &d.Category.Restaurant.Merchant, err
// }

// 优先判断其他错误， 找不到时id返回为0,
func GetMerchantID(account string) (uint, error) {
	var m Merchant
	err := tx.Where("account = ?", account).First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return m.ID, err
}

// 不存在的时候返回nil
func GetMerchant(account string) (*Merchant, error) {
	var m Merchant
	err := tx.Where("account = ?", account).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

// bool: 返回是否created
func CreateMerchant(m *Merchant) (bool, error) {
	attrs := Merchant{
		MerchantName:          m.MerchantName,
		Password:              m.Password,
		PhoneNumber:           m.PhoneNumber,
		MerchantApplicationID: m.MerchantApplicationID,
	}
	merchant := Merchant{}

	res := tx.Where(Merchant{
		Account: m.Account,
	}).Attrs(attrs).FirstOrCreate(&merchant)

	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		// 该申请表已经创建过, 那么就不创建了
		return false, nil
	}

	return res.RowsAffected > 0, res.Error
}

func ExistMerchant(account string) (bool, error) {
	m := Merchant{}
	err := tx.Find(&m, Merchant{Account: account}).Error
	return m.ID != 0, err
}

func GetMerchantByID(id uint) (*Merchant, error) {
	m := Merchant{}
	err := tx.First(&m, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &m, err
}

// 从1开始
func GetMerchants(page_cnt int) ([]Merchant, error) {
	page_size := 10
	offset := max(page_cnt-1, 0) * page_size
	merchants := []Merchant{}
	err := tx.Limit(page_size).Offset(offset).Find(&merchants).Error
	return merchants, err
}

func EnableMerchant(id uint) error {
	err := tx.Model(&Merchant{}).Where("id = ?", id).Update("status", 1).Error
	return err
}

// 禁用merchant账号
func DisableMerchant(id uint) error {
	err := tx.Model(&Merchant{}).Where("id = ?", id).Update("status", 0).Error
	return err
}

func UpdateMerchant(id uint, data map[string]any) error {
	err := tx.Model(&Merchant{}).Where("id = ?", id).Updates(data).Error
	return err
}

func CleanAllMerchants() error {
	log.Info("running merchant cleaning")
	defer log.Info("Deleted merchant have been cleaned")
	res := tx.Unscoped().Where("deleted_at IS NOT NULL").Delete(&Merchant{})
	err := res.Error
	log.Infof("rows affected: [%d]", res.RowsAffected)
	return err
}

func DeleteMerchant(account string) (error, int64) {
	res := tx.Where("account = ?", account).Delete(&Merchant{})
	return res.Error, res.RowsAffected
}
