package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type CustomerRepository struct {
	dao *dao.CustomerDAO
}

func NewCustomerRepository(dao *dao.CustomerDAO) *CustomerRepository {
	return &CustomerRepository{
		dao: dao,
	}

}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, c domain.Customer) error {
	return r.dao.InsertCustomer(ctx, dao.Customer{
		Email:    c.Email,
		Password: c.Password,
	})
	// 在这里操作缓存
}

func (r *CustomerRepository) FindCustomerById(ctx context.Context, id int64) (domain.Customer, error) {
	c, err := r.dao.FindCustomerById(ctx, id)
	if err != nil {
		return domain.Customer{}, err
	}
	return domain.Customer{
		Id:        c.CustomerID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}

func (r *CustomerRepository) FindCustomerByEmail(ctx context.Context, email string) (domain.Customer, error) {
	c, err := r.dao.FindCustomerByEmail(ctx, email)
	if err != nil {
		return domain.Customer{}, err
	}
	return domain.Customer{
		Id:       c.CustomerID,
		Email:    c.Email,
		Password: c.Password,
	}, nil
}

func (r *CustomerRepository) FindCustomerByName(ctx context.Context, name string) ([]domain.Customer, error) {
	customers, err := r.dao.FindCustomerByName(ctx, name)
	if err != nil {
		return nil, err
	}
	var result []domain.Customer
	for _, c := range customers {
		result = append(result, domain.Customer{
			Id:        c.CustomerID,
			Name:      c.Name,
			Email:     c.Email,
			Phone:     c.Phone,
			Address:   c.Address,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return result, nil
}

func (r *CustomerRepository) FindAllCustomers(ctx context.Context) ([]domain.Customer, error) {
	customers, err := r.dao.FindAllCustomers(ctx)
	if err != nil {
		return nil, err
	}
	var result []domain.Customer
	for _, c := range customers {
		result = append(result, domain.Customer{
			Id:        c.CustomerID,
			Name:      c.Name,
			Email:     c.Email,
			Phone:     c.Phone,
			Address:   c.Address,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return result, nil
}

func (r *CustomerRepository) UpdateCustomer(ctx context.Context, c domain.Customer) error {
	return r.dao.UpdateCustomer(ctx, dao.Customer{
		CustomerID: c.Id,
		Name:       c.Name,
		Email:      c.Email,
		Phone:      c.Phone,
		Address:    c.Address,
	})
}

func (r *CustomerRepository) UpdateCustomerPassword(ctx context.Context, c domain.Customer) error {
	return r.dao.UpdateCustomerPassword(ctx, dao.Customer{
		CustomerID: c.Id,
		Password:   c.Password,
	})
}

func (r *CustomerRepository) UpdateCustomerAll(ctx context.Context, c domain.Customer) error {
	return r.dao.UpdateCustomerAll(ctx, dao.Customer{
		CustomerID: c.Id,
		Name:       c.Name,
		Email:      c.Email,
		Password:   c.Password,
		Phone:      c.Phone,
		Address:    c.Address,
	})
}

func (r *CustomerRepository) DeleteCustomer(ctx context.Context, id int64) error {
	return r.dao.DeleteCustomer(ctx, id)
}
