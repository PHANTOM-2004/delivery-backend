package mysql

import (
	"delivery-backend/service/merchant/biz/dal/model"
	"errors"

	"gorm.io/gorm"
)

func CreateMerchant(m *model.Merchant) (bool, error) {
	attrs := model.Merchant{
		MerchantName: m.MerchantName,
		Password:     m.Password,
		PhoneNumber:  m.PhoneNumber,
	}
	merchant := model.Merchant{}

	res := DB.Where(&model.Merchant{
		Account: m.Account,
	}).Attrs(attrs).FirstOrCreate(&merchant)

	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		// 该申请表已经创建过, 那么就不创建了
		return false, nil
	}

	return res.RowsAffected > 0, res.Error
}
