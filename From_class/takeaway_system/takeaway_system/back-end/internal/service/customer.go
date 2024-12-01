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
	ErrUserHasNoPermissionInCustomer    = errors.New("无权限")
	ErrUserDuplicateEmailInCustomer     = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPasswordInCustomer  = errors.New("邮箱或密码错误")
	ErrUserNotFoundInCustomer           = repository.ErrUserNotFound
	ErrPasswordIsWrongInCustomer        = errors.New("密码错误")
	ErrPasswordIsInconsistentInCustomer = errors.New("两次输入的密码不一致")
	ErrFormatForNameInCustomer          = errors.New("用户名格式错误")
	ErrFormatForEmailInCustomer         = errors.New("邮箱格式错误")
	ErrFormatForPasswordInCustomer      = errors.New("密码格式错误")
	ErrFormatForPhoneInCustomer         = errors.New("手机号格式错误")
	ErrFormatForAddressInCustomer       = errors.New("地址格式错误")
	ErrUserListIsEmptyInCustomer        = errors.New("用户列表为空")
)

const (
	// 和用'"'比起来，用'`'看起来更清爽 (不用转译)
	idRegexPattern    = `^[0-9]{1,20}$`
	nameRegexPattern  = `^[\u4e00-\u9fa5a-zA-Z0-9 ]{2,12}$`
	emailRegexPattern = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	// emailRegexPattern = `^[0-9a-zA-Z_]{0,19}@[0-9a-zA-Z]{1,13}\.[com,cn,net]{1,3}$` 约束到了.com之类的后缀，此处先暂时不用
	// 密码包含至少一位数字，字母和特殊字符,且长度8-16
	passwordRegexPattern = `^(?![0-9a-zA-Z]+$)[a-zA-Z0-9~!@#$%^&*?_-]{8,16}$`
	phoneRegexPattern    = `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	addressRegexPattern  = `^[\u4e00-\u9fa5a-zA-Z0-9 ]{2,50}$`
)

type CustomerService struct {
	repo        *repository.CustomerRepository
	idExp       *regexp.Regexp
	nameExp     *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
	addressExp  *regexp.Regexp
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	idExp := regexp.MustCompile(idRegexPattern, regexp.None)
	nameExp := regexp.MustCompile(nameRegexPattern, regexp.None)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	addressExp := regexp.MustCompile(addressRegexPattern, regexp.None)
	return &CustomerService{
		repo:        repo,
		idExp:       idExp,
		nameExp:     nameExp,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
		addressExp:  addressExp,
	}
}

func getCurrentCustomerId(ctx *gin.Context) (int64, error) {
	sess := sessions.Default(ctx)
	if sess.Get("role").(string) != "customer" {
		return -1, errors.New("无权限")
	}
	customerId := sess.Get("id").(int64)
	return customerId, nil
}

func (svc *CustomerService) SignUpCustomer(ctx *gin.Context, email string, password string, confirmPassword string) error {
	ok, err := svc.emailExp.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForEmailInCustomer
	}
	if confirmPassword != password {
		return ErrPasswordIsInconsistentInCustomer
	}
	ok, err = svc.passwordExp.MatchString(password)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPasswordInCustomer
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password = string(hash)
	err = svc.repo.CreateCustomer(ctx, domain.Customer{
		Email:    email,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserDuplicateEmail) {
			return ErrUserDuplicateEmailInCustomer
		} else {
			return err
		}
	}
	return nil
}

func (svc *CustomerService) LogInCustomer(ctx *gin.Context, email string, password string) error {
	// 先找用户
	customer, err := svc.repo.FindCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidUserOrPasswordInCustomer
		} else {
			return err
		}
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err != nil {
		// DEBUG(打日志)
		return ErrInvalidUserOrPasswordInCustomer
	}
	// 设置 session
	sess := sessions.Default(ctx)
	sess.Set("role", "customer")
	sess.Set("id", customer.Id)
	sess.Set("status", "normal")
	err = sess.Save()
	if err != nil {
		return err
	}
	return nil
}

func (svc *CustomerService) EditCustomer(ctx *gin.Context, name string, email string, phone string, address string) error {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInCustomer
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInCustomer
	}
	ok, err = svc.emailExp.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForEmailInCustomer
	}
	ok, err = svc.phoneExp.MatchString(phone)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPhoneInCustomer
	}

	ok, err = svc.addressExp.MatchString(address)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForAddressInCustomer
	}
	err = svc.repo.UpdateCustomer(ctx, domain.Customer{
		Id:      id,
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserDuplicateEmail) {
			return ErrUserDuplicateEmailInCustomer
		} else if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidUserOrPasswordInCustomer
		} else {
			return err
		}
	}
	return nil
}

func (svc *CustomerService) ChangeCustomerPassword(ctx *gin.Context, oldPassword string, newPassword string, confirmPassword string) error {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInCustomer
	}
	if confirmPassword != newPassword {
		return ErrPasswordIsInconsistentInCustomer
	}
	ok, err := svc.passwordExp.MatchString(newPassword)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPasswordInCustomer
	}
	customer, err := svc.repo.FindCustomerPassword(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInCustomer
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(oldPassword))
	if err != nil {
		return ErrPasswordIsWrongInCustomer
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newPassword = string(hash)
	err = svc.repo.UpdateCustomerPassword(ctx, domain.Customer{
		Id:       id,
		Password: newPassword,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *CustomerService) ProfileCustomer(ctx *gin.Context) (domain.Customer, error) {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return domain.Customer{}, ErrUserHasNoPermissionInCustomer
	}
	c, err := svc.repo.FindCustomerById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Customer{}, ErrUserNotFoundInCustomer
		} else {
			return domain.Customer{}, err
		}
	}
	return c, nil
}

func (svc *CustomerService) LogOutCustomer(ctx *gin.Context) error {
	sess := sessions.Default(ctx)
	sess.Clear()
	err := sess.Save()
	if err != nil {
		return err
	}
	return nil
}

func (svc *CustomerService) EditCustomerByAdmin(ctx *gin.Context, id int64, name string, email string, password string, phone string, address string) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInCustomer
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInCustomer
	}
	ok, err = svc.emailExp.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForEmailInCustomer
	}
	ok, err = svc.passwordExp.MatchString(password)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPasswordInCustomer
	}
	ok, err = svc.phoneExp.MatchString(phone)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForPhoneInCustomer
	}

	ok, err = svc.addressExp.MatchString(address)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForAddressInCustomer
	}
	err = svc.repo.UpdateCustomerAll(ctx, domain.Customer{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Phone:    phone,
		Address:  address,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserDuplicateEmail) {
			return ErrUserDuplicateEmailInCustomer
		} else if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInCustomer
		} else {
			return err
		}
	}
	return nil
}

func (svc *CustomerService) GetAllCustomers(ctx *gin.Context) ([]domain.Customer, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInCustomer
	}
	customers, err := svc.repo.FindAllCustomers(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserListIsEmptyInCustomer
		} else {
			return nil, err
		}
	}
	if customers == nil {
		return nil, ErrUserListIsEmptyInCustomer
	}
	return customers, nil
}

func (svc *CustomerService) GetCustomerById(ctx *gin.Context, id int64) (domain.Customer, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return domain.Customer{}, ErrUserHasNoPermissionInCustomer
	}
	customer, err := svc.repo.FindCustomerById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Customer{}, ErrUserNotFoundInCustomer
		} else {
			return domain.Customer{}, err
		}
	}
	return customer, nil
}

func (svc *CustomerService) GetCustomerByName(ctx *gin.Context, name string) ([]domain.Customer, error) {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return nil, ErrUserHasNoPermissionInCustomer
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrFormatForNameInCustomer
	}
	customer, err := svc.repo.FindCustomerByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFoundInCustomer
		} else {
			return nil, err
		}
	}
	if customer == nil {
		return nil, ErrUserListIsEmptyInCustomer
	}
	return customer, nil
}

func (svc *CustomerService) DeleteCustomer(ctx *gin.Context, id int64) error {
	if sessions.Default(ctx).Get("role").(string) != "admin" {
		return ErrUserHasNoPermissionInCustomer
	}
	err := svc.repo.DeleteCustomer(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFoundInCustomer
		} else {
			return err
		}
	}
	return nil
}
