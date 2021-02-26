Paxoi
-----

Assuming that the master is running, to launch a paxoi server type:

    shr-server -paxoi

And to start the improved variation with optimistic execution:

    shr-server -paxoi -optexec

This variation requires custom client:

    shr-client -paxoi -args "-N <number_of_replicas> -pclients 0" -f

For the vanilla paxoi run a simple client with:

    shr-client -f