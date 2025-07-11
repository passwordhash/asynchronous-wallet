package app

import (
	"context"
	"log/slog"

	httpApp "github.com/passwordhash/asynchronous-wallet/internal/app/http"
	"github.com/passwordhash/asynchronous-wallet/internal/config"
	walletSvc "github.com/passwordhash/asynchronous-wallet/internal/service/wallet"
	walletRepo "github.com/passwordhash/asynchronous-wallet/internal/storage/postgres/wallet"
	postgresPkg "github.com/passwordhash/asynchronous-wallet/pkg/postgres"
)

type App struct {
	HTTPSrv *httpApp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	pgPool, err := postgresPkg.NewPool(ctx, cfg.PG.DSN())
	if err != nil {
		panic("failed to create postgres pool: " + err.Error())
	}

	walletRepository := walletRepo.New(pgPool)

	walletService := walletSvc.New(
		log.WithGroup("wallet_service"),
		walletRepository,
	)

	httpSrv := httpApp.New(
		ctx,
		log,
		cfg.HTTP,
		walletService,
	)

	return &App{
		HTTPSrv: httpSrv,
	}
}
