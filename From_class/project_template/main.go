package main

import (
	"GoLandProjects/Works/From_class/haze_detection_system/lib"
	"GoLandProjects/Works/From_class/haze_detection_system/model/connect"
	"GoLandProjects/Works/From_class/haze_detection_system/model/dbinit"
	"GoLandProjects/Works/From_class/haze_detection_system/router"
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
