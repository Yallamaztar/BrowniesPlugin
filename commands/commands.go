package commands

import (
	"fmt"

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
	registerOwnerCommands(cr, bank)
	registerClientCommands(cr, bank)
}
