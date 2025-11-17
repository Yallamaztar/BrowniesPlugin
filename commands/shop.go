package commands

import (
	"fmt"
	"strings"

	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/BrowniesGambling/helpers"
)

func RegisterShopCommands(cr *commandRegister, bank *database.Bank) {
	aliases := map[string]string{
		"kpl":  "killplayer",
		"hd":   "hide",
		"tp":   "teleport",
		"spec": "spectator",
		"st":   "sayto",
		"gw":   "giveweapon",
		"tw":   "takeweapons",
		"frz":  "freeze",
		"ss":   "setspeed",
		"sl":   "slap",
	}

	cr.RegisterCommand("shop", "sh", func(clientNum int, player, xuid string, args []string) {
		items, err := database.ListShopItems(cr.db)
		if err != nil || len(items) == 0 {
			cr.rcon.Tell(clientNum, "Shop empty")
			return
		}
		cr.rcon.Tell(clientNum, "== ^5Shop Items ^7==")
		for _, it := range items {
			cr.rcon.Tell(clientNum, fmt.Sprintf("[^5$%s^7] ^5%s - ^7%s", helpers.FormatMoney(it.Price), it.Name, it.Description))
		}
		cr.rcon.Tell(clientNum, "^7Buy with ^5!buy <item|alias> ^7<player (optional)>")
	})

	cr.RegisterCommand("buy", "bu", func(clientNum int, player, xuid string, args []string) {
		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!buy ^7<item> <player (optional)>")
			return
		}
		nameOrAlias := strings.ToLower(args[0])
		name := nameOrAlias
		if canon, ok := aliases[nameOrAlias]; ok {
			name = canon
		}
		it, err := database.GetShopItem(cr.db, name)
		if err != nil {
			cr.rcon.Tell(clientNum, "Item not found")
			return
		}
		wlt := database.GetWallet(player, xuid, cr.db)
		if wlt == nil {
			cr.rcon.Tell(clientNum, "Wallet not found")
			return
		}
		if wlt.Balance() < it.Price {
			cr.rcon.Tell(clientNum, "Not enough money")
			return
		}
		bank.TransferFromWallet(wlt, it.Price)
		cr.rcon.Tell(clientNum, fmt.Sprintf("Purchased ^5%s ^7for ^5$%s", it.Name, helpers.FormatMoney(it.Price)))
		if it.Command != "" {
			target := player
			if len(args) >= 2 {
				if t := cr.findPlayer(args[1]); t != nil {
					target = t.Name
				} else {
					cr.rcon.Tell(clientNum, "Target player not found")
					return
				}
			}
			cmd := strings.ReplaceAll(it.Command, "%PLAYER%", target)
			cr.rcon.SetDvar("brwns_enabled", "1")
			cr.rcon.SetDvar("brwns_exec_in", cmd)
		}
	})
}
