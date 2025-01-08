package dal

import (
	"delivery-backend/service/merchant/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
