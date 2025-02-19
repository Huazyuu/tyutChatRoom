package jwt

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

var TOKEN string

func InitCore() {
	const ConfigFile = "../../conf/settings.yaml"
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

func TestGenToken(t *testing.T) {
	type args struct {
		user JwtPayLoad
	}
	test := struct {
		name string
		args args
	}{

		name: "GenToken_CASE",
		args: args{
			user: JwtPayLoad{
				Username: "testUser",
				UserID:   "f21a07ff81",
			},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		token, err := GenToken(test.args.user)
		TOKEN = token
		t.Log("token: ", token)
		t.Log("GenToken_CASE SUCCESS")
		if err != nil {
			t.Errorf("GenToken() error = %v", err)
			return
		}
	})

}
