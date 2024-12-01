package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler { // 预编译正则表达式
	return &CustomerHandler{
		svc: svc,
	}
}

func (c *CustomerHandler) ErrOutputForCustomer(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrUserHasNoPermissionInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrUserDuplicateEmailInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 邮箱冲突"})
	} else if errors.Is(err, service.ErrInvalidUserOrPasswordInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 邮箱或密码错误"})
	} else if errors.Is(err, service.ErrFormatForNameInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 姓名格式错误"})
	} else if errors.Is(err, service.ErrFormatForEmailInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 邮箱格式错误"})
	} else if errors.Is(err, service.ErrFormatForPasswordInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 密码包含至少一位数字，字母和特殊字符,且长度8-16"})
	} else if errors.Is(err, service.ErrFormatForPhoneInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 请输入11位的中国大陆地区的手机号"})
	} else if errors.Is(err, service.ErrFormatForAddressInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 地址格式错误"})
	} else if errors.Is(err, service.ErrUserListIsEmptyInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 用户列表为空"})
	} else if errors.Is(err, service.ErrUserNotFoundInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 用户不存在"})
	} else if errors.Is(err, service.ErrPasswordIsWrongInCustomer) {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 密码错误"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "失败, 系统错误"})
	}
}

func (c *CustomerHandler) SignUpCustomer(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirm_password"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "注册失败, JSON字段不匹配"})
		return
	}
	err := c.svc.SignUpCustomer(ctx, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功"}) // 响应
}

func (c *CustomerHandler) LogInCustomer(ctx *gin.Context) {
	type LogInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogInReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "登录失败, JSON字段不匹配"})
		return
	}
	err := c.svc.LogInCustomer(ctx, req.Email, req.Password)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "登录成功"}) // 响应
}

func (c *CustomerHandler) EditCustomer(ctx *gin.Context) {
	type EditReq struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := c.svc.EditCustomer(ctx, req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (c *CustomerHandler) ChangeCustomerPassword(ctx *gin.Context) {
	type ChangePasswordReq struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req ChangePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := c.svc.ChangeCustomerPassword(ctx, req.OldPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (c *CustomerHandler) ProfileCustomer(ctx *gin.Context) {
	customer, err := c.svc.ProfileCustomer(ctx)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerHandler) LogOutCustomer(ctx *gin.Context) {
	err := c.svc.LogOutCustomer(ctx)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func (c *CustomerHandler) EditCustomerByAdmin(ctx *gin.Context) {
	type EditReq struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := c.svc.EditCustomerByAdmin(ctx, req.Id, req.Name, req.Email, req.Password, req.Phone, req.Address)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (c *CustomerHandler) GetAllCustomers(ctx *gin.Context) {
	customers, err := c.svc.GetAllCustomers(ctx)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerHandler) GetCustomerById(ctx *gin.Context) {
	type GetReq struct {
		Id int64 `json:"id"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	customer, err := c.svc.GetCustomerById(ctx, req.Id)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerHandler) GetCustomerByName(ctx *gin.Context) {
	type GetReq struct {
		Name string `json:"name"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	customers, err := c.svc.GetCustomerByName(ctx, req.Name)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerHandler) DeleteCustomer(ctx *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id"`
	}
	var req DeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := c.svc.DeleteCustomer(ctx, req.Id)
	if err != nil {
		c.ErrOutputForCustomer(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
