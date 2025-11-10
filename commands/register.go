package commands

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Yallamaztar/PlutoRCON/rcon"
)

type commandHandler func(clientNum int, player, xuid string, args []string)

type commandRegister struct {
	commands map[string]commandHandler
	logger   *log.Logger
	rcon     *rcon.RCONClient
	db       *sql.DB
}

func New(logger *log.Logger, rc *rcon.RCONClient, db *sql.DB) *commandRegister {
	return &commandRegister{
		commands: make(map[string]commandHandler),
		logger:   logger,
		rcon:     rc,
		db:       db,
	}
}

func (cr *commandRegister) registerClientCommand(name, alias string, handler commandHandler) {
	cr.commands[strings.ToLower(name)] = handler
	if alias != "" {
		cr.commands[strings.ToLower(alias)] = handler
	}
	cr.logger.Printf("Registered command: !%s (!%s)", name, alias)
}

func (cr *commandRegister) Exec(command string, clientNum int, player, xuid string, args []string) bool {
	if handler, exists := cr.commands[strings.ToLower(command)]; exists {
		handler(clientNum, player, xuid, args)
		return true
	}

	return false
}
