setDvarIfUnitialized( d, v ) {
    if (!isDefined(getDvar( d )) || getDvar( d ) == "") {
        setDvar( d, v );
    }
}

RegisterCommand(name, handler) {
    if (!isDefined(level.brwns_cmds)) {
        level.brwns_cmds = [];
    }

    entry = spawnstruct();
    entry.name = toLower(name);
    entry.handler = handler;
    level.brwns_cmds[level.brwns_cmds.size] = entry;
}

findRegisteredCommand(name) {
    if (!isDefined(level.brwns_cmds)) {
        return undefined;
    }

    n = toLower(name);
    for ( i = 0; i < level.brwns_cmds.size; i++ ) {
        e = level.brwns_cmds[i];
        if (isDefined(e) && e.name == n) {
            return e;
        }
    }
    return undefined;
}

dispatchCommand(cmd) {
    if ( !isDefined(cmd) || cmd == "" ) {
        return;
    }

    tokens = strTok(cmd, " ");
    if ( !isDefined(tokens) || tokens.size == 0 ) {
        return;
    }

    cname = toLower(tokens[0]);
    args = [];
    for ( i = 1; i < tokens.size; i++ ) {
        args[args.size] = tokens[i];
    }

    def = findRegisteredCommand(cname);
    if (!isDefined(def)) {
        return;
    }

    thread [[def.handler]](args);
}

commandListenerLoop() {
    level endon("game_ended");
    for (;;) {
        if ( getDvarInt( level.browniesPrefix + "enabled" ) != 1 ) {
            wait 1;
            continue;
        }

        cmd = getDvar( level.browniesPrefix + "exec_in" );
        if ( cmd != "" ) {
            setDvar( level.browniesPrefix + "exec_in", "" );
            thread dispatchCommand(cmd);
        }

        wait 0.01;
    }
}

findPlayerByName(t) {
    v = toLower(t);
    for ( i = 0; i < level.players.size; i++ ) {
        p = level.players[i];
        if ( toLower( p.name ) == v || IsSubStr( toLower( p.name ), v ) )
            return p;
    }
    return undefined;
}