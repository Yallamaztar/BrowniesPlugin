package database

import (
	"context"
	"database/sql"
	"log"
	"math"
	"sync"
)

type Bank struct {
	balance int64
	db      *sql.DB
	logger  *log.Logger
	mu      sync.RWMutex
}

func NewBank(balance int64, db *sql.DB, logger *log.Logger) *Bank {
	schema := `CREATE TABLE IF NOT EXISTS bank (
		total INTEGER NOT NULL
    );`

	if _, err := db.Exec(schema); err != nil {
		return nil
	}

	_, err := db.Exec("INSERT INTO bank (total) SELECT ? WHERE NOT EXISTS (SELECT 1 FROM bank)", balance)
	if err != nil {
		return nil
	}

	return &Bank{balance: balance, db: db, logger: logger}
}

func (b *Bank) Balance() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	err := b.db.QueryRow("SELECT total FROM bank ORDER BY rowid LIMIT 1").Scan(&b.balance)
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

	res, err := b.db.Exec(
		"UPDATE bank SET total = total - ? WHERE rowid = (SELECT rowid FROM bank ORDER BY rowid LIMIT 1) AND total >= ?",
		amount, amount,
	)
	if err != nil {
		return
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		return
	}

	b.balance -= amount
	b.logger.Printf("Withdrew %d from bank (%d)", amount, b.balance)
}

func (b *Bank) TransferToWallet(w *Wallet, amount int64) {
	if amount <= 0 {
		return
	}

	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	res, err := tx.Exec(
		"UPDATE bank SET total = total - ? WHERE rowid = (SELECT rowid FROM bank ORDER BY rowid LIMIT 1) AND total >= ?",
		amount, amount,
	)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		_ = tx.Rollback()
		return
	}

	res, err = tx.Exec("UPDATE wallets SET balance = balance + ? WHERE xuid = ?", amount, w.xuid)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		_ = tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return
	}

	b.mu.Lock()
	w.mu.Lock()

	b.balance -= amount
	w.balance += amount

	w.mu.Unlock()
	b.mu.Unlock()

	b.logger.Printf("Transferred $%d from bank to wallet %s", amount, w.player)
}

func (b *Bank) TransferFromWallet(w *Wallet, amount int64) {
	if amount <= 0 {
		return
	}

	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	res, err := tx.Exec(
		"UPDATE wallets SET balance = balance - ? WHERE xuid = ?",
		amount, w.xuid,
	)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		_ = tx.Rollback()
		return
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return
	}

	w.mu.Lock()
	w.balance -= amount
	w.mu.Unlock()

	credited := false
	ctx2 := context.Background()
	tx2, err := b.db.BeginTx(ctx2, nil)
	if err == nil {
		limit := math.MaxInt64 - amount
		res, err = tx2.Exec(
			"UPDATE bank SET total = total + ? WHERE rowid = (SELECT rowid FROM bank ORDER BY rowid LIMIT 1) AND total <= ?",
			amount, limit,
		)
		if err == nil {
			if rows, _ := res.RowsAffected(); rows == 1 {
				if err := tx2.Commit(); err == nil {
					credited = true
				} else {
					_ = tx2.Rollback()
				}
			} else {
				_ = tx2.Rollback()
			}
		} else {
			_ = tx2.Rollback()
		}
	}

	if credited {
		b.mu.Lock()
		b.balance += amount
		b.mu.Unlock()
		b.logger.Printf("Transferred $%d from wallet %s to bank", amount, w.player)
	} else {
		b.logger.Printf("Tried to transfer $%d from wallet %s (limit reached)", amount, w.player)
	}
}
