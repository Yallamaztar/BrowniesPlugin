package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Yallamaztar/BrowniesPlugin/database"
	"github.com/Yallamaztar/BrowniesPlugin/helpers"
	"github.com/google/shlex"
)

func RegisterAdminCommands(cr *commandRegister, bank *database.Bank) {
	cr.RegisterCommand("jumpheight", "jh", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!jumpheight ^7<player (optional)> <height>")
			return
		}

		var target string
		var heightStr string
		if len(args) == 1 {
			target = player
			heightStr = args[0]
		} else {
			target = args[0]
			heightStr = args[1]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		h, err := strconv.ParseFloat(heightStr, 64)
		if err != nil || h <= 0 {
			cr.rcon.Tell(clientNum, "Invalid height")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("jumpheight %d %d %f", clientNum, t.clientNum, h))
	})

	cr.RegisterCommand("bunnyhop", "bhop", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("bunnyhop %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("freeze", "fz", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("freeze %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("thirdperson", "3rd", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("thirdperson %d %d", clientNum, t.clientNum))
	})

	// cr.RegisterCommand("changeteam", "ct", func(clientNum int, player, xuid string, args []string) {
	// 	isAdmin, _ := database.IsAdmin(cr.db, xuid)
	// 	isOwner, _ := database.IsOwner(cr.db, xuid)
	// 	if !isAdmin && !isOwner {
	// 		cr.rcon.Tell(clientNum, "You do not have permission to use this command")
	// 		return
	// 	}

	// 	if len(args) < 2 {
	// 		cr.rcon.Tell(clientNum, "Usage: ^5!changeteam ^7<player (optional)> <team>")
	// 		return
	// 	}

	// 	var target string
	// 	if len(args) == 0 {
	// 		target = player
	// 	} else {
	// 		target = args[0]
	// 	}

	// 	t := cr.findPlayer(target)
	// 	if t == nil || t.clientNum == -1 {
	// 		cr.rcon.Tell(clientNum, "Player not found")
	// 		return
	// 	}

	// 	team := strings.ToLower(args[len(args)-1])
	// 	if team != "allies" && team != "axis" && team != "spectator" {
	// 		cr.rcon.Tell(clientNum, "Invalid team | Valid teams are: allies, axis, spectator")
	// 		return
	// 	}

	// 	cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("changeteam %s %s %s", player, t.Name, team))
	// })

	cr.RegisterCommand("dropgun", "dg", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("dropgun %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("setgravity", "sg", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!setgravity ^7<player (optional)> <gravity>")
			return
		}

		var target string
		var gravityStr string
		if len(args) == 1 {
			target = player
			gravityStr = args[0]
		} else {
			target = args[0]
			gravityStr = args[1]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		gravity, err := strconv.ParseFloat(gravityStr, 64)
		if err != nil || gravity <= 0 {
			cr.rcon.Tell(clientNum, "Invalid gravity")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("setgravity %d %d %f", clientNum, t.clientNum, gravity))
	})

	cr.RegisterCommand("setspeed", "ss", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!setspeed ^7<player (optional)> <speed>")
			return
		}

		var target string
		var speedStr string
		if len(args) == 1 {
			target = player
			speedStr = args[0]
		} else {
			target = args[0]
			speedStr = args[1]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		speed, err := strconv.ParseFloat(speedStr, 64)
		if err != nil || speed <= 0 {
			cr.rcon.Tell(clientNum, "Invalid speed")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("setspeed %d %d %f", clientNum, t.clientNum, speed))
	})

	cr.RegisterCommand("killplayer", "kpl", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("killplayer %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("hide", "hd", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("hide %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("spectator", "spec", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		var target string
		if len(args) == 0 {
			target = player
		} else {
			target = args[0]
		}

		t := cr.findPlayer(target)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("spectator %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("teleport", "tp", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) == 0 {
			cr.rcon.Tell(clientNum, "Usage: ^5!teleport ^7<player (optional)> <target>")
			return
		}

		if len(args) == 1 {
			t := cr.findPlayer(args[0])
			if t == nil || t.clientNum == -1 {
				cr.rcon.Tell(clientNum, "Player not found")
				return
			}
			cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("teleport %d %d", clientNum, t.clientNum))
			return
		}

		src := cr.findPlayer(args[0])
		dst := cr.findPlayer(args[1])
		if src == nil || dst == nil || src.clientNum == -1 || dst.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}
		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("teleport %d %d %d", clientNum, src.clientNum, dst.clientNum))
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

		cr.rcon.Tell(clientNum, "Fast restart executed successfully")
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

	cr.RegisterCommand("sayas", "says", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!sayas ^7<player> <message> [-d] [-e]")
			return
		}

		var dead, enemy bool
		filtered := make([]string, 0, len(args))

		for _, arg := range args {
			switch arg {
			case "-d":
				dead = true
			case "-e":
				enemy = true
			default:
				filtered = append(filtered, arg)
			}
		}

		if len(filtered) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!sayas ^7<player> <message> [-d] [-e]")
			return
		}

		trgt := filtered[0]
		message := strings.Join(filtered[1:], " ")

		target := cr.findPlayer(trgt)
		if target == nil || target.clientNum == -1 {
			target = &playerInfo{Name: trgt}
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
		cr.rcon.Tell(clientNum, fmt.Sprintf("Message sent as %s", target.Name))
	})

	cr.RegisterCommand("takeweapons", "tw", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) < 1 {
			cr.rcon.Tell(clientNum, "Usage: ^5!takeweapons ^7<player>")
			return
		}

		t := cr.findPlayer(args[0])
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("takeweapons %d %d", clientNum, t.clientNum))
	})

	cr.RegisterCommand("giveweapon", "gw", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) < 2 {
			cr.rcon.Tell(clientNum, "Usage: ^5!giveweapon ^7<player> <weapon+ext>")
			return
		}

		t := cr.findPlayer(args[0])
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("giveweapon %d %d %s", clientNum, t.clientNum, args[1]))
	})

	cr.RegisterCommand("loadout", "ld", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

		if len(args) == 0 {
			cr.rcon.Tell(clientNum, "Usage: ^5!loadout ^7<player> <loadout> ^4or ^5!ld ^7<loadout>")
			return
		}

		var trgt, loadout string
		if len(args) == 1 {
			trgt = player
			loadout = args[0]
		} else {
			trgt = args[0]
			loadout = args[1]
		}

		t := cr.findPlayer(trgt)
		if t == nil || t.clientNum == -1 {
			cr.rcon.Tell(clientNum, "Player not found")
			return
		}

		cr.rcon.SetDvar("brwns_exec_in", fmt.Sprintf("loadout %d %d %s", clientNum, t.clientNum, loadout))
	})

	cr.RegisterCommand("take", "ta", func(clientNum int, player, xuid string, args []string) {
		isAdmin, _ := database.IsAdmin(cr.db, xuid)
		isOwner, _ := database.IsOwner(cr.db, xuid)
		if !isAdmin && !isOwner {
			cr.rcon.Tell(clientNum, "You do not have permission to use this command")
			return
		}

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

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

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

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

		parsed, err := shlex.Split(strings.Join(args, " "))
		if err != nil {
			cr.rcon.Tell(clientNum, "Invalid arguments")
			return
		}
		args = parsed

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
