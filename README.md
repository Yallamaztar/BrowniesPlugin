# BrowniesPlugin

A custom plugin and management tool for Call of Duty: Black Ops 2 (Plutonium T6) servers. Handles economy, administration, and Discord integration.

## Features

*   **Economy System**: Wallets, banking, gambling, and a shop system.
*   **Administration**: Player management (freeze, slap, kill, hide) and map controls.
*   **Discord Bot**: Bridge between Discord and in-game chat/commands.
*   **Custom Commands**: Teleport, godmode, third-person, and fun commands.
*   **IW4M-Admin Support**: Includes scripts to prevent command conflicts with IW4M-Admin.

## Setup

### 1. GSC Scripts
Copy the files from `scripts/mp/` to your server's script directory.

> **Note:** `_ignore_iw4m.gsc` is optional. Install it only if you use IW4M-Admin and want to prevent command conflicts.

### 2. Configuration
The application will generate a `config.json` on first run, or you can create one manually in the root directory:

```json
{
  "servers": [
    {
      "ip": "127.0.0.1",
      "port": "4976",
      "password": "your_rcon_password",
      "logPath": "C:\\path\\to\\games_mp.log"
    }
  ]
}
```

### 3. Add Owner
Before building, open `main.go` and add yourself as an owner to ensure you have full access. Find the `// Add Owners` section:

```go
database.AddOwner(db, "Player", "XUID")
```

### 4. Running
Build and run the application:

```bash
go build .
./BrowniesPlugin.exe
```

## Commands

### Client

| Command | Alias | Description |
| :--- | :--- | :--- |
| `!balance` | `!bal` | Check wallet balance. |
| `!bankbalance` | `!bank` | Check bank balance. |
| `!pay <player> <amount>` | `!pp` | Transfer money to another player. |
| `!gamble <amount>` | `!g` | Gamble money (50% win chance). |
| `!richest` | `!rich` | Show top 5 richest players. |
| `!poorest` | `!poor` | Show bottom 5 poorest players. |
| `!votekick <player>` | `!vk` | Start a votekick against a player. |
| `!discord` | `!dc` | Get the Discord invite link. |
| `!help` | `!?` | Show help menu. |
| `!left` | `!lt` | Toggle left-handed mode (if supported). |

### Admin

| Command | Alias | Description |
| :--- | :--- | :--- |
| `!thirdperson` | `!3rd` | Toggle thirdperson view for a player. |
| `!changeteam <player> <team>` | `!ct` | Change a player's team (allies/axis/spectator). |
| `!dropgun <player>` | `!dg` | Force a player to drop their current weapon. |
| `!friendlyfire <on/off>` | `!ff` | Toggle friendly fire. |
| `!setgravity <player> <value>` | `!sg` | Set a player's gravity. |
| `!setspeed <player> <value>` | `!ss` | Set a player's movement speed. |
| `!killplayer <player>` | `!kpl` | Kill a player. |
| `!hide <player>` | `!hd` | Toggle invisibility for a player. |
| `!spectator <player>` | `!spec` | Move a player to spectator mode. |
| `!teleport <player> <target>` | `!tp` | Teleport a player to another player. |
| `!fast` | `!res` | Perform a fast restart. |
| `!maprot` | `!mapr` | Rotate to the next map. |
| `!gadmins` | `!gay` | List all owners and admins. |
| `!sayas <player> <message>` | `!says` | Send a chat message as another player. |
| `!takeweapons <player>` | `!tw` | Remove all weapons from a player. |
| `!giveweapon <player> <weapon>` | `!gw` | Give a specific weapon to a player. |
| `!loadout <player> <loadout>` | `!ld` | Give a pre-defined loadout to a player. |
| `!take <player> <amount>` | `!ta` | Take money from a player's wallet. |
| `!info <player>` | `!if` | Get player information (XUID, ClientNum). |
| `!giveall <amount>` | `!ga` | Give money to all registered wallets. |
| `!give <player> <amount>` | `!gi` | Give money to a player's wallet. |

### Owner

| Command | Alias | Description |
| :--- | :--- | :--- |
| `!gambling <enable/disable/status>` | `!gmbl` | Manage the gambling system. |
| `!maxbet <amount/0>` | `!mb` | Set the maximum bet limit (0 to disable). |
| `!printmoney <amount>` | `!print` | Add money to your own wallet. |
| `!addowner <player> <xuid>` | `!add` | Add a new owner. |
| `!removeowner <xuid>` | `!remove` | Remove an owner. |
| `!addadmin <player> <xuid>` | `!adda` | Add a new admin. |
| `!removeadmin <xuid>` | `!removea` | Remove an admin. |
