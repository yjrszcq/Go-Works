package domain

import "time"

type Review struct {
	Id           int64
	CustomerID   int64
	CustomerName string
	DishID       int64
	Rating       int
	Comment      string
	ReviewDate   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
