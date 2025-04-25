package main

import (
	"sync"
	"testing"
	"concurrency/wallet"
)
// Пополнения кошелька
func TestRefill(t *testing.T) {
	// Arrange: Подготовка тестовых данных
	userWallet := &wallet.Wallet{} // Создаем новый экземпляр кошелька
	var wg sync.WaitGroup           // Создаем WaitGroup для ожидания завершения горутин
	wg.Add(100000)                  // Устанавливаем количество горутин, которые мы собираемся запустить

	// Act: Выполнение тестируемого действия
	for i := 0; i < 100000; i++ {
		go func() {
			defer wg.Done() // Уменьшаем счетчик WaitGroup после завершения горутины
			wallet.Refill(userWallet, 1) // Пополняем кошелек на 1
		}()
	}

	wg.Wait() // Ожидаем завершения всех горутин

	// Assert: Проверка ожидаемого результата
	expectedBalance := 100000 // Ожидаем, что баланс будет равен 100000
	if got := wallet.GetBalance(userWallet); got != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, got) // Если баланс не совпадает, выводим ошибку
	}
}
//Списания средств из кошелька
func TestWithdrawal(t *testing.T) {
	// Arrange: Подготовка тестовых данных
	userWallet := &wallet.Wallet{} // Создаем новый экземпляр кошелька
	wallet.Refill(userWallet, 10000) // Пополняем кошелек на 10000
	var wg sync.WaitGroup           // Создаем WaitGroup для ожидания завершения горутин
	wg.Add(5000)                    // Устанавливаем количество горутин, которые мы собираемся запустить

	// Act: Выполнение тестируемого действия
	for i := 0; i < 5000; i++ {
		go func() {
			defer wg.Done() // Уменьшаем счетчик WaitGroup после завершения горутины
			if err := wallet.Withdrawal(userWallet, 1); err != nil {
				t.Errorf("unexpected error: %v", err) // Если произошла ошибка, выводим ее
			}
		}()
	}

	wg.Wait() // Ожидаем завершения всех горутин

	// Assert: Проверка ожидаемого результата
	expectedBalance := 10000 - 5000 // Ожидаем, что баланс будет равен 5000 (10000 - 5000)
	if got := wallet.GetBalance(userWallet); got != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, got) // Если баланс не совпадает, выводим ошибку
	}
}