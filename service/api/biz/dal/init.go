package dal

import (
	"delivery-backend/service/api/biz/dal/redis"
)

func Init() {
	redis.Init()
}
