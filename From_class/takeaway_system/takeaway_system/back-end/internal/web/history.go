package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderStatusHistoryHandler struct {
	svc *service.OrderStatusHistoryService
}

func NewOrderStatusHistoryHandler(svc *service.OrderStatusHistoryService) *OrderStatusHistoryHandler {
	return &OrderStatusHistoryHandler{
		svc: svc,
	}
}

func (o *OrderStatusHistoryHandler) ErrOutputForOrderStatusHistory(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrRecordNotFoundInOrderStatusHistory) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无订单状态历史记录"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInOrderStatusHistory) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrRecordNotFoundInOrder) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 订单不存在"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 系统错误"})
	}
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoryByID(ctx *gin.Context) {
	type FindOrderStatusHistoryReq struct {
		Id int64 `json:"id"`
	}
	var req FindOrderStatusHistoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	history, err := o.svc.FindOrderStatusHistoryByID(ctx, req.Id)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, history)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByOrderID(ctx *gin.Context) {
	type FindOrderStatusHistoriesByOrderIDReq struct {
		OrderId int64 `json:"order_id"`
	}
	var req FindOrderStatusHistoriesByOrderIDReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByOrderID(ctx, req.OrderId)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesAllByCustomer(ctx *gin.Context) {
	histories, err := o.svc.FindOrderStatusHistoriesAllByCustomer(ctx)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByStatusByCustomer(ctx *gin.Context) {
	type FindOrderStatusHistoriesByStatusReq struct {
		Status string `json:"status"`
	}
	var req FindOrderStatusHistoriesByStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByStatusByCustomer(ctx, req.Status)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByChangedByIDByCustomer(ctx *gin.Context) {
	type FindOrderStatusHistoriesByChangedByIDReq struct {
		ChangedById int64 `json:"changed_by_id"`
	}
	var req FindOrderStatusHistoriesByChangedByIDReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByChangedByIDByCustomer(ctx, req.ChangedById)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesAllByEmployee(ctx *gin.Context) {
	histories, err := o.svc.FindOrderStatusHistoriesAllByEmployee(ctx)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByStatusByEmployee(ctx *gin.Context) {
	type FindOrderStatusHistoriesByStatusReq struct {
		Status string `json:"status"`
	}
	var req FindOrderStatusHistoriesByStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByStatusByEmployee(ctx, req.Status)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByChangedByIDByEmployee(ctx *gin.Context) {
	type FindOrderStatusHistoriesByChangedByIDReq struct {
		ChangedById int64 `json:"changed_by_id"`
	}
	var req FindOrderStatusHistoriesByChangedByIDReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByChangedByIDByEmployee(ctx, req.ChangedById)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByOrderIDAndStatus(ctx *gin.Context) {
	type FindOrderStatusHistoriesByOrderIDAndStatusReq struct {
		OrderId int64  `json:"order_id"`
		Status  string `json:"status"`
	}
	var req FindOrderStatusHistoriesByOrderIDAndStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByOrderIDAndStatus(ctx, req.OrderId, req.Status)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}

func (o *OrderStatusHistoryHandler) FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx *gin.Context) {
	type FindOrderStatusHistoriesByOrderIDAndChangedByIDReq struct {
		OrderId     int64 `json:"order_id"`
		ChangedById int64 `json:"changed_by_id"`
	}
	var req FindOrderStatusHistoriesByOrderIDAndChangedByIDReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查找失败, JSON字段不匹配"})
		return
	}
	histories, err := o.svc.FindOrderStatusHistoriesByOrderIDAndChangedByID(ctx, req.OrderId, req.ChangedById)
	if err != nil {
		o.ErrOutputForOrderStatusHistory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, histories)
}