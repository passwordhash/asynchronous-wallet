package entity

import "time"

type Wallet struct {
	ID        string
	Balance   int64
	UpdatedAt time.Time
	CreatedAt time.Time
}
