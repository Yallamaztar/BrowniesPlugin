registerAllCommands() {
    // server command to verify that the script is loaded
    scripts\mp\_brownies_core::RegisterCommand("onstart", ::onStart);

    // rcon commands
    scripts\mp\_brownies_core::RegisterCommand("killplayer",   ::impl_killplayer);
    scripts\mp\_brownies_core::RegisterCommand("hide",         ::impl_hideplayer);
    scripts\mp\_brownies_core::RegisterCommand("spectator",    ::impl_setspectator);
    scripts\mp\_brownies_core::RegisterCommand("sayto",        ::impl_sayto);
    scripts\mp\_brownies_core::RegisterCommand("takeweapons",  ::impl_takeweapons);
    scripts\mp\_brownies_core::RegisterCommand("freeze",       ::impl_freezeplayer);
    scripts\mp\_brownies_core::RegisterCommand("slap",         ::impl_slapplayer);
    scripts\mp\_brownies_core::RegisterCommand("loadout",      ::impl_loadout);
    scripts\mp\_brownies_core::RegisterCommand("teleport",     ::impl_teleport);
    scripts\mp\_brownies_core::RegisterCommand("setspeed",     ::impl_setspeed);
    scripts\mp\_brownies_core::RegisterCommand("setgravity",   ::impl_setgravity);
    scripts\mp\_brownies_core::RegisterCommand("friendlyfire", ::impl_setfriendlyfire);
    scripts\mp\_brownies_core::RegisterCommand("dropgun",      ::impl_dropgun);
    scripts\mp\_brownies_core::RegisterCommand("toggleleft",   ::impl_toggleleft);
    // scripts\mp\_brownies_core::RegisterCommand("giveweapon", ::impl_giveWeapon);
}

onStart(args) {
    SetDvar(level.browniesPrefix + "exec_out", "success");
    wait 0.75;
    SetDvar(level.browniesPrefix + "exec_out", "");
}

impl_killplayer(args) {
    if (args.size < 2) return;

    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    if ( !isDefined(t) && !IsAlive(t) ) {
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    t Suicide();
    o IPrintLnBold("Killed ^5" + t.name);
}

impl_hideplayer(args) {
    if (args.size < 2) return;

    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    if (isDefined(t) && IsAlive(t)) {
        if (!isDefined(t.isHidden)) {
            t.isHidden = false;
        }
        if (!t.isHidden) {
            t Hide();
            t.isHidden = true;
            o IPrintLnBold("Hidden ^5" + t.name);
        } else {
            t Show();
            t.isHidden = false;
            o IPrintLnBold("Unhidden ^5" + t.name);
        }

    }
}

impl_teleport(args) {
    if (args.size < 2) return; 
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);

    if (args.size == 3) {
        t1 = scripts\mp\_brownies_core::findPlayerByName(args[1]);
        t2 = scripts\mp\_brownies_core::findPlayerByName(args[2]);

        if (!isDefined(t1) || !isDefined(t2) || !IsAlive(t1) || !IsAlive(t2)) {
            o IPrintLnBold("One or both target players not found or not alive");
            return;
        }

        t1 SetOrigin(t2.origin);
        o IPrintLnBold("Teleported ^5" + t1.name + " ^7to ^5" + t2.name);
        return;
    }

    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);  
    if (!isDefined(o) || !isDefined(t) || !IsAlive(o) || !IsAlive(t)) {
        return;
    }

    o SetOrigin(t.origin);
    o IPrintLnBold("Teleported to ^5" + t.name);
}

impl_setspectator(args) {
    if (args.size < 2) return;

    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    if (!isDefined(t) || t.pers["team"] == "spectator") {
        o IPrintLnBold("Player not found or already in spectator mode");
        return;
    }

    t [[level.spectator]]();
    o IPrintLnBold("Set ^5" + t.name + " ^7to spectator mode");
}

impl_sayto(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    msg = args[1];
    for (i = 2; i < args.size; i++)
        msg += " " + args[i];

    t IPrintLnBold(msg);
    o IPrintLnBold("Sent message to ^5" + t.name);    
}

