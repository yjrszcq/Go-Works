package service

import (
	"back-end/internal/domain"
	"back-end/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrDishInCartNotFoundInCartItem  = repository.ErrCartItemForeignKeyDishConstraintFail
	ErrRecordIsEmptyInCartItem       = errors.New("列表为空")
	ErrRecordNotFoundInCartItem      = repository.ErrCartItemNotFound
	ErrUserHasNoPermissionInCartItem = errors.New("无权限")
	ErrFormatForQuantityInCartItem   = errors.New("数量应大于0, 小于等于99")
)

type CartItemService struct {
	repo *repository.CartItemRepository
}

func NewCartItemService(repo *repository.CartItemRepository) *CartItemService {
	return &CartItemService{
		repo: repo,
	}
}

func (svc *CartItemService) AddCartItem(ctx *gin.Context, dishID int64, quantity int) error {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInCartItem
	}
	if quantity <= 0 || quantity > 99 {
		return ErrFormatForQuantityInCartItem
	}
	cartItem, err := svc.repo.FindCartItemByCustomerIDAndDishID(ctx, id, dishID)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			err = svc.repo.CreateCartItem(ctx, domain.CartItem{
				CustomerID: id,
				DishID:     dishID,
				Quantity:   quantity,
			})
			if err != nil {
				if errors.Is(err, repository.ErrCartItemForeignKeyDishConstraintFail) {
					return ErrDishInCartNotFoundInCartItem
				} else {
					return err
				}
			}
			return nil
		} else {
			return err
		}
	} else {
		err = svc.repo.UpdateCartItemQuantity(ctx, domain.CartItem{
			Id:       cartItem.Id,
			Quantity: quantity + cartItem.Quantity,
		})
		return nil
	}
}

func (svc *CartItemService) FindCartItemByID(ctx *gin.Context, id int64) (domain.CartItem, error) {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return domain.CartItem{}, ErrUserHasNoPermissionInCartItem
	}
	c, err := svc.repo.FindCartItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return domain.CartItem{}, ErrRecordNotFoundInCartItem
		} else {
			return domain.CartItem{}, err
		}
	}
	if customerId != c.CustomerID {
		return domain.CartItem{}, ErrRecordNotFoundInCartItem
	}
	dish, err := GlobalDish.FindDishById(ctx, c.DishID)
	if err == nil {
		c.DishName = dish.Name
	}
	return c, nil
}

func (svc *CartItemService) FindCartItemsByCustomerID(ctx *gin.Context) ([]domain.CartItem, error) {
	id, err := getCurrentCustomerId(ctx)
	if err != nil {
		return nil, ErrUserHasNoPermissionInCartItem
	}
	c, err := svc.repo.FindCartItemByCustomerID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return nil, ErrRecordIsEmptyInCartItem
		}
		return nil, err
	}
	if c == nil {
		return nil, ErrRecordIsEmptyInCartItem
	}
	for i := 0; i < len(c); i++ {
		dish, err := GlobalDish.FindDishById(ctx, c[i].DishID)
		if err == nil {
			c[i].DishName = dish.Name
		}
	}
	return c, nil
}

func (svc *CartItemService) UpdateCartItem(ctx *gin.Context, id int64, quantity int) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInCartItem
	}
	if quantity <= 0 || quantity > 99 {
		return ErrFormatForQuantityInCartItem
	}
	cartItem, err := svc.repo.FindCartItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	if customerId != cartItem.CustomerID {
		return ErrUserHasNoPermissionInCartItem
	}
	err = svc.repo.UpdateCartItemQuantity(ctx, domain.CartItem{
		Id:       id,
		Quantity: quantity,
	})
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	return nil
}

func (svc *CartItemService) DeleteCartItem(ctx *gin.Context, id int64) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInCartItem
	}
	cartItem, err := svc.repo.FindCartItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	if customerId != cartItem.CustomerID {
		return ErrUserHasNoPermissionInCartItem
	}
	err = svc.repo.DeleteCartItem(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	return nil
}

func (svc *CartItemService) DeleteCartItemsByCustomerID(ctx *gin.Context) error {
	customerId, err := getCurrentCustomerId(ctx)
	if err != nil {
		return ErrUserHasNoPermissionInCartItem
	}
	err = svc.repo.DeleteCartItemsByCustomerID(ctx, customerId)
	if err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			return ErrRecordNotFoundInCartItem
		} else {
			return err
		}
	}
	return nil
}
