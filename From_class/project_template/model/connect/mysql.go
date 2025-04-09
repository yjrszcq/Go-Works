package connect

import (
	"GoLandProjects/Works/From_class/haze_detection_system/lib"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(conf *lib.Config) *gorm.DB {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		conf.DbUser,
		conf.DbPassword,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
		conf.DbTimeout)
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		lib.ErrorLog.Fatal(err)
	}
	return DB
}
