package wallet

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/passwordhash/asynchronous-wallet/internal/handler/api/v1/response"
	svcErr "github.com/passwordhash/asynchronous-wallet/internal/service/errors"
)

const (
	depositOperation  = "deposit"
	withdrawOperation = "withdraw"
)

type operationReq struct {
	WalletID      string `json:"walletId" binding:"required,uuid"`
	OperationType string `json:"operationType" binding:"required,oneof=deposit withdraw"`
	Amount        int64  `json:"amount" binding:"required,min=1,gt=0"`
}

type operationResp struct {
	Message string `json:"message"`
}

func (h *Handler) operation(c *gin.Context) {
	var req operationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	switch req.OperationType {
	case depositOperation:
		err := h.walletSvc.Deposit(c.Request.Context(), req.WalletID, req.Amount)
		if isErr := handleServiceError(c, err); isErr {
			return
		}
		response.Success(c, 200, operationResp{Message: "Deposit successful"})
	case withdrawOperation:
		err := h.walletSvc.Withdraw(c.Request.Context(), req.WalletID, req.Amount)
		if isErr := handleServiceError(c, err); isErr {
			return
		}
		response.Success(c, 200, operationResp{Message: "Withdrawal successful"})
	default:
		response.BadRequest(c, response.ErrCodeInvalidRequest, "Invalid operation type", "Must be either 'deposit' or 'withdraw'")
		return
	}
}

func handleServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, svcErr.ErrInvalidParams):
		response.ValidationError(c, "Invalid parameters provided")
	case errors.Is(err, svcErr.ErrWalletNotFound):
		response.NotFound(c, "Wallet not found")
	default:
		response.InternalError(c, "Internal server error")
	}

	return true
}
