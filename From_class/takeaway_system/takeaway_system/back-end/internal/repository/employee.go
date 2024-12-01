package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

type EmployeeRepository struct {
	dao *dao.EmployeeDAO
}

func NewEmployeeRepository(dao *dao.EmployeeDAO) *EmployeeRepository {
	return &EmployeeRepository{
		dao: dao,
	}
}

func (r *EmployeeRepository) CreateEmployee(ctx context.Context, e domain.Employee) error {
	return r.dao.InsertEmployee(ctx, dao.Employee{
		Name:     e.Name,
		Role:     e.Role,
		Email:    e.Email,
		Phone:    e.Phone,
		Password: e.Password,
	})
}

func (r *EmployeeRepository) FindEmployeePassword(ctx context.Context, id int64) (domain.Employee, error) {
	e, err := r.dao.FindEmployeeById(ctx, id)
	if err != nil {
		return domain.Employee{}, err
	}
	return domain.Employee{
		Id:       e.EmployeeID,
		Password: e.Password,
	}, nil
}

func (r *EmployeeRepository) FindEmployeeById(ctx context.Context, id int64) (domain.Employee, error) {
	e, err := r.dao.FindEmployeeById(ctx, id)
	if err != nil {
		return domain.Employee{}, err
	}
	return domain.Employee{
		Id:        e.EmployeeID,
		Name:      e.Name,
		Role:      e.Role,
		Email:     e.Email,
		Phone:     e.Phone,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}, nil
}

func (r *EmployeeRepository) FindEmployeeByEmail(ctx context.Context, email string) (domain.Employee, error) {
	e, err := r.dao.FindEmployeeByEmail(ctx, email)
	if err != nil {
		return domain.Employee{}, err
	}
	return domain.Employee{
		Id:       e.EmployeeID,
		Email:    e.Email,
		Password: e.Password,
		Role:     e.Role,
		Status:   e.Status,
	}, nil
}

func (r *EmployeeRepository) FindEmployeeByName(ctx context.Context, name string) ([]domain.Employee, error) {
	employees, err := r.dao.FindEmployeeByName(ctx, name)
	if err != nil {
		return nil, err
	}
	var result []domain.Employee
	for _, e := range employees {
		result = append(result, domain.Employee{
			Id:        e.EmployeeID,
			Name:      e.Name,
			Role:      e.Role,
			Email:     e.Email,
			Phone:     e.Phone,
			Status:    e.Status,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}
	return result, nil
}

func (r *EmployeeRepository) FindAllEmployees(ctx context.Context) ([]domain.Employee, error) {
	employees, err := r.dao.FindAllEmployees(ctx)
	if err != nil {
		return nil, err
	}
	var result []domain.Employee
	for _, e := range employees {
		result = append(result, domain.Employee{
			Id:        e.EmployeeID,
			Name:      e.Name,
			Role:      e.Role,
			Email:     e.Email,
			Phone:     e.Phone,
			Status:    e.Status,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}
	return result, nil
}

func (r *EmployeeRepository) UpdateEmployee(ctx context.Context, e domain.Employee) error {
	return r.dao.UpdateEmployee(ctx, dao.Employee{
		EmployeeID: e.Id,
		Name:       e.Name,
		Email:      e.Email,
		Phone:      e.Phone,
	})
}

func (r *EmployeeRepository) UpdateEmployeePassword(ctx context.Context, e domain.Employee) error {
	return r.dao.UpdateEmployeePassword(ctx, dao.Employee{
		EmployeeID: e.Id,
		Password:   e.Password,
	})
}

func (r *EmployeeRepository) UpdateEmployeeRole(ctx context.Context, e domain.Employee) error {
	return r.dao.UpdateEmployeeRole(ctx, dao.Employee{
		EmployeeID: e.Id,
		Role:       e.Role,
	})
}

func (r *EmployeeRepository) UpdateEmployeeStatus(ctx context.Context, e domain.Employee) error {
	return r.dao.UpdateEmployeeStatus(ctx, dao.Employee{
		EmployeeID: e.Id,
		Status:     e.Status,
	})
}

func (r *EmployeeRepository) UpdateEmployeeAll(ctx context.Context, e domain.Employee) error {
	return r.dao.UpdateEmployeeAll(ctx, dao.Employee{
		EmployeeID: e.Id,
		Name:       e.Name,
		Role:       e.Role,
		Email:      e.Email,
		Password:   e.Password,
		Phone:      e.Phone,
		Status:     e.Status,
	})
}

func (r *EmployeeRepository) DeleteEmployee(ctx context.Context, id int64) error {
	return r.dao.DeleteEmployee(ctx, id)
}
