package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		svc: svc,
	}
}

func (r *ReviewHandler) ErrOutputForReview(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrRecordNotFoundInReview) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无评价"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInReview) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrCustomerDoNotHaveOrderItem) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 用户没有此订单项, 无法评价"})
	} else if errors.Is(err, service.ErrOrderItemReviewStatusInReview) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 订单项评价状态错误"})
	} else if errors.Is(err, service.ErrUpdateOrderItemStatusInReview) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 更新订单项评价状态失败"})
	} else if errors.Is(err, service.ErrRatingOutOfRangeInReview) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 评分应在1-5之间"})
	} else if errors.Is(err, service.ErrFormatForCommentInReview) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 评价应在200个字符以内"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 系统错误"})
	}
}

func (r *ReviewHandler) CreateReview(ctx *gin.Context) {
	type CreateReviewReq struct {
		OrderItemId int64  `json:"order_item_id"`
		Rating      int    `json:"rating"`
		Comment     string `json:"comment"`
	}
	var req CreateReviewReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "创建失败, JSON字段不匹配"})
		return
	}
	err := r.svc.CreateReview(ctx, req.OrderItemId, req.Rating, req.Comment)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (r *ReviewHandler) GetReviewById(ctx *gin.Context) {
	type GetReviewByIdReq struct {
		Id int64 `json:"id"`
	}
	var req GetReviewByIdReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	review, err := r.svc.FindReviewById(ctx, req.Id)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, review)
}

func (r *ReviewHandler) GetReviewsByDishId(ctx *gin.Context) {
	type GetReviewsByDishIdReq struct {
		DishId int64 `json:"dish_id"`
	}
	var req GetReviewsByDishIdReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	reviews, err := r.svc.FindReviewsByDishId(ctx, req.DishId)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, reviews)
}

func (r *ReviewHandler) GetReviewsByRating(ctx *gin.Context) {
	type GetReviewsByRatingReq struct {
		Rating int `json:"rating"`
	}
	var req GetReviewsByRatingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	reviews, err := r.svc.FindReviewsByRating(ctx, req.Rating)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, reviews)
}

func (r *ReviewHandler) GetReviewsByCustomerId(ctx *gin.Context) {
	reviews, err := r.svc.FindReviewsByCustomerId(ctx)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, reviews)
}

func (r *ReviewHandler) EditReview(ctx *gin.Context) {
	type UpdateReviewReq struct {
		Id      int64  `json:"id"`
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	var req UpdateReviewReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := r.svc.EditReview(ctx, req.Id, req.Rating, req.Comment)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (r *ReviewHandler) DeleteReview(ctx *gin.Context) {
	type DeleteReviewReq struct {
		Id int64 `json:"id"`
	}
	var req DeleteReviewReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := r.svc.DeleteReview(ctx, req.Id)
	if err != nil {
		r.ErrOutputForReview(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
