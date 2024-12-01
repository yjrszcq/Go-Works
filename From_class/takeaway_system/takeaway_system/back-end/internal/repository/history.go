package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrOrderStatusHistoryNotFound = dao.ErrOrderStatusHistoryNotFound
)

type OrderStatusHistoryRepository struct {
	dao *dao.OrderStatusHistoryDAO
}

func NewOrderStatusHistoryRepository(dao *dao.OrderStatusHistoryDAO) *OrderStatusHistoryRepository {
	return &OrderStatusHistoryRepository{
		dao: dao,
	}
}

func (r *OrderStatusHistoryRepository) CreateOrderStatusHistory(ctx context.Context, o domain.OrderStatusHistory) error {
	return r.dao.InsertOrderStatusHistory(ctx, dao.OrderStatusHistory{
		OrderID:     o.OrderID,
		Status:      o.Status,
		ChangedByID: o.ChangedByID,
	})
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoryByID(ctx context.Context, id int64) (domain.OrderStatusHistory, error) {
	o, err := r.dao.FindOrderStatusHistoryByID(ctx, id)
	if err != nil {
		return domain.OrderStatusHistory{}, err
	}
	return domain.OrderStatusHistory{
		HistoryID:   o.HistoryID,
		OrderID:     o.OrderID,
		Status:      o.Status,
		ChangedAt:   o.ChangedAt,
		ChangedByID: o.ChangedByID,
	}, nil
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoriesByOrderID(ctx context.Context, orderId int64) ([]domain.OrderStatusHistory, error) {
	histories, err := r.dao.FindOrderStatusHistoriesByOrderID(ctx, orderId)
	if err != nil {
		return nil, err
	}
	var result []domain.OrderStatusHistory
	for _, h := range histories {
		result = append(result, domain.OrderStatusHistory{
			HistoryID:   h.HistoryID,
			OrderID:     h.OrderID,
			Status:      h.Status,
			ChangedAt:   h.ChangedAt,
			ChangedByID: h.ChangedByID,
		})
	}
	return result, nil
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoriesByStatus(ctx context.Context, status string) ([]domain.OrderStatusHistory, error) {
	histories, err := r.dao.FindOrderStatusHistoriesByStatus(ctx, status)
	if err != nil {
		return nil, err
	}
	var result []domain.OrderStatusHistory
	for _, h := range histories {
		result = append(result, domain.OrderStatusHistory{
			HistoryID:   h.HistoryID,
			OrderID:     h.OrderID,
			Status:      h.Status,
			ChangedAt:   h.ChangedAt,
			ChangedByID: h.ChangedByID,
		})
	}
	return result, nil
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoriesByChangedByID(ctx context.Context, changedById int64) ([]domain.OrderStatusHistory, error) {
	histories, err := r.dao.FindOrderStatusHistoriesByChangedByID(ctx, changedById)
	if err != nil {
		return nil, err
	}
	var result []domain.OrderStatusHistory
	for _, h := range histories {
		result = append(result, domain.OrderStatusHistory{
			HistoryID:   h.HistoryID,
			OrderID:     h.OrderID,
			Status:      h.Status,
			ChangedAt:   h.ChangedAt,
			ChangedByID: h.ChangedByID,
		})
	}
	return result, nil
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoriesByOrderIDAndStatus(ctx context.Context, orderId int64, status string) ([]domain.OrderStatusHistory, error) {
	histories, err := r.dao.FindOrderStatusHistoriesByOrderIDAndStatus(ctx, orderId, status)
	if err != nil {
		return nil, err
	}
	var result []domain.OrderStatusHistory
	for _, h := range histories {
		result = append(result, domain.OrderStatusHistory{
			HistoryID:   h.HistoryID,
			OrderID:     h.OrderID,
			Status:      h.Status,
			ChangedAt:   h.ChangedAt,
			ChangedByID: h.ChangedByID,
		})
	}
	return result, nil
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx context.Context, orderId int64, changedById int64) ([]domain.OrderStatusHistory, error) {
	histories, err := r.dao.FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx, orderId, changedById)
	if err != nil {
		return nil, err
	}
	var result []domain.OrderStatusHistory
	for _, h := range histories {
		result = append(result, domain.OrderStatusHistory{
			HistoryID:   h.HistoryID,
			OrderID:     h.OrderID,
			Status:      h.Status,
			ChangedAt:   h.ChangedAt,
			ChangedByID: h.ChangedByID,
		})
	}
	return result, nil
}

func (r *OrderStatusHistoryRepository) FindOrderStatusHistoriesAll(ctx context.Context) ([]domain.OrderStatusHistory, error) {
	histories, err := r.dao.FindOrderStatusHistoriesAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []domain.OrderStatusHistory
	for _, h := range histories {
		result = append(result, domain.OrderStatusHistory{
			HistoryID:   h.HistoryID,
			OrderID:     h.OrderID,
			Status:      h.Status,
			ChangedAt:   h.ChangedAt,
			ChangedByID: h.ChangedByID,
		})
	}
	return result, nil
}

func (r *OrderStatusHistoryRepository) DeleteOrderStatusHistoryByID(ctx context.Context, id int64) error {
	return r.dao.DeleteOrderStatusHistoryByID(ctx, id)
}

func (r *OrderStatusHistoryRepository) DeleteOrderStatusHistoriesByOrderID(ctx context.Context, orderId int64) error {
	return r.dao.DeleteOrderStatusHistoriesByOrderID(ctx, orderId)
}
