package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/passwordhash/asynchronous-wallet/internal/handler/api/v1/response"
)

const (
	depositOperation  = "deposit"
	withdrawOperation = "withdraw"
)

type operationReq struct {
	WalletID      string `json:"walletId" binding:"required,uuid"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	OperationType string `json:"operationType" binding:"required,oneof=deposit withdraw"`
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
		if err := h.walletSvc.Deposit(c.Request.Context(), req.WalletID, req.Amount); err != nil {
			response.BadRequest(c, response.ErrCodeInvalidRequest, "Failed to deposit amount", err.Error())
			return
		}
		response.Success(c, 200, operationResp{Message: "Deposit successful"})
	case withdrawOperation:
		if err := h.walletSvc.Withdraw(c.Request.Context(), req.WalletID, req.Amount); err != nil {
			response.BadRequest(c, response.ErrCodeInvalidRequest, "Failed to withdraw amount", err.Error())
			return
		}
		response.Success(c, 200, operationResp{Message: "Withdrawal successful"})
	default:
		response.BadRequest(c, response.ErrCodeInvalidRequest, "Invalid operation type", "Must be either 'deposit' or 'withdraw'")
		return
	}
}
