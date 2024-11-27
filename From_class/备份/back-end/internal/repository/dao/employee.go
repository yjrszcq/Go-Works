package dao

import "time"

type Employee struct {
	EmployeeID int64     `gorm:"primary_key;auto_increment;comment:员工唯一标识"`
	Name       string    `gorm:"type:varchar(100);not null;comment:员工姓名"`
	Role       string    `gorm:"type:enum('管理员', '员工', '送餐员');not null;comment:员工角色"`
	Email      string    `gorm:"type:varchar(100);unique;not null;comment:员工邮箱"`
	Phone      string    `gorm:"type:varchar(20);comment:员工电话"`
	Password   string    `gorm:"type:varchar(255);not null;comment:员工密码"`
	Status     string    `gorm:"type:enum('可用', '不可用');default:'可用';comment:员工状态"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}
