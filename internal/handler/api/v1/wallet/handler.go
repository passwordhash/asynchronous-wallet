package wallet

import (
	"context"

	"github.com/gin-gonic/gin"
)

type WalletOperator interface {
	Deposit(ctx context.Context, walletID string, amount int64) error
	Withdraw(ctx context.Context, walletID string, amount int64) error
}

type BalanceProvider interface {
	Balance(ctx context.Context, walletID string) (int64, error)
}

type Handler struct {
	walletSvc       WalletOperator
	balanceProvider BalanceProvider
}

func New(
	walletSvc WalletOperator,
	balanceProvider BalanceProvider,
) *Handler {
	return &Handler{
		walletSvc:       walletSvc,
		balanceProvider: balanceProvider,
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
