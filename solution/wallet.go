package wallet

import (
	"errors"
	"sync"
)

// Ошибка, возникающая при недостатке средств на кошельке
var errInsufficientFunds = errors.New("недостаточно средств")

// Структура Wallet представляет собой кошелек с балансом и мьютексом для синхронизации
type Wallet struct {
	balance int        // Текущий баланс кошелька
	mu      sync.RWMutex // Мьютекс для безопасного доступа к балансу из нескольких горутин
}

// функция списания
func Withdrawal(wallet *Wallet, amount int) error {
	wallet.mu.Lock() // Блокируем мьютекс для записи
	defer wallet.mu.Unlock() // Освобождаем мьютекс после завершения функции
	// Проверяем, достаточно ли средств для списания
	if wallet.balance-amount < 0 {
		return errInsufficientFunds // Возвращаем ошибку, если средств недостаточно
	}
	wallet.balance -= amount // Уменьшаем баланс на указанную сумму
	// fmt.Printf("Было списано %d монет(ы), текущий баланс: %d монет(ы)\n", amount, wallet.balance)
	return nil // Возвращаем nil, если списание прошло успешно
}

// функция пополнения
func Refill(wallet *Wallet, amount int) {
	wallet.mu.Lock() // Блокируем мьютекс для записи
	defer wallet.mu.Unlock() // Освобождаем мьютекс после завершения функции
	wallet.balance += amount // Увеличиваем баланс на указанную сумму
	// fmt.Printf("Было начисленно %d монет(ы), текущий баланс: %d монет(ы)\n", amount, wallet.balance)
}

// функция получения баланса (стандартная)
func GetBalance(wallet *Wallet) int {
	wallet.mu.RLock() // Блокируем мьютекс для чтения
	defer wallet.mu.RUnlock() // Освобождаем мьютекс после завершения функции
	// time.Sleep(100000 * time.Nanosecond) // Задержка для имитации длительной операции
	return wallet.balance // Возвращаем текущий баланс
}
