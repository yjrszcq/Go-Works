package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DishHandler struct {
	svc *service.DishService
}

func NewDishHandler(svc *service.DishService) *DishHandler {
	return &DishHandler{
		svc: svc,
	}
}

func (d *DishHandler) ErrOutputForDish(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrRecordNotFoundInDish) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无菜品"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInDish) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrFormatForDishNameInDish) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 菜品名称格式错误"})
	} else if errors.Is(err, service.ErrFormatForImageUrlInDish) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 图片链接格式错误"})
	} else if errors.Is(err, service.ErrRangeForPriceInDish) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 价格范围错误"})
	} else if errors.Is(err, service.ErrRecordNotFoundInCategory) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 分类不存在"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 系统错误"})
	}
}

func (d *DishHandler) CreateDish(ctx *gin.Context) {
	type CreateReq struct {
		Name       string  `json:"name"`
		ImageURL   string  `json:"image_url"`
		Price      float64 `json:"price"`
		CategoryID int64   `json:"category_id"`
	}
	var req CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "创建失败, JSON字段不匹配"})
		return
	}
	err := d.svc.CreateDish(ctx, req.Name, req.ImageURL, req.Price, req.CategoryID)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (d *DishHandler) GetDishById(ctx *gin.Context) {
	type FindReq struct {
		Id int64 `json:"id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	dish, err := d.svc.FindDishById(ctx, req.Id)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dish)
}

func (d *DishHandler) GetDishByName(ctx *gin.Context) {
	type FindReq struct {
		Name string `json:"name"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	dishes, err := d.svc.FindDishByName(ctx, req.Name)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dishes)
}

func (d *DishHandler) GetDishByCategory(ctx *gin.Context) {
	type FindReq struct {
		CategoryID int64 `json:"category_id"`
	}
	var req FindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	dishes, err := d.svc.FindDishByCategory(ctx, req.CategoryID)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dishes)
}

func (d *DishHandler) GetAllDishes(ctx *gin.Context) {
	dishes, err := d.svc.FindAllDishes(ctx)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dishes)
}

func (d *DishHandler) EditDish(ctx *gin.Context) {
	type EditReq struct {
		Id         int64   `json:"id"`
		Name       string  `json:"name"`
		ImageURL   string  `json:"image_url"`
		Price      float64 `json:"price"`
		CategoryID int64   `json:"category_id"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "编辑失败, JSON字段不匹配"})
		return
	}
	err := d.svc.EditDish(ctx, req.Id, req.Name, req.ImageURL, req.Price, req.CategoryID)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "编辑成功"})
}

func (d *DishHandler) DeleteDish(ctx *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id"`
	}
	var req DeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := d.svc.DeleteDish(ctx, req.Id)
	if err != nil {
		d.ErrOutputForDish(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
