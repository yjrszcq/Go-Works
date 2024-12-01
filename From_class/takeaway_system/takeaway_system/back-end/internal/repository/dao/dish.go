package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDishNotFound = gorm.ErrRecordNotFound
)

type Dish struct {
	DishID     int64     `gorm:"primary_key;auto_increment;comment:菜品唯一标识"`
	Name       string    `gorm:"type:varchar(100);not null;comment:菜品名称"`
	ImageURL   string    `gorm:"type:varchar(255);comment:菜品图片URL"`
	Price      float64   `gorm:"type:decimal(10,2);not null;comment:菜品价格"`
	CategoryID int64     `gorm:"index;comment:分类ID"`
	Category   Category  `gorm:"foreign_key:CategoryID;references:CategoryID"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

type DishDAO struct {
	db *gorm.DB
}

func NewDishDAO(db *gorm.DB) *DishDAO {
	return &DishDAO{db: db}
}

func (dao *DishDAO) InsertDish(ctx context.Context, d Dish) error {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&d).Error
	return err
}

func (dao *DishDAO) FindDishByID(ctx context.Context, id int64) (Dish, error) {
	var d Dish
	err := dao.db.WithContext(ctx).Where("dish_id = ?", id).First(&d).Error
	return d, err
}

func (dao *DishDAO) FindDishByName(ctx context.Context, name string) ([]Dish, error) {
	var dishes []Dish
	err := dao.db.WithContext(ctx).Where("name LIKE ?", "%"+name+"%").Find(&dishes).Error
	return dishes, err
}

func (dao *DishDAO) FindDishByCategoryID(ctx context.Context, categoryID int64) ([]Dish, error) {
	var dishes []Dish
	err := dao.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&dishes).Error
	return dishes, err
}

func (dao *DishDAO) FindAllDishes(ctx context.Context) ([]Dish, error) {
	var dishes []Dish
	err := dao.db.WithContext(ctx).Find(&dishes).Error
	return dishes, err
}

func (dao *DishDAO) UpdateDish(ctx context.Context, d Dish) error {
	now := time.Now()
	d.UpdatedAt = now
	err := dao.db.WithContext(ctx).Model(&Dish{}).Where("dish_id = ?", d.DishID).Updates(Dish{
		Name:       d.Name,
		ImageURL:   d.ImageURL,
		Price:      d.Price,
		CategoryID: d.CategoryID,
		UpdatedAt:  d.UpdatedAt,
	}).Error
	return err
}

func (dao *DishDAO) DeleteDish(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("dish_id = ?", id).Delete(&Dish{}).Error
	return err
}
