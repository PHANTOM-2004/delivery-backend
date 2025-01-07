package dal

import (
	"delivery-backend/service/user/biz/dal/mysql"
	"delivery-backend/service/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
