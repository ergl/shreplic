CURP
----

Assuming that the master is running, to launch a curp server type:

    shr-server -curp -optexec

And to start the improved variation in which replicas send accepts
directly to clients:

    shr-server -curpOpt -optexec

CURP supports only custom clients. Run a CURP client with:

    shr-client -curp -args "-N <number_of_replicas> -pclients 0" -f