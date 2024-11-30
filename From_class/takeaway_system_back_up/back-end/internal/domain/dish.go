package domain

import "time"

type Dish struct {
	Id         int64
	Name       string
	ImageURL   string
	Price      float64
	CategoryID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
