// Optionally, if you are running IW4M-Admin and want to ignore BrowniesPlugin commands 
// from being logged or processed by IW4M-Admin, this script registers those commands 
// so that IW4M-Admin can recognize and ignore them
main() {
    ReplaceFunc(scripts\_integration_t6::RegisterClientCommands, ::RegisterCommands);
}

init() {
    waittillframeend;
    thread moduleSetup();
}

moduleSetup() {
    level waittill( level.notifyTypes.gameFunctionsInitialized );

     // Owner commands
    scripts\_integration_shared::RegisterScriptCommand("Brownies_gambling",     "gambling",     "gmbl",    "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_maxbet",       "maxbet",       "mb",      "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_printmoney",   "printmoney",   "print",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_addowner",     "addowner",     "add",     "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_removeowner",  "removeowner",  "remove",  "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_addadmin",     "addadmin",     "adda",    "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_removeadmin",  "removeadmin",  "remvovea","BrowniesPlugin", "User", "T6", false, ::null );

    // Admin commands
    scripts\_integration_shared::RegisterScriptCommand("Brownies_cheats",       "svcheat",      "svc",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_jump",         "jumpheight",   "jh",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_bunnyhop",     "bunnyhop",     "bh",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_sayas",        "sayas",        "ss",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_left",         "left",         "lt",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_onstart",      "onstart",      "",     "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_killplayer",   "killplayer",   "kpl",  "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_hide",         "hide",         "hd",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_spectator",    "spectator",    "spec", "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_sayto",        "sayto",        "st",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_takeweapons",  "takeweapons",  "tw",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_giveweapon",   "giveweapon",   "gw",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_freeze",       "freeze",       "fz",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_slap",         "slap",         "sl",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_loadout",      "loadout",      "ld",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_teleport",     "teleport",     "tp",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_setspeed",     "setspeed",     "ss",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_setgravity",   "setgravity",   "sg",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_dropgun",      "dropgun",      "dg",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_toggleleft",   "toggleleft",   "tl",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_thirdperson",  "thirdperson",  "3rd",   "BrowniesPlugin", "User", "T6", false, ::null );

    // Client commands
    scripts\_integration_shared::RegisterScriptCommand("Brownies_votekick",     "votekick",     "vk",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_left",         "left",         "lt",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_discord",      "discord",      "dc",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_richest",      "richest",      "rich", "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_poorest",      "poorest",      "poor", "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_pay",          "pay",          "pp",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_help",         "help",         "?",    "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_bankbalance",  "bankbalance",  "bank", "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_balance",      "balance",      "bal",  "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_gamble",       "gamble",       "g",    "BrowniesPlugin", "User", "T6", false, ::null );
}

null() { return; }

RegisterCommands() 
{
    scripts\_integration_base::AddClientCommand( "SwitchTeams",    true,  ::TeamSwitchImpl );
    scripts\_integration_base::AddClientCommand( "Alert",          true,  ::AlertImpl );
    scripts\_integration_base::AddClientCommand( "Goto",           false, ::GotoImpl );
    scripts\_integration_base::AddClientCommand( "PlayerToMe",     true,  ::PlayerToMeImpl );
}
