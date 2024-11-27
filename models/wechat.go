package models

type WechatUser struct {
	Model
	OpenID      string `gorm:"index;size:50" json:"-"`
	Role        uint8  `gorm:"not null;default:1"`
	PhoneNumber string `gorm:"not null;size:20"`
}

const (
	RoleCustomer = 1
	RoleRider    = 2
)
