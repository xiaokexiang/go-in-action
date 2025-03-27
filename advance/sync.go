package main

import (
	"fmt"
	"sync"
	"time"
)

type BankAccount struct {
	balance float64
	mu      sync.Mutex
}

func (ba *BankAccount) Deposit(amount float64) {
	defer ba.mu.Unlock()
	ba.mu.Lock()
	current := ba.balance         // 1. 读取当前余额
	time.Sleep(time.Millisecond)  // 2. 模拟处理延迟
	ba.balance = current + amount // 3. 写入新余额
}

func (ba *BankAccount) Withdraw(amount float64) bool {
	defer ba.mu.Unlock()
	ba.mu.Lock()
	current := ba.balance
	time.Sleep(time.Millisecond)
	if current >= amount {
		ba.balance = current - amount
		return true
	}
	return false
}

func main() {
	account := &BankAccount{balance: 1000}
	var wg sync.WaitGroup

	// 启动100个goroutine并发存款1元
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Deposit(1)
		}()
	}

	// 启动100个goroutine并发取款1元
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Withdraw(1)
		}()
	}

	wg.Wait()
	fmt.Printf("预期余额: 1000, 实际余额: %.2f\n", account.balance)
}
