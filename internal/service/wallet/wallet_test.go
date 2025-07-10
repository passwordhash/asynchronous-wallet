package wallet_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/passwordhash/asynchronous-wallet/internal/entity"
	svcErr "github.com/passwordhash/asynchronous-wallet/internal/service/errors"
	"github.com/passwordhash/asynchronous-wallet/internal/service/wallet"
	"github.com/passwordhash/asynchronous-wallet/internal/service/wallet/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*wallet.Service, *mocks.MockRepository) {
	t.Helper()

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockRepository(ctrl)

	service := wallet.New(log, mockRepo)

	return service, mockRepo
}

func TestDeposit(t *testing.T) {
	t.Parallel()

	validUUID := "11111111-2b2b-4c4c-8d8d-0e0e1f2a3b4c"

	tests := []struct {
		name          string
		walletID      string
		amount        int64
		mockBehavior  func(mock *mocks.MockRepository)
		expectedError error
	}{
		{
			name:     "Ok",
			walletID: validUUID,
			amount:   100,
			mockBehavior: func(mock *mocks.MockRepository) {
				mock.EXPECT().Operation(gomock.Any(), validUUID, int64(100)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Invalid uuid format",
			walletID:      "wallet-id",
			amount:        100,
			mockBehavior:  func(mock *mocks.MockRepository) {},
			expectedError: svcErr.ErrInvalidParams,
		},
		{
			name:          "Amount is zero",
			walletID:      validUUID,
			amount:        0,
			mockBehavior:  func(mock *mocks.MockRepository) {},
			expectedError: svcErr.ErrInvalidParams,
		},
		{
			name:          "Amount is negative",
			walletID:      validUUID,
			amount:        -100,
			mockBehavior:  func(mock *mocks.MockRepository) {},
			expectedError: svcErr.ErrInvalidParams,
		},
		{
			name:     "Wallet not found",
			walletID: validUUID,
			amount:   100,
			mockBehavior: func(mock *mocks.MockRepository) {
				mock.EXPECT().Operation(gomock.Any(), validUUID, int64(100)).Return(svcErr.ErrWalletNotFound)
			},
			expectedError: svcErr.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mockRepo := setupTest(t)

			tt.mockBehavior(mockRepo)

			err := service.Deposit(t.Context(), tt.walletID, tt.amount)

			if tt.expectedError == nil {
				t.Log(err)
				require.NoError(t, err, "expected no error")
			} else {
				require.ErrorIs(t, err, tt.expectedError, "expected error to match")
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	t.Parallel()

	validUUID := "11111111-2b2b-4c4c-8d8d-0e0e1f2a3b4c"

	tests := []struct {
		name          string
		walletID      string
		amount        int64
		mockBehavior  func(mock *mocks.MockRepository)
		expectedError error
	}{
		{
			name:     "Ok",
			walletID: validUUID,
			amount:   100,
			mockBehavior: func(mock *mocks.MockRepository) {
				mock.EXPECT().Operation(gomock.Any(), validUUID, int64(-100)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Invalid uuid format",
			walletID:      "wallet-id",
			amount:        100,
			mockBehavior:  func(mock *mocks.MockRepository) {},
			expectedError: svcErr.ErrInvalidParams,
		},
		{
			name:          "Amount is zero",
			walletID:      validUUID,
			amount:        0,
			mockBehavior:  func(mock *mocks.MockRepository) {},
			expectedError: svcErr.ErrInvalidParams,
		},
		{
			name:          "Amount is negative",
			walletID:      validUUID,
			amount:        -100,
			mockBehavior:  func(mock *mocks.MockRepository) {},
			expectedError: svcErr.ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mockRepo := setupTest(t)

			tt.mockBehavior(mockRepo)

			err := service.Withdraw(t.Context(), tt.walletID, tt.amount)

			if tt.expectedError == nil {
				t.Log(err)
				require.NoError(t, err, "expected no error")
			} else {
				require.ErrorIs(t, err, tt.expectedError, "expected error to match")
			}
		})
	}
}

func TestBalance(t *testing.T) {
	t.Parallel()

	validUUID := "11111111-2b2b-4c4c-8d8d-0e0e1f2a3b4c"

	tests := []struct {
		name            string
		walletID        string
		mockBehavior    func(mock *mocks.MockRepository)
		expectedError   error
		expectedBalance int64
	}{
		{
			name:     "Ok",
			walletID: validUUID,
			mockBehavior: func(mock *mocks.MockRepository) {
				mock.EXPECT().GetByID(gomock.Any(), validUUID).Return(&entity.Wallet{Balance: 100}, nil)
			},
			expectedError:   nil,
			expectedBalance: 100,
		},
		{
			name:            "Invalid uuid format",
			walletID:        "wallet-id",
			mockBehavior:    func(mock *mocks.MockRepository) {},
			expectedError:   svcErr.ErrInvalidParams,
			expectedBalance: 0,
		},
		{
			name:     "Wallet not found",
			walletID: validUUID,
			mockBehavior: func(mock *mocks.MockRepository) {
				mock.EXPECT().GetByID(gomock.Any(), validUUID).Return(nil, svcErr.ErrWalletNotFound)
			},
			expectedError:   svcErr.ErrWalletNotFound,
			expectedBalance: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mockRepo := setupTest(t)

			tt.mockBehavior(mockRepo)

			balance, err := service.Balance(t.Context(), tt.walletID)

			if tt.expectedError == nil {
				t.Log(err)
				require.NoError(t, err, "expected no error")
				require.Equal(t, tt.expectedBalance, balance, "expected balance to match")
			} else {
				require.ErrorIs(t, err, tt.expectedError, "expected error to match")
				require.Equal(t, tt.expectedBalance, balance, "expected balance to be zero on error")
			}
		})
	}
}
