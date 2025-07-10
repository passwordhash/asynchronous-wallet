package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/passwordhash/asynchronous-wallet/internal/app"
	"github.com/passwordhash/asynchronous-wallet/internal/config"
)

// TODO: Logging middleware
// TODO: Pass request ID to the logger (e.g. using context)
// TODO: Stress tests (1000 rps/sec for 1 minute)
// TODO: Swagger documentation

const shutdownTimeout = 5 * time.Second // max time to wait for graceful shutdown

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.App.Env)

	application := app.New(ctx, log, cfg)

	go application.HTTPSrv.MustRun()

	<-ctx.Done()

	log.Info("received signal stop signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	application.HTTPSrv.Stop(shutdownCtx)

	log.Info("application stopped gracefully")
}
