package domain

import "time"

type Order struct {
	Id               int64
	CustomerID       int64
	OrderDate        time.Time
	DeliveryTime     time.Time
	DeliveryLocation string
	Status           string
	PaymentStatus    string
	TotalAmount      float64
	DeliveryPersonID int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
