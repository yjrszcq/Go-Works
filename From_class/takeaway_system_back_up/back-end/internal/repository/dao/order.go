package dao

import "time"

type Order struct {
	OrderID          int64     `gorm:"primary_key;auto_increment;comment:订单唯一标识"`
	CustomerID       int64     `gorm:"index;not null;comment:关联的顾客ID"`
	Customer         Customer  `gorm:"foreign_key:CustomerID;references:CustomerID"`
	OrderDate        time.Time `gorm:"type:datetime;not null;comment:订单创建时间"`
	DeliveryTime     time.Time `gorm:"type:datetime;not null;comment:预定的送餐时间"`
	DeliveryLocation string    `gorm:"type:varchar(255);not null;comment:送餐地点"`
	Status           string    `gorm:"type:enum('确认中', '备餐中', '送餐中', '已送达');default:'确认中';comment:订单当前状态"`
	PaymentStatus    string    `gorm:"type:enum('待支付', '已支付');default:'待支付';comment:支付状态"`
	TotalAmount      float64   `gorm:"type:decimal(10,2);not null;comment:订单总金额"`
	DeliveryPersonID int64     `gorm:"index;comment:关联的送餐员ID"`
	DeliveryPerson   Employee  `gorm:"foreign_key:DeliveryPersonID;references:EmployeeID"`
	CreatedAt        time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt        time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}
