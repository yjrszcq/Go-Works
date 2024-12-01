package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var (
	ErrReviewNotFound = gorm.ErrRecordNotFound
)

type Review struct {
	ReviewID   int64     `gorm:"primary_key;auto_increment;comment:评价唯一标识"`
	CustomerID int64     `gorm:"index;not null;comment:关联的顾客ID"`
	Customer   Customer  `gorm:"foreign_key:CustomerID;references:CustomerID"`
	DishID     int64     `gorm:"index;not null;comment:关联的菜品ID"`
	Dish       Dish      `gorm:"foreign_key:DishID;references:DishID"`
	Rating     int       `gorm:"type:int;not null;check:rating BETWEEN 1 AND 5;comment:评分 (1-5)"`
	Comment    string    `gorm:"type:text;comment:评价内容"`
	ReviewDate time.Time `gorm:"type:datetime;not null;comment:评价时间"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

type ReviewDAO struct {
	db *gorm.DB
}

func NewReviewDAO(db *gorm.DB) *ReviewDAO {
	return &ReviewDAO{db: db}
}

func (dao *ReviewDAO) InsertReview(ctx context.Context, r Review) error {
	r.ReviewDate = time.Now()
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	err := dao.db.WithContext(ctx).Create(&r).Error
	return err
}

func (dao *ReviewDAO) FindReviewByID(ctx context.Context, id int64) (Review, error) {
	var r Review
	err := dao.db.WithContext(ctx).Where("review_id = ?", id).First(&r).Error
	return r, err
}

func (dao *ReviewDAO) FindReviewsByCustomerID(ctx context.Context, customerId int64) ([]Review, error) {
	var reviews []Review
	err := dao.db.WithContext(ctx).Where("customer_id = ?", customerId).Find(&reviews).Error
	return reviews, err
}

func (dao *ReviewDAO) FindReviewsByDishID(ctx context.Context, dishId int64) ([]Review, error) {
	var reviews []Review
	err := dao.db.WithContext(ctx).Where("dish_id = ?", dishId).Find(&reviews).Error
	return reviews, err
}

func (dao *ReviewDAO) FindReviewsByRating(ctx context.Context, rating int) ([]Review, error) {
	var reviews []Review
	err := dao.db.WithContext(ctx).Where("rating = ?", rating).Find(&reviews).Error
	return reviews, err
}

func (dao *ReviewDAO) UpdateReview(ctx context.Context, r Review) error {
	r.UpdatedAt = time.Now()
	err := dao.db.WithContext(ctx).Model(&Review{}).Where("review_id = ?", r.ReviewID).Updates(Review{
		Rating:  r.Rating,
		Comment: r.Comment,
	}).Error
	return err
}

func (dao *ReviewDAO) DeleteReview(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("review_id = ?", id).Delete(&Review{}).Error
	return err
}
