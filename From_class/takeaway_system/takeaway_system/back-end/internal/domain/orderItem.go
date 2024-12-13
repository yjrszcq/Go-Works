package domain

import "time"

type OrderItem struct {
	Id           int64
	OrderID      int64
	DishID       int64
	DishName     string
	Quantity     int
	UnitPrice    float64
	TotalPrice   float64
	ReviewStatus string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
