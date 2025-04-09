package dbinit

import (
	"GoLandProjects/Works/From_class/haze_detection_system/model"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error { // 自动创建表
	return db.AutoMigrate(
		&model.User{})
}
