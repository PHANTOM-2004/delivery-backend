package common

import "gorm.io/plugin/soft_delete"

type Model struct {
	// 不使用uint64, 我们也用不到那么多数据
	ID        uint32 `gorm:"primaryKey" json:"id"`
	CreatedAt uint64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt uint64 `gorm:"autoUpdateTime" json:"updated_at"`
	// 仿照gorm模型添加索引
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}

type MerchantApplication struct {
	Model
	// 在告诉gorm默认值的时候gorm才知道默认值，否则这里会插入一个0
	Status      uint8  `gorm:"default:2;not null" json:"status"` // 1:不通过审核，2:未审核，3:代表通过审核
	Description string `gorm:"size:300;not null" json:"description"`
	License     string `gorm:"size:200;not null" json:"license"`
	Email       string `gorm:"size:50;not null" json:"email"`
	PhoneNumber string `gorm:"size:30;not null" json:"phone_number"`
	Name        string `gorm:"size:20;not null" json:"name"`
	EmailStatus uint8  `gorm:"default:0;not null" json:"email_status"`
}

const (
	EmailNotSent   = 0
	EmailSent      = 1
	EmailSentError = 2
)

const (
	ApplicationDisapproved = 1
	ApplicationToBeViewed  = 2
	ApplicationApproved    = 3
)
