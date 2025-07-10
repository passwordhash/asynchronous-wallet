package model

import (
	"time"

	"github.com/passwordhash/asynchronous-wallet/internal/entity"
)

type Wallet struct {
	ID        string    `db:"id"`
	Balance   int64     `db:"balance"`
	UpdatedAt time.Time `db:"updated_at"`
	CreateAt  time.Time `db:"created_at"`
}

func (w Wallet) ToEntity() *entity.Wallet {
	return &entity.Wallet{
		ID:        w.ID,
		Balance:   w.Balance,
		UpdatedAt: w.UpdatedAt,
		CreatedAt: w.CreateAt,
	}
}
