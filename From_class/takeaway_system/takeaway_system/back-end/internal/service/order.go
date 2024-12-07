package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	ErrRecordIsEmptyInOrder                  = errors.New("列表为空")
	ErrRecordNotFoundInOrder                 = repository.ErrOrderNotFound
	ErrUserHasNoPermissionInOrder            = errors.New("无权限")
	ErrFormatForDeliveryLocationInOrder      = errors.New("送餐地址格式错误")
	ErrFormatForDeliveryTimeInOrder          = errors.New("送餐时间格式错误")
	ErrEmptyDeliveryLocationInOrder          = errors.New("未设置送餐地址")
	ErrDeleteInvalidOrderItemInOrder         = errors.New("无效订单项删除失败, 订单项删除失败")
	ErrDeleteInvalidOrderInOrder             = errors.New("无效订单删除失败")
	ErrCreateNewOrderInOrder                 = errors.New("创建新订单失败")
	ErrEmployeeHasNotOrderCanTake            = errors.New("暂无待接订单")
	ErrDeliverymanHasNotOrderCanTake         = errors.New("暂无待送订单")
	ErrDeliverymanHasNotDeliveringOrder      = errors.New("暂无送餐中订单")
	ErrDeliverymanHasNotDeliveredOrder       = errors.New("暂无已送达订单")
	ErrOrderStatusErrorInOrder               = errors.New("订单状态错误")
	ErrCanNotBindDeliverymanInOrder          = errors.New("无法绑定送餐员")
	ErrUpdateOrderItemStatusInOrder          = errors.New("更新订单项状态失败")
	ErrUpdateDeliveryTimeInOrder             = errors.New("更新送达时间失败")
	ErrCancelPaidOrderByCustomerInOrder      = errors.New("订单已支付，无法取消")
	ErrCancelCanceledOrderInOrder            = errors.New("订单不能重复取消")
	ErrCancelConfirmedOrderByEmployeeInOrder = errors.New("只有未备餐的订单可以取消")
	ErrUpdateOrderPaymentStatusInOrder       = errors.New("订单支付状态更新失败")
	ErrCreateHistoryInOrder                  = errors.New("创建订单状态历史记录失败")
	ErrDeleteInvalidHistoriesInOrder         = errors.New("无效历史记录删除失败")
	ErrOrderItemReviewedInOrder              = errors.New("含有已评价订单项的订单不可删除")
	ErrDeleteInvalidCartItemInOrder          = errors.New("无效购物车项删除失败")
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (svc *OrderService) CreateOrder(ctx *gin.Context, deliveryLocation string, deliveryTimeString string, cartItemIds []int64) error {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "无权限"})
		return ErrUserHasNoPermissionInOrder
	}
	var deliveryTime time.Time
	if deliveryTimeString == "" {
		deliveryTime = time.Now()
	} else {
		deliveryTime, err = time.Parse("2006-01-02 15:04:05", deliveryTimeString)
		if err != nil {
			return ErrFormatForDeliveryTimeInOrder
		}
		if deliveryTime.Before(time.Now()) {
			deliveryTime = time.Now()
		}
	}
	if len(cartItemIds) == 0 {
		return ErrRecordNotFoundInCartItem
	}
	if deliveryLocation == "" {
		customer, err := GlobalCustomer.FindCustomerById(ctx, id)
		if err != nil {
			return err
		}
		if customer.Address == "" {
			return ErrEmptyDeliveryLocationInOrder
		}
		deliveryLocation = customer.Address
	}
	cartItems := make([]domain.CartItem, 0)
	unitPrice := make([]float64, 0)
	totalAmount := 0.0
	for _, v := range cartItemIds {
		cartItem, err := GlobalCartItem.FindCartItemByID(ctx, v)
		if err != nil {
			if errors.Is(err, repository.ErrCartItemNotFound) {
				return ErrRecordNotFoundInCartItem
			} else {
				return err
			}
		}
		cartItems = append(cartItems, cartItem)
		dish, err := GlobalDish.FindDishById(ctx, cartItem.DishID)
		if err != nil {
			if errors.Is(err, repository.ErrDishNotFound) {
				return ErrRecordNotFoundInDish
			} else {
				return err
			}
		}
		unitPrice = append(unitPrice, dish.Price)
		totalAmount += float64(cartItem.Quantity) * dish.Price
	}
	orderId, err := svc.repo.CreateOrder(ctx, domain.Order{
		CustomerID:       id,
		DeliveryLocation: deliveryLocation,
		DeliveryTime:     deliveryTime,
		TotalAmount:      totalAmount,
	})
	if err != nil {
		return err
	}
	// 根据cartItems创建订单项
	for k, v := range cartItems {
		err = GlobalOrderItem.CreateOrderItem(ctx, domain.OrderItem{
			OrderID:   orderId,
			DishID:    v.DishID,
			Quantity:  v.Quantity,
			UnitPrice: unitPrice[k],
		})
		if err != nil {
			// 删除已经创建的订单项
			err = GlobalOrderItem.DeleteOrderItemsByOrderId(ctx, orderId)
			if err != nil {
				if !errors.Is(err, repository.ErrOrderItemNotFound) {
					return ErrDeleteInvalidOrderItemInOrder
				}
			}
			err = GlobalOrder.DeleteOrderById(ctx, orderId)
			if err != nil {
				return ErrDeleteInvalidOrderInOrder
			}
			return ErrCreateNewOrderInOrder
		}
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID: orderId,
		Status:  "未支付",
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) PayTheOrder(ctx *gin.Context, id int64) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.CustomerID != customerId {
		return ErrRecordNotFoundInOrder
	}
	if order.PaymentStatus != "待支付" {
		return ErrOrderStatusErrorInOrder
	}
	err = svc.repo.UpdateOrderPaymentStatus(ctx, domain.Order{
		Id:            id,
		PaymentStatus: "已支付",
	})
	if err != nil {
		return ErrUpdateOrderPaymentStatusInOrder
	}
	orderItems, err := GlobalOrderItem.FindOrderItemsByOrderId(ctx, order.Id)
	for _, v := range orderItems {
		err := GlobalCartItem.DeleteCartItemByCustomerIDAndDishID(ctx, customerId, v.DishID)
		if err != nil {
			return ErrDeleteInvalidCartItemInOrder
		}
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID: id,
		Status:  "确认中",
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) GetOrderById(ctx *gin.Context, id int64) (domain.Order, error) {
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return domain.Order{}, ErrRecordNotFoundInOrder
		} else {
			return domain.Order{}, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		if order.CustomerID != customerId {
			return domain.Order{}, ErrRecordNotFoundInOrder
		}
	}
	return order, nil
}

