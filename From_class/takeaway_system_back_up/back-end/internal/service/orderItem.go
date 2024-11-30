package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ErrOrderItemForeignKeyDishConstraintFail = repository.ErrOrderItemForeignKeyDishConstraintFail
	ErrRecordNotFoundInOrderItem             = repository.ErrOrderItemNotFound
	ErrUserHasNoPermissionInOrderItem        = errors.New("无权限")
)

type OrderItemService struct {
	repo *repository.OrderItemRepository
}

func NewOrderItemService(repo *repository.OrderItemRepository) *OrderItemService {
	return &OrderItemService{
		repo: repo,
	}
}

func (svc *OrderItemService) GetOrderItemById(ctx *gin.Context, id int64) (domain.OrderItem, error) {
	orderItem, err := svc.repo.FindOrderItemById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return domain.OrderItem{}, ErrRecordNotFoundInOrderItem
		} else {
			return domain.OrderItem{}, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		order, err := GlobalOrder.FindOrderById(ctx, orderItem.OrderID)
		if err != nil {
			if errors.Is(err, repository.ErrOrderNotFound) {
				return domain.OrderItem{}, ErrRecordNotFoundInOrder
			} else {
				return domain.OrderItem{}, err
			}
		}
		if order.CustomerID != customerId {
			return domain.OrderItem{}, ErrRecordNotFoundInOrderItem
		}
	}
	return orderItem, nil
}

func (svc *OrderItemService) GetOrderItemsByOrderId(ctx *gin.Context, orderId int64) ([]domain.OrderItem, error) {
	id, err := getCurrentCustomerId(ctx)
	if err == nil {
		order, err := GlobalOrder.FindOrderById(ctx, orderId)
		if err != nil {
			if errors.Is(err, repository.ErrOrderNotFound) {
				return nil, ErrRecordNotFoundInOrder
			} else {
				return nil, err
			}
		}
		if order.CustomerID != id {
			return nil, ErrRecordNotFoundInOrderItem
		}
	}
	orderItems, err := svc.repo.FindOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return nil, ErrRecordNotFoundInOrderItem
		} else {
			return nil, err
		}
	}
	return orderItems, nil
}

func (svc *OrderItemService) GetOrderItemsByDishId(ctx *gin.Context, dishId int64) ([]domain.OrderItem, error) {
	if sessions.Default(ctx).Get("role") != "employee" {
		return nil, ErrUserHasNoPermissionInOrderItem
	}
	orderItems, err := svc.repo.FindOrderItemsByDishId(ctx, dishId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return nil, ErrRecordNotFoundInOrderItem
		} else {
			return nil, err
		}
	}
	return orderItems, nil
}

func (svc *OrderItemService) GetOrderItemsByDishIdByCustomer(ctx *gin.Context, dishId int64) ([]domain.OrderItem, error) {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrderItem
	}
	orders, err := GlobalOrder.FindOrdersByCustomerId(ctx, id)
	orderItems := make([]domain.OrderItem, 0)
	for _, order := range orders {
		tempOrderItems, err := svc.repo.FindOrderItemsByOrderId(ctx, order.Id)
		if err != nil {
			if !errors.Is(err, repository.ErrOrderItemNotFound) {
				return nil, err
			}
			continue
		}
		for _, temp := range tempOrderItems {
			if temp.DishID == dishId {
				orderItems = append(orderItems, temp)
			}
		}
	}
	if len(orderItems) == 0 {
		return nil, ErrRecordNotFoundInOrderItem
	}
	return orderItems, nil
}

func (svc *OrderItemService) GetOrderItemsByReviewStatus(ctx *gin.Context, status string) ([]domain.OrderItem, error) {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrderItem
	}
	orders, err := GlobalOrder.FindOrdersByCustomerId(ctx, id)
	orderItems := make([]domain.OrderItem, 0)
	for _, order := range orders {
		tempOrderItems, err := svc.repo.FindOrderItemsByOrderId(ctx, order.Id)
		if err != nil {
			if !errors.Is(err, repository.ErrOrderItemNotFound) {
				return nil, err
			}
			continue
		}
		for _, temp := range tempOrderItems {
			if temp.ReviewStatus == status {
				orderItems = append(orderItems, temp)
			}
		}
	}
	if len(orderItems) == 0 {
		return nil, ErrRecordNotFoundInOrderItem
	}
	return orderItems, nil
}

func (svc *OrderItemService) GetAllOrderItemsByCustomer(ctx *gin.Context) ([]domain.OrderItem, error) {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrderItem
	}
	orders, err := GlobalOrder.FindOrdersByCustomerId(ctx, id)
	orderItems := make([]domain.OrderItem, 0)
	for _, order := range orders {
		tempOrderItems, err := svc.repo.FindOrderItemsByOrderId(ctx, order.Id)
		if err != nil {
			if !errors.Is(err, repository.ErrOrderItemNotFound) {
				return nil, err
			}
			continue
		}
		orderItems = append(orderItems, tempOrderItems...)
	}
	if len(orderItems) == 0 {
		return nil, ErrRecordNotFoundInOrderItem
	}
	return orderItems, nil
}

func (svc *OrderItemService) GetAllOrderItemsByEmployee(ctx *gin.Context) ([]domain.OrderItem, error) {
	if sessions.Default(ctx).Get("role") != "employee" {
		return nil, ErrUserHasNoPermissionInOrderItem
	}
	orderItems, err := svc.repo.FindAllOrderItems(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return nil, ErrRecordNotFoundInOrderItem
		} else {
			return nil, err
		}
	}
	return orderItems, nil
}
