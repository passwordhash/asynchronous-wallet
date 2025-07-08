package wallet

import (
	"context"

	"github.com/gin-gonic/gin"
)

type WalletOperator interface {
	Deposit(ctx context.Context, walletID string, amount int64) error
	Withdraw(ctx context.Context, walletID string, amount int64) error
}

type Handler struct {
	walletSvc WalletOperator
}

func New(walletSvc WalletOperator) *Handler {
	return &Handler{
		walletSvc: walletSvc,
	}
}

func (h *Handler) RegisterRoutes(base *gin.RouterGroup) {
	walletGroup := base.Group("/wallet")
	{
		walletGroup.POST("", h.operation)
	}
}
