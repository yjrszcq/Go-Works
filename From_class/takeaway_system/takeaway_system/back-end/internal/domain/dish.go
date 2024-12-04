package domain

import "time"

type Dish struct {
	Id          int64
	Name        string
	ImageURL    string
	Price       float64
	Description string
	CategoryID  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
