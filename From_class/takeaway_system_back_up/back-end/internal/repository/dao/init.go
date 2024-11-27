package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&Customer{},
		&Employee{},
		&Category{},
		&Dish{},
		&Order{},
		&OrderItem{},
		&Payment{},
		&Review{},
		&OrderStatusHistory{})
}
