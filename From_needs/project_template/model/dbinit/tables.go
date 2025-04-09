package dbinit

import (
	"gorm.io/gorm"
	"project/model"
)

func InitTable(db *gorm.DB) error { // 自动创建表
	return db.AutoMigrate(
		&model.User{})
}
