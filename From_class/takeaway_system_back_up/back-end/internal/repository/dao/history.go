package dao

import "time"

type OrderStatusHistory struct {
	HistoryID   int64     `gorm:"primary_key;auto_increment;comment:状态历史唯一标识"`
	OrderID     int64     `gorm:"index;not null;comment:关联的订单ID"`
	Order       Order     `gorm:"foreign_key:OrderID;references:OrderID"`
	Status      string    `gorm:"type:enum('确认中', '备餐中', '送餐中', '已送达');not null;comment:订单状态"`
	ChangedAt   time.Time `gorm:"type:datetime;not null;comment:状态变更时间"`
	ChangedByID int64     `gorm:"index;comment:操作员工ID"`
	ChangedBy   Employee  `gorm:"foreign_key:ChangedByID;references:EmployeeID"`
}
