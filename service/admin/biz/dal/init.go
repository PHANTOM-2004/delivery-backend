package dal

import (
	"delivery-backend/service/admin/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
