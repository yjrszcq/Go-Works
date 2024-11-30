package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrReviewNotFound = dao.ErrReviewNotFound
)

type ReviewRepository struct {
	dao *dao.ReviewDAO
}

func NewReviewRepository(dao *dao.ReviewDAO) *ReviewRepository {
	return &ReviewRepository{
		dao: dao,
	}
}

func (repo *ReviewRepository) CreateReview(ctx context.Context, r domain.Review) error {
	return repo.dao.InsertReview(ctx, dao.Review{
		CustomerID: r.CustomerID,
		DishID:     r.DishID,
		Rating:     r.Rating,
		Comment:    r.Comment,
	})
}

func (repo *ReviewRepository) FindReviewById(ctx context.Context, id int64) (domain.Review, error) {
	r, err := repo.dao.FindReviewByID(ctx, id)
	if err != nil {
		return domain.Review{}, err
	}
	return domain.Review{
		Id:         r.ReviewID,
		CustomerID: r.CustomerID,
		DishID:     r.DishID,
		Rating:     r.Rating,
		Comment:    r.Comment,
		ReviewDate: r.ReviewDate,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}, nil
}

func (repo *ReviewRepository) FindReviewsByCustomerId(ctx context.Context, customerId int64) ([]domain.Review, error) {
	reviews, err := repo.dao.FindReviewsByCustomerID(ctx, customerId)
	if err != nil {
		return nil, err
	}
	var domainReviews []domain.Review
	for _, r := range reviews {
		domainReviews = append(domainReviews, domain.Review{
			Id:         r.ReviewID,
			CustomerID: r.CustomerID,
			DishID:     r.DishID,
			Rating:     r.Rating,
			Comment:    r.Comment,
			ReviewDate: r.ReviewDate,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}
	return domainReviews, nil
}

func (repo *ReviewRepository) FindReviewsByDishId(ctx context.Context, dishId int64) ([]domain.Review, error) {
	reviews, err := repo.dao.FindReviewsByDishID(ctx, dishId)
	if err != nil {
		return nil, err
	}
	var domainReviews []domain.Review
	for _, r := range reviews {
		domainReviews = append(domainReviews, domain.Review{
			Id:         r.ReviewID,
			CustomerID: r.CustomerID,
			DishID:     r.DishID,
			Rating:     r.Rating,
			Comment:    r.Comment,
			ReviewDate: r.ReviewDate,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}
	return domainReviews, nil
}

func (repo *ReviewRepository) FindReviewsByRating(ctx context.Context, rating int) ([]domain.Review, error) {
	reviews, err := repo.dao.FindReviewsByRating(ctx, rating)
	if err != nil {
		return nil, err
	}
	var domainReviews []domain.Review
	for _, r := range reviews {
		domainReviews = append(domainReviews, domain.Review{
			Id:         r.ReviewID,
			CustomerID: r.CustomerID,
			DishID:     r.DishID,
			Rating:     r.Rating,
			Comment:    r.Comment,
			ReviewDate: r.ReviewDate,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}
	return domainReviews, nil
}

func (repo *ReviewRepository) UpdateReview(ctx context.Context, r domain.Review) error {
	return repo.dao.UpdateReview(ctx, dao.Review{
		ReviewID: r.Id,
		Rating:   r.Rating,
		Comment:  r.Comment,
	})
}

func (repo *ReviewRepository) DeleteReview(ctx context.Context, id int64) error {
	return repo.dao.DeleteReview(ctx, id)
}
