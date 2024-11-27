package web

import (
	"back-end/internal/domain"
	"back-end/internal/service"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmployeeHandler struct {
	svc         *service.EmployeeService
	nameExp     *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
}

func NewEmployeeHandler(svc *service.EmployeeService) *EmployeeHandler {
	nameExp := regexp.MustCompile(nameRegexPattern, regexp.None)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	return &EmployeeHandler{
		svc:         svc,
		nameExp:     nameExp,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
	}
}

func getCurrentEmployeeId(ctx *gin.Context) int64 {
	sess := sessions.Default(ctx)
	if sess.Get("role").(string) != "employee" {
		return -1
	}
	employeeId := sess.Get("employeeId").(int64)
	return employeeId
}

func (e *EmployeeHandler) SignUpEmployee(ctx *gin.Context) {
	type SignUpReq struct {
		Name            string `json:"name"`
		Phone           string `json:"phone"`
		Email           string `json:"email"`
		Role            string `json:"role"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignUpReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err := e.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	ok, err = e.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码包含至少一位数字，字母和特殊字符,且长度8-16")
		return
	}

	ok, err = e.nameExp.MatchString(req.Name)
	if err != nil {
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "姓名不合法")
		return
	}

	ok, err = e.phoneExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "注册失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "请输入11位的中国大陆地区的手机号")
		return
	}

	err = e.svc.SignUpEmployee(ctx.Request.Context(), domain.Employee{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
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

func (e *EmployeeHandler) LogInEmployee(ctx *gin.Context) {
	type LogInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogInReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "登录失败, 系统错误")
		return
	}
	employee, err := e.svc.LogInEmployee(ctx, domain.Employee{
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
	sess.Set("role", "employee")
	sess.Set("employeeId", employee.Id)
	switch employee.Role {
	case "管理员":
		sess.Set("employeeRole", "admin")
	case "员工":
		sess.Set("employeeRole", "employee")
	case "送餐员":
		sess.Set("employeeRole", "deliveryman")
	}
	switch employee.Status {
	case "可用":
		sess.Set("employeeStatus", "available")
	case "不可用":
		sess.Set("employeeStatus", "unavailable")
	}
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (e *EmployeeHandler) EditEmployee(ctx *gin.Context) {
	type EditReq struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	id := getCurrentEmployeeId(ctx)
	if id == -1 {
		ctx.String(http.StatusOK, "无权限")
		return
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}

	ok, err := e.nameExp.MatchString(req.Name)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "姓名格式错误")
	}

	ok, err = e.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	ok, err = e.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码包含至少一位数字，字母和特殊字符,且长度8-16")
		return
	}

	ok, err = e.phoneExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "修改失败, 系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "请输入11位的中国大陆地区的手机号")
		return
	}

	if req.Role != "管理员" && req.Role != "员工" && req.Role != "送餐员" {
		ctx.String(http.StatusOK, "修改失败, 角色错误")
		return
	}

	if req.Status != "可用" && req.Status != "不可用" {
		ctx.String(http.StatusOK, "修改失败, 状态错误")
		return
	}

	err = e.svc.EditEmployee(ctx.Request.Context(), domain.Employee{
		Id:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   req.Status,
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

func (e *EmployeeHandler) ProfileEmployee(ctx *gin.Context) {
	id := getCurrentEmployeeId(ctx)
	if id == -1 {
		ctx.String(http.StatusOK, "无权限")
		return
	}
	employee, err := e.svc.GetProfileEmployee(ctx, domain.Employee{
		Id: id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (e *EmployeeHandler) LogOutEmployee(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Delete("role")
	sess.Delete("employeeId")
	sess.Delete("employeeRole")
	sess.Delete("employeeStatus")
	err := sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "退出成功")
}

func (e *EmployeeHandler) GetAllEmployees(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	role := sess.Get("employeeRole").(string)
	if role != "admin" {
		ctx.String(http.StatusOK, "无权限")
		return
	}
	employees, err := e.svc.GetAllEmployees(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) GetEmployeeById(ctx *gin.Context) {
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

	employee, err := e.svc.GetProfileEmployee(ctx, domain.Employee{
		Id: req.Id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (e *EmployeeHandler) DeleteEmployee(ctx *gin.Context) {
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

	err := e.svc.DeleteEmployee(ctx, domain.Employee{
		Id: req.Id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "删除成功")
}
