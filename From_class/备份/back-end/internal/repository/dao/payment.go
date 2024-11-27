package dao

import "time"

type Payment struct {
	PaymentID   int64     `gorm:"primary_key;auto_increment;comment:支付唯一标识"`
	OrderID     int64     `gorm:"index;not null;comment:关联的订单ID"`
	Order       Order     `gorm:"foreign_key:OrderID;references:OrderID"`
	PaymentDate time.Time `gorm:"type:datetime;not null;comment:支付时间"`
	Amount      float64   `gorm:"type:decimal(10,2);not null;comment:支付金额"`
	Method      string    `gorm:"type:varchar(50);comment:支付方式"`
	Status      string    `gorm:"type:enum('待支付', '已完成', '失败');default:'待支付';comment:支付状态"`
	CreatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}
