package wallet

import (
	"context"

	"github.com/gin-gonic/gin"
)

type WalletService interface {
	Deposit(ctx context.Context, walletID string, amount int64) error
	Withdraw(ctx context.Context, walletID string, amount int64) error
	Balance(ctx context.Context, walletID string) (int64, error)
}

type Handler struct {
	walletSvc WalletService
}

func New(
	walletSvc WalletService,
) *Handler {
	return &Handler{
		walletSvc: walletSvc,
	}
}

func (h *Handler) RegisterRoutes(base *gin.RouterGroup) {
	walletGroup := base.Group("/wallet")
	{
		walletGroup.POST("", h.operation)
	}

	walletsGroup := base.Group("/wallets")
	{
		walletIDGroup := walletsGroup.Group("/:id")
		{
			walletIDGroup.GET("", h.balance)
		}
	}
}
