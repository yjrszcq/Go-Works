package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrOrderNotFound = dao.ErrOrderNotFound
)

type OrderRepository struct {
	dao *dao.OrderDAO
}

func NewOrderRepository(dao *dao.OrderDAO) *OrderRepository {
	return &OrderRepository{
		dao: dao,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, o domain.Order) (int64, error) {
	return r.dao.InsertOrder(ctx, dao.Order{
		CustomerID:       o.CustomerID,
		DeliveryLocation: o.DeliveryLocation,
		TotalAmount:      o.TotalAmount,
		DeliveryTime:     o.DeliveryTime,
	})
}

func (r *OrderRepository) FindOrderById(ctx context.Context, id int64) (domain.Order, error) {
	o, err := r.dao.FindOrderByID(ctx, id)
	if err != nil {
		return domain.Order{}, err
	}
	return domain.Order{
		Id:               o.OrderID,
		CustomerID:       o.CustomerID,
		DeliveryLocation: o.DeliveryLocation,
		Status:           o.Status,
		PaymentStatus:    o.PaymentStatus,
		TotalAmount:      o.TotalAmount,
		OrderDate:        o.OrderDate,
		DeliveryTime:     o.DeliveryTime,
		DeliveryPersonID: o.DeliveryPersonID,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}, nil
}

func (r *OrderRepository) FindOrdersByCustomerId(ctx context.Context, customerId int64) ([]domain.Order, error) {
	orders, err := r.dao.FindOrdersByCustomerID(ctx, customerId)
	if err != nil {
		return nil, err
	}
	var domainOrders []domain.Order
	for _, o := range orders {
		domainOrders = append(domainOrders, domain.Order{
			Id:               o.OrderID,
			CustomerID:       o.CustomerID,
			DeliveryLocation: o.DeliveryLocation,
			Status:           o.Status,
			PaymentStatus:    o.PaymentStatus,
			TotalAmount:      o.TotalAmount,
			OrderDate:        o.OrderDate,
			DeliveryTime:     o.DeliveryTime,
			DeliveryPersonID: o.DeliveryPersonID,
			CreatedAt:        o.CreatedAt,
			UpdatedAt:        o.UpdatedAt,
		})
	}
	return domainOrders, nil
}

func (r *OrderRepository) FindOrdersByDeliveryPersonId(ctx context.Context, deliveryPersonId int64) ([]domain.Order, error) {
	orders, err := r.dao.FindOrdersByDeliveryPersonID(ctx, deliveryPersonId)
	if err != nil {
		return nil, err
	}
	var domainOrders []domain.Order
	for _, o := range orders {
		domainOrders = append(domainOrders, domain.Order{
			Id:               o.OrderID,
			CustomerID:       o.CustomerID,
			DeliveryLocation: o.DeliveryLocation,
			Status:           o.Status,
			PaymentStatus:    o.PaymentStatus,
			TotalAmount:      o.TotalAmount,
			OrderDate:        o.OrderDate,
			DeliveryTime:     o.DeliveryTime,
			DeliveryPersonID: o.DeliveryPersonID,
			CreatedAt:        o.CreatedAt,
			UpdatedAt:        o.UpdatedAt,
		})
	}
	return domainOrders, nil
}

func (r *OrderRepository) FindOrdersByStatus(ctx context.Context, status string) ([]domain.Order, error) {
	orders, err := r.dao.FindOrdersByStatus(ctx, status)
	if err != nil {
		return nil, err
	}
	var domainOrders []domain.Order
	for _, o := range orders {
		domainOrders = append(domainOrders, domain.Order{
			Id:               o.OrderID,
			CustomerID:       o.CustomerID,
			DeliveryLocation: o.DeliveryLocation,
			Status:           o.Status,
			PaymentStatus:    o.PaymentStatus,
			TotalAmount:      o.TotalAmount,
			OrderDate:        o.OrderDate,
			DeliveryTime:     o.DeliveryTime,
			DeliveryPersonID: o.DeliveryPersonID,
			CreatedAt:        o.CreatedAt,
			UpdatedAt:        o.UpdatedAt,
		})
	}
	return domainOrders, nil
}

func (r *OrderRepository) FindOrdersByPaymentStatus(ctx context.Context, paymentStatus string) ([]domain.Order, error) {
	orders, err := r.dao.FindOrdersByPaymentStatus(ctx, paymentStatus)
	if err != nil {
		return nil, err
	}
	var domainOrders []domain.Order
	for _, o := range orders {
		domainOrders = append(domainOrders, domain.Order{
			Id:               o.OrderID,
			CustomerID:       o.CustomerID,
			DeliveryLocation: o.DeliveryLocation,
			Status:           o.Status,
			PaymentStatus:    o.PaymentStatus,
			TotalAmount:      o.TotalAmount,
			OrderDate:        o.OrderDate,
			DeliveryTime:     o.DeliveryTime,
			DeliveryPersonID: o.DeliveryPersonID,
			CreatedAt:        o.CreatedAt,
			UpdatedAt:        o.UpdatedAt,
		})
	}
	return domainOrders, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, o domain.Order) error {
	return r.dao.UpdateOrderStatus(ctx, dao.Order{
		OrderID: o.Id,
		Status:  o.Status,
	})
}

func (r *OrderRepository) UpdateOrderPaymentStatus(ctx context.Context, o domain.Order) error {
	return r.dao.UpdateOrderPaymentStatus(ctx, dao.Order{
		OrderID:       o.Id,
		PaymentStatus: o.PaymentStatus,
	})
}

func (r *OrderRepository) UpdateOrderDeliveryPerson(ctx context.Context, o domain.Order) error {
	return r.dao.UpdateOrderDeliveryPerson(ctx, dao.Order{
		OrderID:          o.Id,
		DeliveryPersonID: o.DeliveryPersonID,
	})
}

func (r *OrderRepository) UpdateOrderDeliveryTime(ctx context.Context, id int64) error {
	return r.dao.UpdateOrderDeliveryTime(ctx, id)
}

func (r *OrderRepository) DeleteOrderById(ctx context.Context, id int64) error {
	return r.dao.DeleteOrderById(ctx, id)
}
