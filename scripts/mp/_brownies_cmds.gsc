registerAllCommands() {
    // server command to verify that the script is loaded
    scripts\mp\_brownies_core::RegisterCommand("onstart", ::onStart);

    // rcon commands
    scripts\mp\_brownies_core::RegisterCommand("killplayer",  ::impl_killplayer);
    scripts\mp\_brownies_core::RegisterCommand("hide",        ::impl_hideplayer);
    scripts\mp\_brownies_core::RegisterCommand("spectator",   ::impl_setspectator);
    scripts\mp\_brownies_core::RegisterCommand("sayto",       ::impl_sayto);
    scripts\mp\_brownies_core::RegisterCommand("takeweapons", ::impl_takeweapons);
    scripts\mp\_brownies_core::RegisterCommand("freeze",      ::impl_freezeplayer);
    scripts\mp\_brownies_core::RegisterCommand("slap",        ::impl_slapplayer);
    scripts\mp\_brownies_core::RegisterCommand("loadout",     ::impl_loadout);
    scripts\mp\_brownies_core::RegisterCommand("teleport",    ::impl_teleport);
    scripts\mp\_brownies_core::RegisterCommand("setspeed",    ::impl_setspeed);
    scripts\mp\_brownies_core::RegisterCommand("setgravity",  ::impl_setgravity);
    // scripts\mp\_brownies_core::RegisterCommand("giveweapon", ::impl_giveWeapon);
}

onStart(args) {
    SetDvar(level.browniesPrefix + "exec_out", "success");
    wait 0.75;
    SetDvar(level.browniesPrefix + "exec_out", "");
}

impl_killplayer(args) {
    if ( !isDefined(args) || args.size < 1 ) {
        return;
    }
        
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if ( isDefined(p) && IsAlive(p) ) {
        p Suicide();
    }
}

impl_hideplayer(args) {
    if ( !isDefined(args) || args.size < 1 ) {
        return;
    }
        
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (isDefined(p) && IsAlive(p)) {
        if (!isDefined(p.isHidden)) {
            p.isHidden = false;
        }
        if (!p.isHidden) {
            p Hide();
            p.isHidden = true;
        } else {
            p Show();
            p.isHidden = false;
        }

    }
}

impl_teleport(args) {
    if (args.size < 2) return; 
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
    if (!isDefined(p) || !isDefined(t) || !IsAlive(p) || !IsAlive(t)) {
        return;
    }
    p SetOrigin(t.origin);
}

impl_setspectator(args) {
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);

    if (!isDefined(p) || p.pers["team"] == "spectator") {
        return;
    }

    p [[level.spectator]]();
}

impl_sayto(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p)) return;

    msg = args[1];
    for (i = 2; i < args.size; i++)
        msg += " " + args[i];

    p IPrintLnBold(msg);
}

impl_giveweapon(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    weapon = args[1];
    p GiveWeapon(weapon);
    p SwitchToWeapon(weapon);
}

impl_takeweapons(args) {
    if (args.size < 1) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    p TakeAllWeapons();
}

impl_freezeplayer(args) {
    if (args.size < 1) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p)) return;

	if (!isDefined(p.isFrozen)) p.isFrozen = false;
    if (!p.isFrozen) {
        p FreezeControls(true);
        p.isFrozen = true;
    } else {
        p FreezeControls(false);
        p.isFrozen = false;
    }
}

impl_setspeed(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    speed = float(args[1]);
    p SetMoveSpeedScale(speed);
}

impl_slapplayer(args) {
    if (args.size < 1) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    dir = (randomInt(300) - 100, randomInt(500) - 100, 200);
    p SetVelocity(dir);
}

impl_loadout(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    if (args[1] == "ballista" || args[1] == "bal") {
        p TakeAllWeapons();
        wait 0.1;
        p GiveWeapon("ballista_mp", 0, randomIntRange(1, 45));
        p SwitchToWeapon("ballista_mp");
    } else if (args[1] == "dsr50" || args[1] == "dsr") {
        p TakeAllWeapons();
        wait 0.1;
        p GiveWeapon("dsr50_mp", 0, randomIntRange(1, 45));
        p SwitchToWeapon("dsr50_mp");
    } else {
        p TakeAllWeapons();
        wait 0.1;
        p GiveWeapon("dsr50_mp", 0, randomIntRange(1, 45));
        p SwitchToWeapon("dsr50_mp");
    }
}

impl_setgravity(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    gravity = float(args[1]);
    p SetClientDvar("bg_gravity", gravity);
}