package commands

import (
	"fmt"

	"github.com/Yallamaztar/BrowniesGambling/database"
)

func RegisterOwnerCommands(cr *commandRegister) {
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
