package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc: svc,
	}
}

func (c *CategoryHandler) ErrOutputForCategory(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrRecordIsEmptyInCategory) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 暂无分类"})
	} else if errors.Is(err, service.ErrRecordNotFoundInCategory) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "失败, 分类不存在"})
	} else if errors.Is(err, service.ErrDuplicateNameInCategory) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "失败, 分类名称重复"})
	} else if errors.Is(err, service.ErrUserHasNoPermissionInCategory) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrFormatForNameInCategory) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 分类名称应小于20个字符"})
	} else if errors.Is(err, service.ErrFormatForDescriptionInCategory) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 分类描述应小于200个字符"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 系统错误"})
	}
}

func (c *CategoryHandler) CreateCategory(ctx *gin.Context) {
	type CreateReq struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var req CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "创建失败, JSON字段不匹配"})
		return
	}
	err := c.svc.CreateCategory(ctx, req.Name, req.Description)
	if err != nil {
		c.ErrOutputForCategory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (c *CategoryHandler) GetCategoryById(ctx *gin.Context) {
	type CreateReq struct {
		Id int64 `json:"id"`
	}
	var req CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	category, err := c.svc.FindCategoryByID(ctx, req.Id)
	if err != nil {
		c.ErrOutputForCategory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryHandler) GetCategoryByName(ctx *gin.Context) {
	type CreateReq struct {
		Name string `json:"name"`
	}
	var req CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	category, err := c.svc.FindCategoryByName(ctx, req.Name)
	if err != nil {
		c.ErrOutputForCategory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryHandler) GetAllCategories(ctx *gin.Context) {
	categories, err := c.svc.FindAllCategories(ctx)
	if err != nil {
		c.ErrOutputForCategory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

func (c *CategoryHandler) EditCategory(ctx *gin.Context) {
	type UpdateReq struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "更新失败, JSON字段不匹配"})
		return
	}
	err := c.svc.UpdateCategory(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		c.ErrOutputForCategory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (c *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id"`
	}
	var req DeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := c.svc.DeleteCategory(ctx, req.Id)
	if err != nil {
		c.ErrOutputForCategory(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
