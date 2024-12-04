package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrDishNotFound = dao.ErrDishNotFound
)

type DishRepository struct {
	dao *dao.DishDAO
}

func NewDishRepository(dao *dao.DishDAO) *DishRepository {
	return &DishRepository{
		dao: dao,
	}
}

func (r *DishRepository) CreateDish(ctx context.Context, d domain.Dish) error {
	return r.dao.InsertDish(ctx, dao.Dish{
		Name:        d.Name,
		ImageURL:    d.ImageURL,
		Price:       d.Price,
		Description: d.Description,
		CategoryID:  d.CategoryID,
	})
}

func (r *DishRepository) FindDishById(ctx context.Context, id int64) (domain.Dish, error) {
	d, err := r.dao.FindDishByID(ctx, id)
	if err != nil {
		return domain.Dish{}, err
	}
	return domain.Dish{
		Id:          d.DishID,
		Name:        d.Name,
		ImageURL:    d.ImageURL,
		Price:       d.Price,
		Description: d.Description,
		CategoryID:  d.CategoryID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

func (r *DishRepository) FindDishByName(ctx context.Context, name string) ([]domain.Dish, error) {
	ds, err := r.dao.FindDishByName(ctx, name)
	if err != nil {
		return nil, err
	}
	var dishes []domain.Dish
	for _, d := range ds {
		dishes = append(dishes, domain.Dish{
			Id:          d.DishID,
			Name:        d.Name,
			ImageURL:    d.ImageURL,
			Price:       d.Price,
			Description: d.Description,
			CategoryID:  d.CategoryID,
		})
	}
	return dishes, nil
}

func (r *DishRepository) FindDishByCategoryID(ctx context.Context, categoryID int64) ([]domain.Dish, error) {
	ds, err := r.dao.FindDishByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	var dishes []domain.Dish
	for _, d := range ds {
		dishes = append(dishes, domain.Dish{
			Id:          d.DishID,
			Name:        d.Name,
			ImageURL:    d.ImageURL,
			Price:       d.Price,
			Description: d.Description,
			CategoryID:  d.CategoryID,
		})
	}
	return dishes, nil
}

func (r *DishRepository) FindAllDishes(ctx context.Context) ([]domain.Dish, error) {
	ds, err := r.dao.FindAllDishes(ctx)
	if err != nil {
		return nil, err
	}
	var dishes []domain.Dish
	for _, d := range ds {
		dishes = append(dishes, domain.Dish{
			Id:          d.DishID,
			Name:        d.Name,
			ImageURL:    d.ImageURL,
			Price:       d.Price,
			Description: d.Description,
			CategoryID:  d.CategoryID,
		})
	}
	return dishes, nil
}

func (r *DishRepository) UpdateDish(ctx context.Context, d domain.Dish) error {
	return r.dao.UpdateDish(ctx, dao.Dish{
		DishID:      d.Id,
		Name:        d.Name,
		ImageURL:    d.ImageURL,
		Price:       d.Price,
		Description: d.Description,
		CategoryID:  d.CategoryID,
	})
}

func (r *DishRepository) DeleteDish(ctx context.Context, id int64) error {
	return r.dao.DeleteDish(ctx, id)
}
