package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrOrderItemForeignKeyDishConstraintFail = errors.New("菜品不存在")
	ErrOrderItemNotFound                     = gorm.ErrRecordNotFound
)

type OrderItem struct {
	OrderItemID  int64     `gorm:"primary_key;auto_increment;comment:订单项唯一标识"`
	OrderID      int64     `gorm:"index;not null;comment:关联的订单ID"`
	Order        Order     `gorm:"foreign_key:OrderID;references:OrderID"`
	DishID       int64     `gorm:"index;not null;comment:关联的菜品ID"`
	Dish         Dish      `gorm:"foreign_key:DishID;references:DishID"`
	Quantity     int       `gorm:"type:int;not null;comment:菜品数量"`
	UnitPrice    float64   `gorm:"type:decimal(10,2);not null;comment:菜品单价"`
	TotalPrice   float64   `gorm:"type:decimal(10,2);not null;comment:菜品总价"`
	ReviewStatus string    `gorm:"type:enum('不可评价', '未评价', '已评价');default:'不可评价';comment:评价状态"`
	CreatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

type OrderItemDAO struct {
	db *gorm.DB
}

func NewOrderItemDAO(db *gorm.DB) *OrderItemDAO {
	return &OrderItemDAO{db: db}
}

func (dao *OrderItemDAO) InsertOrderItem(ctx context.Context, o OrderItem) error {
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	o.ReviewStatus = "不可评价"
	o.TotalPrice = float64(o.Quantity) * o.UnitPrice
	err := dao.db.WithContext(ctx).Create(&o).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const foreignKeyConstraintFailErrNo uint16 = 1452
		if mysqlErr.Number == foreignKeyConstraintFailErrNo {
			return ErrOrderItemForeignKeyDishConstraintFail
		}
	}
	return err
}

func (dao *OrderItemDAO) FindOrderItemByID(ctx context.Context, id int64) (OrderItem, error) {
	var o OrderItem
	err := dao.db.WithContext(ctx).Where("order_item_id = ?", id).First(&o).Error
	return o, err
}

func (dao *OrderItemDAO) FindOrderItemsByOrderID(ctx context.Context, orderID int64) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := dao.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&orderItems).Error
	return orderItems, err
}

func (dao *OrderItemDAO) FindOrderItemsByDishID(ctx context.Context, dishID int64) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := dao.db.WithContext(ctx).Where("dish_id = ?", dishID).Find(&orderItems).Error
	return orderItems, err
}

func (dao *OrderItemDAO) FindOrderItemsByReviewStatus(ctx context.Context, status string) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := dao.db.WithContext(ctx).Where("review_status = ?", status).Find(&orderItems).Error
	return orderItems, err
}

func (dao *OrderItemDAO) FindAllOrderItems(ctx context.Context) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := dao.db.WithContext(ctx).Find(&orderItems).Error
	return orderItems, err
}

func (dao *OrderItemDAO) UpdateOrderItemReviewStatus(ctx context.Context, o OrderItem) error {
	err := dao.db.WithContext(ctx).Model(&OrderItem{}).Where("order_item_id = ?", o.OrderItemID).Updates(OrderItem{
		ReviewStatus: o.ReviewStatus,
		UpdatedAt:    time.Now(),
	}).Error
	return err
}

func (dao *OrderItemDAO) DeleteOrderItem(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("order_item_id = ?", id).Delete(&OrderItem{}).Error
	return err
}

func (dao *OrderItemDAO) DeleteOrderItemsByOrderID(ctx context.Context, orderID int64) error {
	err := dao.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&OrderItem{}).Error
	return err
}