func (svc *OrderService) GetOrdersByCustomerId(ctx *gin.Context) ([]domain.Order, error) {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByCustomerId(ctx, customerId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrRecordIsEmptyInOrder
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrRecordIsEmptyInOrder
	}
	return orders, nil
}

func (svc *OrderService) GetOrdersByCustomerIdByEmployee(ctx *gin.Context, customerId int64) ([]domain.Order, error) {
	_, err := getCurrentEmployeeId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByCustomerId(ctx, customerId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrRecordIsEmptyInOrder
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrRecordIsEmptyInOrder
	}
	return orders, nil
}

func (svc *OrderService) GetOrdersByDeliveryPersonId(ctx *gin.Context, deliverymanId int64) ([]domain.Order, error) {
	_, err := getCurrentEmployeeId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByDeliveryPersonId(ctx, deliverymanId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrRecordIsEmptyInOrder
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrRecordIsEmptyInOrder
	}
	return orders, nil
}

func (svc *OrderService) GetOrdersByStatus(ctx *gin.Context, status string) ([]domain.Order, error) {
	orders, err := svc.repo.FindOrdersByStatus(ctx, status)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrRecordIsEmptyInOrder
		} else {
			return nil, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		customerOrder := make([]domain.Order, 0)
		for _, v := range orders {
			if v.CustomerID == customerId {
				customerOrder = append(customerOrder, v)
			}
		}
		if len(customerOrder) == 0 {
			return nil, ErrRecordIsEmptyInOrder
		}
		if customerOrder == nil {
			return nil, ErrRecordIsEmptyInOrder
		}
		return customerOrder, nil
	}
	if orders == nil {
		return nil, ErrRecordIsEmptyInOrder
	}
	return orders, nil
}

