package models

import (
	"delivery-backend/internal/setting"
	"fmt"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/soft_delete"
)

type Model struct {
	// 不使用uint64, 我们也用不到那么多数据
	ID        uint   `gorm:"primaryKey" json:"id"`
	CreatedAt uint64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt uint64 `gorm:"autoUpdateTime" json:"updated_at"`
	// 仿照gorm模型添加索引
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}

// 用于复用的transaction
var tx *gorm.DB

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

	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			// 设置默认的db table handler
			NamingStrategy: schema.NamingStrategy{
				// table name prefix, table for `User` would be `t_users`
				TablePrefix: setting.DatabaseSetting.TablePrefix,
				// use singular table name, table for `User` would be `user` with this option enabled
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(setting.DatabaseSetting.GetLogLevel()),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin")

	err = db.AutoMigrate(
		&Admin{},
		&MerchantApplication{},
		&Merchant{},
		&Restaurant{},
		&RestaurantTime{},
		&Category{},
		&Flavor{},
		&Dish{},
	)

	log.Info("tables created")
	if err != nil {
		log.Fatal(err)
	}

	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)

	// NOTE:一定要注意如何reuse
	// https://gorm.io/zh_CN/docs/method_chaining.html
	tx = db.Session(&gorm.Session{})
}
