package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Yallamaztar/BrowniesPlugin/bot"
	"github.com/Yallamaztar/BrowniesPlugin/commands"
	"github.com/Yallamaztar/BrowniesPlugin/config"
	"github.com/Yallamaztar/BrowniesPlugin/database"
	"github.com/Yallamaztar/BrowniesPlugin/helpers"
	"github.com/Yallamaztar/BrowniesPlugin/rcon"
)

func setupDatabase(logger *log.Logger) (*sql.DB, *database.Bank, error) {
	db, err := database.Open("brownies.db")
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

	if err := database.EnsureSettings(db); err != nil {
		logger.Fatalf("Failed to ensure settings table: %v", err)
	}

	database.EnableGambling(db, true)

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

	// Add Owners
	database.AddOwner(db, "budiworld", "2045030")
	database.AddOwner(db, "SsugonmaA", "4544213")
	database.AddOwner(db, "icmp", "4840217")

	// Add Admins
	database.AddAdmin(db, "[XRP]OGRuntz", "4941187")
	database.AddAdmin(db, "[kitty]xAkame", "5968446")
	database.AddAdmin(db, "B R I K", "4002521")
	database.AddAdmin(db, "Larry Funk", "2538213")
	database.AddAdmin(db, "F A Y A Z", "4096583")
	database.AddAdmin(db, "bakedzzz", "5975331")
	database.AddAdmin(db, "zPoi", "1230466")
	database.AddAdmin(db, "ssy", "5358787")
	database.AddAdmin(db, "palma.", "5648548")
	database.AddAdmin(db, "yungraven", "4790962")

	config.InitConfig()

	var wg sync.WaitGroup
	var rcs []*rcon.RCONClient
	cfg, err := config.Load("config.json")
	if err == nil && cfg != nil && len(cfg.Servers) > 0 {
		for _, s := range cfg.Servers {
			slogger := log.New(os.Stdout, fmt.Sprintf("[%s:%s][Gambling] ", s.IP, s.Port), log.LstdFlags)
			rc, err := setupRCON(s.IP, s.Port, s.Password, slogger)
			if err != nil {
				slogger.Printf("RCON setup failed for %s:%s: %v", s.IP, s.Port, err)
				continue
			}
			reg := commands.New(slogger, rc, db)
			reg.RegisterCommands(db, bdb, rc)

			wg.Add(1)
			go func(logPath string, rc *rcon.RCONClient, slogger *log.Logger) {
				defer wg.Done()
				commands.HandleEvents(logPath, ctx, rc, slogger, db, bdb, reg)
			}(s.LogPath, rc, slogger)

			rcs = append(rcs, rc)
		}
	}

	if token := os.Getenv("BOT_TOKEN"); token != "" && len(rcs) > 0 {
		go bot.RunDiscordBotMulti(ctx, token, rcs, logger)
	}

	helpers.OnReadyWebhook()

	<-ctx.Done()
	logger.Println("Shutting down...")

	stop()

	for _, rc := range rcs {
		_ = rc.Close()
	}

	wg.Wait()
	helpers.OnShutdownWebhook()
	_ = db.Close()
}