func (svc *OrderService) GetOrdersByPaymentStatus(ctx *gin.Context, paymentStatus string) ([]domain.Order, error) {
	orders, err := svc.repo.FindOrdersByPaymentStatus(ctx, paymentStatus)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrRecordIsEmptyInOrder
		} else {
			return nil, err
		}
	}
	customerId, err := getCurrentCustomerId(ctx)
	if err == nil {
		customerOrder := make([]domain.Order, 0)
		for _, v := range orders {
			if v.CustomerID == customerId {
				customerOrder = append(customerOrder, v)
			}
		}
		if len(customerOrder) == 0 {
			return nil, ErrRecordIsEmptyInOrder
		}
		if customerOrder == nil {
			return nil, ErrRecordIsEmptyInOrder
		}
		return customerOrder, nil
	}
	if orders == nil {
		return nil, ErrRecordIsEmptyInOrder
	}
	return orders, nil
}

func (svc *OrderService) EmployeeGetOrders(ctx *gin.Context) ([]domain.Order, error) {
	if sessions.Default(ctx).Get("role") != "employee" {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByPaymentStatus(ctx, "已支付")
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrEmployeeHasNotOrderCanTake
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrEmployeeHasNotOrderCanTake
	}
	return orders, nil
}

func (svc *OrderService) DeliverymanGetOrdersWaitingForDelivery(ctx *gin.Context) ([]domain.Order, error) {
	if sessions.Default(ctx).Get("role") != "deliveryman" {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByStatus(ctx, "待送餐")
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrDeliverymanHasNotOrderCanTake
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrDeliverymanHasNotOrderCanTake
	}
	return orders, nil
}

func (svc *OrderService) DeliverymanGetOrdersDelivering(ctx *gin.Context) ([]domain.Order, error) {
	if sessions.Default(ctx).Get("role") != "deliveryman" {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByStatus(ctx, "送餐中")
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrDeliverymanHasNotDeliveringOrder
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrDeliverymanHasNotDeliveringOrder
	}
	employeeId, _ := getCurrentEmployeeId(ctx)
	deliverymanOrders := make([]domain.Order, 0)
	for _, v := range orders {
		if v.DeliveryPersonID == employeeId {
			deliverymanOrders = append(deliverymanOrders, v)
		}
	}
	if len(deliverymanOrders) == 0 {
		return nil, ErrDeliverymanHasNotDeliveringOrder
	}
	return deliverymanOrders, nil
}

func (svc *OrderService) DeliverymanGetOrdersDelivered(ctx *gin.Context) ([]domain.Order, error) {
	if sessions.Default(ctx).Get("role") != "deliveryman" {
		return nil, ErrUserHasNoPermissionInOrder
	}
	orders, err := svc.repo.FindOrdersByStatus(ctx, "已送达")
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrDeliverymanHasNotDeliveredOrder
		} else {
			return nil, err
		}
	}
	if orders == nil {
		return nil, ErrDeliverymanHasNotDeliveredOrder
	}
	employeeId, _ := getCurrentEmployeeId(ctx)
	deliverymanOrders := make([]domain.Order, 0)
	for _, v := range orders {
		if v.DeliveryPersonID == employeeId {
			deliverymanOrders = append(deliverymanOrders, v)
		}
	}
	if len(deliverymanOrders) == 0 {
		return nil, ErrDeliverymanHasNotDeliveredOrder
	}
	return deliverymanOrders, nil
}

