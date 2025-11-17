init() {
    level.browniesPrefix = "brwns_";
    level.brwns_cmds = [];

    scripts\mp\_brownies_core::setDvarIfUnitialized(level.browniesPrefix + "enabled", 1 );
    scripts\mp\_brownies_core::setDvarIfUnitialized(level.browniesPrefix + "exec_in", "" );
    scripts\mp\_brownies_core::setDvarIfUnitialized(level.browniesPrefix + "exec_out", "" );

    scripts\mp\_brownies_cmds::registerAllCommands();
    level thread scripts\mp\_brownies_core::commandListenerLoop();
}
