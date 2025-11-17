package commands

import (
	"fmt"
	"strconv"

	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/BrowniesGambling/helpers"
)

func RegisterOwnerCommands(cr *commandRegister, bank *database.Bank) {
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

	cr.RegisterCommand("removeadmin", "remvovea", func(clientNum int, player, xuid string, args []string) {
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
