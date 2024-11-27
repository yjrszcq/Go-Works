package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

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

type EmployeeDAO struct {
	db *gorm.DB
}

func NewEmployeeDAO(db *gorm.DB) *EmployeeDAO {
	return &EmployeeDAO{
		db: db,
	}
}

func (dao *EmployeeDAO) InsertEmployee(ctx context.Context, e Employee) error {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	e.Status = "可用"
	err := dao.db.WithContext(ctx).Create(&e).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueIndexConflictsErrNo uint16 = 1062 // 唯一索引冲突报错号
		if mysqlErr.Number == uniqueIndexConflictsErrNo {
			// 邮箱冲突（因为目前含唯一索引的就只有邮箱）
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *EmployeeDAO) FindEmployeeByEmail(ctx context.Context, email string) (Employee, error) {
	var e Employee
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&e).Error
	return e, err
}

func (dao *EmployeeDAO) FindEmployeeById(ctx context.Context, id int64) (Employee, error) {
	var e Employee
	err := dao.db.WithContext(ctx).Where("employee_id = ?", id).First(&e).Error
	return e, err
}

func (dao *EmployeeDAO) FindAllEmployees(ctx context.Context) ([]Employee, error) {
	var employees []Employee
	err := dao.db.WithContext(ctx).Find(&employees).Error
	return employees, err
}

func (dao *EmployeeDAO) UpdateEmployee(ctx context.Context, e Employee) error {
	now := time.Now()
	e.UpdatedAt = now
	err := dao.db.WithContext(ctx).Model(&Employee{}).Where("employee_id = ?", e.EmployeeID).Updates(Employee{
		Name:      e.Name,
		Role:      e.Role,
		Email:     e.Email,
		Phone:     e.Phone,
		Status:    e.Status,
		UpdatedAt: e.UpdatedAt,
	}).Error
	if err != nil {
		return err
	}
	if e.Password != "" {
		err = dao.db.WithContext(ctx).Model(&Employee{}).Where("employee_id = ?", e.EmployeeID).Updates(Employee{
			Password: e.Password,
		}).Error
	}
	return err
}

func (dao *EmployeeDAO) DeleteEmployee(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Delete(&Employee{}, id).Error
}
