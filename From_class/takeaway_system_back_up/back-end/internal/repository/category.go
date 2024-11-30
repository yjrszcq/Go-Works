package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrCategoryDuplicateName = dao.ErrCategoryDuplicateName
	ErrCategoryNotFound      = dao.ErrCategoryNotFound
)

type CategoryRepository struct {
	dao *dao.CategoryDAO
}

func NewCategoryRepository(dao *dao.CategoryDAO) *CategoryRepository {
	return &CategoryRepository{
		dao: dao,
	}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, c domain.Category) error {
	return r.dao.InsertCategory(ctx, dao.Category{
		Name:        c.Name,
		Description: c.Description,
	})
}

func (r *CategoryRepository) FindCategoryByID(ctx context.Context, id int64) (domain.Category, error) {
	c, err := r.dao.FindCategoryByID(ctx, id)
	if err != nil {
		return domain.Category{}, err
	}
	return domain.Category{
		Id:          c.CategoryID,
		Name:        c.Name,
		Description: c.Description,
	}, nil
}

func (r *CategoryRepository) FindCategoryByName(ctx context.Context, name string) (domain.Category, error) {
	c, err := r.dao.FindCategoryByName(ctx, name)
	if err != nil {
		return domain.Category{}, err
	}
	return domain.Category{
		Id:          c.CategoryID,
		Name:        c.Name,
		Description: c.Description,
	}, nil
}

func (r *CategoryRepository) FindAllCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.dao.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	var result []domain.Category
	for _, c := range categories {
		result = append(result, domain.Category{
			Id:          c.CategoryID,
			Name:        c.Name,
			Description: c.Description,
		})
	}
	return result, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, c domain.Category) error {
	return r.dao.UpdateCategory(ctx, dao.Category{
		CategoryID:  c.Id,
		Name:        c.Name,
		Description: c.Description,
	})
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	return r.dao.DeleteCategory(ctx, id)
}
