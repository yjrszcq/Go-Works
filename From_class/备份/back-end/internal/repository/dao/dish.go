package dao

import "time"

type Dish struct {
	DishID     int64     `gorm:"primary_key;auto_increment;comment:菜品唯一标识"`
	Name       string    `gorm:"type:varchar(100);not null;comment:菜品名称"`
	ImageURL   string    `gorm:"type:varchar(255);comment:菜品图片URL"`
	Price      float64   `gorm:"type:decimal(10,2);not null;comment:菜品价格"`
	CategoryID int64     `gorm:"index;comment:分类ID"`
	Category   Category  `gorm:"foreign_key:CategoryID;references:CategoryID"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}
