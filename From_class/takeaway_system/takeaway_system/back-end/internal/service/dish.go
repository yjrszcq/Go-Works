package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrRecordNotFoundInDish      = repository.ErrDishNotFound
	ErrUserHasNoPermissionInDish = errors.New("无权限")
	ErrFormatForDishNameInDish   = errors.New("菜品名称格式错误")
	ErrFormatForImageUrlInDish   = errors.New("图片链接格式错误")
	ErrRangeForPriceInDish       = errors.New("价格范围错误")
)

const (
	dishNameRegexPattern = `^[a-zA-Z0-9\u4e00-\u9fa5 ]{1,20}$`
	imageURLRegexPattern = `^https?://.*\.(?:png|jpg|jpeg|webp)$`
	dishPriceMin         = 0.01
	dishPriceMax         = 9999.99
)

type DishService struct {
	repo        *repository.DishRepository
	dishExp     *regexp.Regexp
	imageURLExp *regexp.Regexp
}

func NewDishService(repo *repository.DishRepository) *DishService {
	dishExp := regexp.MustCompile(dishNameRegexPattern, regexp.None)
	imageURLExp := regexp.MustCompile(imageURLRegexPattern, regexp.None)
	return &DishService{
		repo:        repo,
		dishExp:     dishExp,
		imageURLExp: imageURLExp,
	}
}

func (svc *DishService) CreateDish(ctx *gin.Context, name string, imageUrl string, price float64, categoryId int64) error {
	role := sessions.Default(ctx).Get("role").(string)
	if role != "employee" && role != "admin" {
		return ErrUserHasNoPermissionInDish
	}
	ok, err := svc.dishExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForDishNameInDish
	}
	ok, err = svc.imageURLExp.MatchString(imageUrl)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForImageUrlInDish
	}
	if price < dishPriceMin || price > dishPriceMax {
		return ErrRangeForPriceInDish
	}
	_, err = GlobalCategory.FindCategoryByID(ctx, categoryId)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	err = svc.repo.CreateDish(ctx, domain.Dish{
		Name:       name,
		ImageURL:   imageUrl,
		Price:      price,
		CategoryID: categoryId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *DishService) FindDishById(ctx *gin.Context, id int64) (domain.Dish, error) {
	d, err := svc.repo.FindDishById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return domain.Dish{}, ErrRecordNotFoundInDish
		} else {
			return domain.Dish{}, err
		}
	}
	return d, nil
}

func (svc *DishService) FindDishByName(ctx *gin.Context, name string) ([]domain.Dish, error) {
	d, err := svc.repo.FindDishByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return nil, ErrRecordNotFoundInDish
		} else {
			return nil, err
		}
	}
	if d == nil {
		return nil, ErrRecordNotFoundInDish
	}
	return d, nil
}

func (svc *DishService) FindDishByCategory(ctx *gin.Context, categoryID int64) ([]domain.Dish, error) {
	d, err := svc.repo.FindDishByCategoryID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return nil, ErrRecordNotFoundInDish
		} else {
			return nil, err
		}
	}
	if d == nil {
		return nil, ErrRecordNotFoundInDish
	}
	return d, nil
}

func (svc *DishService) FindAllDishes(ctx *gin.Context) ([]domain.Dish, error) {
	d, err := svc.repo.FindAllDishes(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return nil, ErrRecordNotFoundInDish
		} else {
			return nil, err
		}
	}
	if d == nil {
		return nil, ErrRecordNotFoundInDish
	}
	return d, nil
}

func (svc *DishService) EditDish(ctx *gin.Context, id int64, name string, imageUrl string, price float64, categoryId int64) error {
	role := sessions.Default(ctx).Get("role").(string)
	if role != "employee" && role != "admin" {
		return ErrUserHasNoPermissionInDish
	}
	ok, err := svc.dishExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForDishNameInDish
	}
	ok, err = svc.imageURLExp.MatchString(imageUrl)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForImageUrlInDish
	}
	if price < dishPriceMin || price > dishPriceMax {
		return ErrRangeForPriceInDish
	}
	_, err = GlobalCategory.FindCategoryByID(ctx, categoryId)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	err = svc.repo.UpdateDish(ctx, domain.Dish{
		Id:         id,
		Name:       name,
		ImageURL:   imageUrl,
		Price:      price,
		CategoryID: categoryId,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return ErrRecordNotFoundInDish
		} else {
			return err
		}
	}
	return nil
}

func (svc *DishService) DeleteDish(ctx *gin.Context, id int64) error {
	role := sessions.Default(ctx).Get("role").(string)
	if role != "employee" && role != "admin" {
		ctx.JSON(http.StatusOK, gin.H{"message": "无权限"})
		return ErrUserHasNoPermissionInDish
	}
	err := svc.repo.DeleteDish(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return ErrRecordNotFoundInDish
		} else {
			return err
		}
	}
	return nil
}