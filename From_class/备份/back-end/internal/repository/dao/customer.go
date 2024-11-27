package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

// User 直接对应数据库表结构

type Customer struct {
	CustomerID int64     `gorm:"primary_key;auto_increment;comment:顾客唯一标识"`
	Name       string    `gorm:"type:varchar(100);default:'新用户';not null;comment:顾客姓名"`
	Email      string    `gorm:"type:varchar(100);unique;not null;comment:顾客邮箱"`
	Password   string    `gorm:"type:varchar(255);not null;comment:顾客密码"`
	Phone      string    `gorm:"type:varchar(20);comment:顾客电话"`
	Address    string    `gorm:"type:varchar(255);comment:顾客地址"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

type CustomerDAO struct {
	db *gorm.DB
}

func NewCustomerDAO(db *gorm.DB) *CustomerDAO {
	return &CustomerDAO{
		db: db,
	}
}

func (dao *CustomerDAO) Insert(ctx context.Context, c Customer) error {
	//now := time.Now().UnixMilli() //返回 int64 毫秒数
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&c).Error
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

func (dao *CustomerDAO) FindByEmail(ctx context.Context, email string) (Customer, error) {
	var c Customer
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&c).Error // 只找第一条记录
	return c, err
}

func (dao *CustomerDAO) FindById(ctx context.Context, id int64) (Customer, error) {
	var c Customer
	err := dao.db.WithContext(ctx).Where("customer_id = ?", id).First(&c).Error // 只找第一条记录
	return c, err
}

func (dao *CustomerDAO) Update(ctx context.Context, c Customer) error {
	now := time.Now()
	c.UpdatedAt = now
	err := dao.db.WithContext(ctx).Save(&c).Error
	return err
}
