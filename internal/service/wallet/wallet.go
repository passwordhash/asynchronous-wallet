package wallet

import (
	"context"
	"log/slog"
)

type WalletRepository interface {
	Deposit(ctx context.Context, walletID string, amount int64) error
	Withdraw(ctx context.Context, walletID string, amount int64) error
	Balance(ctx context.Context, walletID string) (int64, error)
}

type Service struct {
	log  *slog.Logger
	repo WalletRepository
}

func New(
	log *slog.Logger,
	repo WalletRepository,
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

	err := s.repo.Deposit(ctx, walletID, amount)
	// TODO: handler specific errors
	if err != nil {
		log.Error("failed to deposit", "err", err)

		return err
	}

	return nil
}

func (s *Service) Withdraw(ctx context.Context, walletID string, amount int64) error {
	const op = "service.wallet.Withdraw"

	log := s.log.With(
		"op", op,
		"walletID", walletID,
		"amount", amount,
	)

	err := s.repo.Withdraw(ctx, walletID, amount)
	// TODO: handler specific errors
	if err != nil {
		log.Error("failed to withdraw", "err", err)

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

	balance, err := s.repo.Balance(ctx, walletID)
	if err != nil {
		log.Error("failed to get balance", "err", err)

		return 0, err
	}

	log.Info("balance retrieved successfully")

	return balance, nil
}
