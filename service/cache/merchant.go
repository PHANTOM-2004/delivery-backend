package cache

import (
	"delivery-backend/internal/gredis"
	"strconv"
)

// store in redis
type MerchantBlacklist struct {
	key string
}

func NewMerchantBlacklist(merchant_id uint) *MerchantBlacklist {
	key := "MERCH_BLACKLIST_" + strconv.Itoa(int(merchant_id))
	return &MerchantBlacklist{key: key}
}

func (m *MerchantBlacklist) Add() error {
	err := gredis.Set(m.key, "", 0)
	return err
}

func (m *MerchantBlacklist) Exists() (bool, error) {
	exist, err := gredis.Exists(m.key)
	return exist, err
}

func (m *MerchantBlacklist) Remove() error {
	err := gredis.Delete(m.key)
	return err
}
