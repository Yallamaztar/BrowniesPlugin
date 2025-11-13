package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Yallamaztar/BrowniesGambling/database"
)

func RegisterAdminCommands(cr *commandRegister, bank *database.Bank) {
	cr.RegisterCommand("gadmins", "gay", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		owners, err := database.ListOwners(cr.db)
		if err != nil || len(owners) == 0 {
			cr.rcon.Tell(clientNum, "No owners found")
			return
		}

		admins, err := database.ListAdmins(cr.db)
		if err != nil || len(admins) == 0 {
			cr.rcon.Tell(clientNum, "No admins found")
			return
		}

		for _, owner := range owners {
			cr.rcon.Tell(clientNum, fmt.Sprintf("[^5owner^7] %s | ^5XUID: ^7%s", owner.Player, owner.XUID))
		}

		for _, admin := range admins {
			cr.rcon.Tell(clientNum, fmt.Sprintf("[^5admin^7] %s | ^5XUID: ^7%s", admin.Player, admin.XUID))
		}
	})
	cr.RegisterCommand("sayas", "sayas", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!sayas ^7<player> <message> [-d] [-e]")
			return
		}

		var dead bool
		var enemy bool

		for _, arg := range args {
			switch arg {
			case "-d":
				dead = true
			case "-e":
				enemy = true
			}
		}

		filtered := make([]string, 0, len(args))
		for _, arg := range args {
			if arg != "-d" && arg != "-e" {
				filtered = append(filtered, arg)
			}
		}

		if len(filtered) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!sayas ^7<player> <message> [-d] [-e]")
			return
		}

		message := strings.Join(filtered[1:], " ")

		target := cr.findPlayer(filtered[0])
		if target == nil || target.clientNum == -1 {
			target = &playerInfo{
				Name: filtered[0],
			}
		}

		prefix := "^2"
		channel := "[Playing-All]"
		if enemy {
			prefix = "^1"
		}

		if dead {
			channel = "[Dead-All]"
		}

		formatted := fmt.Sprintf("%s%s %s: ^7%s", prefix, target.Name, channel, message)
		cr.rcon.SayRaw(formatted)
	})

	cr.RegisterCommand("loadout", "loadout", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!loadout ^7<player> <loadout>")
			return
		}

		t := cr.findPlayer(args[0])
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_enabled", "1")
		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("loadout %s %s", args[0], args[1]))
	})

	cr.RegisterCommand("take", "ta", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!take ^7<player> <amount>")
			return
		}

		t := cr.findPlayer(args[0])
		if t == nil {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}
		amount, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil || amount <= 0 {
			cr.rcon.Tell(clientNum, "Invalid amount")
			return
		}

		owlt := database.GetWallet(player, xuid, cr.db)
		twlt := database.GetWallet(t.Name, t.XUID, cr.db)
		if twlt == nil {
			cr.rcon.Tell(clientNum, "Target wallet not found")
			return
		}

		if isOwner {
			bank.TransferFromWallet(twlt, amount)
		} else if isAdmin {
			database.TransferFromWalletToWallet(twlt, owlt, amount)
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Took ^5$%d ^7from %s", amount, t.Name))
		if t.clientNum != -1 {
			cr.rcon.Tell(t.clientNum, fmt.Sprintf("^5%s ^7took ^5$%d from you", player, amount))
		}

		cr.logger.Printf("%s took $%d from %s", player, amount, t.Name)
	})

	cr.RegisterCommand("info", "if", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!xuid ^7<player>")
			return
		}

		t := cr.findPlayer(args[0])
		if t == nil {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}
		cr.rcon.Tell(clientNum, fmt.Sprintf("Player ^5%s ^4| ^7XUID: ^5%s ^4| ^7ClientNum: ^5%d", t.Name, t.XUID, t.clientNum))
	})

	cr.RegisterCommand("giveall", "ga", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!giveall ^7<amount>")
			return
		}

		amount, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil || amount <= 0 {
			cr.rcon.Tell(clientNum, "Invalid amount")
			return
		}

		if !isOwner {
			atotal, _ := database.AdminGiveTotal(cr.db, xuid)
			max := bank.Balance() / 20
			if atotal+amount > max {
				cr.rcon.Tell(clientNum, fmt.Sprintf("Admins can only give up to $%d total (you have given $%d)", max, atotal))
				return
			}
			database.IncrementAdminGiveTotal(cr.db, xuid, amount)
		}

		count, err := database.GiveAllWallets(cr.db, amount)
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to give all wallets")
			return
		}
		cr.rcon.Say(fmt.Sprintf("Gave ^5$%d ^7to ^5%d ^7wallets", amount, count))
	})

	cr.RegisterCommand("give", "gi", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!give ^7<player> <amount>")
			return
		}

		t := cr.findPlayer(args[0])
		if t == nil {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		wlt := database.GetWallet(t.Name, t.XUID, cr.db)

		amount, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil || amount <= 0 {
			cr.rcon.Tell(clientNum, "Invalid amount")
			return
		}

		if !isOwner {
			atotal, _ := database.AdminGiveTotal(cr.db, xuid)
			max := bank.Balance() / 20
			if atotal+amount > max {
				cr.rcon.Tell(clientNum, fmt.Sprintf("Admins can only give up to $%d total (you have given $%d)", max, atotal))
				return
			}
			database.IncrementAdminGiveTotal(cr.db, xuid, amount)
		}

		bank.TransferToWallet(wlt, amount)
		cr.rcon.Tell(clientNum, fmt.Sprintf("Gave ^5$%d ^7to %s", amount, t.Name))
		if tcn, ok := cr.GetClientNum(t.XUID); ok {
			cr.rcon.Tell(tcn, fmt.Sprintf("You received ^5$%d ^7from %s", amount, player))
		}
		cr.logger.Printf("%s gave $%d to %s from bank", player, amount, t.Name)
	})
}
