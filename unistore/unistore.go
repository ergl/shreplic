package unistore

import (
	"time"

	"github.com/vonaka/shreplic/server/smr"
	"github.com/vonaka/shreplic/tools/fastrpc"
)

type Replica struct {
	*smr.Replica
	cs CommunicationSupply
	//      ...
}

type CommunicationSupply struct {
	maxLatency time.Duration

	pingChan chan fastrpc.Serializable
	pongChan chan fastrpc.Serializable

	pingRPC uint8
	pongRPC uint8
}

func NewReplica(replicaId int, addrs []string, f int, exec, drep bool, args string, proxy_addrs map[string]struct{}) *Replica {
	r := &Replica{
		Replica: smr.NewReplica(replicaId, f, addrs, false, exec, false, drep, proxy_addrs),
		cs: CommunicationSupply{
			pingChan: make(chan fastrpc.Serializable, smr.CHAN_BUFFER_SIZE),
			pongChan: make(chan fastrpc.Serializable, smr.CHAN_BUFFER_SIZE),

			pingRPC: 0,
			pongRPC: 0,
		},
		//...
	}

	r.cs.pingRPC = r.RPC.Register(new(MPing), r.cs.pingChan)
	r.cs.pongRPC = r.RPC.Register(new(MPong), r.cs.pongChan)

	go r.run()

	return r
}

func (r *Replica) run() {
	r.ConnectToPeers()
	latencies := r.ComputeClosestPeers()
	for _, l := range latencies {
		d := time.Duration(l*1000*1000) * time.Nanosecond
		if d > r.cs.maxLatency {
			r.cs.maxLatency = d
		}
	}

	go r.WaitForClientConnections()

	for !r.Shutdown {
		select {
		case propose := <-r.ProposeChan:
			r.handlePropose(propose)

		case m := <-r.cs.pingChan:
			ping := m.(*MPing)
			r.handlePing(ping)

		case m := <-r.cs.pongChan:
			pong := m.(*MPong)
			r.handlePong(pong)
		}
	}
}

func (r *Replica) handlePropose(propose *smr.GPropose) {
	//...
}

func (r *Replica) handlePing(msg *MPing) {
	//...
}

func (r *Replica) handlePong(msg *MPong) {
	//...
}

func (m *MPing) New() fastrpc.Serializable {
	return new(MPing)
}

func (m *MPong) New() fastrpc.Serializable {
	return new(MPong)
}
