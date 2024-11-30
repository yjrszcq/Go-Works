package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var (
	ErrOrderNotFound = gorm.ErrRecordNotFound
)

type Order struct {
	OrderID          int64     `gorm:"primary_key;auto_increment;comment:订单唯一标识"`
	CustomerID       int64     `gorm:"index;not null;comment:关联的顾客ID"`
	Customer         Customer  `gorm:"foreign_key:CustomerID;references:CustomerID"`
	OrderDate        time.Time `gorm:"type:datetime;not null;comment:订单创建时间"`
	DeliveryTime     time.Time `gorm:"type:datetime;not null;comment:预定的送餐时间"`
	DeliveryLocation string    `gorm:"type:varchar(255);not null;comment:送餐地点"`
	Status           string    `gorm:"type:enum('已取消', '确认中', '备餐中', '待送餐', '送餐中', '已送达');default:'确认中';comment:订单当前状态"`
	PaymentStatus    string    `gorm:"type:enum('待支付', '已支付');default:'待支付';comment:支付状态"`
	TotalAmount      float64   `gorm:"type:decimal(10,2);not null;comment:订单总金额"`
	DeliveryPersonID int64     `gorm:"index;comment:关联的送餐员ID"`
	DeliveryPerson   Employee  `gorm:"foreign_key:DeliveryPersonID;references:EmployeeID"`
	CreatedAt        time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt        time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录更新时间"`
}

type OrderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) *OrderDAO {
	return &OrderDAO{db: db}
}

func (dao *OrderDAO) InsertOrder(ctx context.Context, o Order) (int64, error) {
	now := time.Now()
	o.OrderDate = now
	o.CreatedAt = now
	o.UpdatedAt = now
	if o.DeliveryTime.Before(now) {
		o.DeliveryTime = now
	}
	o.Status = "确认中"
	o.PaymentStatus = "待支付"
	err := dao.db.WithContext(ctx).Create(&o).Error
	if err != nil {
		return 0, err
	}
	return o.OrderID, err
}

func (dao *OrderDAO) FindOrderByID(ctx context.Context, id int64) (Order, error) {
	var o Order
	err := dao.db.WithContext(ctx).Where("order_id = ?", id).First(&o).Error
	return o, err
}

func (dao *OrderDAO) FindOrdersByCustomerID(ctx context.Context, customerID int64) ([]Order, error) {
	var orders []Order
	err := dao.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&orders).Error
	return orders, err
}

func (dao *OrderDAO) FindOrdersByDeliveryPersonID(ctx context.Context, deliveryPersonID int64) ([]Order, error) {
	var orders []Order
	err := dao.db.WithContext(ctx).Where("delivery_person_id = ?", deliveryPersonID).Find(&orders).Error
	return orders, err
}

func (dao *OrderDAO) FindOrdersByStatus(ctx context.Context, status string) ([]Order, error) {
	var orders []Order
	err := dao.db.WithContext(ctx).Where("status = ?", status).Find(&orders).Error
	return orders, err
}

func (dao *OrderDAO) FindOrdersByPaymentStatus(ctx context.Context, paymentStatus string) ([]Order, error) {
	var orders []Order
	err := dao.db.WithContext(ctx).Where("payment_status = ?", paymentStatus).Find(&orders).Error
	return orders, err
}

func (dao *OrderDAO) UpdateOrderStatus(ctx context.Context, o Order) error {
	err := dao.db.WithContext(ctx).Model(&Order{}).Where("order_id = ?", o.OrderID).Updates(Order{
		Status:    o.Status,
		UpdatedAt: time.Now(),
	}).Error
	return err
}

func (dao *OrderDAO) UpdateOrderDeliveryTime(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Model(&Order{}).Where("order_id = ?", id).Updates(Order{
		DeliveryTime: time.Now(),
	}).Error
	return err
}

func (dao *OrderDAO) UpdateOrderPaymentStatus(ctx context.Context, o Order) error {
	err := dao.db.WithContext(ctx).Model(&Order{}).Where("order_id = ?", o.OrderID).Updates(Order{
		PaymentStatus: o.PaymentStatus,
		UpdatedAt:     time.Now(),
	}).Error
	return err
}

func (dao *OrderDAO) UpdateOrderDeliveryPerson(ctx context.Context, o Order) error {
	err := dao.db.WithContext(ctx).Model(&Order{}).Where("order_id = ?", o.OrderID).Updates(Order{
		DeliveryPersonID: o.DeliveryPersonID,
		UpdatedAt:        time.Now(),
	}).Error
	return err
}

func (dao *OrderDAO) DeleteOrderById(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("order_id = ?", id).Delete(&Order{}).Error
	return err
}
