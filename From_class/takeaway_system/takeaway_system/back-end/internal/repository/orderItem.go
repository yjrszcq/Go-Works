package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrOrderItemForeignKeyDishConstraintFail = dao.ErrOrderItemForeignKeyDishConstraintFail
	ErrOrderItemNotFound                     = dao.ErrOrderItemNotFound
)

type OrderItemRepository struct {
	dao *dao.OrderItemDAO
}

func NewOrderItemRepository(dao *dao.OrderItemDAO) *OrderItemRepository {
	return &OrderItemRepository{
		dao: dao,
	}
}

func (r *OrderItemRepository) CreateOrderItem(ctx context.Context, o domain.OrderItem) error {
	return r.dao.InsertOrderItem(ctx, dao.OrderItem{
		OrderID:   o.OrderID,
		DishID:    o.DishID,
		Quantity:  o.Quantity,
		UnitPrice: o.UnitPrice,
	})
}

func (r *OrderItemRepository) FindOrderItemById(ctx context.Context, id int64) (domain.OrderItem, error) {
	o, err := r.dao.FindOrderItemByID(ctx, id)
	if err != nil {
		return domain.OrderItem{}, err
	}
	return domain.OrderItem{
		Id:           o.OrderItemID,
		OrderID:      o.OrderID,
		DishID:       o.DishID,
		Quantity:     o.Quantity,
		UnitPrice:    o.UnitPrice,
		TotalPrice:   o.TotalPrice,
		ReviewStatus: o.ReviewStatus,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}, nil
}

func (r *OrderItemRepository) FindOrderItemsByOrderId(ctx context.Context, orderId int64) ([]domain.OrderItem, error) {
	orderItems, err := r.dao.FindOrderItemsByOrderID(ctx, orderId)
	if err != nil {
		return nil, err
	}
	var domainOrderItems []domain.OrderItem
	for _, o := range orderItems {
		domainOrderItems = append(domainOrderItems, domain.OrderItem{
			Id:           o.OrderItemID,
			OrderID:      o.OrderID,
			DishID:       o.DishID,
			Quantity:     o.Quantity,
			UnitPrice:    o.UnitPrice,
			TotalPrice:   o.TotalPrice,
			ReviewStatus: o.ReviewStatus,
			CreatedAt:    o.CreatedAt,
			UpdatedAt:    o.UpdatedAt,
		})
	}
	return domainOrderItems, nil
}

func (r *OrderItemRepository) FindOrderItemsByDishId(ctx context.Context, dishId int64) ([]domain.OrderItem, error) {
	orderItems, err := r.dao.FindOrderItemsByDishID(ctx, dishId)
	if err != nil {
		return nil, err
	}
	var domainOrderItems []domain.OrderItem
	for _, o := range orderItems {
		domainOrderItems = append(domainOrderItems, domain.OrderItem{
			Id:           o.OrderItemID,
			OrderID:      o.OrderID,
			DishID:       o.DishID,
			Quantity:     o.Quantity,
			UnitPrice:    o.UnitPrice,
			TotalPrice:   o.TotalPrice,
			ReviewStatus: o.ReviewStatus,
			CreatedAt:    o.CreatedAt,
			UpdatedAt:    o.UpdatedAt,
		})
	}
	return domainOrderItems, nil
}

func (r *OrderItemRepository) FindOrderItemsByReviewStatus(ctx context.Context, status string) ([]domain.OrderItem, error) {
	orderItems, err := r.dao.FindOrderItemsByReviewStatus(ctx, status)
	if err != nil {
		return nil, err
	}
	var domainOrderItems []domain.OrderItem
	for _, o := range orderItems {
		domainOrderItems = append(domainOrderItems, domain.OrderItem{
			Id:           o.OrderItemID,
			OrderID:      o.OrderID,
			DishID:       o.DishID,
			Quantity:     o.Quantity,
			UnitPrice:    o.UnitPrice,
			TotalPrice:   o.TotalPrice,
			ReviewStatus: o.ReviewStatus,
			CreatedAt:    o.CreatedAt,
			UpdatedAt:    o.UpdatedAt,
		})
	}
	return domainOrderItems, nil
}

func (r *OrderItemRepository) FindAllOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	orderItems, err := r.dao.FindAllOrderItems(ctx)
	if err != nil {
		return nil, err
	}
	var domainOrderItems []domain.OrderItem
	for _, o := range orderItems {
		domainOrderItems = append(domainOrderItems, domain.OrderItem{
			Id:           o.OrderItemID,
			OrderID:      o.OrderID,
			DishID:       o.DishID,
			Quantity:     o.Quantity,
			UnitPrice:    o.UnitPrice,
			TotalPrice:   o.TotalPrice,
			ReviewStatus: o.ReviewStatus,
			CreatedAt:    o.CreatedAt,
			UpdatedAt:    o.UpdatedAt,
		})
	}
	return domainOrderItems, nil
}

func (r *OrderItemRepository) UpdateOrderItemReviewStatus(ctx context.Context, o domain.OrderItem) error {
	return r.dao.UpdateOrderItemReviewStatus(ctx, dao.OrderItem{
		OrderItemID:  o.Id,
		ReviewStatus: o.ReviewStatus,
	})
}

func (r *OrderItemRepository) DeleteOrderItem(ctx context.Context, id int64) error {
	return r.dao.DeleteOrderItem(ctx, id)
}

func (r *OrderItemRepository) DeleteOrderItemsByOrderId(ctx context.Context, orderId int64) error {
	return r.dao.DeleteOrderItemsByOrderID(ctx, orderId)
}
