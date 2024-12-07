package models

import "gorm.io/gorm"

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

type Rider struct {
	Model
	StudentName  string `gorm:"not null;size:30"`
	StudentNo    string `gorm:"not null;size:7"`
	StudentCard  string `gorm:"not null;size:100"`
	WechatUserID uint   `gorm:"not null"`
}

type RiderApplication struct {
	Model
	StudentName  string `gorm:"not null;size:30" json:"student_name"`
	StudentNo    string `gorm:"not null;size:7" json:"student_no"`
	StudentCard  string `gorm:"not null;size:100" json:"student_card"`
	WechatUserID uint   `gorm:"not null" json:"-"`
	Status       uint8  `gorm:"not null;default:1" json:"status"`
}

const (
	RiderApplicationToBeViewed  = 1
	RiderApplicationApproved    = 2
	RiderApplicationDisapproved = 3
)

func CreateRiderApplication(r *RiderApplication) error {
	err := tx.Create(r).Error
	return err
}

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

func ApproveRider(application_id uint) (bool, error) {
	succ := false
	err := tx.Transaction(func(ftx *gorm.DB) error {
		app := RiderApplication{}
		err := ftx.Find(&app, application_id).Error
		if err != nil {
			return err
		}
		if app.ID == 0 {
			return nil
		}
		if app.Status == RiderApplicationApproved {
			return nil
		}
		err = ftx.Model(&RiderApplication{}).UpdateColumn("status", RiderApplicationApproved).Error
		if err != nil {
			return err
		}
		// 考虑更改身份
		err = ftx.Model(&WechatUser{}).UpdateColumn("role", RoleRider).Error
		if err != nil {
			return err
		}
		// 创建信息
		rider := Rider{
			StudentName:  app.StudentName,
			StudentCard:  app.StudentCard,
			StudentNo:    app.StudentNo,
			WechatUserID: app.WechatUserID,
		}
		err = ftx.Create(&rider).Error
		if err != nil {
			return err
		}
		succ = true
		return nil
	})
	return succ, err
}

func DisapproveRider(application_id uint) (bool, error) {
	succ := false
	err := tx.Transaction(func(ftx *gorm.DB) error {
		app := RiderApplication{}
		err := ftx.Find(&app, application_id).Error
		if err != nil {
			return err
		}
		if app.ID == 0 {
			return nil
		}
		if app.Status != RiderApplicationToBeViewed {
			return nil
		}
		err = ftx.Model(&RiderApplication{}).UpdateColumn("status", RiderApplicationDisapproved).Error
		if err != nil {
			return err
		}
		succ = true
		return nil
	})
	return succ, err
}

func GetRiderApplications() ([]RiderApplication, error) {
	apps := []RiderApplication{}
	err := tx.Find(&apps).Error
	return apps, err
}
