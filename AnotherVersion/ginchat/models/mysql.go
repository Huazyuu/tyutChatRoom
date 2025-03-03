package models

import (
	"ginchat/conf"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var ChatDB *gorm.DB

func InitDB() {
	// 配置日志
	var dbLogger logger.Interface
	if conf.GlobalConf.App.Mode == "dev" {
		dbLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 写入标准输出
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      true,        // 彩色日志
			},
		)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	dsn := conf.GlobalConf.Mysql.Dsn
	zap.L().Info(dsn)
	var err error
	ChatDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		zap.L().Fatal("gorm连接出错", zap.Error(err))
		return
	} else {
		zap.L().Info("gorm init success")
	}
	return
}
