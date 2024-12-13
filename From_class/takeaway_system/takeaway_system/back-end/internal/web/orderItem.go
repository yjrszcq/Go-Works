package web

import (
	"back-end/internal/service"
	"back-end/internal/web/web_log"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderItemHandler struct {
	svc *service.OrderItemService
}

func NewOrderItemHandler(svc *service.OrderItemService) *OrderItemHandler {
	return &OrderItemHandler{
		svc: svc,
	}
}

func (o *OrderItemHandler) ErrOutputForOrderItem(ctx *gin.Context, err error) {
	web_log.WebLogger.ErrorLogger.Println(err)
	if errors.Is(err, service.ErrOrderItemForeignKeyDishConstraintFail) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 菜品不存在"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInOrderItem) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrRecordIsEmptyInOrderItem) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无订单项"})
	} else if errors.Is(err, service.ErrRecordNotFoundInOrderItem) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "失败, 订单项不存在"})
	} else if errors.Is(err, service.ErrRecordNotFoundInOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 订单不存在"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 系统错误"})
	}
}

func (o *OrderItemHandler) GetOrderItemById(ctx *gin.Context) {
	type FindReq struct {
		Id int64 `json:"id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orderItem, err := o.svc.GetOrderItemById(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("订单项ID %d 查找成功", req.Id)
	ctx.JSON(http.StatusOK, orderItem)
}

func (o *OrderItemHandler) GetOrderItemsByOrderId(ctx *gin.Context) {
	type FindReq struct {
		OrderId int64 `json:"order_id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orderItems, err := o.svc.GetOrderItemsByOrderId(ctx, req.OrderId)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("订单ID %d 的订单项 查找成功", req.OrderId)
	ctx.JSON(http.StatusOK, orderItems)
}

func (o *OrderItemHandler) GetOrderItemsByDishId(ctx *gin.Context) {

	type FindReq struct {
		DishId int64 `json:"dish_id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orderItems, err := o.svc.GetOrderItemsByDishId(ctx, req.DishId)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("菜品ID %d 的订单项 查找成功", req.DishId)
	ctx.JSON(http.StatusOK, orderItems)
}

func (o *OrderItemHandler) GetOrderItemsByReviewStatus(ctx *gin.Context) {
	type FindReq struct {
		Status string `json:"status"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orderItems, err := o.svc.GetOrderItemsByReviewStatus(ctx, req.Status)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("订单项状态 %s 的订单项 查找成功", req.Status)
	ctx.JSON(http.StatusOK, orderItems)
}

func (o *OrderItemHandler) GetAllOrderItemsByCustomer(ctx *gin.Context) {
	orderItems, err := o.svc.GetAllOrderItemsByCustomer(ctx)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("用户的订单项 查找成功")
	ctx.JSON(http.StatusOK, orderItems)
}

func (o *OrderItemHandler) GetAllOrderItemsByEmployee(ctx *gin.Context) {
	orderItems, err := o.svc.GetAllOrderItemsByEmployee(ctx)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("员工对 订单项 查找成功")
	ctx.JSON(http.StatusOK, orderItems)
}

func (o *OrderItemHandler) GetOrderItemsByDishIdByCustomer(ctx *gin.Context) {
	type FindReq struct {
		DishId int64 `json:"dish_id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	orderItems, err := o.svc.GetOrderItemsByDishIdByCustomer(ctx, req.DishId)
	if err != nil {
		o.ErrOutputForOrderItem(ctx, err)
		return
	}
	web_log.WebLogger.InfoLogger.Printf("菜品ID %d 的订单项 查找成功", req.DishId)
	ctx.JSON(http.StatusOK, orderItems)
}
