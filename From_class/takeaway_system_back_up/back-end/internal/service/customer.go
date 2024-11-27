package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码错误")
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (svc *CustomerService) SignUpCustomer(ctx context.Context, c domain.Customer) error {
	// 要考虑加密放在哪里，然后就是存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hash)

	return svc.repo.CreateCustomer(ctx, c)
}

func (svc *CustomerService) LogInCustomer(ctx context.Context, c domain.Customer) (domain.Customer, error) {
	// 先找用户
	customer, err := svc.repo.FindCustomerByEmail(ctx, c.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Customer{}, ErrInvalidUserOrPassword
		} else {
			return domain.Customer{}, err
		}
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(c.Password))
	if err != nil {
		// DEBUG(打日志)
		return domain.Customer{}, ErrInvalidUserOrPassword
	}
	return customer, nil
}

func (svc *CustomerService) EditCustomer(ctx context.Context, c domain.Customer) error {
	if c.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		c.Password = string(hash)
	}
	return svc.repo.UpdateCustomer(ctx, c)
}

func (svc *CustomerService) GetProfileCustomer(ctx context.Context, c domain.Customer) (domain.Customer, error) {
	return svc.repo.FindCustomerById(ctx, c.Id)
}

func (svc *CustomerService) GetAllCustomers(ctx context.Context) ([]domain.Customer, error) {
	return svc.repo.FindAllCustomers(ctx)
}

func (svc *CustomerService) DeleteCustomer(ctx context.Context, c domain.Customer) error {
	return svc.repo.DeleteCustomer(ctx, c.Id)
}
