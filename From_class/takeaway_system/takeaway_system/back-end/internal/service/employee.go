package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserHasNoPermissionInEmployee               = errors.New("无权限")
	ErrUserDuplicateEmailInEmployee                = repository.ErrUserDuplicateEmail
	ErrPasswordIsWrongInEmployee                   = errors.New("密码错误")
	ErrInvalidUserOrPasswordInEmployee             = errors.New("邮箱或密码错误")
	ErrUserNotFoundInEmployee                      = repository.ErrUserNotFound
	ErrPasswordIsInconsistentInEmployee            = errors.New("两次输入的密码不一致")
	ErrFormatForNameInEmployee                     = errors.New("用户名格式错误")
	ErrFormatForEmailInEmployee                    = errors.New("邮箱格式错误")
	ErrFormatForPasswordInEmployee                 = errors.New("密码格式错误")
	ErrFormatForPhoneInEmployee                    = errors.New("手机号格式错误")
	ErrUserListIsEmptyInEmployee                   = errors.New("用户列表为空")
	ErrRoleInputInEmployee                         = errors.New("角色输入错误")
	ErrStatusInputInEmployee                       = errors.New("状态输入错误")
	ErrUnassignedEmployeeMustUnavailableInEmployee = errors.New("未分配员工不可用")
)

type EmployeeService struct {
	repo        *repository.EmployeeRepository
	idExp       *regexp.Regexp
	nameExp     *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
}

func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	idExp := regexp.MustCompile(idRegexPattern, regexp.None)
	nameExp := regexp.MustCompile(nameRegexPattern, regexp.None)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	return &EmployeeService{
		repo:        repo,
		idExp:       idExp,
		nameExp:     nameExp,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
	}
}

func getCurrentEmployeeId(ctx *gin.Context) (int64, error) {
	sess := sessions.Default(ctx)
	role := sess.Get("role").(string)
	if role != "employee" && role != "deliveryman" {
		return -1, ErrUserHasNoPermissionInEmployee
	}
	employeeId := sess.Get("id").(int64)
	return employeeId, nil
}

func (svc *EmployeeService) SignUpEmployee(ctx *gin.Context, name string, phone string, email string, password string, confirmPassword string) error {
	if confirmPassword != password {
		return ErrPasswordIsInconsistentInEmployee
	}
	ok, err := svc.emailExp.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForEmailInEmployee
	}
	ok, err = svc.passwordExp.MatchString(password)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPasswordInEmployee
	}
	ok, err = svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInEmployee
	}
	ok, err = svc.phoneExp.MatchString(phone)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPhoneInEmployee
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password = string(hash)
	err = svc.repo.CreateEmployee(ctx, domain.Employee{
		Name:     name,
		Email:    email,
		Password: password,
		Phone:    phone,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserDuplicateEmail) {
			return ErrUserDuplicateEmailInEmployee
		} else {
			return err
		}
	}
	return nil
}

func (svc *EmployeeService) LogInEmployee(ctx *gin.Context, email string, password string) (domain.Employee, error) {
	employee, err := svc.repo.FindEmployeeByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Employee{}, ErrInvalidUserOrPasswordInEmployee
		} else {
			return domain.Employee{}, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password))
	if err != nil {
		// DEBUG(打日志)
		return domain.Employee{}, ErrInvalidUserOrPasswordInEmployee
	}
	// 设置 session
	sess := sessions.Default(ctx)
	switch employee.Role {
	case "员工":
		sess.Set("role", "employee")
	case "送餐员":
		sess.Set("role", "deliveryman")
	case "未分配":
		sess.Set("role", "unassigned")
	}
	sess.Set("id", employee.Id)
	switch employee.Status {
	case "可用":
		sess.Set("status", "available")
	case "不可用":
		sess.Set("status", "unavailable")
	}
	err = sess.Save()
	if err != nil {
		return domain.Employee{}, err
	}
	employee.Password = ""
	return employee, nil
}

func (svc *EmployeeService) EditEmployee(ctx *gin.Context, name string, phone string, email string) error {
	id, err := getCurrentEmployeeId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInEmployee
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInEmployee
	}
	ok, err = svc.emailExp.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForEmailInEmployee
	}
	ok, err = svc.phoneExp.MatchString(phone)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPhoneInEmployee
	}
	err = svc.repo.UpdateEmployee(ctx, domain.Employee{
		Id:    id,
		Name:  name,
		Email: email,
		Phone: phone,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserDuplicateEmail) {
			return ErrUserDuplicateEmailInEmployee
		} else if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidUserOrPasswordInEmployee
		} else {
			return err
		}
	}
	return nil
}

func (svc *EmployeeService) ChangeEmployeePassword(ctx *gin.Context, oldPassword string, newPassword string, confirmPassword string) error {
	id, err := getCurrentEmployeeId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInEmployee
	}
	if newPassword != confirmPassword {
		return ErrPasswordIsInconsistentInEmployee
	}
	ok, err := svc.passwordExp.MatchString(newPassword)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPasswordInEmployee
	}
	employee, err := svc.repo.FindEmployeePassword(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(oldPassword))
	if err != nil {
		return ErrPasswordIsWrongInEmployee
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newPassword = string(hash)
	err = svc.repo.UpdateEmployeePassword(ctx, domain.Employee{
		Id:       id,
		Password: newPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	return nil
}

func (svc *EmployeeService) ProfileEmployee(ctx *gin.Context) (domain.Employee, error) {
	id, err := getCurrentEmployeeId(ctx)
	if err != nil {
		return domain.Employee{}, ErrUserHasNoPermissionInEmployee
	}
	employee, err := svc.repo.FindEmployeeById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Employee{}, ErrUserNotFoundInEmployee
		} else {
			return domain.Employee{}, err
		}
	}
	return employee, nil
}

