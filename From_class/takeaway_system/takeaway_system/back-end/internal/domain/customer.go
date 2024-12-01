package domain

import "time"

type Customer struct {
	Id        int64
	Name      string
	Email     string
	Password  string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
