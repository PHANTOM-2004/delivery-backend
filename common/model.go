package common

import "gorm.io/plugin/soft_delete"

type Model struct {
	// 不使用uint64, 我们也用不到那么多数据
	ID        uint   `gorm:"primaryKey" json:"id"`
	CreatedAt uint64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt uint64 `gorm:"autoUpdateTime" json:"updated_at"`
	// 仿照gorm模型添加索引
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}
