package initialize

import (
	"douyin/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitORM(conn string) {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	// 1. 开启数据库连接
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: conn,
		}),
		&gorm.Config{
			Logger: ormLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 表明不加s
			},
		},
	)
	if err != nil {
		panic(err)
	}

	// 2. 连接池设置
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	repository.ORM = db
}
