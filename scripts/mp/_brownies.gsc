init() {
    level.browniesPrefix = "";
    level.brwns_cmds = [];

    scripts\mp\_brownies_core::setDvarIfUnitialized("brwns_enabled", 1 );
    scripts\mp\_brownies_core::setDvarIfUnitialized("brwns_exec", "" );

    scripts\mp\_brownies_cmds::registerAllCommands();
    level thread scripts\mp\_brownies_core::commandListenerLoop();
}

