package wallet

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/passwordhash/asynchronous-wallet/internal/entity"
	svcErr "github.com/passwordhash/asynchronous-wallet/internal/service/errors"
	repoErr "github.com/passwordhash/asynchronous-wallet/internal/storage/errors"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks github.com/passwordhash/asynchronous-wallet/internal/service/wallet Repository
type Repository interface {
	Operation(ctx context.Context, walletID string, amount int64) error
	GetByID(ctx context.Context, walletID string) (*entity.Wallet, error)
}

type Service struct {
	log  *slog.Logger
	repo Repository
}

func New(
	log *slog.Logger,
	repo Repository,
) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) Deposit(ctx context.Context, walletID string, amount int64) error {
	const op = "service.wallet.Deposit"

	log := s.log.With(
		"op", op,
		"walletID", walletID,
		"amount", amount,
	)

	if err := validate(walletID, amount); err != nil {
		log.Error("invalid parameters", "err", err)

		return svcErr.ErrInvalidParams
	}

	err := s.repo.Operation(ctx, walletID, amount)
	if errors.Is(err, repoErr.ErrWalletNotFound) {
		log.Warn("wallet not found", "err", err)

		return svcErr.ErrWalletNotFound
	}
	if err != nil {
		log.Error("failed to update balance", "err", err)

		return err
	}

	log.Info("deposit successful")

	return nil
}

func (s *Service) Withdraw(ctx context.Context, walletID string, amount int64) error {
	const op = "service.wallet.Withdraw"

	log := s.log.With(
		"op", op,
		"walletID", walletID,
		"amount", amount,
	)

	if err := validate(walletID, amount); err != nil {
		log.Error("invalid parameters", "err", err)

		return svcErr.ErrInvalidParams
	}

	err := s.repo.Operation(ctx, walletID, -amount)
	if errors.Is(err, repoErr.ErrWalletNotFound) {
		log.Warn("wallet not found", "err", err)

		return svcErr.ErrWalletNotFound
	}
	if err != nil {
		log.Error("failed to update balance", "err", err)

		return err
	}

	log.Info("withdrawal successful")

	return nil
}

func (s *Service) Balance(ctx context.Context, walletID string) (int64, error) {
	const op = "service.wallet.Balance"

	log := s.log.With(
		"op", op,
		"walletID", walletID,
	)

	if uuid.Validate(walletID) != nil {
		log.Warn("invalid wallet ID format", "walletID", walletID)

		return 0, svcErr.ErrInvalidParams
	}

	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		log.Error("failed to get balance", "err", err)

		return 0, err
	}

	log.Info("wallet balance retrieved")

	return wallet.Balance, nil
}

func validate(walletID string, amount int64) error {
	if uuid.Validate(walletID) != nil {
		return svcErr.ErrInvalidParams
	}
	if amount <= 0 {
		return svcErr.ErrInvalidParams
	}
	return nil
}
