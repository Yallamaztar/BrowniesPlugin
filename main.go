package main

import (
	"context"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Yallamaztar/BrowniesGambling/commands"
	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/PlutoRCON/rcon"
)

var (
	dbPath  = "brownies.db"
	logPath = filepath.Join(os.Getenv("LOCALAPPDATA"), "Plutonium2", "storage", "t6", "main", "logs", "games_mp3.log")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.New(os.Stdout, "[Gambling] ", log.LstdFlags)

	db, err := database.Open(dbPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	bdb := database.NewBank(math.MaxInt64, db)
	if bdb == nil {
		logger.Fatalf("Failed to initialize bank database")
	}
	logger.Printf("Bank initialized with balance: %d", bdb.Balance())

	if err := database.NewOwners(db); err != nil {
		logger.Fatalf("Failed to ensure owners table: %v", err)
	}

	if err := database.EnsureShop(db); err != nil {
		logger.Fatalf("Failed to ensure shop table: %v", err)
	}
	if err := database.SeedShop(db); err != nil {
		logger.Printf("Shop seeding warning: %v", err)
	}

	database.AddOwner(db, "budiworld", "2045030")

	rc, err := rcon.New(
		os.Getenv("RCON_IP"),
		os.Getenv("RCON_PORT"),
		os.Getenv("RCON_PASSWORD"),
	)

	if err != nil {
		logger.Fatalf("Failed to connect to RCON: %v", err)
	}

	rc.SetDvar("brwns_enabled", "1")
	rc.SetDvar("brwns_exec", "hide budigp rgo")
	defer rc.Close()

	reg := commands.New(logger, rc, db)
	reg.RegisterCommands(bdb)

	go commands.HandleEvents(logPath, ctx, rc, logger, db, bdb, reg)

	<-ctx.Done()
	rc.Close()
	db.Close()
	stop()
}
