package domain

import "time"

type CartItem struct {
	Id         int64
	CustomerID int64
	DishID     int64
	Quantity   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
