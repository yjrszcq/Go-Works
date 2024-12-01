package dao

import "gorm.io/gorm"

// 自动建表(如果表不存在才会建表，存在则不会覆盖，但会更新表结构)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&Customer{},
		&Employee{},
		&Category{},
		&Dish{},
		&CartItem{},
		&Order{},
		&OrderItem{},
		&Review{},
		&OrderStatusHistory{})
}
