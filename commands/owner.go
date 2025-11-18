package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Yallamaztar/BrowniesPlugin/database"
	"github.com/Yallamaztar/BrowniesPlugin/helpers"
)

func RegisterOwnerCommands(cr *commandRegister, bank *database.Bank) {
	cr.RegisterCommand("godmode", "god", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!godmode ^7<player (optional)>")
			return
		}

		var target string
		if len(args) >= 1 {
			target = args[0]
		} else {
			target = player
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("godmode %s", t.Name))
	})

	cr.RegisterCommand("gambling", "gmbl", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!gambling ^7<enable|disable|status>")
			return
		}

		switch strings.ToLower(args[0]) {
		case "enable":
			err := database.EnableGambling(cr.db, true)
			if err != nil {
				cr.rcon.Tell(clientNum, "Failed to enable gambling")
				return
			}
			cr.rcon.Say("Gambling has been ^5enabled^7")
			cr.logger.Printf("%s enabled gambling", player)

		case "disable":
			err := database.EnableGambling(cr.db, false)
			if err != nil {
				cr.rcon.Tell(clientNum, "Failed to disable gambling")
				return
			}
			cr.rcon.Say("Gambling has been ^5disabled^7")
			cr.logger.Printf("%s disabled gambling", player)

		case "status":
			enabled := database.IsGamblingEnabled(cr.db)
			status := "^1disabled^7"
			if enabled {
				status = "^2enabled^7"
			}
			cr.rcon.Say(fmt.Sprintf("Gambling is currently %s", status))

		default:
			cr.rcon.Tell(clientNum, "Usage: ^5!gambling ^7<enable|disable|status>")
		}
	})

	cr.RegisterCommand("maxbet", "mb", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!maxbet ^7<amount|0 to disable>")
			return
		}

		if strings.ToLower(args[0]) == "status" {
			max := database.GetMaxBet(cr.db)
			if max == 0 {
				cr.rcon.Tell(clientNum, "Max bet is currently ^5disabled^7")
			} else {
				cr.rcon.Tell(clientNum, fmt.Sprintf("Max bet is currently set to ^5$%s^7", helpers.FormatMoney(max)))
			}
			return
		}

		amount, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil || amount < 0 {
			cr.rcon.Tell(clientNum, "Invalid amount")
			return
		}

		err = database.SetMaxBet(cr.db, amount)
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to set max bet")
			return
		}

		if amount == 0 {
			cr.rcon.Say("Max bet has been ^5disabled^7")
			cr.logger.Printf("%s disabled max bet", player)
		} else {
			cr.rcon.Say(fmt.Sprintf("Max bet set to ^5$%s^7", helpers.FormatMoney(amount)))
			cr.logger.Printf("%s set max bet to $%d", player, amount)
		}
	})

	cr.RegisterCommand("printmoney", "print", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!printmoney ^7<amount>")
			return
		}

		amount, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil || amount <= 0 {
			cr.rcon.Tell(clientNum, "Invalid amount")
			return
		}

		bank.TransferToWallet(database.GetWallet(player, xuid, cr.db), amount)
		cr.rcon.Tell(clientNum, fmt.Sprintf("Printed ^5$%s ^7to your wallet", helpers.FormatMoney(amount)))
		cr.logger.Printf("%s printed $%d to their wallet", player, amount)
	})

	cr.RegisterCommand("addowner", "add", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!addowner ^7<player> <xuid>")
			return
		}

		err = database.AddOwner(cr.db, args[0], args[1])
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to add owner")
			return
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Added ^5%s ^7as owner", args[0]))
		cr.logger.Printf("%s added %s as owner", player, args[0])
	})

	cr.RegisterCommand("removeowner", "remove", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!removeowner ^7<xuid>")
			return
		}

		err = database.RemoveOwner(cr.db, args[0])
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to remove owner")
			return
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Removed ^5%s ^7from owners", args[0]))
		cr.logger.Printf("%s removed %s from owners", player, args[0])
	})

	cr.RegisterCommand("addadmin", "adda", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!addadmin ^7<player> <xuid>")
			return
		}

		err = database.AddAdmin(cr.db, args[0], args[1])
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to add admin")
			return
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Added ^5%s ^7as admin", args[0]))
		cr.logger.Printf("%s added %s as admin", player, args[0])
	})

	cr.RegisterCommand("removeadmin", "removea", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!removeadmin ^7<xuid>")
			return
		}

		err = database.RemoveAdmin(cr.db, args[0])
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to remove admin")
			return
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Removed ^5%s ^7from admins", args[0]))
		cr.logger.Printf("%s removed %s from admins", player, args[0])
	})
}
