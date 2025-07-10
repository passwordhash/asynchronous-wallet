package wallet

import (
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/passwordhash/asynchronous-wallet/internal/entity"
	repoErr "github.com/passwordhash/asynchronous-wallet/internal/storage/errors"
	"github.com/stretchr/testify/require"
)

var walletColumns = []string{"id", "balance", "updated_at", "created_at"}

type mockBehavior func(mock pgxmock.PgxPoolIface)

func setupTest(t *testing.T) (pgxmock.PgxPoolIface, *Repository) {
	t.Helper()

	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	require.NoError(t, err)

	repo := New(mock)

	return mock, repo
}

func TestOperation(t *testing.T) {
	t.Parallel()

	const getQuery = `SELECT.*FROM wallets WHERE id = \$1 FOR UPDATE`
	const updateQuery = `UPDATE wallets SET balance = \$1, updated_at = NOW\(\) WHERE id = \$2`

	tests := []struct {
		name          string
		walletID      string
		amount        int64
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:     "Deposit",
			walletID: "test-wallet-id",
			amount:   100,
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(getQuery).
					WithArgs("test-wallet-id").
					WillReturnRows(pgxmock.NewRows(walletColumns).
						AddRow("test-wallet-id", 100, time.Time{}, time.Time{}))
				mock.ExpectExec(updateQuery).
					WithArgs(int64(200), "test-wallet-id").
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:     "Withdraw",
			walletID: "test-wallet-id",
			amount:   -50,
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(getQuery).
					WithArgs("test-wallet-id").
					WillReturnRows(pgxmock.NewRows(walletColumns).
						AddRow("test-wallet-id", 100, time.Time{}, time.Time{}))
				mock.ExpectExec(updateQuery).
					WithArgs(int64(50), "test-wallet-id").
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:     "WalletNotFound",
			walletID: "non-existent-wallet-id",
			amount:   50,
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(getQuery).
					WithArgs("non-existent-wallet-id").
					WillReturnRows(pgxmock.NewRows(walletColumns))
				mock.ExpectRollback()
			},
			expectedError: repoErr.ErrWalletNotFound,
		},
		{
			name:     "UpdateError",
			walletID: "test-wallet-id",
			amount:   50,
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(getQuery).
					WithArgs("test-wallet-id").
					WillReturnRows(pgxmock.NewRows(walletColumns).
						AddRow("test-wallet-id", 100, time.Time{}, time.Time{}))
				mock.ExpectExec(updateQuery).
					WithArgs(int64(150), "test-wallet-id").
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, repo := setupTest(t)

			tt.mockBehavior(mock)

			err := repo.Operation(t.Context(), tt.walletID, tt.amount)

			require.NoError(t, mock.ExpectationsWereMet(), "expectations were not met")
			if tt.expectedError == nil {
				require.NoError(t, err, "expected no error")
			} else {
				require.ErrorIs(t, err, tt.expectedError, "expected error to match")
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	const query = `SELECT.*FROM wallets WHERE id = \$1$`

	tests := []struct {
		name           string
		walletID       string
		mockBehavior   mockBehavior
		expectedWallet *entity.Wallet
		expectedError  error
	}{
		{
			name:     "Ok",
			walletID: "test-wallet-id",
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(query).
					WithArgs("test-wallet-id").
					WillReturnRows(pgxmock.NewRows(walletColumns).
						AddRow("test-wallet-id", 100, time.Time{}, time.Time{}))
			},
			expectedWallet: &entity.Wallet{
				ID:        "test-wallet-id",
				Balance:   100,
				UpdatedAt: time.Time{},
				CreatedAt: time.Time{},
			},
		},
		{
			name:     "NotFound",
			walletID: "non-existent-wallet-id",
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(query).
					WithArgs("non-existent-wallet-id").
					WillReturnRows(pgxmock.NewRows(walletColumns))
			},
			expectedWallet: nil,
			expectedError:  repoErr.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, repo := setupTest(t)

			tt.mockBehavior(mock)

			wallet, err := repo.GetByID(t.Context(), tt.walletID)

			require.NoError(t, mock.ExpectationsWereMet(), "expectations were not met")
			if tt.expectedError == nil {
				require.NotNil(t, wallet, "expected wallet to be returned")
				require.Equal(t, tt.expectedWallet, wallet, "expected wallet to match")
			} else {
				require.ErrorIs(t, err, tt.expectedError, "expected error to match")
				require.Nil(t, wallet, "expected wallet to be nil")
			}
		})
	}
}
