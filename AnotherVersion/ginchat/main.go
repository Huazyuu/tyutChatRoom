package main

import (
	"fmt"
	"ginchat/conf"
	"ginchat/log"
	"ginchat/models"
	"ginchat/routes"
	"go.uber.org/zap"
)

func main() {
	conf.InitCore("")
	log.InitLog(conf.GlobalConf.App.Mode)
	models.InitDB()
	router := routes.InitRoute()
	addr := fmt.Sprintf("127.0.0.1:%s", conf.GlobalConf.App.Port)
	err := router.Run(addr)
	if err != nil {
		zap.L().Error("路由注册出错", zap.Error(err))
		return
	} else {
		zap.L().Info("监听:", zap.String("addr", addr))
	}

}
