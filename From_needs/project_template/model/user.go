package model

import (
	"GoLandProjects/Works/From_class/haze_detection_system/lib"
	"GoLandProjects/Works/From_class/haze_detection_system/model/connect"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrorUserDuplicateEmail = errors.New("邮箱冲突")
	ErrorUserNotFound       = gorm.ErrRecordNotFound
)

type User struct {
	UserId    string    `gorm:"type:varchar(50);primary_key;comment:用户唯一标识"`
	Name      string    `gorm:"type:varchar(100);default:'新用户';not null;comment:用户姓名"`
	Email     string    `gorm:"type:varchar(100);unique;not null;comment:用户邮箱"`
	Password  string    `gorm:"type:varchar(255);not null;comment:用户密码"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

func InsertUser(ctx context.Context, u *User) (string, error) {
	//now := time.Now().UnixMilli() //返回 int64 毫秒数
	now := time.Now()
	u.UserId = lib.GenerateUUID()
	u.CreatedAt = now
	u.UpdatedAt = now
	err := connect.DB.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueIndexConflictsErrNo uint16 = 1062 // 唯一索引冲突报错号
		if mysqlErr.Number == uniqueIndexConflictsErrNo {
			// 邮箱冲突（因为目前含唯一索引的就只有邮箱）
			return "", ErrorUserDuplicateEmail
		}
	}
	return u.UserId, nil
}

func SelectUser(ctx context.Context, u User) (User, error) {
	var user User
	var err error
	if u.UserId != "" { // 只找第一条记录
		err = connect.DB.WithContext(ctx).Where("user_id = ?", u.UserId).First(&user).Error
	} else if u.Email != "" {
		err = connect.DB.WithContext(ctx).Where("email = ?", u.Email).First(&user).Error
	}
	return user, err
}

func UpdateUser(ctx context.Context, u User) error {
	return connect.DB.WithContext(ctx).Model(&User{}).Where("user_id = ?", u.UserId).Updates(User{
		Name:      u.Name,
		Email:     u.Email,
		UpdatedAt: time.Now(),
	}).Error
}

func UpdateUserPassword(ctx context.Context, u User) error {
	return connect.DB.WithContext(ctx).Model(&User{}).Where("user_id = ?", u.UserId).Updates(User{
		Password:  u.Password,
		UpdatedAt: time.Now(),
	}).Error
}
