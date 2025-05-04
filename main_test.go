package main

import (
	wallet "concurrency/solution" // Импортируем пакет wallet из папки solution
	"reflect"
	"sync"
	"testing"
)

// isRWMutexUsed проверяет, содержит ли переданная структура или указатель на структуру поле типа sync.RWMutex
func isRWMutexUsed(v interface{}) bool {
	typ := reflect.TypeOf(v) // Получаем тип переданного значения

	// Если передан указатель, получаем тип элемента (структуры), чтобы проверить её поля
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Проверяем, что тип — структура
	if typ.Kind() == reflect.Struct {
		// Перебираем все поля структуры
		for i := 0; i < typ.NumField(); i++ {
			// Если тип поля совпадает с "sync.RWMutex", возвращаем true
			if typ.Field(i).Type.String() == "sync.RWMutex" {
				return true
			}
		}
	}

	// Если поле sync.RWMutex не найдено, возвращаем false
	return false
}

// TestIsRWMutexUsed данный тест проверяет, что структура Wallet содержит поле sync.RWMutex
func TestIsRWMutexUsed(t *testing.T) {
	// Arrange: подготовить тестовые данные — создаём указатель на Wallet
	userWallet := &wallet.Wallet{}
	// Act: выполняем действие — вызываем функцию isRWMutexUsed, передавая указатель
	result := isRWMutexUsed(userWallet)
	// Assert: проверяем результат — ожидаем, что поле sync.RWMutex есть (true)
	if !result {
		t.Errorf("expected Wallet to have sync.RWMutex, but it is missing")
	}
}

// TestRefill тестирует пополнение кошелька
func TestRefill(t *testing.T) {
	// Arrange: Подготовка тестовых данных
	userWallet := &wallet.Wallet{} // Создаем новый экземпляр кошелька
	var wg sync.WaitGroup          // Создаем WaitGroup для ожидания завершения горутин
	wg.Add(100000)                 // Устанавливаем количество горутин, которые мы собираемся запустить

	// Act: Выполнение тестируемого действия
	for i := 0; i < 100000; i++ {
		go func() {
			defer wg.Done()              // Уменьшаем счетчик WaitGroup после завершения горутины
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

// TestWithdrawal тестирует списание средств из кошелька
func TestWithdrawal(t *testing.T) {
	// Arrange: Подготовка тестовых данных
	userWallet := &wallet.Wallet{}   // Создаем новый экземпляр кошелька
	wallet.Refill(userWallet, 10000) // Пополняем кошелек на 10000
	var wg sync.WaitGroup            // Создаем WaitGroup для ожидания завершения горутин
	wg.Add(5000)                     // Устанавливаем количество горутин, которые мы собираемся запустить

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
