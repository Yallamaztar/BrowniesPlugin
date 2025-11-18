// Optionally, if you are running IW4M-Admin and want to ignore BrowniesPlugin commands 
// from being logged or processed by IW4M-Admin, this script registers those commands 
// so that IW4M-Admin can recognize and ignore them

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
    scripts\_integration_shared::RegisterScriptCommand("Brownies_sayas",        "sayas", "ss", "say as another player", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_left",         "left", "l", "say as another player", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_onstart",      "onstart",      "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_killplayer",   "killplayer",   "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_hide",         "hide",         "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_spectator",    "spectator",    "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_sayto",        "sayto",        "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_takeweapons",  "takeweapons",  "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_giveweapon",   "giveweapon",   "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_freeze",       "freeze",       "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_slap",         "slap",         "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_loadout",      "loadout",      "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_teleport",     "teleport",     "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_setspeed",     "setspeed",     "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_setgravity",   "setgravity",   "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_dropgun",      "dropgun",      "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_toggleleft",   "toggleleft",   "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_changeteam",   "changeteam",   "",   "BrowniesPlugin", "User", "T6", false, ::null );
    scripts\_integration_shared::RegisterScriptCommand("Brownies_thirdperson",  "thirdperson",  "",   "BrowniesPlugin", "User", "T6", false, ::null );

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