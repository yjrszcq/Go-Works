package repository

import (
	"back-end/internal/domain"
	"back-end/internal/repository/dao"
	"context"
)

var (
	ErrCartItemForeignKeyDishConstraintFail = dao.ErrCartItemForeignKeyDishConstraintFail
	ErrCartItemNotFound                     = dao.ErrCartItemNotFound
)

type CartItemRepository struct {
	dao *dao.CartItemDAO
}

func NewCartItemRepository(dao *dao.CartItemDAO) *CartItemRepository {
	return &CartItemRepository{
		dao: dao,
	}
}

func (r *CartItemRepository) CreateCartItem(ctx context.Context, cartItem domain.CartItem) error {
	return r.dao.InsertCartItem(ctx, dao.CartItem{
		CustomerID: cartItem.CustomerID,
		DishID:     cartItem.DishID,
		Quantity:   cartItem.Quantity,
	})
}

func (r *CartItemRepository) FindCartItemByID(ctx context.Context, id int64) (domain.CartItem, error) {
	c, err := r.dao.FindCartItemByID(ctx, id)
	if err != nil {
		return domain.CartItem{}, err
	}
	return domain.CartItem{
		Id:         c.CartItemID,
		CustomerID: c.CustomerID,
		DishID:     c.DishID,
		Quantity:   c.Quantity,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}, nil
}

func (r *CartItemRepository) FindCartItemByCustomerIDAndDishID(ctx context.Context, customerID, dishID int64) (domain.CartItem, error) {
	c, err := r.dao.FindCartItemByCustomerIDAndDishID(ctx, customerID, dishID)
	if err != nil {
		return domain.CartItem{}, err
	}
	return domain.CartItem{
		Id:         c.CartItemID,
		CustomerID: c.CustomerID,
		DishID:     c.DishID,
		Quantity:   c.Quantity,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}, nil
}

func (r *CartItemRepository) FindCartItemByCustomerID(ctx context.Context, customerID int64) ([]domain.CartItem, error) {
	cartItems, err := r.dao.FindCartItemsByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	var res []domain.CartItem
	for _, c := range cartItems {
		res = append(res, domain.CartItem{
			Id:         c.CartItemID,
			CustomerID: c.CustomerID,
			DishID:     c.DishID,
			Quantity:   c.Quantity,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		})
	}
	return res, nil
}

func (r *CartItemRepository) UpdateCartItemQuantity(ctx context.Context, c domain.CartItem) error {
	return r.dao.UpdateCartItemQuantity(ctx, dao.CartItem{
		CartItemID: c.Id,
		Quantity:   c.Quantity,
	})
}

func (r *CartItemRepository) DeleteCartItem(ctx context.Context, id int64) error {
	return r.dao.DeleteCartItem(ctx, id)
}

func (r *CartItemRepository) DeleteCartItemsByCustomerID(ctx context.Context, customerID int64) error {
	return r.dao.DeleteCartItemsByCustomerID(ctx, customerID)
}

func (r *CartItemRepository) DeleteCartItemsByDishID(ctx context.Context, dishID int64) error {
	return r.dao.DeleteCartItemsByDishID(ctx, dishID)
}

func (r *CartItemRepository) DeleteCartItemByCustomerIDAndDishID(ctx context.Context, customerID, dishID int64) error {
	return r.dao.DeleteCartItemByCustomerIDAndDishID(ctx, customerID, dishID)
}
