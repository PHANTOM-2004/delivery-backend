package dal

import (
	"delivery-backend/service/admin/biz/dal/mysql"
	"delivery-backend/service/admin/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
