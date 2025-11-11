package commands

import (
	"fmt"
	"strconv"

	"github.com/Yallamaztar/BrowniesGambling/database"
)

func registerOwnerCommands(cr *commandRegister, bank *database.Bank) {
	cr.registerCommand("take", "ta", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
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

		wlt := database.GetWallet(t.Name, t.XUID, cr.db)
		bank.TransferFromWallet(wlt, amount)

		cr.rcon.Tell(clientNum, fmt.Sprintf("Took ^5$%d ^7from %s", amount, t.Name))
		if t.clientNum != -1 {
			cr.rcon.Tell(t.clientNum, fmt.Sprintf("^5%s ^7took ^5$%d from you", player, amount))
		}

		cr.logger.Printf("%s took $%d from %s", player, amount, t.Name)
	})

	cr.registerCommand("info", "if", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
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

	cr.registerCommand("delbanker", "del", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!delbanker ^7<xuid>")
			return
		}

		err = database.RemoveOwner(cr.db, args[0])
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to remove banker")
			return
		}
	})

	cr.registerCommand("addbanker", "add", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
			return
		}

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!addbanker ^7<player> <xuid>")
			return
		}

		err = database.AddOwner(cr.db, args[0], args[1])
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to add banker")
			return
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Added ^5%s ^7as banker", args[0]))
		cr.logger.Printf("%s added %s as banker", player, args[0])
	})

	cr.registerCommand("giveall", "ga", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
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

		count, err := database.GiveAllWallets(cr.db, amount)
		if err != nil {
			cr.rcon.Tell(clientNum, "Failed to give all wallets")
			return
		}
		cr.rcon.Say(fmt.Sprintf("[^5Gambling^7] Gave ^5$%d ^7to ^5%d ^7wallets", amount, count))
	})

	cr.registerCommand("give", "gi", func(clientNum int, player, xuid string, args []string) {
		owner, err := database.IsOwner(cr.db, xuid)
		if err != nil || !owner {
			cr.rcon.Tell(clientNum, "You dont have permission to use this command")
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

		bank.TransferToWallet(wlt, amount)
		cr.rcon.Tell(clientNum, fmt.Sprintf("Gave ^5$%d ^7to %s", amount, t.Name))
		if tcn, ok := cr.GetClientNum(t.XUID); ok {
			cr.rcon.Tell(tcn, fmt.Sprintf("You received ^5$%d ^7from %s", amount, player))
		}
		cr.logger.Printf("%s gave $%d to %s from bank", player, amount, t.Name)
	})
}
