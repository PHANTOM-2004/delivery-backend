package mysql

import (
	"delivery-backend/service/user/biz/dal/model"
	"delivery-backend/service/user/conf"
	"fmt"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	defer klog.Info("User Service: Mysql Init Done")

	dsn := fmt.Sprintf(
		conf.GetConf().MySQL.DSN,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	)
	klog.Info("User Service: Mysql DSN: ", dsn)

	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	// init tables

	err = DB.AutoMigrate(
		&model.MerchantApplication{},
		&model.Merchant{},
	)
	if err != nil {
		panic(err)
	}
}
