package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Yallamaztar/BrowniesGambling/database"
)

type PlayerInfo struct {
	Name      string
	XUID      string
	clientNum int
}

func (cr *commandRegister) findPlayer(partialName string) *PlayerInfo {
	var playerName, xuid string
	query := "SELECT player, xuid FROM wallets WHERE player LIKE ? ORDER BY created_at DESC LIMIT 1"
	err := cr.db.QueryRow(query, "%"+partialName+"%").Scan(&playerName, &xuid)
	if err != nil {
		fmt.Println("Error finding player by name:", err)
		return nil
	}

	if cn, ok := cr.GetClientNum(xuid); ok {
		return &PlayerInfo{
			Name:      playerName,
			XUID:      xuid,
			clientNum: cn,
		}
	}

	return &PlayerInfo{
		Name:      playerName,
		XUID:      xuid,
		clientNum: -1,
	}
}

func (cr *commandRegister) RegisterCommands(bank *database.Bank) {
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

	cr.registerCommand("pay", "pp", func(clientNum int, player, xuid string, args []string) {
		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!pay ^7<player> <amount>")
			return
		}

		wlt := database.GetWallet(player, xuid, cr.db)
		t := cr.findPlayer(args[0])
		if t == nil {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		if t.XUID == xuid {
			cr.rcon.Tell(clientNum, "Cant pay yourself")
			return
		}

		twlt := database.GetWallet(t.Name, t.XUID, cr.db)

		amount, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil || amount <= 0 {
			cr.rcon.Tell(clientNum, "Invalid amount")
			return
		}

		err = database.TransferFromWalletToWallet(wlt, twlt, amount)
		if err != nil {
			cr.rcon.Tell(clientNum, fmt.Sprintf("Failed to pay %s", t.Name))
			return
		}

		cr.rcon.Tell(clientNum, fmt.Sprintf("Paid ^5$%d ^7to %s", amount, t.Name))
		if targetCN, ok := cr.GetClientNum(t.XUID); ok {
			cr.rcon.Tell(targetCN, fmt.Sprintf("You received ^5$%d ^7from %s", amount, player))
		}
	})

	cr.registerCommand("help", "?", func(clientNum int, player, xuid string, args []string) {
		cr.rcon.Tell(clientNum, "^3Available commands:")
		cr.rcon.Tell(clientNum, "^5!balance ^7[player] (!bal) - Check wallet balance")
		cr.rcon.Tell(clientNum, "^5Usage: !balance, !balance PlayerName, !balance 5, !balance @xuid")
		cr.rcon.Tell(clientNum, "^5!gamble ^7<amount> (!g) - Place a bet")
		cr.rcon.Tell(clientNum, "^5!help ^7(!?) - Show this help")
	})

	cr.registerCommand("bankbalance", "bankbal", func(clientNum int, player, xuid string, args []string) {
		cr.rcon.Tell(clientNum, fmt.Sprintf("Bank ^5balance: ^7$%d", bank.Balance()))
	})

	cr.registerCommand("balance", "bal", func(clientNum int, player, xuid string, args []string) {
		if len(args) == 0 {
			wlt := database.GetWallet(player, xuid, cr.db)
			if wlt != nil {
				cr.rcon.Tell(clientNum, fmt.Sprintf("Your wallet ^5balance: ^7$%d", wlt.Balance()))
			}
		} else {
			target := args[0]
			t := cr.findPlayer(target)
			wlt := database.GetWallet(t.Name, t.XUID, cr.db)
			if wlt != nil {
				if wlt.Balance() < 0 {
					cr.rcon.Tell(clientNum, fmt.Sprintf("%s ^5balance: ^7%d$", t.Name, wlt.Balance()))
					return
				}
				cr.rcon.Tell(clientNum, fmt.Sprintf("%s ^5balance: ^7$%d", t.Name, wlt.Balance()))
			} else {
				cr.rcon.Tell(clientNum, "Player wallet not found")
			}
		}
	})

	cr.registerCommand("gamble", "g", func(clientNum int, player, xuid string, args []string) {
		if len(args) == 0 {
			cr.rcon.Tell(clientNum, "Usage: !gamble <amount>")
			return
		}

		wlt := database.GetWallet(player, xuid, cr.db)

		if strings.ToLower(args[0]) == "all" || strings.ToLower(args[0]) == "a" {
			bet := wlt.Balance()
			if bet <= 0 {
				cr.rcon.Say(fmt.Sprintf("%s is ^F^1Gay n Poor", player))
				return
			}

			win := (len(player) % 2) == 0
			if win {
				bank.TransferToWallet(wlt, bet)
				cr.rcon.Tell(clientNum, fmt.Sprintf("you ^5won! ^7$%d", bet))
			} else {
				bank.TransferFromWallet(wlt, bet)
				cr.rcon.Tell(clientNum, fmt.Sprintf("you ^1lost! ^7$%d", bet))
			}

		} else if strings.ToLower(args[0]) == "half" || strings.ToLower(args[0]) == "h" {
			bet := wlt.Balance() / 2
			if bet <= 0 {
				cr.rcon.Say(fmt.Sprintf("%s is ^F^1Gay n Poor", player))
				return
			}

			win := (len(player) % 2) == 0
			if win {
				bank.TransferToWallet(wlt, bet)
				cr.rcon.Tell(clientNum, fmt.Sprintf("you ^5won! ^7$%d", bet))
			} else {
				bank.TransferFromWallet(wlt, bet)
				cr.rcon.Tell(clientNum, fmt.Sprintf("you ^1lost! ^7$%d", bet))
			}

		} else {
			bet, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || bet <= 0 {
				cr.rcon.Tell(clientNum, "Invalid amount")
				return
			}

			if wlt.Balance() < bet {
				cr.rcon.Tell(clientNum, "Insufficient wallet balance")
				return
			}

			win := (len(player) % 2) == 0
			if win {
				bank.TransferToWallet(wlt, bet)
				cr.rcon.Tell(clientNum, fmt.Sprintf("you ^5won! ^7$%d", bet))
			} else {
				bank.TransferFromWallet(wlt, bet)
				cr.rcon.Tell(clientNum, fmt.Sprintf("you ^1lost! ^7$%d", bet))
			}
		}
	})
}
