package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateNameInCategory        = repository.ErrCategoryDuplicateName
	ErrRecordIsEmptyInCategory        = errors.New("列表为空")
	ErrRecordNotFoundInCategory       = repository.ErrCategoryNotFound
	ErrUserHasNoPermissionInCategory  = errors.New("无权限")
	ErrFormatForNameInCategory        = errors.New("分类名称应小于20个字符")
	ErrFormatForDescriptionInCategory = errors.New("分类描述应小于200个字符")
)

const (
	categoryNameRegexPattern = `^[a-zA-Z0-9\u4e00-\u9fa5 ]{1,20}$`
	categoryDescRegexPattern = `^[a-zA-Z0-9\u4e00-\u9fa5,\.\u3002\uff0c ]{0,200}$`
)

type CategoryService struct {
	repo    *repository.CategoryRepository
	nameExp *regexp.Regexp
	descExp *regexp.Regexp
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	nameExp := regexp.MustCompile(categoryNameRegexPattern, regexp.None)
	descExp := regexp.MustCompile(categoryDescRegexPattern, regexp.None)
	return &CategoryService{
		repo:    repo,
		nameExp: nameExp,
		descExp: descExp,
	}
}

func (svc *CategoryService) CreateCategory(ctx *gin.Context, name string, description string) error {
	role := sessions.Default(ctx).Get("role")
	if role != "employee" && role != "admin" {
		return ErrUserHasNoPermissionInCategory
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInCategory
	}
	ok, err = svc.descExp.MatchString(description)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForDescriptionInCategory
	}
	err = svc.repo.CreateCategory(ctx, domain.Category{
		Name:        name,
		Description: description,
	})
	if err != nil {
		if errors.Is(err, repository.ErrCategoryDuplicateName) {
			return ErrDuplicateNameInCategory
		} else {
			return err
		}
	}
	return nil
}

func (svc *CategoryService) FindCategoryByID(ctx *gin.Context, id int64) (domain.Category, error) {
	c, err := svc.repo.FindCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return domain.Category{}, ErrRecordNotFoundInCategory
		} else {
			return domain.Category{}, err
		}
	}
	return c, nil
}

func (svc *CategoryService) FindCategoryByName(ctx *gin.Context, name string) (domain.Category, error) {
	c, err := svc.repo.FindCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return domain.Category{}, ErrRecordNotFoundInCategory
		} else {
			return domain.Category{}, err
		}
	}
	return c, nil
}

func (svc *CategoryService) FindAllCategories(ctx *gin.Context) ([]domain.Category, error) {
	c, err := svc.repo.FindAllCategories(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return nil, ErrRecordIsEmptyInCategory
		} else {
			return nil, err
		}
	}
	if c == nil {
		return nil, ErrRecordIsEmptyInCategory
	}
	return c, nil
}

func (svc *CategoryService) UpdateCategory(ctx *gin.Context, id int64, name string, description string) error {
	role := sessions.Default(ctx).Get("role")
	if role != "employee" && role != "admin" {
		return ErrUserHasNoPermissionInCategory
	}
	ok, err := svc.nameExp.MatchString(name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForNameInCategory
	}
	ok, err = svc.descExp.MatchString(description)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFormatForDescriptionInCategory
	}
	err = svc.repo.UpdateCategory(ctx, domain.Category{
		Id:          id,
		Name:        name,
		Description: description,
	})
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return ErrRecordNotFoundInCategory
		} else {
			return err
		}
	}
	return nil
}

func (svc *CategoryService) DeleteCategory(ctx *gin.Context, id int64) error {
	role := sessions.Default(ctx).Get("role")
	if role != "employee" && role != "admin" {
		return ErrUserHasNoPermissionInCategory
	}
	err := svc.repo.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return ErrRecordNotFoundInCategory
		} else {
			return err
		}
	}
	return nil
}
