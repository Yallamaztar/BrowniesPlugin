package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/PlutoRCON/rcon"
	"github.com/Yallamaztar/events/events"
)

func HandleEvents(ctx context.Context, rc *rcon.RCONClient, logger *log.Logger, db *sql.DB, bdb *database.Bank, reg *commandRegister) {
	ch := make(chan events.Event, 128)
	go func() {
		if err := events.TailFileContext(ctx, "games_mp3.log", true, ch); err != nil {
			logger.Fatalf("Failed to tail log file: %v", err)
		}
		close(ch)
	}()

	for e := range ch {
		switch t := e.(type) {
		case *events.KillEvent:
			reg.SetClientNum(t.AttackerXUID, t.AttackerClientNum)
			reg.SetClientNum(t.VictimXUID, t.VictimClientNum)
			awlt := database.GetWallet(t.VictimName, t.VictimXUID, db)
			bdb.TransferToWallet(awlt, 150)
			rc.Tell(t.VictimClientNum, fmt.Sprintf("Kill Reward: ^5$%d", 150))

			vwlt := database.GetWallet(t.AttackerName, t.AttackerXUID, db)
			bdb.TransferFromWallet(vwlt, 200)
			rc.Tell(t.AttackerClientNum, fmt.Sprintf("Death Penalty: ^1$%d", 200))

		case *events.PlayerEvent:
			reg.SetClientNum(t.XUID, t.Flag)
			if t.Command == "J" {
				wlt := database.GetWallet(t.Player, t.XUID, db)
				if wlt != nil {
					bdb.TransferToWallet(wlt, 50)
				}
			}

			if t.Command == "say" || t.Command == "sayteam" {
				if after, ok := strings.CutPrefix(t.Message, "!"); ok {
					cmd := after
					parts := strings.Fields(cmd)
					if len(parts) > 0 {
						args := []string{}

						if len(parts) > 1 {
							args = parts[1:]
						}

						if reg.Exec(parts[0], t.Flag, t.Player, t.XUID, args) {
							logger.Printf("%s: !%s %v", t.Player, parts[0], args)
						}
					}
				}
			}
		}
	}
}
