package models

type WechatUser struct {
	Model
	OpenID          string `gorm:"not null;index;size:50" json:"-"`
	Role            uint8  `gorm:"not null;default:1" json:"role"`
	PhoneNumber     string `gorm:"size:20" json:"phone_number"`
	ProfileImageURL string `gorm:"size:200" json:"profile_image_url"`
	NickName        string `gorm:"size:50" json:"nickname"`
}

// 默认创建顾客
const (
	RoleCustomer = 1
	RoleRider    = 2
)

// 如果找不到这个wechat user， 那么就创建这个wechat user, 利用openid
// 因此必然返回一个对象， 除了error发生
// bool值代表是否创建了一个user, 用于指示是否是新用户
func GetOrCreateWechatUser(openid string) (*WechatUser, bool, error) {
	user := WechatUser{}
	res := tx.FirstOrCreate(&user, WechatUser{OpenID: openid})
	return &user, res.RowsAffected != 0, res.Error
}

func UpdateWechatUser(id uint, w *WechatUser) error {
	err := tx.Model(&WechatUser{}).Where("id = ?", id).Updates(*w).Error
	return err
}
