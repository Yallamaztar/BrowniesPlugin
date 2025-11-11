package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type Wallet struct {
	player  string
	xuid    string
	balance int64
	db      *sql.DB
	mu      sync.RWMutex
}

func NewWallet(player, xuid string, balance int64, db *sql.DB) *Wallet {
	schema := `CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player TEXT NOT NULL,
		xuid TEXT NOT NULL UNIQUE,
		balance INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil
	}

	_, err := db.Exec("INSERT OR IGNORE INTO wallets (player, xuid, balance) VALUES (?, ?, ?)", player, xuid, balance)
	if err != nil {
		return nil
	}

	var actualBalance int64
	err = db.QueryRow("SELECT balance FROM wallets WHERE xuid = ?", xuid).Scan(&actualBalance)
	if err != nil {
		actualBalance = balance
	}

	return &Wallet{player: player, xuid: xuid, balance: actualBalance, db: db}
}

func (w *Wallet) Balance() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var balance int64
	err := w.db.QueryRow("SELECT balance FROM wallets WHERE xuid = ?", w.xuid).Scan(&balance)
	if err == nil {
		w.balance = balance
	}

	return w.balance
}

func (w *Wallet) AddBalance(amount int64) error {
	tx, err := w.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err := w.AddBalanceTx(tx, amount); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	w.mu.Lock()
	w.balance += amount
	w.mu.Unlock()
	return nil
}

func (w *Wallet) SubtractBalance(amount int64) error {
	tx, err := w.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err := w.SubtractBalanceTx(tx, amount); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	w.mu.Lock()
	w.balance -= amount
	w.mu.Unlock()
	return nil
}

func (w *Wallet) AddBalanceTx(tx *sql.Tx, amount int64) error {
	if amount <= 0 {
		return nil
	}

	res, err := tx.Exec("UPDATE wallets SET balance = balance + ? WHERE xuid = ?", amount, w.xuid)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return sql.ErrNoRows
	}
	return nil
}

func (w *Wallet) SubtractBalanceTx(tx *sql.Tx, amount int64) error {
	if amount <= 0 {
		return nil
	}

	res, err := tx.Exec("UPDATE wallets SET balance = balance - ? WHERE xuid = ? AND balance >= ?", amount, w.xuid, amount)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return fmt.Errorf("insufficient funds or wallet not found")
	}
	return nil
}

func GetWallet(player, xuid string, db *sql.DB) *Wallet {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM wallets WHERE xuid = ?)", xuid).Scan(&exists)

	if err != nil || !exists {
		return NewWallet(player, xuid, 1000, db)
	}

	var balance int64
	err = db.QueryRow("SELECT balance FROM wallets WHERE xuid = ?", xuid).Scan(&balance)
	if err != nil {
		balance = 1000
	}

	return &Wallet{player: player, xuid: xuid, balance: balance, db: db}
}

func TransferFromWalletToWallet(from, to *Wallet, amount int64) error {
	if from == nil || to == nil {
		return fmt.Errorf("wallet is nil")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if from.xuid == to.xuid {
		return fmt.Errorf("cannot transfer to the same wallet")
	}
	if from.db != to.db {
		return fmt.Errorf("wallets use different database handles")
	}

	tx, err := from.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err := from.SubtractBalanceTx(tx, amount); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := to.AddBalanceTx(tx, amount); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	first, second := from, to
	debitFirst := true
	if to.xuid < from.xuid {
		first, second = to, from
		debitFirst = false
	}

	first.mu.Lock()
	second.mu.Lock()
	if debitFirst {
		from.balance -= amount
		to.balance += amount
	} else {
		to.balance += amount
		from.balance -= amount
	}
	second.mu.Unlock()
	first.mu.Unlock()

	return nil
}
