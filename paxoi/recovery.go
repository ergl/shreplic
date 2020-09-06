package paxoi

import (
	"sync"
)

func (r *Replica) handleNewLeader(nl *MNewLeader) {
	if r.ballot >= nl.Ballot {
		return
	}

	r.status = RECOVERING
	r.ballot = nl.Ballot
	r.stopDescs()
	// TODO: send NewLeaderAck
	// stop processing normal channels:
	for r.status == RECOVERING {
		select {
		case m := <-r.cs.newLeaderChan:
			newLeader := m.(*MNewLeader)
			r.handleNewLeader(newLeader)

		case m := <-r.cs.newLeaderAckChan:
			newLeaderAck := m.(*MNewLeaderAck)
			r.handleNewLeaderAck(newLeaderAck)

		case m := <-r.cs.syncChan:
			sync := m.(*MSync)
			r.handleSync(sync)
		}
	}
}

func (r *Replica) stopDescs() {
	var wg sync.WaitGroup
	r.cmdDescs.IterCb(func(_ string, v interface{}) {
		desc := v.(*commandDesc)
		if !desc.seq {
			go func() {
				wg.Add(1)
				desc.stopChan <- &wg
			}()
		}
	})
	wg.Wait()
}
