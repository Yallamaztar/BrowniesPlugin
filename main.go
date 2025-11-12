package main

import (
	"context"
	"database/sql"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Yallamaztar/BrowniesGambling/commands"
	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/BrowniesGambling/rcon"
)

var (
	dbPath  = "brownies.db"
	logPath = filepath.Join(os.Getenv("LOCALAPPDATA"), "Plutonium2", "storage", "t6", "main", "logs", "games_mp3.log")
)

func setupDatabase(logger *log.Logger) (*sql.DB, *database.Bank, error) {
	db, err := database.Open(dbPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	// Do NOT close here; caller (main) owns the lifetime of db

	bdb := database.NewBank(math.MaxInt64, db, logger)
	if bdb == nil {
		logger.Fatalf("Failed to initialize bank database")
	}
	logger.Printf("Bank loaded with balance: %d", bdb.Balance())

	if err := database.NewOwners(db); err != nil {
		logger.Fatalf("Failed to ensure owners table: %v", err)
	}

	if err := database.NewAdmins(db); err != nil {
		logger.Fatalf("Failed to ensure admins table: %v", err)
	}

	if err := database.NewAdminLimits(db); err != nil {
		logger.Fatalf("Failed to ensure admin limits table: %v", err)
	}

	if err := database.EnsureShop(db); err != nil {
		logger.Fatalf("Failed to ensure shop table: %v", err)
	}
	if err := database.SeedShop(db); err != nil {
		logger.Printf("Shop seeding warning: %v", err)
	}

	return db, bdb, nil
}

func setupRCON(ip, port, password string, logger *log.Logger) (*rcon.RCONClient, error) {
	rc, err := rcon.New(ip, port, password)
	if err != nil {
		return nil, err
	}

	for {
		logger.Println("Attempting to detect Brownies mod on the server")
		rc.SetDvar("brwns_enabled", "1")
		rc.SetDvar("brwns_exec_in", "onstart")

		d, err := rc.GetDvar("brwns_exec_out")
		if err != nil || d[22:][:7] != "success" {
			logger.Println("brownies mod not detected on the server")
			time.Sleep(1 * time.Second)
			continue
		} else if d[22:][:7] == "success" {
			break
		}
	}

	logger.Println("Brownies mod detected on the server")
	return rc, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.New(os.Stdout, "[Gambling] ", log.LstdFlags)

	logger.Println("Setting up database")
	db, bdb, err := setupDatabase(logger)
	if err != nil {
		logger.Fatalf("Database setup failed: %v", err)
	}
	logger.Println("Database setup complete")

	database.AddOwner(db, "budiworld", "2045030")

	rc, err := setupRCON(
		os.Getenv("RCON_IP"),
		os.Getenv("RCON_PORT"),
		os.Getenv("RCON_PASSWORD"),
		logger,
	)

	if err != nil {
		logger.Fatalf("RCON setup failed: %v", err)
	}

	reg := commands.New(logger, rc, db)
	reg.RegisterCommands(db, bdb, rc)

	go commands.HandleEvents(logPath, ctx, rc, logger, db, bdb, reg)

	<-ctx.Done()
	rc.Close()
	db.Close()
}
