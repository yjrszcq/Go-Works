package controller

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/domain"
	"project/lib"
	"project/service"
)

func ErrorOutputForUser(ctx *gin.Context, err error) {
	lib.Logger.Error.Println(err)
	if errors.Is(err, service.ErrorUserHasNoPermission) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrorUserIsUnavailable) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 用户不可用"})
	} else if errors.Is(err, service.ErrorUserDuplicateEmail) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "失败, 邮箱冲突"})
	} else if errors.Is(err, service.ErrorUserInvalidUserOrPassword) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "失败, 邮箱或密码错误"})
	} else if errors.Is(err, service.ErrorUserFormatOfName) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 姓名格式错误"})
	} else if errors.Is(err, service.ErrorUserFormatOfEmail) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 邮箱格式错误"})
	} else if errors.Is(err, service.ErrorUserFormatOfPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 密码包含至少一位数字，字母和特殊字符,且长度8-16"})
	} else if errors.Is(err, service.ErrorUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "失败, 用户不存在"})
	} else if errors.Is(err, service.ErrorUserPasswordIsWrong) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "失败, 密码错误"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 系统错误"})
	}

}

func SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "注册失败, JSON字段不匹配"})
		return
	}
	user, err := service.SignUp(ctx, &domain.User{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		ErrorOutputForUser(ctx, err)
		return
	}
	lib.Logger.Info.Printf("邮箱 %s 注册成功", req.Email)
	ctx.JSON(http.StatusOK, user)
}

func Login(ctx *gin.Context) {
	type LogInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogInReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "登录失败, JSON字段不匹配"})
		return
	}
	user, err := service.LogIn(ctx, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ErrorOutputForUser(ctx, err)
		return
	}
	lib.Logger.Info.Println("用户ID ", sessions.Default(ctx).Get("id"), " 登录成功")
	ctx.JSON(http.StatusOK, user) // 响应
}

func LogOut(ctx *gin.Context) {
	err := service.LogOut(ctx)
	if err != nil {
		ErrorOutputForUser(ctx, err)
		return
	}
	lib.Logger.Info.Println("登出成功")
	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func Profile(ctx *gin.Context) {
	user, err := service.Profile(ctx)
	if err != nil {
		ErrorOutputForUser(ctx, err)
		return
	}
	lib.Logger.Info.Println("用户ID ", sessions.Default(ctx).Get("id"), " 获取个人信息成功")
	ctx.JSON(http.StatusOK, user)
}

func ChangeProfile(ctx *gin.Context) {
	type EditReq struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := service.ChangeProfile(ctx, &domain.User{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		ErrorOutputForUser(ctx, err)
		return
	}
	lib.Logger.Info.Println("用户ID ", sessions.Default(ctx).Get("id"), " 修改成功")
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func ChangePassword(ctx *gin.Context) {
	type ChangePasswordReq struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req ChangePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := service.ChangePassword(ctx, &domain.User{
		Password:        req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		ErrorOutputForUser(ctx, err)
		return
	}
	lib.Logger.Info.Println("用户ID ", sessions.Default(ctx).Get("id"), " 修改密码成功")
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}
