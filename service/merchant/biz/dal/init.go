package dal

import (
	"delivery-backend/service/merchant/biz/dal/mysql"
	"delivery-backend/service/merchant/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
