package merchant_service

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
	handler "delivery-backend/service"
	"delivery-backend/service/cache"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 过滤商家黑名单
func MerchantBlacklistFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		merchant_id := handler.NewMerchInfoHanlder(c).GetID()
		in_blacklist, err := MerchantInBlacklist(merchant_id)
		if err != nil {
			app.ResponseInternalError(c, err)
			return
		}

		if in_blacklist {
			app.Response(c, http.StatusUnauthorized,
				ecode.ERROR_MERCHANT_ACCOUNT_BANNED, nil)
			c.Abort()
			return
		}

		log.Debug("pass: merchant account not in blacklist")
		c.Next()
	}
}

func EnableMerchant(merchant_id uint) error {
	b := cache.NewMerchantBlacklist(merchant_id)
	err := b.Remove()
	if err != nil {
		return err
	}

	// 写入数据库
	err = models.EnableMerchant(merchant_id)
	return err
}

func DisableMerchant(merchant_id uint) error {
	b := cache.NewMerchantBlacklist(merchant_id)
	err := b.Add()
	if err != nil {
		return err
	}

	// 然后设置数据库商家的状态为禁用状态
	err = models.DisableMerchant(merchant_id)
	return err
}

// 只从redis中查询状态
func MerchantInBlacklist(merchant_id uint) (bool, error) {
	b := cache.NewMerchantBlacklist(merchant_id)
	exist, err := b.Exists()
	return exist, err
}
