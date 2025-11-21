package bus

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Yallamaztar/BrowniesPlugin/database"
	"github.com/Yallamaztar/BrowniesPlugin/rcon"
	"github.com/Yallamaztar/events/events"
)

type commandExecutor interface {
	SetClientNum(xuid string, clientNum int)
	Exec(command string, clientNum int, player, xuid string, args []string) bool
}

func HandleEvents(logPath string, ctx context.Context, rc *rcon.RCONClient, logger *log.Logger, db *sql.DB, bdb *database.Bank, cr commandExecutor) {
	ch := make(chan events.Event, 128)
	go func() {
		if err := events.TailFileContext(ctx, logPath, true, ch); err != nil {
			logger.Fatalf("Failed to tail log file: %v", err)
		}
		close(ch)
	}()

	for e := range ch {
		if !database.IsGamblingEnabled(db) {
			continue
		}

		switch t := e.(type) {
		case *events.KillEvent:

			cr.SetClientNum(t.AttackerXUID, t.AttackerClientNum)
			cr.SetClientNum(t.VictimXUID, t.VictimClientNum)

			awlt := database.GetWallet(t.VictimName, t.VictimXUID, db)
			bdb.TransferToWallet(awlt, 150)
			rc.Tell(t.VictimClientNum, fmt.Sprintf("Kill Reward: ^5$%d", 150))

			vwlt := database.GetWallet(t.AttackerName, t.AttackerXUID, db)
			bdb.TransferFromWallet(vwlt, 200)
			rc.Tell(t.AttackerClientNum, fmt.Sprintf("Death Penalty: ^1$%d", 200))

		case *events.PlayerEvent:
			cr.SetClientNum(t.XUID, t.Flag)

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

						if cr.Exec(parts[0], t.Flag, t.Player, t.XUID, args) {
							logger.Printf("%s: !%s %v", t.Player, parts[0], args)
						}
					}
				}
			}
		}
	}
}
