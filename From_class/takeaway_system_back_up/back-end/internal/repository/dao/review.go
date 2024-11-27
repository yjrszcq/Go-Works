package dao

import "time"

type Review struct {
	ReviewID   int64     `gorm:"primary_key;auto_increment;comment:评价唯一标识"`
	OrderID    int64     `gorm:"index;not null;comment:关联的订单ID"`
	Order      Order     `gorm:"foreign_key:OrderID;references:OrderID"`
	CustomerID int64     `gorm:"index;not null;comment:关联的顾客ID"`
	Customer   Customer  `gorm:"foreign_key:CustomerID;references:CustomerID"`
	Rating     int       `gorm:"type:int;not null;check:rating BETWEEN 1 AND 5;comment:评分 (1-5)"`
	Comment    string    `gorm:"type:text;comment:评价内容"`
	ReviewDate time.Time `gorm:"type:datetime;not null;comment:评价时间"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}
