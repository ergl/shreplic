shreplic
========

Shreplic is a framework that simplifies the implementation of
broadcast protocols for SMR-replication in message-passing systems
where any minority of processes can fail by crashing. Based on message
definitions it automatically generates *"boring"* communication and
marshaling code and embeds it in the main application.

Originally fork of [Pierre Sutra's fork][otrack] of [epaxos], shreplic
improves on both of its predecessors by fixing several bugs and adding
new protocols and functionality.

Installation
------------

    go get github.com/vonaka/shreplic
    go install github.com/vonaka/shreplic
    shreplic -i

Supported protocols
-------------------

|  Name                  | Comments                                    |
|------------------------|---------------------------------------------|
| [Paxos][paxos_src]     | Sutra's version with some minor fixes.      |
| [N2Paxos][n2paxos_src] | All-to-all variant of Paxos.                |
| [EPaxos][epaxos_src]   | Slightly [improved][epaxos_fix] version of Sutra's fork<br />in which read operations are excluded from dependencies<br />of other read requests performed on the same key. |
| [Paxoi][paxoi_src]     | -                                           |
| [CURP][curp_src]       | -                                           |

Add new protocol
----------------

To add a new protocol named `shmaxos`, first create a go file with the
message definitions, let's call it `def.go`. Each message is a go
structure, which name starts with a capital `M`:

```go
package shmaxos

import "github.com/vonaka/shreplic/state"

type M2A struct {
	Replica int32
	Ballot  int32
	Cmd     state.Command
	CmdSlot int
}

type M2B struct {
	Replica int32
	Ballot  int32
	CmdSlot int
}
```

Then pass this definition to shreplic:

    shreplic -n shmaxos -msg def.go

This will generate the following tree:

    shmaxos
    ├── defs.go
    ├── Makefile
    ├── proto.go
    └── shmaxos.go

After editing `shmaxos.go` run `shreplic -i`.

Usage
-----

First start master, which coordinates communications between clients
and (in this case three) servers:

    shr-master -N 3

Run each server with the appropriate options:

    shr-server -shmaxos

Once all servers are ready start a client:

    shr-client -q 100

[otrack]: https://github.com/otrack/epaxos
[epaxos]: https://github.com/efficient/epaxos
[epaxos_fix]: https://github.com/vonaka/shreplic/commit/5e4dcb5736dd3c4d3e87aeb18f67c4371e3c429c
[paxos_src]: https://github.com/vonaka/shreplic/tree/master/paxos
[n2paxos_src]: https://github.com/vonaka/shreplic/tree/master/n2paxos
[epaxos_src]: https://github.com/vonaka/shreplic/tree/master/epaxos
[paxoi_src]: https://github.com/vonaka/shreplic/tree/master/paxoi
[curp_src]: https://github.com/vonaka/shreplic/tree/master/curp