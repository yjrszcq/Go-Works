package domain

import "time"

type OrderStatusHistory struct {
	HistoryID   int64
	OrderID     int64
	Status      string
	ChangedAt   time.Time
	ChangedByID int64
}
