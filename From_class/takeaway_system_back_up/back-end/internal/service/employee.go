package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService struct {
	repo *repository.EmployeeRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		repo: repo,
	}
}

func (svc *EmployeeService) SignUpEmployee(ctx context.Context, e domain.Employee) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hash)
	return svc.repo.CreateEmployee(ctx, e)
}

func (svc *EmployeeService) LogInEmployee(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	// 先找用户
	employee, err := svc.repo.FindEmployeeByEmail(ctx, e.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.Employee{}, ErrInvalidUserOrPassword
		} else {
			return domain.Employee{}, err
		}
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(e.Password))
	if err != nil {
		// DEBUG(打日志)
		return domain.Employee{}, ErrInvalidUserOrPassword
	}
	return employee, nil
}

func (svc *EmployeeService) EditEmployee(ctx context.Context, e domain.Employee) error {
	if e.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		e.Password = string(hash)
	}
	return svc.repo.UpdateEmployee(ctx, e)
}

func (svc *EmployeeService) GetProfileEmployee(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	return svc.repo.FindEmployeeById(ctx, e.Id)
}

func (svc *EmployeeService) GetAllEmployees(ctx context.Context) ([]domain.Employee, error) {
	return svc.repo.FindAllEmployees(ctx)
}

func (svc *EmployeeService) DeleteEmployee(ctx context.Context, e domain.Employee) error {
	return svc.repo.DeleteEmployee(ctx, e.Id)
}
