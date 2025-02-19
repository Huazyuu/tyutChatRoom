package core

import (
	"context"
	"gin-gorilla/global"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"time"
)

func InitRedis() *redis.Client {
	return InitRedisDB(global.Config.Redis.DB)
}
func InitRedisDB(db int) *redis.Client {
	redisConf := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.GetAddr(),
		Password: redisConf.Password,
		DB:       db,
		PoolSize: redisConf.PoolSize,
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		global.Log.Errorf("连接redis失败%s", redisConf.GetAddr())
		return nil
	}
	logrus.Info("redis init success")
	return rdb
}
