package core

import (
	"fmt"
	"gin-gorilla/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitGorm() *gorm.DB {
	if global.Config.Mysql.Host == "" {
		logrus.Warnln("未配置mysql 取消gorm连接")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()

	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		mysqlLogger = logger.Default.LogMode(logger.Info) // 打印所有
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error) // 只打印错误的
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   mysqlLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalf(fmt.Sprintf("[%s] mysql 连接失败", dsn))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(global.Config.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.Config.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(global.Config.Mysql.MaxConnLifeTime))

	logrus.Info("mysql init success")
	global.Log.Info("db init success")
	return db
}
