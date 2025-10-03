package main

import (
	"errors"
	"fmt"
	"sync"
)

type BankAccount struct {
	accountNumber string
	holderName    string
	balance       float64
	mux           sync.RWMutex
}

func (b *BankAccount) Deposit(amount float64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.balance += amount
}

func (b *BankAccount) Withdraw(amount float64) error {
	b.mux.Lock()
	defer b.mux.Unlock()

	if amount > b.balance {
		return errors.New("недостаточно средств")
	}
	b.balance -= amount
	return nil
}

func (b *BankAccount) GetBalance() float64 {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.balance
}

func main() {
	account := BankAccount{
		accountNumber: "24354236",
		holderName:    "Петр Михайлов",
		balance:       1000.0,
	}

	fmt.Printf("Баланс: %.2f\n", account.GetBalance())

	account.Deposit(500.0)
	fmt.Printf("Пополнения: %.2f\n", account.GetBalance())

	err := account.Withdraw(300.0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Списание: %.2f\n", account.GetBalance())
	}

	err = account.Withdraw(2000.0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
}
