package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrCartItemForeignKeyDishConstraintFail = errors.New("菜品不存在")
	ErrCartItemNotFound                     = gorm.ErrRecordNotFound
)

type CartItem struct {
	CartItemID int64     `gorm:"primary_key;auto_increment;comment:购物车项唯一标识"`
	CustomerID int64     `gorm:"index;not null;comment:关联的顾客ID"`
	Customer   Customer  `gorm:"foreign_key:CustomerID;references:CustomerID"`
	DishID     int64     `gorm:"index;not null;comment:关联的菜品ID"`
	Dish       Dish      `gorm:"foreign_key:DishID;references:DishID"`
	Quantity   int       `gorm:"not null;comment:数量"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

type CartItemDAO struct {
	db *gorm.DB
}

func NewCartItemDAO(db *gorm.DB) *CartItemDAO {
	return &CartItemDAO{db: db}
}

func (dao *CartItemDAO) InsertCartItem(ctx context.Context, c CartItem) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&c).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const foreignKeyConstraintFailErrNo uint16 = 1452
		if mysqlErr.Number == foreignKeyConstraintFailErrNo {
			return ErrCartItemForeignKeyDishConstraintFail
		}
	}
	return err
}

func (dao *CartItemDAO) FindCartItemByID(ctx context.Context, id int64) (CartItem, error) {
	var c CartItem
	err := dao.db.WithContext(ctx).Where("cart_item_id = ?", id).First(&c).Error
	return c, err
}

func (dao *CartItemDAO) FindCartItemByCustomerIDAndDishID(ctx context.Context, customerID, dishID int64) (CartItem, error) {
	var c CartItem
	err := dao.db.WithContext(ctx).Where("customer_id = ? AND dish_id = ?", customerID, dishID).First(&c).Error
	return c, err
}

func (dao *CartItemDAO) FindCartItemsByCustomerID(ctx context.Context, customerID int64) ([]CartItem, error) {
	var cartItems []CartItem
	err := dao.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&cartItems).Error
	return cartItems, err
}

func (dao *CartItemDAO) UpdateCartItemQuantity(ctx context.Context, c CartItem) error {
	if c.Quantity <= 0 {
		return dao.DeleteCartItem(ctx, c.CartItemID)
	}
	err := dao.db.WithContext(ctx).Model(&CartItem{}).Where("cart_item_id = ?", c.CartItemID).Updates(CartItem{
		Quantity:  c.Quantity,
		UpdatedAt: time.Now(),
	}).Error
	return err
}

func (dao *CartItemDAO) DeleteCartItem(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("cart_item_id = ?", id).Delete(&CartItem{}).Error
	return err
}

func (dao *CartItemDAO) DeleteCartItemsByCustomerID(ctx context.Context, customerID int64) error {
	err := dao.db.WithContext(ctx).Where("customer_id = ?", customerID).Delete(&CartItem{}).Error
	return err
}

func (dao *CartItemDAO) DeleteCartItemsByDishID(ctx context.Context, dishID int64) error {
	err := dao.db.WithContext(ctx).Where("dish_id = ?", dishID).Delete(&CartItem{}).Error
	return err
}

func (dao *CartItemDAO) DeleteCartItemByCustomerIDAndDishID(ctx context.Context, customerID, dishID int64) error {
	err := dao.db.WithContext(ctx).Where("customer_id = ? AND dish_id = ?", customerID, dishID).Delete(&CartItem{}).Error
	return err
}
