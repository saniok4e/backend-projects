package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Account struct {
	mu      sync.Mutex
	ID      int
	Balance int
}

func (from *Account) Transfer(to *Account, wg *sync.WaitGroup, amount int) bool {
	defer wg.Done()

	if from.ID == to.ID {
		from.mu.Lock()
		defer from.mu.Unlock()
	} else if from.ID < to.ID {
		from.mu.Lock()
		to.mu.Lock()
		defer from.mu.Unlock()
		defer to.mu.Unlock()
	} else {
		to.mu.Lock()
		from.mu.Lock()
		defer to.mu.Unlock()
		defer from.mu.Unlock()
	}

	if from.Balance >= amount {
		from.Balance -= amount
		to.Balance += amount
		return true
	}

	return false
}

func main() {
	accounts := make([]Account, 10)

	for i := 0; i < 10; i++ {
		accounts[i] = Account{ID: i, Balance: 1000}
	}

	var success, failed int
	var statsMu sync.Mutex
	wg := sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			from := &accounts[rand.Intn(len(accounts))]
			to := &accounts[rand.Intn(len(accounts))]
			amount := rand.Intn(500) + 1
			if from.Transfer(to, &wg, amount) {
				statsMu.Lock()
				success++
				statsMu.Unlock()
			} else {
				statsMu.Lock()
				failed++
				statsMu.Unlock()
			}

		}()
	}

	wg.Wait()

	sum := 0
	for _, acc := range accounts {
		sum += acc.Balance
	}

	fmt.Printf("✅Успешных переводов: %d\n", success)
	fmt.Printf("❌Неудачных переводов: %d\n", failed)
	fmt.Printf("Всего денег: %d\n", sum)
	fmt.Printf("%v\n", accounts)
}
