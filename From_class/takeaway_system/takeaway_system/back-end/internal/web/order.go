package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{
		svc: svc,
	}
}

func (o *OrderHandler) ErrOutputForOrder(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrRecordIsEmptyInOrder) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无订单"})
	} else if errors.Is(err, service.ErrRecordNotFoundInOrder) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "失败, 订单不存在"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInOrder) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrFormatForDeliveryLocationInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 送餐地址格式错误"})
	} else if errors.Is(err, service.ErrFormatForDeliveryTimeInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 送餐时间格式错误"})
	} else if errors.Is(err, service.ErrEmptyDeliveryLocationInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 未设置送餐地址"})
	} else if errors.Is(err, service.ErrDeleteInvalidOrderItemInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 无效订单项删除失败, 订单项删除失败"})
	} else if errors.Is(err, service.ErrDeleteInvalidOrderInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 无效订单删除失败"})
	} else if errors.Is(err, service.ErrCreateNewOrderInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 创建新订单失败"})
	} else if errors.Is(err, service.ErrEmployeeHasNotOrderCanTake) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无待接订单"})
	} else if errors.Is(err, service.ErrDeliverymanHasNotOrderCanTake) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无待送订单"})
	} else if errors.Is(err, service.ErrDeliverymanHasNotDeliveringOrder) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无送餐中订单"})
	} else if errors.Is(err, service.ErrDeliverymanHasNotDeliveredOrder) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无已送达订单"})
	} else if errors.Is(err, service.ErrOrderStatusErrorInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 订单状态错误"})
	} else if errors.Is(err, service.ErrCanNotBindDeliverymanInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 无法绑定送餐员"})
	} else if errors.Is(err, service.ErrRecordNotFoundInOrderItem) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 订单项不存在"})
	} else if errors.Is(err, service.ErrUpdateOrderItemStatusInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 订单项状态更新失败"})
	} else if errors.Is(err, service.ErrUpdateDeliveryTimeInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 送达时间更新失败"})
	} else if errors.Is(err, service.ErrCancelPaidOrderByCustomerInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 已支付订单无法取消"})
	} else if errors.Is(err, service.ErrCancelCanceledOrderInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 订单不能重复取消"})
	} else if errors.Is(err, service.ErrCancelConfirmedOrderByEmployeeInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 只有未备餐的订单可以取消"})
	} else if errors.Is(err, service.ErrUpdateOrderPaymentStatusInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 订单支付状态更新失败"})
	} else if errors.Is(err, service.ErrCreateHistoryInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "成功, 但是创建订单状态历史记录失败"})
	} else if errors.Is(err, service.ErrDeleteInvalidHistoriesInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 无效历史记录删除失败"})
	} else if errors.Is(err, service.ErrOrderItemReviewedInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 含有已评价订单项的订单不可删除"})
	} else if errors.Is(err, service.ErrDeleteInvalidCartItemInOrder) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 无效购物车项删除失败"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 系统错误"})
	}
}

func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	type CreateOrderReq struct {
		DeliveryLocation string  `json:"delivery_location"`
		DeliveryTime     string  `json:"delivery_time"` // 格式 2006-01-02 15:04:05
		CartItemId       []int64 `json:"cart_item_id"`
	}
	var req CreateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "创建失败, JSON字段不匹配"})
		return
	}
	order, err := o.svc.CreateOrder(ctx, req.DeliveryLocation, req.DeliveryTime, req.CartItemId)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) PayTheOrder(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "支付失败, JSON字段不匹配"})
		return
	}
	err := o.svc.PayTheOrder(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "支付成功"})
}

func (o *OrderHandler) GetOrderById(ctx *gin.Context) {
	type FindReq struct {
		Id int64 `json:"id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	order, err := o.svc.GetOrderById(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) GetOrdersByCustomerId(ctx *gin.Context) {
	orders, err := o.svc.GetOrdersByCustomerId(ctx)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) GetOrdersByCustomerIdByEmployee(ctx *gin.Context) {
	type GetReq struct {
		CustomerId int64 `json:"customer_id"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orders, err := o.svc.GetOrdersByCustomerIdByEmployee(ctx, req.CustomerId)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) GetOrdersByDeliveryPersonId(ctx *gin.Context) {
	type GetReq struct {
		DeliverymanId int64 `json:"deliveryman_id"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orders, err := o.svc.GetOrdersByDeliveryPersonId(ctx, req.DeliverymanId)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) GetOrdersByStatus(ctx *gin.Context) {
	type GetReq struct {
		Status string `json:"status"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orders, err := o.svc.GetOrdersByStatus(ctx, req.Status)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) GetOrdersByPaymentStatus(ctx *gin.Context) {
	type FindReq struct {
		PaymentStatus string `json:"payment_status"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orders, err := o.svc.GetOrdersByPaymentStatus(ctx, req.PaymentStatus)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) EmployeeGetOrders(ctx *gin.Context) {
	order, err := o.svc.EmployeeGetOrders(ctx)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) DeliverymanGetOrdersWaitingForDelivery(ctx *gin.Context) {
	order, err := o.svc.DeliverymanGetOrdersWaitingForDelivery(ctx)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) DeliverymanGetOrdersDelivering(ctx *gin.Context) {
	order, err := o.svc.DeliverymanGetOrdersDelivering(ctx)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) DeliverymanGetOrdersDelivered(ctx *gin.Context) {
	order, err := o.svc.DeliverymanGetOrdersDelivered(ctx)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) ConfirmTheOrder(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := o.svc.ConfirmTheOrder(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (o *OrderHandler) MealPreparationCompleted(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := o.svc.MealPreparationCompleted(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (o *OrderHandler) DeliverTheFood(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := o.svc.DeliverTheFood(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (o *OrderHandler) FoodDelivered(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := o.svc.FoodDelivered(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (o *OrderHandler) CancelTheOrder(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "取消失败, JSON字段不匹配"})
		return
	}
	err := o.svc.CancelTheOrder(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "订单取消成功"})
}

func (o *OrderHandler) DeleteTheOrder(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := o.svc.DeleteTheOrder(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "订单删除成功"})
}

func (o *OrderHandler) CancelTheOrderByEmployee(ctx *gin.Context) {
	type UpdateReq struct {
		Id int64 `json:"id"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "取消失败, JSON字段不匹配"})
		return
	}
	err := o.svc.CancelTheOrderByEmployee(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "订单取消成功"})
}
