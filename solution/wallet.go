package wallet

import (
	"errors"
	"sync"
)

// ошибка, возникающая при недостатке средств на кошельке
var errInsufficientFunds = errors.New("недостаточно средств")

// структура Кошелёк
type Wallet struct {
	balance int   // банас кошелька
	log     []int // лог операций кошелька
	mu      sync.RWMutex
}

// функция списания
func (wallet *Wallet) Withdrawal(amount int) error {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	// проверка, достаточно ли средств для списания
	if wallet.balance-amount < 0 {
		return errInsufficientFunds
	}
	wallet.balance -= amount
	wallet.log = append(wallet.log, -amount)
	return nil
}

// функция пополнения
func (wallet *Wallet) Refill(amount int) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	wallet.balance += amount
	wallet.log = append(wallet.log, amount)
}

// функция получения баланса
func (wallet *Wallet) GetBalance() int {
	wallet.mu.RLock()
	defer wallet.mu.RUnlock()
	return wallet.balance
}

func (wallet *Wallet) GetLog() []int {
	wallet.mu.RLock()
	defer wallet.mu.RUnlock()
	return wallet.log
}
