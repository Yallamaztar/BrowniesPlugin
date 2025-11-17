package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/BrowniesGambling/helpers"
)

func RegisterAdminCommands(cr *commandRegister, bank *database.Bank) {
	cr.RegisterCommand("setgravity", "sg", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!setgravity ^7<gravity>")
			return
		}

		gravity, err := strconv.ParseFloat(args[0], 64)
		if err != nil || gravity <= 0 {
			cr.rcon.Tell(clientNum, "Invalid gravity")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("setgravity %s %f", player, gravity))
	})

	cr.RegisterCommand("setspeed", "ss", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!setspeed ^7<speed>")
			return
		}

		speed, err := strconv.ParseFloat(args[0], 64)
		if err != nil || speed <= 0 {
			cr.rcon.Tell(clientNum, "Invalid speed")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("setspeed %s %f", player, speed))
	})

	cr.RegisterCommand("killplayer", "kpl", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!killplayer ^7<player>")
			return
		}

		var t, l string
		if len(args) == 1 {
			t = player
			l = args[0]
		} else {
			l = args[1]
			t := cr.findPlayer(args[0])
			if t == nil || t.clientNum == -1 {
				cr.rcon.Tell(clientNum, "Player not found")
				return
			}
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("killplayer %s %s", t, l))
	})

	cr.RegisterCommand("hide", "hd", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!hide ^7<player>")
			return
		}

		var t, l string
		if len(args) == 1 {
			t = player
			l = args[0]
		} else {
			l = args[1]
			t := cr.findPlayer(args[0])
			if t == nil || t.clientNum == -1 {
				cr.rcon.Tell(clientNum, "Player not found")
				return
			}
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("hide %s %s", t, l))
	})

	cr.RegisterCommand("spectator", "spec", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!spectator ^7<player>")
			return
		}

		var t, l string
		if len(args) == 1 {
			t = player
			l = args[0]
		} else {
			l = args[1]
			t := cr.findPlayer(args[0])
			if t == nil || t.clientNum == -1 {
				cr.rcon.Tell(clientNum, "Player not found")
				return
			}
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("spectator %s %s", t, l))
	})

	cr.RegisterCommand("teleport", "tp", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!teleport ^7<player> <target>")
			return
		}

		var t, l string
		if len(args) == 2 {
			t = args[0]
			l = args[1]
		} else {
			l = args[1]
			t := cr.findPlayer(args[0])
			if t == nil || t.clientNum == -1 {
				cr.rcon.Tell(clientNum, "Player not found")
				return
			}
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("teleport %s %s", t, l))
	})

	cr.RegisterCommand("fast", "res", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		err := cr.rcon.FastRestart()
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to fast restart")
			return
		}

	})

	cr.RegisterCommand("maprot", "mapr", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		err := cr.rcon.MapRotate()
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to rotate map")
			return
		}
		cr.rcon.Tell(clientNum, "Map rotated successfully")
	})

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

		var dead, enemy bool
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

		if len(args) == 0 {
			cr.rcon.Tell(clientNum, "Usage: ^5!loadout ^7<player> <loadout> ^4or ^5!loadout ^7<loadout>")
			return
		}

		var t, l string
		if len(args) == 1 {
			t = player
			l = args[0]
		} else {
			l = args[1]
			t := cr.findPlayer(args[0])
			if t == nil || t.clientNum == -1 {
				cr.rcon.Tell(clientNum, "Player not found")
				return
			}
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("loadout %s %s", t, l))
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

		cr.rcon.Tell(clientNum, fmt.Sprintf("Took ^5$%s ^7from %s", helpers.FormatMoney(amount), t.Name))
		if t.clientNum != -1 {
			cr.rcon.Tell(t.clientNum, fmt.Sprintf("^5%s ^7took ^5$%s from you", player, helpers.FormatMoney(amount)))
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
				cr.rcon.Tell(clientNum, fmt.Sprintf("Admins can only give up to $%s total (you have given $%s)", helpers.FormatMoney(max), helpers.FormatMoney(atotal)))
				return
			}
			database.IncrementAdminGiveTotal(cr.db, xuid, amount)
		}

		count, err := database.GiveAllWallets(cr.db, amount)
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to give all wallets")
			return
		}
		cr.rcon.Say(fmt.Sprintf("Gave ^5$%s ^7to ^5%d ^7wallets", helpers.FormatMoney(amount), count))
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
				cr.rcon.Tell(clientNum, fmt.Sprintf("Admins can only give up to $%s total (you have given $%s)", helpers.FormatMoney(max), helpers.FormatMoney(atotal)))
				return
			}
			database.IncrementAdminGiveTotal(cr.db, xuid, amount)
		}

		bank.TransferToWallet(wlt, amount)
		cr.rcon.Tell(clientNum, fmt.Sprintf("Gave ^5$%s ^7to %s", helpers.FormatMoney(amount), t.Name))
		if tcn, ok := cr.GetClientNum(t.XUID); ok {
			cr.rcon.Tell(tcn, fmt.Sprintf("You received ^5$%s ^7from %s", helpers.FormatMoney(amount), player))
		}
		cr.logger.Printf("%s gave $%s to %s from bank", player, helpers.FormatMoney(amount), t.Name)
	})
}
