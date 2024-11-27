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

func (r *CustomerRepository) Create(ctx context.Context, c domain.Customer) error {
	return r.dao.Insert(ctx, dao.Customer{
		Email:    c.Email,
		Password: c.Password,
	})
	// 在这里操作缓存
}

func (r *CustomerRepository) FindById(ctx context.Context, id int64) (domain.Customer, error) {
	c, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.Customer{}, err
	}
	return domain.Customer{
		Id:       c.CustomerID,
		Email:    c.Email,
		Password: c.Password,
	}, nil
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (domain.Customer, error) {
	c, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.Customer{}, err
	}
	return domain.Customer{
		Id:       c.CustomerID,
		Email:    c.Email,
		Password: c.Password,
	}, nil
}

func (r *CustomerRepository) Update(ctx context.Context, c domain.Customer) error {
	return r.dao.Update(ctx, dao.Customer{
		CustomerID: c.Id,
		Name:       c.Name,
		Email:      c.Email,
		Password:   c.Password,
		Phone:      c.Phone,
		Address:    c.Address,
	})
}
