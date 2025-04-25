package wallet

import (
	"errors"
)

// Ошибка, возникающая при недостатке средств на кошельке
var errInsufficientFunds = errors.New("недостаточно средств")

// Структура Wallet представляет собой кошелек с балансом и мьютексом для синхронизации
type Wallet struct {
	balance int // Текущий баланс кошелька
}

// функция списания
func Withdrawal(wallet *Wallet, amount int) error {
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
	wallet.balance += amount // Увеличиваем баланс на указанную сумму
	// fmt.Printf("Было начисленно %d монет(ы), текущий баланс: %d монет(ы)\n", amount, wallet.balance)
}

// функция получения баланса (стандартная)
func GetBalance(wallet *Wallet) int {
	return wallet.balance // Возвращаем текущий баланс
}
