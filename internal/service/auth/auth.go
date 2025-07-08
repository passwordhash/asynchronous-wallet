package auth

import (
	"log/slog"
)

type Service struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) *Service {
	return &Service{
		log: log,
	}
}