func (svc *OrderService) ConfirmTheOrder(ctx *gin.Context, id int64) error {
	if sessions.Default(ctx).Get("role") != "employee" {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.Status != "确认中" {
		return ErrOrderStatusErrorInOrder
	}
	err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
		Id:     id,
		Status: "备餐中",
	})
	if err != nil {
		return err
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID:     id,
		Status:      "备餐中",
		ChangedByID: sessions.Default(ctx).Get("id").(int64),
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) MealPreparationCompleted(ctx *gin.Context, id int64) error {
	if sessions.Default(ctx).Get("role") != "employee" {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.Status != "备餐中" {
		return ErrOrderStatusErrorInOrder
	}
	err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
		Id:     id,
		Status: "待送餐",
	})
	if err != nil {
		return err
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID:     id,
		Status:      "待送餐",
		ChangedByID: sessions.Default(ctx).Get("id").(int64),
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) DeliverTheFood(ctx *gin.Context, id int64) error {
	if sessions.Default(ctx).Get("role") != "deliveryman" {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.Status != "待送餐" {
		return ErrOrderStatusErrorInOrder
	}
	err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
		Id:     id,
		Status: "送餐中",
	})
	if err != nil {
		return err
	}
	employeeId, _ := getCurrentEmployeeId(ctx)
	err = svc.repo.UpdateOrderDeliveryPerson(ctx, domain.Order{
		Id:               id,
		DeliveryPersonID: employeeId,
	})
	if err != nil {
		return ErrCanNotBindDeliverymanInOrder
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID:     id,
		Status:      "送餐中",
		ChangedByID: sessions.Default(ctx).Get("id").(int64),
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) FoodDelivered(ctx *gin.Context, id int64) error {
	if sessions.Default(ctx).Get("role") != "deliveryman" {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.Status != "送餐中" {
		return ErrOrderStatusErrorInOrder
	}
	err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
		Id:     id,
		Status: "已送达",
	})
	if err != nil {
		return err
	}
	orderItems, err := GlobalOrderItem.FindOrderItemsByOrderId(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return ErrRecordNotFoundInOrderItem
		} else {
			return err
		}
	}
	for _, v := range orderItems {
		err = GlobalOrderItem.UpdateOrderItemReviewStatus(ctx, domain.OrderItem{
			Id:           v.Id,
			ReviewStatus: "未评价",
		})
		if err != nil {
			err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
				Id:     id,
				Status: "送餐中",
			})
			if err != nil {
				return err
			}
			return ErrUpdateOrderItemStatusInOrder
		}
	}
	err = svc.repo.UpdateOrderDeliveryTime(ctx, id)
	if err != nil {
		return ErrUpdateDeliveryTimeInOrder
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID:     id,
		Status:      "已送达",
		ChangedByID: sessions.Default(ctx).Get("id").(int64),
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) CancelTheOrder(ctx *gin.Context, id int64) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.CustomerID != customerId {
		return ErrUserHasNoPermissionInOrder
	}
	if order.PaymentStatus != "待支付" {
		return ErrCancelPaidOrderByCustomerInOrder
	}
	if order.Status == "已取消" {
		return ErrCancelCanceledOrderInOrder
	}
	err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
		Id:     id,
		Status: "已取消",
	})
	if err != nil {
		return ErrUpdateOrderPaymentStatusInOrder
	}
	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID: id,
		Status:  "已取消",
	})
	if err != nil {
		return ErrCreateHistoryInOrder
	}
	return nil
}

func (svc *OrderService) DeleteTheOrder(ctx *gin.Context, id int64) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.CustomerID != customerId {
		return ErrUserHasNoPermissionInOrder
	}
	if order.Status != "已送达" && order.Status != "已取消" {
		return ErrOrderStatusErrorInOrder
	}
	orderItems, err := GlobalOrderItem.FindOrderItemsByOrderId(ctx, id)
	for _, v := range orderItems {
		if v.ReviewStatus == "已评价" {
			return ErrOrderItemReviewedInOrder
		}
	}
	err = GlobalOrderItem.DeleteOrderItemsByOrderId(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return ErrRecordNotFoundInOrderItem
		} else {
			return ErrDeleteInvalidOrderItemInOrder
		}
	}
	err = GlobalOrderStatusHistory.DeleteOrderStatusHistoriesByOrderID(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrOrderStatusHistoryNotFound) {
			return ErrDeleteInvalidHistoriesInOrder
		}
	}
	err = svc.repo.DeleteOrderById(ctx, id)
	if err != nil {
		return ErrDeleteInvalidOrderInOrder
	}
	return nil
}

func (svc *OrderService) CancelTheOrderByEmployee(ctx *gin.Context, id int64) error {
	employeeId, err := getCurrentEmployeeId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInOrder
	}
	if sessions.Default(ctx).Get("role") == "deliveryman" {
		return ErrUserHasNoPermissionInOrder
	}
	order, err := svc.repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrRecordNotFoundInOrder
		} else {
			return err
		}
	}
	if order.Status != "确认中" {
		return ErrCancelConfirmedOrderByEmployeeInOrder
	}
	err = svc.repo.UpdateOrderStatus(ctx, domain.Order{
		Id:     id,
		Status: "已取消",
	})
	if err != nil {
		return ErrUpdateOrderPaymentStatusInOrder
	}

	err = GlobalOrderStatusHistory.CreateOrderStatusHistory(ctx, domain.OrderStatusHistory{
		OrderID:     id,
		Status:      "已取消",
		ChangedByID: employeeId,
	})
	return nil
}
