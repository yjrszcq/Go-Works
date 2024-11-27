package web

import (
	"back-end/internal/domain"
	"back-end/internal/servies"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2" // 支持复杂正则表达式的解析 (go自带的仅支持简单正则表达式的解析)
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	// 和用'"'比起来，用'`'看起来更清爽 (不用转译)
	emailRegexPattern = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	// emailRegexPattern = `^[0-9a-zA-Z_]{0,19}@[0-9a-zA-Z]{1,13}\.[com,cn,net]{1,3}$` 约束到了.com之类的后缀，此处先暂时不用
	// 密码包含至少一位数字，字母和特殊字符,且长度8-16
	passwordRegexPattern = `^(?![0-9a-zA-Z]+$)[a-zA-Z0-9~!@#$%^&*?_-]{8,16}$`
)

// UserHandler 我准备在它上面定义跟用户有关的路由

type CustomerHandler struct {
	svc         *servies.CustomerService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewCustomerHandler(svc *servies.CustomerService) *CustomerHandler { // 预编译正则表达式
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &CustomerHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (c *CustomerHandler) SignUp(ctx *gin.Context) {
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
	err = c.svc.SignUp(ctx.Request.Context(), domain.Customer{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, servies.ErrUserDuplicateEmail) {
			ctx.String(http.StatusOK, "邮箱冲突")
		} else {
			ctx.String(http.StatusOK, "系统错误")
		}
		return
	}
	ctx.String(http.StatusOK, "注册成功") // 响应
}

func (c *CustomerHandler) LogIn(ctx *gin.Context) {
	type LogInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogInReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "登录失败, 系统错误")
		return
	}

	customer, err := c.svc.LogIn(ctx, domain.Customer{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, servies.ErrInvalidUserOrPassword) {
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
	sess.Set("userId", customer.Id)
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *CustomerHandler) Edit(ctx *gin.Context) {}

func (u *CustomerHandler) Profile(ctx *gin.Context) {}
