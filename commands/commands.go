package commands

import (
	"database/sql"
	"strings"

	"github.com/Yallamaztar/BrowniesPlugin/database"
	"github.com/Yallamaztar/BrowniesPlugin/rcon"
)

type playerInfo struct {
	Name      string
	XUID      string
	clientNum int
}

func (cr *commandRegister) findPlayer(partial string) *playerInfo {
	partial = strings.ToLower(strings.TrimSpace(partial))

	status, err := cr.rcon.Status()
	if err != nil {
		return nil
	}

	var (
		exactMatch   *playerInfo
		partialMatch *playerInfo
	)

	for _, p := range status.Players {
		if strings.ToLower(p.Name) == partial {
			exactMatch = &playerInfo{
				Name:      p.Name,
				XUID:      p.GUID,
				clientNum: p.ClientNum,
			}
			break
		}

		if strings.Contains(strings.ToLower(p.Name), partial) && partialMatch == nil {
			partialMatch = &playerInfo{
				Name:      p.Name,
				XUID:      p.GUID,
				clientNum: p.ClientNum,
			}
		}
	}

	if exactMatch != nil {
		return exactMatch
	}

	if partialMatch != nil {
		return partialMatch
	}

	var (
		dbName string
		dbXUID string
	)

	query := `
		SELECT player, xuid 
		FROM wallets 
		WHERE LOWER(player) LIKE ? 
		ORDER BY created_at DESC 
		LIMIT 1
	`

	err = cr.db.QueryRow(query, "%"+partial+"%").Scan(&dbName, &dbXUID)
	if err != nil {
		return nil
	}

	if cn, ok := cr.GetClientNum(dbXUID); ok {
		return &playerInfo{
			Name:      dbName,
			XUID:      dbXUID,
			clientNum: cn,
		}
	}

	return &playerInfo{
		Name:      dbName,
		XUID:      dbXUID,
		clientNum: -1,
	}
}

func (cr *commandRegister) RegisterCommands(db *sql.DB, bank *database.Bank, rc *rcon.RCONClient) {
	RegisterOwnerCommands(cr, bank)
	RegisterAdminCommands(cr, bank)
	RegisterClientCommands(cr, bank)
	RegisterShopCommands(cr, bank)
}
