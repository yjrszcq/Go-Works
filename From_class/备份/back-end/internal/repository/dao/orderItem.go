package dao

import "time"

type OrderItem struct {
	OrderItemID int64     `gorm:"primary_key;auto_increment;comment:订单项唯一标识"`
	OrderID     int64     `gorm:"index;not null;comment:关联的订单ID"`
	Order       Order     `gorm:"foreign_key:OrderID;references:OrderID"`
	DishID      int64     `gorm:"index;not null;comment:关联的菜品ID"`
	Dish        Dish      `gorm:"foreign_key:DishID;references:DishID"`
	Quantity    int       `gorm:"type:int;not null;comment:菜品数量"`
	UnitPrice   float64   `gorm:"type:decimal(10,2);not null;comment:菜品单价"`
	TotalPrice  float64   `gorm:"type:decimal(10,2);not null;comment:菜品总价"`
	CreatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}
