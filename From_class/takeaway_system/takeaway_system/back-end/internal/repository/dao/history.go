package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var (
	ErrOrderStatusHistoryNotFound = gorm.ErrRecordNotFound
)

type OrderStatusHistory struct {
	HistoryID   int64     `gorm:"primary_key;auto_increment;comment:状态历史唯一标识"`
	OrderID     int64     `gorm:"index;not null;comment:关联的订单ID"`
	Order       Order     `gorm:"foreign_key:OrderID;references:OrderID"`
	Status      string    `gorm:"type:enum('已取消', '未支付', '确认中', '备餐中', '待送餐', '送餐中', '已送达');not null;comment:订单状态"`
	ChangedAt   time.Time `gorm:"type:datetime;not null;comment:状态变更时间"`
	ChangedByID int64     `gorm:"index;default:null;comment:操作员工ID"`
	ChangedBy   Employee  `gorm:"foreign_key:ChangedByID;references:EmployeeID"`
}

type OrderStatusHistoryDAO struct {
	db *gorm.DB
}

func NewOrderStatusHistoryDAO(db *gorm.DB) *OrderStatusHistoryDAO {
	return &OrderStatusHistoryDAO{db: db}
}

func (dao *OrderStatusHistoryDAO) InsertOrderStatusHistory(ctx context.Context, o OrderStatusHistory) error {
	o.ChangedAt = time.Now()
	err := dao.db.WithContext(ctx).Create(&o).Error
	return err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoryByID(ctx context.Context, id int64) (OrderStatusHistory, error) {
	var o OrderStatusHistory
	err := dao.db.WithContext(ctx).Where("history_id = ?", id).First(&o).Error
	return o, err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoriesByOrderID(ctx context.Context, orderId int64) ([]OrderStatusHistory, error) {
	var histories []OrderStatusHistory
	err := dao.db.WithContext(ctx).Where("order_id = ?", orderId).Find(&histories).Error
	return histories, err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoriesByStatus(ctx context.Context, status string) ([]OrderStatusHistory, error) {
	var histories []OrderStatusHistory
	err := dao.db.WithContext(ctx).Where("status = ?", status).Find(&histories).Error
	return histories, err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoriesByChangedByID(ctx context.Context, changedById int64) ([]OrderStatusHistory, error) {
	var histories []OrderStatusHistory
	err := dao.db.WithContext(ctx).Where("changed_by_id = ?", changedById).Find(&histories).Error
	return histories, err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoriesByOrderIDAndStatus(ctx context.Context, orderId int64, status string) ([]OrderStatusHistory, error) {
	var histories []OrderStatusHistory
	err := dao.db.WithContext(ctx).Where("order_id = ? AND status = ?", orderId, status).Find(&histories).Error
	return histories, err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx context.Context, orderId int64, changedById int64) ([]OrderStatusHistory, error) {
	var histories []OrderStatusHistory
	err := dao.db.WithContext(ctx).Where("order_id = ? AND changed_by_id = ?", orderId, changedById).Find(&histories).Error
	return histories, err
}

func (dao *OrderStatusHistoryDAO) FindOrderStatusHistoriesAll(ctx context.Context) ([]OrderStatusHistory, error) {
	var histories []OrderStatusHistory
	err := dao.db.WithContext(ctx).Find(&histories).Error
	return histories, err
}

func (dao *OrderStatusHistoryDAO) DeleteOrderStatusHistoryByID(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("history_id = ?", id).Delete(&OrderStatusHistory{}).Error
	return err
}

func (dao *OrderStatusHistoryDAO) DeleteOrderStatusHistoriesByOrderID(ctx context.Context, orderId int64) error {
	err := dao.db.WithContext(ctx).Where("order_id = ?", orderId).Delete(&OrderStatusHistory{}).Error
	return err
}
