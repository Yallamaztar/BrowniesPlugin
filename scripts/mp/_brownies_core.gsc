setDvarIfUnitialized( d, v ) {
    if ( !isDefined(getDvar( d )) || getDvar( d ) == "" )
        setDvar( d, v );
}

registerCommand(name, alias, handler, minArgs)
{
    if ( !isDefined(level.brwns_cmds) )
        level.brwns_cmds = [];

    entry = spawnstruct();
    entry.name = toLower(name);
    entry.alias = toLower(alias);
    entry.handler = handler;
    entry.minArgs = minArgs;
    level.brwns_cmds[level.brwns_cmds.size] = entry;
}

findRegisteredCommand(name)
{
    if ( !isDefined(level.brwns_cmds) )
        return undefined;

    n = toLower(name);
    for ( i = 0; i < level.brwns_cmds.size; i++ )
    {
        e = level.brwns_cmds[i];
        if ( isDefined(e) && (e.name == n || (isDefined(e.alias) && e.alias == n)) )
            return e;
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

    if ( isDefined(def.minArgs) && args.size < def.minArgs ) {
        return;
    }

    thread [[def.handler]](args);
}

commandListenerLoop() {
    level endon("game_ended");
    for (;;) {
        if ( getDvarInt( "brwns_enabled" ) != 1 ) {
            wait 1;
            continue;
        }

        cmd = getDvar( "brwns_exec" );
        if ( cmd != "" ) {
            setDvar( "brwns_exec", "" );
            thread dispatchCommand(cmd);
        }

        wait 1;
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