func (svc *EmployeeService) LogOutEmployee(ctx *gin.Context) error {
	sess := sessions.Default(ctx)
	sess.Clear()
	err := sess.Save()
	if err != nil {
		return err
	}
	return nil
}

func (svc *EmployeeService) EditEmployeeByAdmin(ctx *gin.Context, id int64, name string, email string, password string, phone string, role string, status string) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInEmployee
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInEmployee
	}
	ok, err = svc.emailExp.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForEmailInEmployee
	}
	ok, err = svc.passwordExp.MatchString(password)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPasswordInEmployee
	}
	ok, err = svc.phoneExp.MatchString(phone)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPhoneInEmployee
	}
	if role != "员工" && role != "送餐员" {
		return ErrRoleInputInEmployee
	}
	if status != "可用" && status != "不可用" {
		return ErrStatusInputInEmployee
	}
	employee, err := svc.repo.FindEmployeeById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	if employee.Role == "未分配" && status == "可用" {
		return ErrUnassignedEmployeeMustUnavailableInEmployee
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password = string(hash)
	err = svc.repo.UpdateEmployeeAll(ctx, domain.Employee{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Phone:    phone,
		Role:     role,
		Status:   status,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserDuplicateEmail) {
			return ErrUserDuplicateEmailInEmployee
		} else if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	return nil
}

func (svc *EmployeeService) EditEmployeeRole(ctx *gin.Context, id int64, role string) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInEmployee
	}
	if role != "员工" && role != "送餐员" {
		return ErrRoleInputInEmployee
	}
	err := svc.repo.UpdateEmployeeRole(ctx, domain.Employee{
		Id:   id,
		Role: role,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	return nil
}

func (svc *EmployeeService) EditEmployeeStatus(ctx *gin.Context, id int64, status string) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInEmployee
	}
	if status != "可用" && status != "不可用" {
		return ErrStatusInputInEmployee
	}
	employee, err := svc.repo.FindEmployeeById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	if employee.Role == "未分配" && status == "可用" {
		return ErrUnassignedEmployeeMustUnavailableInEmployee
	}
	err = svc.repo.UpdateEmployeeStatus(ctx, domain.Employee{
		Id:     id,
		Status: status,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *EmployeeService) InitEmployeePassword(ctx *gin.Context, id int64, password string) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInEmployee
	}
	if password == "" {
		password = GlobalDefaultPassword
	} else {
		ok, err := svc.passwordExp.MatchString(password)
		if err != nil {
			return err
		}
		if !ok {
			return ErrFormatForPasswordInEmployee
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password = string(hash)
	err = svc.repo.UpdateEmployeePassword(ctx, domain.Employee{
		Id:       id,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	return nil
}

func (svc *EmployeeService) GetAllEmployees(ctx *gin.Context) ([]domain.Employee, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInEmployee
	}
	e, err := svc.repo.FindAllEmployees(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserListIsEmptyInEmployee
		} else {
			return nil, err
		}
	}
	if e == nil {
		return nil, ErrUserListIsEmptyInEmployee
	}
	return e, nil
}

func (svc *EmployeeService) GetEmployeeById(ctx *gin.Context, id int64) (domain.Employee, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return domain.Employee{}, ErrUserHasNoPermissionInEmployee
	}
	e, err := svc.repo.FindEmployeeById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Employee{}, ErrUserNotFoundInEmployee
		} else {
			return domain.Employee{}, err
		}
	}
	return e, nil
}

func (svc *EmployeeService) GetEmployeeByName(ctx *gin.Context, name string) ([]domain.Employee, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInEmployee
	}
	e, err := svc.repo.FindEmployeeByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFoundInEmployee
		} else {
			return nil, err
		}
	}
	if e == nil {
		return nil, ErrUserNotFoundInEmployee
	}
	return e, nil
}

func (svc *EmployeeService) GetEmployeeByRole(ctx *gin.Context, role string) ([]domain.Employee, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInEmployee
	}
	if role != "管理员" && role != "员工" && role != "送餐员" && role != "未分配" {
		return nil, ErrRoleInputInEmployee
	}
	e, err := svc.repo.FindEmployeeByRole(ctx, role)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserListIsEmptyInEmployee
		} else {
			return nil, err
		}
	}
	if e == nil {
		return nil, ErrUserListIsEmptyInEmployee
	}
	return e, nil
}

func (svc *EmployeeService) GetEmployeeByStatus(ctx *gin.Context, status string) ([]domain.Employee, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInEmployee
	}
	if status != "可用" && status != "不可用" {
		return nil, ErrStatusInputInEmployee
	}
	e, err := svc.repo.FindEmployeeByStatus(ctx, status)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserListIsEmptyInEmployee
		} else {
			return nil, err
		}
	}
	if e == nil {
		return nil, ErrUserListIsEmptyInEmployee
	}
	return e, nil
}

func (svc *EmployeeService) GetNewEmployees(ctx *gin.Context) ([]domain.Employee, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInEmployee
	}
	e, err := svc.repo.FindEmployeeByRole(ctx, "未分配")
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserListIsEmptyInEmployee
		} else {
			return nil, err
		}
	}
	if e == nil {
		return nil, ErrUserListIsEmptyInEmployee
	}
	return e, nil
}

func (svc *EmployeeService) DeleteEmployee(ctx *gin.Context, id int64) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInEmployee
	}
	err := svc.repo.DeleteEmployee(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInEmployee
		} else {
			return err
		}
	}
	return nil
}
