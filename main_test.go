package main

import (
	"concurrency/wallet"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

// проверяет, содержит ли структура поле типа sync.RWMutex
func isRWMutexUsed(v interface{}) bool {
	var typ = reflect.TypeOf(v) // Получаем тип переданного значения
	// если передан указатель, получаем тип элемента (структуры), чтобы проверить её поля
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// проверяем, что тип — структура
	if typ.Kind() == reflect.Struct {
		// перебираем все поля структуры
		for i := 0; i < typ.NumField(); i++ {
			// если тип поля совпадает с "sync.RWMutex", возвращаем true
			if typ.Field(i).Type.String() == "sync.RWMutex" {
				return true
			}
		}
	}
	return false
}

// проверка того, что структура Wallet содержит поле sync.RWMutex
func TestIsRWMutexUsed(t *testing.T) {
	var userWallet = &wallet.Wallet{}
	if !isRWMutexUsed(userWallet) {
		t.Errorf("Ожидалось, что структура Wallet содержит sync.RWMutex, однако такого поля нет\n")
	}
}

// тест конкрурентного пополнения кошелька
func TestRefill(t *testing.T) {
	var userWallet = &wallet.Wallet{}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			userWallet.Refill(1)
		}()
	}
	wg.Wait()
	var expectedBalance = 10000
	// тест не пройден, если полученное значение баланса не равно ожидаемому
	if got := userWallet.GetBalance(); got != expectedBalance {
		t.Errorf("Ожидался баланс %d, получено %d\n", expectedBalance, got)
	}
}

// тест конкурентного списания средств из кошелька
func TestWithdrawal(t *testing.T) {
	var userWallet = &wallet.Wallet{}
	userWallet.Refill(10000)
	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := userWallet.Withdrawal(1); err != nil {
				t.Errorf("Ошибка списания: %v", err) // тест провален, если произошла ошибка списания (нехватка средств)
			}
		}()
	}
	wg.Wait()
	var expectedBalance = 10000 - 5000 // ожидаемое значение баланса = начальное состояние (10000) - сумма списаний (5000)
	// тест не пройден, если полученное значение баланса не равно ожидаемому
	if got := userWallet.GetBalance(); got != expectedBalance {
		t.Errorf("Ожидался баланс %d, получено %d\n", expectedBalance, got)
	}
}

// проверка корректности отображения баланса при конкурентном доступе
func TestGetBalance(t *testing.T) {
	var userWallet = &wallet.Wallet{}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Nanosecond * time.Duration(rand.Intn(1000)))
			if i%2 == 0 {
				userWallet.Refill(1)
			} else {
				userWallet.Withdrawal(1)
			}
		}()
	}
	wg.Wait()
	var log = userWallet.GetLog()
	var logSum int
	for i := 0; i < len(log); i++ {
		logSum += log[i]
	}
	// тест не пройден, если баланс, полученный двумя разными способами не совпадает
	if logSum != userWallet.GetBalance() {
		t.Errorf("Ожидалось, что баланс, полученный из переменной баланса, и баланс, полученный в результате анализа логов списания и пополнения, будут совпадать, но получены различные значения: баланс(логи) %d, баланс (переменная) %d\n", logSum, userWallet.GetBalance())
	}
}
