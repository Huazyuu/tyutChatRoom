package global

import (
	"gin-gorilla/conf"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Config *conf.Config
	DB     *gorm.DB
	Log    *logrus.Logger
	Redis  *redis.Client
)
