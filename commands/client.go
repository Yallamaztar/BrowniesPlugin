package commands

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/BrowniesGambling/helpers"
)

func RegisterClientCommands(cr *commandRegister, bank *database.Bank) {
	cr.RegisterCommand("discord", "dc", func(clientNum int, player, xuid string, args []string) {
		cr.rcon.Tell(clientNum, "Join our Discord: ^5dsc.gg/browner")
	})

	cr.RegisterCommand("richest", "rich", func(clientNum int, player, xuid string, args []string) {
		wallets, err := database.Top5RichestWallets(cr.db)
		if err != nil || len(wallets) == 0 {
			cr.rcon.Tell(clientNum, "No wallets found")
			return
		}

		for i, rw := range wallets {
			cr.rcon.Say(fmt.Sprintf("[^5#%d^7] %s ^7- ^5$%s", i+1, rw.Player, helpers.FormatMoney(rw.Balance)))
		}
	})

	cr.RegisterCommand("poorest", "poor", func(clientNum int, player, xuid string, args []string) {
		wallets, err := database.Bottom5PoorestWallets(cr.db)
		if err != nil || len(wallets) == 0 {
			cr.rcon.Tell(clientNum, "No wallets found")
			return
		}

		for i, rw := range wallets {
			cr.rcon.Say(fmt.Sprintf("[^5#%d^7] %s ^7- ^5$%d", i+1, rw.Player, rw.Balance))
		}
	})

	cr.RegisterCommand("pay", "pp", func(clientNum int, player, xuid string, args []string) {
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

		cr.rcon.Tell(clientNum, fmt.Sprintf("Paid ^5$%s ^7to %s", helpers.FormatMoney(amount), t.Name))
		if targetCN, ok := cr.GetClientNum(t.XUID); ok {
			cr.rcon.Tell(targetCN, fmt.Sprintf("You received ^5$%s ^7from %s", helpers.FormatMoney(amount), player))
		}
	})

	cr.RegisterCommand("help", "?", func(clientNum int, player, xuid string, args []string) {
		page := ""
		if len(args) > 0 {
			page = args[0]
		}

		switch page {
		case "2":
			cr.rcon.Tell(clientNum, "Available commands:")
			cr.rcon.Tell(clientNum, "^5!balance ^7[player] (^5!bal^7) - Check wallet balance")
			cr.rcon.Tell(clientNum, "^5!pay ^7<player> <amount> (^5!pp^7) - Pay another player")
			cr.rcon.Tell(clientNum, "^5!bankbalance (^5!bank^7) - Check bank balance")
			cr.rcon.Tell(clientNum, "^5!help 3 ^7- More commands")

		case "3":
			cr.rcon.Tell(clientNum, "Available commands:")
			cr.rcon.Tell(clientNum, "^5!richest (^5!rich^7) - Show top 5 richest players")
			cr.rcon.Tell(clientNum, "^5!poorest (^5!poor^7) - Show bottom 5 poorest players")
			cr.rcon.Tell(clientNum, "^5!discord (^5!dc^7) - Get Discord invite link")

		default:
			cr.rcon.Tell(clientNum, "Available commands:")
			cr.rcon.Tell(clientNum, "^5!gamble ^7<amount> (^5!g^7) - Place a bet")
			cr.rcon.Tell(clientNum, "^5!shop (^5!sh^7) - View shop items")
			cr.rcon.Tell(clientNum, "^5!buy ^7<item|alias> <player (optional)> (^5!bu^7) - Buy an item from the shop")
			cr.rcon.Tell(clientNum, "^5!help 2 ^7- More commands")
		}

	})

	cr.RegisterCommand("bankbalance", "bank", func(clientNum int, player, xuid string, args []string) {
		cr.rcon.Tell(clientNum, fmt.Sprintf("Bank ^5balance: ^7$%s", helpers.FormatMoney(bank.Balance())))
	})

	cr.RegisterCommand("balance", "bal", func(clientNum int, player, xuid string, args []string) {
		if len(args) == 0 {
			wlt := database.GetWallet(player, xuid, cr.db)
			if wlt != nil {
				cr.rcon.Tell(clientNum, fmt.Sprintf("Your wallet ^5balance: ^7$%s", helpers.FormatMoney(wlt.Balance())))
			}
		} else {
			t := cr.findPlayer(args[0])
			wlt := database.GetWallet(t.Name, t.XUID, cr.db)
			if wlt != nil {
				if wlt.Balance() < 0 {
					cr.rcon.Tell(clientNum, fmt.Sprintf("%s ^5balance: ^7%s", t.Name, helpers.FormatMoney(wlt.Balance())))
					return
				}
				cr.rcon.Tell(clientNum, fmt.Sprintf("%s ^5balance: ^7$%s", t.Name, helpers.FormatMoney(wlt.Balance())))
			} else {
				cr.rcon.Tell(clientNum, "Player wallet not found")
			}
		}
	})

	cr.RegisterCommand("gamble", "g", func(clientNum int, player, xuid string, args []string) {
		if !database.IsGamblingEnabled(cr.db) {
			cr.rcon.Tell(clientNum, "Gambling is currently ^1disabled")
			return
		}

		if len(args) == 0 {
			cr.rcon.Tell(clientNum, "Usage: !gamble <amount>")
			return
		}

		wlt := database.GetWallet(player, xuid, cr.db)
		balance := wlt.Balance()

		if balance <= 0 {
			cr.rcon.Say(fmt.Sprintf("%s is ^F^1Gay n Poor", player))
			return
		}

		var bet int64
		arg := strings.ToLower(args[0])

		switch arg {
		case "all", "a":
			bet = balance
		case "half", "h":
			bet = balance / 2
		default:
			amt, err := strconv.ParseInt(arg, 10, 64)
			if err != nil || amt <= 0 {
				cr.rcon.Tell(clientNum, "Invalid amount")
				return
			}
			bet = amt
		}

		if max := database.GetMaxBet(cr.db); max > 0 && bet > max {
			cr.rcon.Tell(clientNum, fmt.Sprintf("Max bet is ^5$%s^7", helpers.FormatMoney(max)))
			return
		}

		if bet > balance {
			cr.rcon.Tell(clientNum, "Insufficient wallet balance")
			return
		}

		if rand.Float64() < 0.45 {
			bank.TransferToWallet(wlt, bet)
			cr.rcon.Tell(clientNum, fmt.Sprintf("You ^5won! ^7$%s", helpers.FormatMoney(bet)))
			cr.rcon.Say(fmt.Sprintf("%s just ^5won ^7$%s in gambling!", player, helpers.FormatMoney(bet)))
		} else {
			bank.TransferFromWallet(wlt, bet)
			cr.rcon.Tell(clientNum, fmt.Sprintf("You ^1lost! ^7$%s", helpers.FormatMoney(bet)))
			cr.rcon.Say(fmt.Sprintf("%s just ^1lost ^7$%s in gambling!", player, helpers.FormatMoney(bet)))
		}
	})
}
