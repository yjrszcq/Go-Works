package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrCategoryDuplicateName = errors.New("分类名称重复")
	ErrCategoryNotFound      = gorm.ErrRecordNotFound
)

type Category struct {
	CategoryID  int64  `gorm:"primary_key;auto_increment;comment:分类唯一标识"`
	Name        string `gorm:"type:varchar(100);unique;not null;comment:分类名称"`
	Description string `gorm:"type:text;comment:分类描述"`
}

type CategoryDAO struct {
	db *gorm.DB
}

func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

func (dao *CategoryDAO) InsertCategory(ctx context.Context, c Category) error {
	err := dao.db.WithContext(ctx).Create(&c).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueIndexConflictsErrNo uint16 = 1062 // 唯一索引冲突报错号
		if mysqlErr.Number == uniqueIndexConflictsErrNo {
			return ErrCategoryDuplicateName
		}
	}
	return err
}

func (dao *CategoryDAO) FindCategoryByID(ctx context.Context, id int64) (Category, error) {
	var c Category
	err := dao.db.WithContext(ctx).Where("category_id = ?", id).First(&c).Error
	return c, err
}

func (dao *CategoryDAO) FindCategoryByName(ctx context.Context, name string) (Category, error) {
	var c Category
	err := dao.db.WithContext(ctx).Where("name = ?", name).First(&c).Error
	return c, err
}

func (dao *CategoryDAO) FindAllCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := dao.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (dao *CategoryDAO) UpdateCategory(ctx context.Context, c Category) error {
	err := dao.db.WithContext(ctx).Save(&c).Error
	return err
}

func (dao *CategoryDAO) DeleteCategory(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("category_id = ?", id).Delete(&Category{}).Error
	return err
}
