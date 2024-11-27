package web

import (
	"back-end/internal/domain"
	"back-end/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2" // 支持复杂正则表达式的解析 (go自带的仅支持简单正则表达式的解析)
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	// 和用'"'比起来，用'`'看起来更清爽 (不用转译)
	nameRegexPattern  = `^[\u4e00-\u9fa5a-zA-Z0-9 ]{2,12}$`
	emailRegexPattern = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	// emailRegexPattern = `^[0-9a-zA-Z_]{0,19}@[0-9a-zA-Z]{1,13}\.[com,cn,net]{1,3}$` 约束到了.com之类的后缀，此处先暂时不用
	// 密码包含至少一位数字，字母和特殊字符,且长度8-16
	passwordRegexPattern = `^(?![0-9a-zA-Z]+$)[a-zA-Z0-9~!@#$%^&*?_-]{8,16}$`
	phoneRegexPattern    = `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	addressRegexPattern  = `^[\u4e00-\u9fa5a-zA-Z0-9 ]{2,50}$`
)

// UserHandler 我准备在它上面定义跟用户有关的路由

type CustomerHandler struct {
	svc         *service.CustomerService
	nameExp     *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
	addressExp  *regexp.Regexp
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler { // 预编译正则表达式
	nameExp := regexp.MustCompile(nameRegexPattern, regexp.None)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	addressExp := regexp.MustCompile(addressRegexPattern, regexp.None)
	return &CustomerHandler{
		svc:         svc,
		nameExp:     nameExp,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
		addressExp:  addressExp,
	}
}

func getCurrentCustomerId(ctx *gin.Context) int64 {
	sess := sessions.Default(ctx)
	if sess.Get("role").(string) != "customer" {
		return -1
	}
	customerId := sess.Get("customerId").(int64)
	return customerId
}

func (c *CustomerHandler) SignUpCustomer(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirm_password"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}

	ok, err := c.emailExp.MatchString(req.Email)
	if err != nil { // 这里一般超时才会报err
		// 记录日志
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = c.passwordExp.MatchString(req.Password)
	if err != nil { // 这里一般超时才会报err
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码包含至少一位数字，字母和特殊字符,且长度8-16")
		return
	}
	fmt.Printf("%v", req)
	// 在下面进行数据库操作(调用一下 svc 的方法)
	err = c.svc.SignUpCustomer(ctx.Request.Context(), domain.Customer{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserDuplicateEmail) {
			ctx.String(http.StatusOK, "邮箱冲突")
		} else {
			ctx.String(http.StatusOK, "系统错误")
		}
		return
	}
	ctx.String(http.StatusOK, "注册成功") // 响应
}

func (c *CustomerHandler) LogInCustomer(ctx *gin.Context) {
	type LogInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogInReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "登录失败, 系统错误")
		return
	}

	customer, err := c.svc.LogInCustomer(ctx, domain.Customer{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			ctx.String(http.StatusOK, "用户名或密码错误")
			return
		} else {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
	}
	// 在这里登录成功
	// 设置 session
	sess := sessions.Default(ctx)
	// 可以随便设置值了
	// 要放在 session 里面的值
	sess.Set("role", "customer")
	sess.Set("customerId", customer.Id)
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (c *CustomerHandler) EditCustomer(ctx *gin.Context) {
	type EditReq struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}

	id := getCurrentCustomerId(ctx)
	if id == -1 {
		ctx.String(http.StatusOK, "无权限")
		return
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}

	ok, err := c.nameExp.MatchString(req.Name)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "姓名格式错误")
	}

	ok, err = c.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	ok, err = c.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码包含至少一位数字，字母和特殊字符,且长度8-16")
		return
	}

	ok, err = c.phoneExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "请输入11位的中国大陆地区的手机号")
		return
	}

	ok, err = c.addressExp.MatchString(req.Address)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "地址格式错误")
		return
	}

	err = c.svc.EditCustomer(ctx, domain.Customer{
		Id:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserDuplicateEmail) {
			ctx.String(http.StatusOK, "邮箱冲突")
		} else if errors.Is(err, service.ErrInvalidUserOrPassword) {
			ctx.String(http.StatusOK, "用户不存在")
		} else {
			ctx.String(http.StatusOK, "系统错误")
		}
		return
	}
	ctx.String(http.StatusOK, "修改成功") // 响应
}

func (c *CustomerHandler) ProfileCustomer(ctx *gin.Context) {
	id := getCurrentCustomerId(ctx)
	if id == -1 {
		ctx.String(http.StatusOK, "无权限")
		return
	}
	customer, err := c.svc.GetProfileCustomer(ctx, domain.Customer{
		Id: id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerHandler) LogOutCustomer(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Delete("role")
	sess.Delete("customerId")
	err := sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "退出成功")
}

func (c *CustomerHandler) GetAllCustomers(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	role := sess.Get("employeeRole").(string)
	if role != "admin" {
		ctx.String(http.StatusOK, "无权限")
		return
	}
	customers, err := c.svc.GetAllCustomers(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerHandler) GetCustomerById(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	if sess.Get("role") != "admin" {
		ctx.String(http.StatusOK, "无权限")
		return
	}
	type GetReq struct {
		Id int64 `json:"id"`
	}

	var req GetReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "查询失败, 系统错误")
		return
	}
	customer, err := c.svc.GetProfileCustomer(ctx, domain.Customer{
		Id: req.Id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerHandler) DeleteCustomer(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	if sess.Get("role") != "admin" {
		ctx.String(http.StatusOK, "无权限")
		return
	}
	type DeleteReq struct {
		Id int64 `json:"id"`
	}

	var req DeleteReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "删除失败, 系统错误")
		return
	}
	err := c.svc.DeleteCustomer(ctx, domain.Customer{
		Id: req.Id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "删除成功")
}
