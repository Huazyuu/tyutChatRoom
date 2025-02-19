package model

import (
	"fmt"
	"gin-gorilla/conf"
	"gin-gorilla/global"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

func InitCore() {
	const ConfigFile = "../conf/settings.yaml"
	c := &conf.Config{}
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error : %s", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("conf init Unmarshal: %v", err)
	}
	global.Config = c
	logrus.Info("conf init success")
}

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
	return db
}

func init() {
	InitCore()
	global.DB = InitGorm()
}

func TestUserModel(t *testing.T) {
	user := UserModel{
		UserName: "testUser",
		Password: "testPassword",
		Email:    "test2@example.com",
	}
	result := global.DB.Create(&user)
	if result.Error != nil {
		println("创建用户记录失败:", result.Error.Error())
	} else {
		println("创建用户记录成功，UserID:", user.UserID)
	}
}
func TestDBModel(t *testing.T) {
	chat := ChatModel{
		UserID:   "f21a07ff81",
		TargetID: "b818804dbf",
		Content:  "test",
		IP:       "127.0.0.1",
		Addr:     "内网地址",
		IsGroup:  false,
		MsgType:  0,
	}
	result := global.DB.Create(&chat)
	if result.Error != nil {
		println("创建对话失败:", result.Error.Error())
	} else {
		println("创建对话记录成功，UserID:", chat.Content)
	}
}
