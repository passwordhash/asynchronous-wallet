package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/passwordhash/asynchronous-wallet/internal/entity"
	repoErr "github.com/passwordhash/asynchronous-wallet/internal/storage/errors"
	"github.com/passwordhash/asynchronous-wallet/internal/storage/postgres/wallet/model"
	postgresPkg "github.com/passwordhash/asynchronous-wallet/pkg/postgres"
)

type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type Repository struct {
	db DB
}

func New(db DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Operation is a method that performs a deposit or withdrawal operation on a wallet.
// If amount is positive, it performs a deposit; if negative, it performs a withdrawal.
// If wallet with the given ID does not exist, it returns [repoErr.ErrWalletNotFound].
// Note: If balance becomes negative after the operation, it will still be updated.
// Safe for concurrent use, as it uses a transaction to ensure atomicity.
func (r *Repository) Operation(ctx context.Context, walletID string, amount int64) (err error) {
	const op = "repository.wallet.Deposit"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("%s: rollback failed: %v, original error: %w", op, rbErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("%s: commit failed: %w", op, commitErr)
			}
		}
	}()

	wallet, err := r.getByID(ctx, tx, walletID, true)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	balance := wallet.Balance + amount

	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(ctx, query, balance, walletID)
	if err != nil {
		return fmt.Errorf("%s: failed to update wallet balance: %w", op, err)
	}

	return nil
}

// GetByID is a method that retrieves a wallet by its ID.
// If the wallet is not found, it returns [repoErr.ErrWalletNotFound].
func (r *Repository) GetByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	const op = "repository.wallet.Balance"

	wallet, err := r.getByID(ctx, r.db, walletID, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return wallet, nil
}

// getByID is a helper method that retrieves a wallet by its ID.
// If the wallet is not found, it returns [repoErr.ErrWalletNotFound].
// If the query is executed within a transaction, it locks the row for update.
func (r *Repository) getByID(
	ctx context.Context,
	q postgresPkg.Queryer,
	walletID string,
	isForUpdating bool,
) (*entity.Wallet, error) {
	query := `SELECT * FROM wallets WHERE id = $1`

	if isForUpdating {
		query += " FOR UPDATE"
	}

	rows, err := q.Query(ctx, query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallet, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Wallet])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repoErr.ErrWalletNotFound
	}
	if err != nil {
		return nil, err
	}

	return wallet.ToEntity(), nil
}
