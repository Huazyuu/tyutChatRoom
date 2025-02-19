package main

import (
	"gin-gorilla/global"
	"gin-gorilla/global/core"
	"gin-gorilla/router"
)

func main() {
	core.InitCore()
	global.Log = core.InitLog()
	global.DB = core.InitGorm()
	global.Redis = core.InitRedis()
	r := router.InitRouter()
	global.Log.Infof("[gin-ws]  backend 运行在 http://%s:%d/api", global.Config.System.Host, global.Config.System.Port)
	err := r.Run(global.Config.System.Addr())
	if err != nil {
		global.Log.Error(err.Error())
	}
}
