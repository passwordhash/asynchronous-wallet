package model

import (
	"database/sql"
	"time"

	"github.com/passwordhash/asynchronous-wallet/internal/entity"
)

type Wallet struct {
	ID        string       `db:"id"`
	Balance   int64        `db:"balance"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	CreateAt  time.Time    `db:"created_at"`
}

func (w Wallet) ToEntity() *entity.Wallet {
	return &entity.Wallet{
		ID:        w.ID,
		Balance:   w.Balance,
		UpdatedAt: w.UpdatedAt.Time,
		CreateAt:  w.CreateAt,
	}
}
