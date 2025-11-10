package main

import (
	"context"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yallamaztar/BrowniesGambling/commands"
	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/PlutoRCON/rcon"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.New(os.Stdout, "[Gambling] ", log.LstdFlags)

	db, err := database.Open("brownies.db")
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

	database.AddOwner(db, "budiworld", "1F3466")

	rc, err := rcon.New(
		os.Getenv("RCON_IP"),
		os.Getenv("RCON_PORT"),
		os.Getenv("RCON_PASSWORD"),
	)

	if err != nil {
		logger.Fatalf("Failed to connect to RCON: %v", err)
	}

	defer rc.Close()

	reg := commands.New(logger, rc, db)
	reg.RegisterCommands(bdb)

	go commands.HandleEvents(ctx, rc, logger, db, bdb, reg)

	<-ctx.Done()
	rc.Close()
	db.Close()
	stop()
}
