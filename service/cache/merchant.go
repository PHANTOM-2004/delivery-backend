package cache

import (
	"delivery-backend/internal/gredis"
	"strconv"
)

type RedisMerchantBlacklist struct {
	key string
}

func NewRedisMerchantBlacklist(merchant_id uint) *RedisMerchantBlacklist {
	key := "MERCH_BLACKLIST_" + strconv.Itoa(int(merchant_id))
	return &RedisMerchantBlacklist{key: key}
}

func (m *RedisMerchantBlacklist) Add() error {
	err := gredis.Set(m.key, "", 0)
	return err
}

func (m *RedisMerchantBlacklist) Exists() (bool, error) {
	exist, err := gredis.Exists(m.key)
	return exist, err
}

func (m *RedisMerchantBlacklist) Remove() error {
	err := gredis.Delete(m.key)
	return err
}
