package wallet_test

import (
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var u = url.URL{
	Scheme: "http",
	Host:   "localhost:8080",
	Path:   "/api/v1/",
}

const (
	walletID = "11111111-2b2b-4c4c-8d8d-0e0e1f2a3b4c" // Existing wallet ID for testing

	depositOperation  = "deposit"
	withdrawOperation = "withdraw"
)

type operationResp struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
	Error   any  `json:"error,omitempty"`
}

func TestBalanceOperation_Ok(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	t.Run("Deposit Operation", func(t *testing.T) {
		initialBalance := getBalance(t, e, walletID)
		amount := int64(1000)
		expectedBalance := initialBalance + amount

		mustSuccessOperationReq(t, e, walletID, depositOperation, amount)

		newBalance := getBalance(t, e, walletID)

		if newBalance != expectedBalance {
			t.Errorf("Expected balance %d, got %d", expectedBalance, newBalance)
		}
	})

	t.Run("Withdraw Operation", func(t *testing.T) {
		initialBalance := getBalance(t, e, walletID)
		amount := int64(500)
		expectedBalance := initialBalance - amount

		mustSuccessOperationReq(t, e, walletID, withdrawOperation, amount)

		newBalance := getBalance(t, e, walletID)

		if newBalance != expectedBalance {
			t.Errorf("Expected balance %d, got %d", expectedBalance, newBalance)
		}
	})

	t.Run("Deposit and Withdraw", func(t *testing.T) {
		initialBalance := getBalance(t, e, walletID)
		depositAmount := int64(2000)
		withdrawAmount := int64(800)
		expectedBalance := initialBalance + depositAmount - withdrawAmount

		mustSuccessOperationReq(t, e, walletID, depositOperation, depositAmount)
		mustSuccessOperationReq(t, e, walletID, withdrawOperation, withdrawAmount)

		newBalance := getBalance(t, e, walletID)

		if newBalance != expectedBalance {
			t.Errorf("Expected balance %d, got %d", expectedBalance, newBalance)
		}
	})
}

func TestBalanceOperation_Error(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	t.Run("Deposit operation with negative amount", func(t *testing.T) {
		t.Parallel()

		operationReq(e, walletID, depositOperation, -1000).
			Status(400).
			JSON().
			Object().
			HasValue("success", false).
			ContainsKey("error")
	})

	t.Run("Withdraw operation with negative amount", func(t *testing.T) {
		t.Parallel()

		operationReq(e, walletID, withdrawOperation, -500).
			Status(400).
			JSON().
			Object().
			HasValue("success", false).
			ContainsKey("error")
	})

	t.Run("Wallet not found", func(t *testing.T) {
		t.Parallel()

		nonExistentWalletID := "00000000-0000-0000-0000-000000000000"

		operationReq(e, nonExistentWalletID, depositOperation, 1000).
			Status(404).
			JSON().
			Object().
			HasValue("success", false).
			ContainsKey("error")
	})

	t.Run("Zero amount operation", func(t *testing.T) {
		t.Parallel()

		operationReq(e, walletID, depositOperation, 0).
			Status(400).
			JSON().
			Object().
			HasValue("success", false).
			ContainsKey("error")

		operationReq(e, walletID, withdrawOperation, 0).
			Status(400).
			JSON().
			Object().
			HasValue("success", false).
			ContainsKey("error")
	})

	t.Run("Unsupported operation type", func(t *testing.T) {
		t.Parallel()

		operationReq(e, walletID, "unsupported_operation", 1000).
			Status(400).
			JSON().
			Object().
			HasValue("success", false).
			ContainsKey("error")
	})
}

func mustSuccessOperationReq(t *testing.T, e *httpexpect.Expect, walletID, operationType string, amount int64) {
	var resp operationResp
	operationReq(e, walletID, operationType, amount).
		Status(200).
		JSON().
		Object().
		HasValue("success", true).
		NotContainsKey("error").
		Decode(&resp)

	if resp.Error != nil {
		t.Fatalf("Operation %s returned error: %v", operationType, resp.Error)
	}
}

func operationReq(e *httpexpect.Expect, walletID, operationType string, amount int64) *httpexpect.Response {
	type operationReq struct {
		WalletID      string `json:"walletId"`
		OperationType string `json:"operationType"`
		Amount        int64  `json:"amount"`
	}

	return e.POST("/wallet", walletID).
		WithJSON(operationReq{
			WalletID:      walletID,
			OperationType: operationType,
			Amount:        amount,
		}).
		Expect()
}

func getBalance(t *testing.T, e *httpexpect.Expect, walletID string) int64 {
	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Balance int64 `json:"balance"`
		} `json:"data,omitempty"`
		Error any `json:"error,omitempty"`
	}

	e.GET("/wallets/{id}", walletID).
		Expect().
		Status(200).
		JSON().
		Object().
		HasValue("success", true).
		NotContainsKey("error").
		Decode(&resp)
	return resp.Data.Balance
}
