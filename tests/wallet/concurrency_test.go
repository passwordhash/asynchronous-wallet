package wallet_test

import (
	"sync"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentOperations(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	// Получаем начальный баланс
	initialBalance := getBalance(t, e, walletID)

	// Количество параллельных операций
	const numOperations = 10

	// Суммы для депозита и снятия
	const depositAmount = 100
	const withdrawAmount = 50

	// Ожидаемое изменение баланса после всех операций
	expectedBalanceChange := numOperations * (depositAmount - withdrawAmount)

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup
	wg.Add(numOperations * 2) // Умножаем на 2, т.к. у нас 2 операции (депозит и снятие) для каждой итерации

	// Запускаем параллельные операции
	for i := 0; i < numOperations; i++ {
		// Запускаем депозит
		go func() {
			defer wg.Done()
			operationReq(e, walletID, depositOperation, depositAmount).
				Status(200).
				JSON().
				Object().
				HasValue("success", true)
		}()

		// Запускаем снятие
		go func() {
			defer wg.Done()
			operationReq(e, walletID, withdrawOperation, withdrawAmount).
				Status(200).
				JSON().
				Object().
				HasValue("success", true)
		}()
	}

	// Ждем завершения всех операций
	wg.Wait()

	// Получаем итоговый баланс
	finalBalance := getBalance(t, e, walletID)

	// Проверяем, что итоговый баланс соответствует ожидаемому
	expectedBalance := initialBalance + int64(expectedBalanceChange)
	assert.Equal(t, expectedBalance, finalBalance,
		"Final balance should be initial balance + expected change after concurrent operations")
}
