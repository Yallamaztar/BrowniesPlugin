package database

import (
	"database/sql"
	"log"
	"sync"
)

type Bank struct {
	balance int64
	db      *sql.DB
	mu      sync.RWMutex
}

func NewBank(balance int64, db *sql.DB) *Bank {
	schema := `CREATE TABLE IF NOT EXISTS bank (
        total INTEGER NOT NULL
    );`

	if _, err := db.Exec(schema); err != nil {
		return nil
	}

	_, err := db.Exec("INSERT INTO bank (total) VALUES (?)", balance)
	if err != nil {
		return nil
	}

	return &Bank{balance: balance, db: db}
}

func (b *Bank) Balance() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	err := b.db.QueryRow("SELECT total FROM bank").Scan(&b.balance)
	if err != nil {
		return 0
	}

	return b.balance
}

func (b *Bank) Withdraw(amount int64) {
	if amount <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.balance -= amount
	_, err := b.db.Exec("UPDATE bank SET total = ?", b.balance)
	if err != nil {
		return
	}

	log.Printf("Withdrew %d from bank (%d)", amount, b.balance)
}

func (b *Bank) Deposit(amount int64) {
	if amount <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.balance += amount
	_, err := b.db.Exec("UPDATE bank SET total = ?", b.balance)
	if err != nil {
		return
	}

	log.Printf("Deposited %d to bank (%d)", amount, b.balance)
}

func (b *Bank) TransferToWallet(w *Wallet, amount int64) {
	if amount <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	if b.balance < amount {
		return
	}

	b.balance -= amount
	w.balance += amount

	_, err := b.db.Exec("UPDATE bank SET total = ?", b.balance)
	if err != nil {
		return
	}
	_, err = w.db.Exec("UPDATE wallets SET balance = ? WHERE xuid = ?", w.balance, w.xuid)
	if err != nil {
		return
	}

	log.Printf("Transferred $%d from bank to wallet %s ($%d)", amount, w.player, b.balance)
}

func (b *Bank) TranserFromWallet(w *Wallet, amount int64) {
	if amount <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	b.balance += amount
	w.balance -= amount
	_, err := b.db.Exec("UPDATE bank SET total = ?", b.balance)
	if err != nil {
		return
	}

	_, err = w.db.Exec("UPDATE wallets SET balance = ? WHERE xuid = ?", w.balance, w.xuid)
	if err != nil {
		return
	}

	log.Printf("Transferred $%d from wallet %s to bank ($%d)", amount, w.player, b.balance)
}
