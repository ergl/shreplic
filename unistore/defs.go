package unistore

import "github.com/vonaka/shreplic/state"

type MPing struct {
	Replica int32
	Ballot  int32
	Cmd     state.Command
	CmdSlot int
}

type MPong struct {
	Replica int32
	Ballot  int32
	CmdSlot int
}
