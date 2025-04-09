package main

import (
	"project/lib"
	"project/model/connect"
	"project/model/dbinit"
	"project/router"
	"strconv"
)

func main() {
	cfg := lib.LoadConfig()
	lib.InitLog(cfg)
	lib.InitRegExp()
	if err := dbinit.InitTable(connect.InitDB(cfg)); err != nil {
		lib.ErrorLog.Fatal(err)
	}
	if err := router.SetupRoute(cfg).Run(":" + strconv.Itoa(cfg.HttpPort)); err != nil {
		lib.ErrorLog.Fatal("服务器启动失败: ", err)
	}
}
