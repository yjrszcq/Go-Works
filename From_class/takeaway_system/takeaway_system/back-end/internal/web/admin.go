package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{
		svc: svc,
	}
}

func (a *AdminHandler) ErrOutputForAdmin(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrPasswordIsWrongInAdmin) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "登录失败, 用户名或密码错误"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "登录失败, 系统错误"})
	}
}

func (a *AdminHandler) LogInAdmin(ctx *gin.Context) {
	type LogInReq struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	var req LogInReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "登录失败, JSON字段不匹配"})
		return
	}
	err := a.svc.LogInAdmin(ctx, req.Name, req.Password)
	if err != nil {
		a.ErrOutputForAdmin(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func (a *AdminHandler) LogOutAdmin(ctx *gin.Context) {
	err := a.svc.LogOutAdmin(ctx)
	if err != nil {
		a.ErrOutputForAdmin(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}