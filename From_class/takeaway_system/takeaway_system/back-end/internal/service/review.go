package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

var (
	ErrRecordNotFoundInReview        = repository.ErrReviewNotFound
	ErrUserHasNoPermissionInReview   = errors.New("无权限")
	ErrCustomerDoNotHaveOrderItem    = errors.New("用户没有此订单项，无法评价")
	ErrOrderItemReviewStatusInReview = errors.New("订单项评价状态错误")
	ErrUpdateOrderItemStatusInReview = errors.New("更新订单项评价状态失败")
	ErrRatingOutOfRangeInReview      = errors.New("评分应在1-5之间")
	ErrFormatForCommentInReview      = errors.New("评价应在200个字符以内")
)

const (
	commentRegexPattern = `^[a-zA-Z0-9\u4e00-\u9fa5,\.\u3002\uff0c ]{0,200}$`
)

type ReviewService struct {
	repo       *repository.ReviewRepository
	commentExp *regexp.Regexp
}

func NewReviewService(repo *repository.ReviewRepository) *ReviewService {
	commentExp := regexp.MustCompile(commentRegexPattern, regexp.None)
	return &ReviewService{
		repo:       repo,
		commentExp: commentExp,
	}
}

func (svc *ReviewService) CreateReview(ctx *gin.Context, orderItemId int64, rating int, comment string) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInReview
	}
	orderItem, err := GlobalOrderItem.FindOrderItemById(ctx, orderItemId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderItemNotFound) {
			return ErrCustomerDoNotHaveOrderItem
		} else {
			return err
		}
	}
	order, err := GlobalOrder.FindOrderById(ctx, orderItem.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrCustomerDoNotHaveOrderItem
		} else {
			return err
		}
	}
	if order.CustomerID != customerId {
		return ErrCustomerDoNotHaveOrderItem
	}
	if orderItem.ReviewStatus != "未评价" {
		return ErrOrderItemReviewStatusInReview
	}
	if rating < 1 || rating > 5 {
		return ErrRatingOutOfRangeInReview
	}
	ok, err := svc.commentExp.MatchString(comment)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForCommentInReview
	}
	err = GlobalOrderItem.UpdateOrderItemReviewStatus(ctx, domain.OrderItem{
		Id:           orderItemId,
		ReviewStatus: "已评价",
	})
	if err != nil {
		return ErrUpdateOrderItemStatusInReview
	}
	err = svc.repo.CreateReview(ctx, domain.Review{
		CustomerID: customerId,
		DishID:     orderItem.DishID,
		Rating:     rating,
		Comment:    comment,
	})
	if err != nil {
		err = GlobalOrderItem.UpdateOrderItemReviewStatus(ctx, domain.OrderItem{
			Id:           orderItemId,
			ReviewStatus: "未评价",
		})
		if err != nil {
			return ErrUpdateOrderItemStatusInReview
		}
		return err
	}
	return nil
}

func (svc *ReviewService) FindReviewById(ctx *gin.Context, id int64) (domain.Review, error) {
	review, err := svc.repo.FindReviewById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return domain.Review{}, ErrRecordNotFoundInReview
		} else {
			return domain.Review{}, err
		}
	}
	customer, err := GlobalCustomer.FindCustomerById(ctx, review.CustomerID)
	if err != nil {
		return review, nil
	}
	review.CustomerName = customer.Name
	return review, nil
}

func (svc *ReviewService) FindReviewsByCustomerId(ctx *gin.Context) ([]domain.Review, error) {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInReview
	}
	reviews, err := svc.repo.FindReviewsByCustomerId(ctx, customerId)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return nil, ErrRecordNotFoundInReview
		} else {
			return nil, err
		}
	}
	if reviews == nil {
		return nil, ErrRecordNotFoundInReview
	}
	return reviews, nil
}

func (svc *ReviewService) FindReviewsByRating(ctx *gin.Context, rating int) ([]domain.Review, error) {
	reviews, err := svc.repo.FindReviewsByRating(ctx, rating)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return nil, ErrRecordNotFoundInReview
		} else {
			return nil, err
		}
	}
	if reviews == nil {
		return nil, ErrRecordNotFoundInReview
	}
	for i := 0; i < len(reviews); i++ {
		customer, err := GlobalCustomer.FindCustomerById(ctx, reviews[i].CustomerID)
		if err != nil {
			continue
		}
		reviews[i].CustomerName = customer.Name
	}
	return reviews, nil
}

func (svc *ReviewService) FindReviewsByDishId(ctx *gin.Context, dishId int64) ([]domain.Review, error) {
	reviews, err := svc.repo.FindReviewsByDishId(ctx, dishId)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return nil, ErrRecordNotFoundInReview
		} else {
			return nil, err
		}
	}
	if reviews == nil {
		return nil, ErrRecordNotFoundInReview
	}
	for i := 0; i < len(reviews); i++ {
		customer, err := GlobalCustomer.FindCustomerById(ctx, reviews[i].CustomerID)
		if err != nil {
			continue
		}
		reviews[i].CustomerName = customer.Name
	}
	return reviews, nil
}

func (svc *ReviewService) EditReview(ctx *gin.Context, id int64, rating int, comment string) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInReview
	}
	review, err := svc.repo.FindReviewById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return ErrRecordNotFoundInReview
		} else {
			return err
		}
	}
	if review.CustomerID != customerId {
		return ErrUserHasNoPermissionInReview
	}
	if rating < 1 || rating > 5 {
		return ErrRatingOutOfRangeInReview
	}
	ok, err := svc.commentExp.MatchString(comment)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForCommentInReview
	}
	err = svc.repo.UpdateReview(ctx, domain.Review{
		Id:      id,
		Rating:  rating,
		Comment: comment,
	})
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return ErrRecordNotFoundInReview
		} else {
			return err
		}
	}
	return nil
}

func (svc *ReviewService) DeleteReview(ctx *gin.Context, id int64) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInReview
	}
	review, err := svc.repo.FindReviewById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return ErrRecordNotFoundInReview
		} else {
			return err
		}
	}
	if review.CustomerID != customerId {
		return ErrUserHasNoPermissionInReview
	}
	err = svc.repo.DeleteReview(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return ErrRecordNotFoundInReview
		} else {
			return err
		}
	}
	return nil
}
