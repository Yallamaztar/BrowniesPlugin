package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Yallamaztar/BrowniesGambling/rcon"
	"github.com/bwmarrin/discordgo"
)

func RunDiscordBot(ctx context.Context, token string, rc *rcon.RCONClient, logger *log.Logger) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("discord session error: %w", err)
	}

	if err := dg.Open(); err != nil {
		return fmt.Errorf("discord open error: %w", err)
	}

	logger.Println("Discord bot is now running")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				status, err := rc.Status()
				if err != nil {
					logger.Println("RCON status error:", err)
					_ = dg.UpdateGameStatus(0, "RCON Offline")
				} else {
					presence := fmt.Sprintf("%d players on %s", len(status.Players), status.Map)
					_ = dg.UpdateGameStatus(0, presence)
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()

	<-ctx.Done()
	logger.Println("Shutting down Discord bot...")
	return dg.Close()
}
