N2Paxos
-------

N2Paxos is the variation of Paxos in which `2b` message is broadcast
to all replicas. For a clients nearby the replicas, N2Paxos improves
latency by one message delay. Assuming that the master is running, to
launch a n2paxos server type:

    shr-server -n2paxos

Run a simple client with:

    shr-client -f