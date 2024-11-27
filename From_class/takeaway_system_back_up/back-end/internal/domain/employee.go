package domain

import "time"

type Employee struct {
	Id        int64
	Name      string
	Role      string
	Email     string
	Phone     string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
