package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ErrRecordNotFoundInOrderStatusHistory      = repository.ErrOrderStatusHistoryNotFound
	ErrUserHasNoPermissionInOrderStatusHistory = errors.New("无权限")
)

type OrderStatusHistoryService struct {
	repo *repository.OrderStatusHistoryRepository
}

func NewOrderStatusHistoryService(repo *repository.OrderStatusHistoryRepository) *OrderStatusHistoryService {
	return &OrderStatusHistoryService{
		repo: repo,
	}
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoryByID(ctx *gin.Context, id int64) (domain.OrderStatusHistory, error) {
	history, err := svc.repo.FindOrderStatusHistoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return domain.OrderStatusHistory{}, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return domain.OrderStatusHistory{}, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		order, err := GlobalOrder.FindOrderById(ctx, history.OrderID)
		if err != nil {
			if errors.Is(err, repository.ErrOrderNotFound) {
				return domain.OrderStatusHistory{}, ErrRecordNotFoundInOrder
			} else {
				return domain.OrderStatusHistory{}, err
			}
		}
		if order.CustomerID != customerId {
			return domain.OrderStatusHistory{}, ErrRecordNotFoundInOrderStatusHistory
		}
	}
	return history, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByOrderID(ctx *gin.Context, orderId int64) ([]domain.OrderStatusHistory, error) {
	histories, err := svc.repo.FindOrderStatusHistoriesByOrderID(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return nil, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		order, err := GlobalOrder.FindOrderById(ctx, orderId)
		if err != nil {
			if errors.Is(err, repository.ErrOrderNotFound) {
				return nil, ErrRecordNotFoundInOrder
			} else {
				return nil, err
			}
		}
		if order.CustomerID != customerId {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		}
	}
	return histories, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesAllByCustomer(ctx *gin.Context) ([]domain.OrderStatusHistory, error) {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrderStatusHistory
	}
	orders, err := GlobalOrder.FindOrdersByCustomerId(ctx, customerId)
	var result []domain.OrderStatusHistory
	for _, order := range orders {
		histories, err := svc.repo.FindOrderStatusHistoriesByOrderID(ctx, order.Id)
		if err != nil {
			if !errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
				return nil, err
			}
			continue
		}
		result = append(result, histories...)
	}
	return result, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByStatusByCustomer(ctx *gin.Context, status string) ([]domain.OrderStatusHistory, error) {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrderStatusHistory
	}
	orders, err := GlobalOrder.FindOrdersByCustomerId(ctx, customerId)
	var result []domain.OrderStatusHistory
	for _, order := range orders {
		histories, err := svc.repo.FindOrderStatusHistoriesByOrderIDAndStatus(ctx, order.Id, status)
		if err != nil {
			if !errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
				return nil, err
			}
			continue
		}
		result = append(result, histories...)
	}
	return result, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByChangedByIDByCustomer(ctx *gin.Context, changedById int64) ([]domain.OrderStatusHistory, error) {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrderStatusHistory
	}
	orders, err := GlobalOrder.FindOrdersByCustomerId(ctx, customerId)
	var result []domain.OrderStatusHistory
	for _, order := range orders {
		histories, err := svc.repo.FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx, order.Id, changedById)
		if err != nil {
			if !errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
				return nil, err
			}
			continue
		}
		result = append(result, histories...)
	}
	return result, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesAllByEmployee(ctx *gin.Context) ([]domain.OrderStatusHistory, error) {
	role := sessions.Default(ctx).Get("role")
	if role != "employee" && role != "admin" {
		return nil, ErrUserHasNoPermissionInOrderStatusHistory
	}
	histories, err := svc.repo.FindOrderStatusHistoriesAll(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return nil, err
		}
	}
	return histories, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByStatusByEmployee(ctx *gin.Context, status string) ([]domain.OrderStatusHistory, error) {
	role := sessions.Default(ctx).Get("role")
	if role != "employee" && role != "admin" {
		return nil, ErrUserHasNoPermissionInOrderStatusHistory
	}
	histories, err := svc.repo.FindOrderStatusHistoriesByStatus(ctx, status)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return nil, err
		}
	}
	return histories, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByChangedByIDByEmployee(ctx *gin.Context, changedById int64) ([]domain.OrderStatusHistory, error) {
	role := sessions.Default(ctx).Get("role")
	if role != "employee" && role != "admin" {
		return nil, ErrUserHasNoPermissionInOrderStatusHistory
	}
	histories, err := svc.repo.FindOrderStatusHistoriesByChangedByID(ctx, changedById)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return nil, err
		}
	}
	return histories, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByOrderIDAndStatus(ctx *gin.Context, orderId int64, status string) ([]domain.OrderStatusHistory, error) {
	histories, err := svc.repo.FindOrderStatusHistoriesByOrderIDAndStatus(ctx, orderId, status)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return nil, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		order, err := GlobalOrder.FindOrderById(ctx, orderId)
		if err != nil {
			if errors.Is(err, repository.ErrOrderNotFound) {
				return nil, ErrRecordNotFoundInOrder
			} else {
				return nil, err
			}
		}
		if order.CustomerID != customerId {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		}
	}
	return histories, nil
}

func (svc *OrderStatusHistoryService) FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx *gin.Context, orderId int64, changedById int64) ([]domain.OrderStatusHistory, error) {
	histories, err := svc.repo.FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx, orderId, changedById)
	if err != nil {
		if errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		} else {
			return nil, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		order, err := GlobalOrder.FindOrderById(ctx, orderId)
		if err != nil {
			if errors.Is(err, repository.ErrOrderNotFound) {
				return nil, ErrRecordNotFoundInOrder
			} else {
				return nil, err
			}
		}
		if order.CustomerID != customerId {
			return nil, ErrRecordNotFoundInOrderStatusHistory
		}
	}
	return histories, nil
}
