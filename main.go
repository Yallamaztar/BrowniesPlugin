package main

import (
	"context"
	"database/sql"
	"fmt"
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
	dbPath = "brownies.db"
	path   = filepath.Join(os.Getenv("LOCALAPPDATA"), "Plutonium2", "storage", "t6", "mods", "mp_brownies", "logs", "games_mp3.log") // CHANGE ME: Update the log path if necessary
)

func setupDatabase(logger *log.Logger) (*sql.DB, *database.Bank, error) {
	db, err := database.Open(dbPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}

	// Initialize the bank with a balance close to MaxInt64 to prevent overflow
	bdb := database.NewBank((math.MaxInt64 - 9223372003), db, logger)
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

	const maxAttempts = 5

	for i := 1; i <= maxAttempts; i++ {
		logger.Printf("Attempt %d/%d: Attempting to detect Brownies mod on the server\n", i, maxAttempts)

		rc.SetDvar("brwns_enabled", "1")
		rc.SetDvar("brwns_exec_in", "onstart")

		d, err := rc.GetDvar("brwns_exec_out")
		if err != nil {
			logger.Println("Error reading brwns_exec_out:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(d) >= 29 && d[22:29] == "success" {
			logger.Println("Brownies mod detected on the server")
			return rc, nil
		}

		logger.Println("Brownies mod not detected")
		time.Sleep(750 * time.Millisecond)
	}

	logger.Fatal("Brownies mod not detected after 5 attempts")
	return nil, fmt.Errorf("brownies mod not detected")
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
	database.AddAdmin(db, "[XRP]OGRuntz", "4941187")
	database.AddAdmin(db, "[kitty]xAkame", "5968446")
	database.AddAdmin(db, "B R I K", "4002521")
	database.AddAdmin(db, "Larry Funk", "2538213")

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

	go commands.HandleEvents(path, ctx, rc, logger, db, bdb, reg)

	<-ctx.Done()
	rc.Close()
	db.Close()
}
