package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CartItemHandler struct {
	svc *service.CartItemService
}

func NewCartItemHandler(svc *service.CartItemService) *CartItemHandler {
	return &CartItemHandler{
		svc: svc,
	}
}

func (c *CartItemHandler) ErrOutputForCartItem(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrRecordIsEmptyInCartItem) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 购物车为空"})
	} else if errors.Is(err, service.ErrRecordNotFoundInCartItem) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "失败, 购物车项不存在"})
	} else if errors.Is(err, service.ErrDishInCartNotFoundInCartItem) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 菜品不存在"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInCartItem) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrFormatForQuantityInCartItem) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 数量应大于0, 小于等于99"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 系统错误"})
	}
}

func (c *CartItemHandler) AddCartItem(ctx *gin.Context) {
	type CreateCartItemReq struct {
		DishID   int64 `json:"dish_id"`
		Quantity int   `json:"quantity"`
	}
	var req CreateCartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "添加失败, JSON字段不匹配"})
		return
	}
	err := c.svc.AddCartItem(ctx, req.DishID, req.Quantity)
	if err != nil {
		c.ErrOutputForCartItem(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

func (c *CartItemHandler) GetCartItemsByCustomerId(ctx *gin.Context) {
	cartItems, err := c.svc.FindCartItemsByCustomerID(ctx)
	if err != nil {
		c.ErrOutputForCartItem(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, cartItems)
}

func (c *CartItemHandler) GetCartItemById(ctx *gin.Context) {
	type GetReq struct {
		CartItemID int64 `json:"id"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "获取失败, JSON字段不匹配"})
		return
	}
	cartItem, err := c.svc.FindCartItemByID(ctx, req.CartItemID)
	if err != nil {
		c.ErrOutputForCartItem(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, cartItem)
}

func (c *CartItemHandler) EditCartItem(ctx *gin.Context) {
	type UpdateCartItemReq struct {
		CartItemID int64 `json:"id"`
		Quantity   int   `json:"quantity"`
	}
	var req UpdateCartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := c.svc.UpdateCartItem(ctx, req.CartItemID, req.Quantity)
	if err != nil {
		c.ErrOutputForCartItem(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (c *CartItemHandler) DeleteCartItem(ctx *gin.Context) {
	type DeleteCartItemReq struct {
		CartItemID int64 `json:"id"`
	}
	var req DeleteCartItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := c.svc.DeleteCartItem(ctx, req.CartItemID)
	if err != nil {
		c.ErrOutputForCartItem(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (c *CartItemHandler) DeleteCartItemsByCustomerId(ctx *gin.Context) {
	err := c.svc.DeleteCartItemsByCustomerID(ctx)
	if err != nil {
		c.ErrOutputForCartItem(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
