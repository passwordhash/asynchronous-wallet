package httpapp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/passwordhash/asynchronous-wallet/internal/config"
	walletHandler "github.com/passwordhash/asynchronous-wallet/internal/handler/api/v1/wallet"
)

type App struct {
	log *slog.Logger

	port         int
	readTimeout  time.Duration
	writeTimeout time.Duration

	server *http.Server
}

func New(
	_ context.Context,
	log *slog.Logger,
	cfg config.HttpConfig,
) *App {
	return &App{
		log: log,

		port:         cfg.Port,
		readTimeout:  cfg.ReadTimeout,
		writeTimeout: cfg.WriteTimeout,
	}
}

// MustRun starts the HTTP server and panics if it fails to start.
func (a *App) MustRun() {
	err := a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic("failed to run HTTP server: " + err.Error())
	}
}

// Run starts the HTTP server and listens on the specified port.
func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	walletHandler := walletHandler.New()
	app := gin.New()

	api := app.Group("/api")
	v1 := api.Group("/api/v1")

	walletHandler.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(a.port),
		Handler:      app,
		ReadTimeout:  a.readTimeout,
		WriteTimeout: a.writeTimeout,
	}
	a.server = srv

	log.Info("Starting HTTP server")

	return srv.ListenAndServe()
}

// Stop gracefully stops the HTTP server.
func (a *App) Stop(ctx context.Context) {
	const op = "httpapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("Stopping HTTP server")

	// Shutdown stops receiving new requests and waits for existing requests to finish.
	if err := a.server.Shutdown(ctx); err != nil {
		log.Error("Failed to gracefully stop HTTP server", slog.Any("error", err))
	} else {
		log.Info("HTTP server stopped gracefully")
	}
}
