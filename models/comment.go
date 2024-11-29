package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Comment struct {
	Model
	Rating         uint8            `gorm:"not null;default:0" json:"rating"`
	Content        string           `gorm:"not null;size:300" json:"content"`
	OrderID        uint             `gorm:"not null" json:"-"`
	WechatUserID   uint             `gorm:"not null" json:"-"`
	WechatUser     WechatUser       `json:"wechat_user"`
	RestaurantID   uint             `gorm:"not null" json:"-"`
	CommentDetails []*CommentDetail `json:"comment_details"`
	Reply          []*Reply         `json:"replies"`
}

// 既然已经有明细，那么说明必然存在图片, 否则明细没有意义
type CommentDetail struct {
	Model
	ImagePath string `gorm:"size:100;not null" json:"image_path"`
	CommentID uint   `gorm:"not null" json:"-"`
}

// 回复表是评论下的所有回复, 他的id可能是商家
// reply属于某一个comment下
type Reply struct {
	Model
	CommentID    uint
	From         uint
	To           uint
	FromMerchant bool
	ToMerchant   bool
}

// func CreateComment
func CreateComment(c *Comment, path []string) error {
	err := tx.Transaction(
		func(ftx *gorm.DB) error {
			var err error
			err = ftx.Create(c).Error
			if err != nil {
				return err
			}
			if len(path) == 0 {
				log.Trace("comment has no pictures, skipped")
				return nil
			}

			// 创建评论明细
			comment_id := c.ID
			details := make([]CommentDetail, len(path))
			for i := range details {
				details[i].ImagePath = path[i]
				details[i].CommentID = comment_id
			}
			log.Trace("prepared comment details:")
			log.Trace(details)

			err = ftx.Create(&details).Error
			return err
		},
	)
	return err
}