impl_giveweapon(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    weapon = args[1];
    t GiveWeapon(weapon);
    t SwitchToWeapon(weapon);

    o IPrintLnBold("Gave ^5" + t.name + " ^7weapon: " + weapon);
}

impl_takeweapons(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    t TakeAllWeapons();
    o IPrintLnBold("Took weapons from ^5" + t.name);
}

impl_freezeplayer(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

	if (!isDefined(t.isFrozen)) t.isFrozen = false;
    if (!t.isFrozen) {
        t FreezeControls(true);
        t.isFrozen = true;
        o IPrintLnBold("Frozen ^5" + t.name);
    } else {
        t FreezeControls(false);
        t.isFrozen = false;
        o IPrintLnBold("Unfrozen ^5" + t.name);
    }
}

impl_setspeed(args) {
    if (args.size < 3) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    speed = float(args[2]);
    t SetMoveSpeedScale(speed);

    o IPrintLnBold("Set ^5" + t.name + " ^7speed to " + speed);
}

impl_slapplayer(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    dir = (randomInt(300) - 100, randomInt(500) - 100, 200);
    t SetVelocity(dir);

    o IPrintLnBold("Slapped ^5" + t.name);
}

impl_loadout(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    if (args[2] == "ballista" || args[2] == "bal" || args[2] == "1") {
        t TakeAllWeapons();
        wait 0.1;
        t GiveWeapon("ballista_mp", 0, randomIntRange(1, 45));
        t SwitchToWeapon("ballista_mp");
    } else if (args[2] == "dsr50" || args[2] == "dsr" || args[2] == "2") {
        t TakeAllWeapons();
        wait 0.1;
        t GiveWeapon("dsr50_mp", 0, randomIntRange(1, 45));
        t SwitchToWeapon("dsr50_mp");
    } else {
        t TakeAllWeapons();
        wait 0.1;
        t GiveWeapon("dsr50_mp", 0, randomIntRange(1, 45));
        t SwitchToWeapon("dsr50_mp");
    }

    o iPrintlnBold("Gave ^5" + t.name + " ^7loadout: " + args[2]);
}

impl_setgravity(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    if (!isDefined(t) || !IsAlive(t)) { 
        o IPrintLnBold("Player not found or not alive");
        return;
    }

    gravity = float(args[2]);
    t SetClientDvar("bg_gravity", gravity);

    o IPrintLnBold("Set ^5" + t.name + " ^7gravity to " + gravity);
}

impl_dropgun(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    w = t GetCurrentWeapon();
    t DropItem(w);

    o IPrintLnBold("Dropped ^5" + t.name + " ^7weapon");
}

impl_setfriendlyfire(args) {
    if (args.size < 2) return;
    
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);

    SetDvar("scr_team_fftype", args[1]);
    o IPrintLnBold("Set friendly fire to ^5" + args[1]);
}

impl_toggleleft(args) {
    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(o) || !IsAlive(o)) {
        return;
    }

    if (!isDefined(o.pers["lefttoggle"])) {
        o.pers["lefttoggle"] = false;
    }

    if(o.pers["lefttoggle"] == 1) {
        o SetClientDvar("cg_gun_y", "7");
        o.pers["lefttoggle"] = 0;
    } else {
        o SetClientDvar("cg_gun_y", "0");
        o.pers["lefttoggle"] = 1;
    }
}

impl_changeteam(args) {
    if (args.size < 3) return;

    o = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);

    if (!isDefined(o) || !isDefined(t)) {
        return;
    }

    team = ToLower(args[2]);
    if (team != "allies" && team != "axis" && team != "spectator") {
        o IPrintLnBold("Invalid team: allies / axis / spectator");
        return;
    }

    if (t.pers["team"] == team) {
        o IPrintLnBold("^5" + t.name + " ^7is already on team ^5" + team);
        return;
    }

    if (team == "allies") {
        t [[level.allies]]();
    } else if (team == "axis") {
        t [[level.axis]]();
    } else if (team == "spectator") {
        t [[level.spectator]]();
    }
}