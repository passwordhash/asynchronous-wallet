package wallet_test

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var u url.URL

func init() {
	port := "8080"
	if p := os.Getenv("APP_OUT_PORT"); p != "" {
		port = p
	}

	u = url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port),
		Path:   "/api/v1/",
	}
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

		assert.Equal(t, expectedBalance, newBalance, "Balance should be increased by the deposit amount")
	})

	t.Run("Withdraw Operation", func(t *testing.T) {
		initialBalance := getBalance(t, e, walletID)
		amount := int64(500)
		expectedBalance := initialBalance - amount

		mustSuccessOperationReq(t, e, walletID, withdrawOperation, amount)

		newBalance := getBalance(t, e, walletID)

		assert.Equal(t, expectedBalance, newBalance, "Balance should be decreased by the withdrawal amount")
	})

	t.Run("Deposit and Withdraw", func(t *testing.T) {
		initialBalance := getBalance(t, e, walletID)
		depositAmount := int64(2000)
		withdrawAmount := int64(800)
		expectedBalance := initialBalance + depositAmount - withdrawAmount

		mustSuccessOperationReq(t, e, walletID, depositOperation, depositAmount)
		mustSuccessOperationReq(t, e, walletID, withdrawOperation, withdrawAmount)

		newBalance := getBalance(t, e, walletID)

		assert.Equal(t, expectedBalance, newBalance, "Balance should reflect both deposit and withdrawal")
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

func TestBalanceOperation_EdgeCases(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	t.Run("Minimum operation amount", func(t *testing.T) {
		initialBalance := getBalance(t, e, walletID)
		minAmount := int64(1)

		mustSuccessOperationReq(t, e, walletID, depositOperation, minAmount)
		newBalance := getBalance(t, e, walletID)
		assert.Equal(t, initialBalance+minAmount, newBalance, "Balance should be updated with minimum deposit")

		mustSuccessOperationReq(t, e, walletID, withdrawOperation, minAmount)
		newBalance = getBalance(t, e, walletID)
		assert.Equal(t, initialBalance, newBalance, "Balance should be restored after minimum withdrawal")
	})

	// TODO: Add more edge cases as needed as it business logic evolves
}

func TestBalanceOperation_Concurrent(t *testing.T) {
	// Test 1000 rps/sec
	e := httpexpect.Default(t, u.String())

	const numOperationPairs = 500

	depositAmount := int64(11)
	withdrawAmount := int64(10)

	balanceOnStart := getBalance(t, e, walletID)
	difference := atomic.Int64{}
	expectedBalance := balanceOnStart + int64(numOperationPairs)*(depositAmount-withdrawAmount)

	wg := sync.WaitGroup{}

	wg.Add(numOperationPairs * 2) // Each pair consists of a deposit and a withdraw operation
	for range numOperationPairs {
		go func() {
			defer wg.Done()
			mustSuccessOperationReq(t, e, walletID, depositOperation, depositAmount)
			difference.Add(depositAmount)
		}()

		go func() {
			defer wg.Done()
			mustSuccessOperationReq(t, e, walletID, withdrawOperation, withdrawAmount)
			difference.Add(-withdrawAmount)
		}()
	}

	wg.Wait()

	require.Equal(t, expectedBalance, balanceOnStart+difference.Load(),
		"Final balance should match expected after concurrent operations")
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

	require.Nil(t, resp.Error, "Operation %s should not return an error", operationType)
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
