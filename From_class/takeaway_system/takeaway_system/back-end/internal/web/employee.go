package web

import (
	"back-end/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmployeeHandler struct {
	svc *service.EmployeeService
}

func NewEmployeeHandler(svc *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		svc: svc,
	}
}

func (e *EmployeeHandler) ErrOutputForEmployee(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrUserHasNoPermissionInEmployee) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "失败, 无权限"})
	} else if errors.Is(err, service.ErrUserDuplicateEmailInEmployee) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "失败, 邮箱冲突"})
	} else if errors.Is(err, service.ErrInvalidUserOrPasswordInEmployee) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "失败, 邮箱或密码错误"})
	} else if errors.Is(err, service.ErrFormatForNameInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 姓名格式错误"})
	} else if errors.Is(err, service.ErrFormatForEmailInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 邮箱格式错误"})
	} else if errors.Is(err, service.ErrFormatForPasswordInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 密码包含至少一位数字，字母和特殊字符,且长度8-16"})
	} else if errors.Is(err, service.ErrFormatForPhoneInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 请输入11位的中国大陆地区的手机号"})
	} else if errors.Is(err, service.ErrUserListIsEmptyInEmployee) {
		ctx.JSON(http.StatusOK, gin.H{"message": "成功, 用户列表为空"})
	} else if errors.Is(err, service.ErrUserNotFoundInEmployee) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "失败, 用户不存在"})
	} else if errors.Is(err, service.ErrRoleInputInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 角色不存在"})
	} else if errors.Is(err, service.ErrStatusInputInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 状态不存在"})
	} else if errors.Is(err, service.ErrUnassignedEmployeeMustUnavailableInEmployee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "失败, 未分配员工不可用"})
	} else if errors.Is(err, service.ErrPasswordIsWrongInEmployee) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "失败, 密码错误"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "失败, 系统错误"})
	}
}

func (e *EmployeeHandler) SignUpEmployee(ctx *gin.Context) {
	type SignUpReq struct {
		Name            string `json:"name"`
		Phone           string `json:"phone"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "注册失败, JSON字段不匹配"})
		return
	}
	err := e.svc.SignUpEmployee(ctx, req.Name, req.Phone, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功"}) // 响应
}

func (e *EmployeeHandler) LogInEmployee(ctx *gin.Context) {
	type LogInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LogInReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "登录失败, JSON字段不匹配"})
		return
	}
	err := e.svc.LogInEmployee(ctx, req.Email, req.Password)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func (e *EmployeeHandler) EditEmployee(ctx *gin.Context) {
	type EditReq struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := e.svc.EditEmployee(ctx, req.Name, req.Phone, req.Email)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (e *EmployeeHandler) ChangeEmployeePassword(ctx *gin.Context) {
	type ChangeReq struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req ChangeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := e.svc.ChangeEmployeePassword(ctx, req.OldPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (e *EmployeeHandler) ProfileEmployee(ctx *gin.Context) {
	employee, err := e.svc.ProfileEmployee(ctx)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (e *EmployeeHandler) LogOutEmployee(ctx *gin.Context) {
	err := e.svc.LogOutEmployee(ctx)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "退出成功"})
}

func (e *EmployeeHandler) EditEmployeeByAdmin(ctx *gin.Context) {
	type EditReq struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := e.svc.EditEmployeeByAdmin(ctx, req.Id, req.Name, req.Email, req.Password, req.Phone, req.Role, req.Status)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (e *EmployeeHandler) EditEmployeeRole(ctx *gin.Context) {
	type EditReq struct {
		Id   int64  `json:"id"`
		Role string `json:"role"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := e.svc.EditEmployeeRole(ctx, req.Id, req.Role)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (e *EmployeeHandler) EditEmployeeStatus(ctx *gin.Context) {
	type EditReq struct {
		Id     int64  `json:"id"`
		Status string `json:"status"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "修改失败, JSON字段不匹配"})
		return
	}
	err := e.svc.EditEmployeeStatus(ctx, req.Id, req.Status)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "修改成功"}) // 响应
}

func (e *EmployeeHandler) InitEmployeePassword(ctx *gin.Context) {
	type InitReq struct {
		Id       int64  `json:"id"`
		Password string `json:"password"`
	}
	var req InitReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "初始化失败, JSON字段不匹配"})
		return
	}
	err := e.svc.InitEmployeePassword(ctx, req.Id, req.Password)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "初始化成功"})
}

func (e *EmployeeHandler) GetAllEmployees(ctx *gin.Context) {
	employees, err := e.svc.GetAllEmployees(ctx)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) GetEmployeeById(ctx *gin.Context) {
	type GetReq struct {
		Id int64 `json:"id"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	employee, err := e.svc.GetEmployeeById(ctx, req.Id)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (e *EmployeeHandler) GetEmployeeByName(ctx *gin.Context) {
	type GetReq struct {
		Name string `json:"name"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	employees, err := e.svc.GetEmployeeByName(ctx, req.Name)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) GetEmployeeByRole(ctx *gin.Context) {
	type GetReq struct {
		Role string `json:"role"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	employees, err := e.svc.GetEmployeeByRole(ctx, req.Role)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) GetEmployeeByStatus(ctx *gin.Context) {
	type GetReq struct {
		Status string `json:"status"`
	}
	var req GetReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询失败, JSON字段不匹配"})
		return
	}
	employees, err := e.svc.GetEmployeeByStatus(ctx, req.Status)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) GetNewEmployees(ctx *gin.Context) {
	employees, err := e.svc.GetNewEmployees(ctx)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeHandler) DeleteEmployee(ctx *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id"`
	}
	var req DeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除失败, JSON字段不匹配"})
		return
	}
	err := e.svc.DeleteEmployee(ctx, req.Id)
	if err != nil {
		e.ErrOutputForEmployee(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
