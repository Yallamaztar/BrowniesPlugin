package commands

import (
	"database/sql"

	"github.com/Yallamaztar/BrowniesGambling/database"
	"github.com/Yallamaztar/BrowniesGambling/rcon"
)

type playerInfo struct {
	Name      string
	XUID      string
	clientNum int
}

func (cr *commandRegister) findPlayer(partialName string) *playerInfo {
	var playerName, xuid string
	query := "SELECT player, xuid FROM wallets WHERE player LIKE ? ORDER BY created_at DESC LIMIT 1"
	err := cr.db.QueryRow(query, "%"+partialName+"%").Scan(&playerName, &xuid)
	if err != nil {
		return nil
	}

	if cn, ok := cr.GetClientNum(xuid); ok {
		return &playerInfo{
			Name:      playerName,
			XUID:      xuid,
			clientNum: cn,
		}
	}

	return &playerInfo{
		Name:      playerName,
		XUID:      xuid,
		clientNum: -1,
	}
}

func (cr *commandRegister) RegisterCommands(db *sql.DB, bank *database.Bank, rc *rcon.RCONClient) {
	RegisterOwnerCommands(cr, bank)
	RegisterAdminCommands(cr, bank)
	RegisterClientCommands(cr, bank)
	RegisterShopCommands(cr, bank)
}
