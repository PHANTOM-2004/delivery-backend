package models

type WechatUser struct {
	Model
	OpenID      string `gorm:"index;size:50" json:"-"`
	Role        uint8  `gorm:"not null;default:1"`
	PhoneNumber string `gorm:"size:20"`
}

//默认创建顾客
const (
	RoleCustomer = 1
	RoleRider    = 2
)

// 如果找不到这个wechat user， 那么就创建这个wechat user, 利用openid
// 因此必然返回一个对象， 除了error发生
func GetOrCreateWechatUser(openid string) (*WechatUser, error) {
	res := WechatUser{}
	err := tx.FirstOrCreate(&res, WechatUser{OpenID: openid}).Error
	return &res, err
}

