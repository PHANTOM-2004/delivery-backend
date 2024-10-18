package models

import (
	"delivery-backend/internal/setting"
	"fmt"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func SetUp() {
	defer log.Info("DB connection initialized")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)
  log.Info("initializing Database with Setting:")
	log.Info(dsn, setting.DatabaseSetting.Type)

	var err error

	db, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			// 设置默认的db table handler
			NamingStrategy: schema.NamingStrategy{
				// table name prefix, table for `User` would be `t_users`
				TablePrefix: setting.DatabaseSetting.TablePrefix,
				// use singular table name, table for `User` would be `user` with this option enabled
				SingularTable: true,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
}
