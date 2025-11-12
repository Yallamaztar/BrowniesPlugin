package commands

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/Yallamaztar/BrowniesGambling/rcon"
)

type commandHandler func(clientNum int, player, xuid string, args []string)

type commandRegister struct {
	commands map[string]commandHandler
	logger   *log.Logger
	rcon     *rcon.RCONClient
	db       *sql.DB
	clients  map[string]int
	mu       sync.RWMutex
}

func New(logger *log.Logger, rc *rcon.RCONClient, db *sql.DB) *commandRegister {
	return &commandRegister{
		commands: make(map[string]commandHandler),
		logger:   logger,
		rcon:     rc,
		db:       db,
		clients:  make(map[string]int),
	}
}

func (cr *commandRegister) RegisterCommand(name, alias string, handler commandHandler) {
	cr.commands[strings.ToLower(name)] = handler
	if alias != "" {
		cr.commands[strings.ToLower(alias)] = handler
	}
	cr.logger.Printf("Registered command: !%s (!%s)", name, alias)
}

func (cr *commandRegister) Exec(command string, clientNum int, player, xuid string, args []string) bool {
	cr.mu.Lock()
	cr.clients[xuid] = clientNum
	cr.mu.Unlock()
	if handler, exists := cr.commands[strings.ToLower(command)]; exists {
		handler(clientNum, player, xuid, args)
		return true
	}

	return false
}

func (cr *commandRegister) SetClientNum(xuid string, clientNum int) {
	cr.mu.Lock()
	cr.clients[xuid] = clientNum
	cr.mu.Unlock()
}

func (cr *commandRegister) GetClientNum(xuid string) (int, bool) {
	cr.mu.RLock()
	cn, ok := cr.clients[xuid]
	cr.mu.RUnlock()
	return cn, ok
}
