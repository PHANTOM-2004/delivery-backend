package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"

	"github.com/gin-gonic/gin"
)

type RiderInfoReq struct {
	StuName string `json:"student_name" validate:"max=20"`
	StuNo   string `json:"student_no" validate:"max=7"`
	StuCard string `json:"student_card"`
}

func UploadRiderApplication(c *gin.Context) {
	req := RiderInfoReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	session := wechat.DefaultSession(c)
	info, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	err = models.CreateRiderApplication(&models.RiderApplication{
		WechatUserID: info.ID,
		StudentName:  req.StuName,
		StudentNo:    req.StuNo,
		StudentCard:  req.StuCard,
	})
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}
