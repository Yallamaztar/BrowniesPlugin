package discord

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Yallamaztar/BrowniesPlugin/rcon"
	"github.com/bwmarrin/discordgo"
)

func RunDiscordBot(ctx context.Context, token string, rcs []*rcon.RCONClient, logger *log.Logger) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("discord session error: %w", err)
	}

	if err := dg.Open(); err != nil {
		return fmt.Errorf("discord open error: %w", err)
	}

	logger.Println("Discord bot is now running")
	// Add your own logic here im not bothered rn atleast
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				totalPlayers := 0
				liveServers := 0
				for _, rc := range rcs {
					if rc == nil {
						continue
					}
					if status, err := rc.Status(); err == nil && status != nil {
						totalPlayers += len(status.Players)
						liveServers++
					}
				}
				presence := fmt.Sprintf("Watching %d gambler(s) across %d servers", totalPlayers, liveServers)
				_ = dg.UpdateGameStatus(0, presence)

				time.Sleep(5 * time.Second)
			}
		}
	}()

	<-ctx.Done()
	logger.Println("Shutting down Discord bot (multi)...")
	return dg.Close()
}
