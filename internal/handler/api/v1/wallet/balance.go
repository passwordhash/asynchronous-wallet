package wallet

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/passwordhash/asynchronous-wallet/internal/handler/api/v1/response"
	svcErr "github.com/passwordhash/asynchronous-wallet/internal/service/errors"
)

type balanceReq struct {
	WalletID string `uri:"id" binding:"required"`
}

type balanceResp struct {
	WalletID string `json:"walletId"`
	Balance  int64  `json:"balance"`
}

func (h *Handler) balance(c *gin.Context) {
	var req balanceReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	amount, err := h.walletSvc.Balance(c.Request.Context(), req.WalletID)
	if errors.Is(err, svcErr.ErrWalletNotFound) {
		response.NotFound(c, "Wallet not found")
	}
	if err != nil {
		response.BadRequest(c, response.ErrCodeInvalidRequest, "Failed to retrieve balance", err.Error())
		return
	}

	response.Success(c, 200, balanceResp{
		WalletID: req.WalletID,
		Balance:  amount,
	})
}
