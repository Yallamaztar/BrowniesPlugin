registerAllCommands()
{
    scripts\mp\_brownies_core::registerCommand("killplayer",  "kpl",  ::killPlayer,    1);
    scripts\mp\_brownies_core::registerCommand("hide",        "hd",   ::hidePlayer,    1);
    scripts\mp\_brownies_core::registerCommand("teleport",    "tp",   ::teleport,      2);
    scripts\mp\_brownies_core::registerCommand("spectator",   "spec", ::setSpectator,  1);
    scripts\mp\_brownies_core::registerCommand("sayto",       "st",   ::sayTo,         2);
    scripts\mp\_brownies_core::registerCommand("giveweapon",  "gw",   ::giveWeapon,    2);
    scripts\mp\_brownies_core::registerCommand("takeweapons", "tw",   ::takeWeapons,   1);
    scripts\mp\_brownies_core::registerCommand("freeze",      "frz",  ::freezePlayer,  1);
    scripts\mp\_brownies_core::registerCommand("setspeed",    "ss",   ::setSpeed,      2);
    scripts\mp\_brownies_core::registerCommand("slap",        "sl",   ::slapPlayer,    1);
}

killPlayer(args) {
    if ( !isDefined(args) || args.size < 1 ) {
        return;
    }
        
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if ( isDefined(p) && IsAlive(p) ) {
        p Suicide();
    }
}

hidePlayer(args) {
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

teleport(args) {
    if (args.size == 2) {
        p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
        t = scripts\mp\_brownies_core::findPlayerByName(args[1]);
    
        if (!isDefined(p) || !isDefined(t) || !IsAlive(p) || !IsAlive(t)) {
            return;
        }
        p SetOrigin(t.origin);
        
    } else if (args.size >= 4) {
        x = int(args[1]);
        y = int(args[2]);
        z = int(args[3]);

        p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
        if (!isDefined(p) || !IsAlive(p)) {
            return;
        }
        p SetOrigin( (x, y, z) );
    }
}

setSpectator(args) {
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);

    if (!isDefined(p) || p.pers["team"] == "spectator") {
        return;
    }

    p [[level.spectator]]();
}

sayTo(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p)) return;

    msg = args[1];
    for (i = 2; i < args.size; i++)
        msg += " " + args[i];

    p IPrintLnBold(msg);
}

giveWeapon(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    weapon = args[1];
    p GiveWeapon(weapon);
    p SwitchToWeapon(weapon);
}

takeWeapons(args) {
    if (args.size < 1) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    p TakeAllWeapons();
}

freezePlayer(args) {
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

setSpeed(args) {
    if (args.size < 2) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    speed = float(args[1]);
    p SetMoveSpeedScale(speed);
}

slapPlayer(args) {
    if (args.size < 1) return;
    p = scripts\mp\_brownies_core::findPlayerByName(args[0]);
    if (!isDefined(p) || !IsAlive(p)) return;

    dir = (randomInt(300) - 100, randomInt(500) - 100, 200);
    p SetVelocity(dir);
}