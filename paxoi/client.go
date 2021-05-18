package paxoi

import (
	"flag"
	"log"
	"strings"

	"github.com/vonaka/shreplic/client/base"
	"github.com/vonaka/shreplic/server/smr"
	"github.com/vonaka/shreplic/state"
)

type Client struct {
	*base.SimpleClient

	val       state.Value
	ready     chan struct{}
	ballot    int32
	delivered map[CommandId]struct{}

	SQ        smr.QuorumI
	FQ        smr.QuorumI
	slowPathH *smr.MsgSet
	fastPathH *smr.MsgSet

	fixedMajority bool

	slowPaths   int
	alreadySlow map[CommandId]struct{}

	cs CommunicationSupply
}

func NewClient(maddr, collocated string, mport, reqNum, writes, psize, conflict int,
	fast, lread, leaderless, verbose bool, logger *log.Logger, args string) *Client {

	// args must be of the form "-N <rep_num>"
	f := flag.NewFlagSet("custom Paxoi arguments", flag.ExitOnError)
	repNum := f.Int("N", -1, "Number of replicas")
	f.Int("pclients", 0, "Number of clients already running on other machines")
	f.Parse(strings.Fields(args))
	if *repNum == -1 {
		f.Usage()
		return nil
	}

	c := &Client{
		SimpleClient: base.NewSimpleClient(maddr, collocated, mport, reqNum, writes,
			psize, conflict, fast, lread, leaderless, verbose, logger),

		val:       nil,
		ready:     make(chan struct{}, 1),
		ballot:    -1,
		delivered: make(map[CommandId]struct{}),

		SQ: smr.NewMajorityOf(*repNum),
		FQ: smr.NewThreeQuartersOf(*repNum),

		fixedMajority: true,

		slowPaths:   0,
		alreadySlow: make(map[CommandId]struct{}),
	}

	c.ReadTable = true
	c.WaitResponse = func() error {
		<-c.ready
		return nil
	}

	if c.fixedMajority {
		c.FQ = smr.NewMajorityOf(*repNum)
	}

	initCs(&c.cs, c.RPC)
	c.reinitFastAndSlowAcks()

	go c.handleMsgs()

	return c
}

func (c *Client) reinitFastAndSlowAcks() {
	accept := func(msg, leaderMsg interface{}) bool {
		if leaderMsg == nil {
			return true
		}
		leaderFastAck := leaderMsg.(*MFastAck)
		fastAck := msg.(*MFastAck)
		return fastAck.Checksum == nil || SHashesEq(leaderFastAck.Checksum, fastAck.Checksum)
	}

	free := func(msg interface{}) {
		switch f := msg.(type) {
		case *MFastAck:
			releaseFastAck(f)
		}
	}

	c.slowPathH.Free()
	c.slowPathH = c.slowPathH.ReinitMsgSet(c.SQ, accept, free, c.handleFastAndSlowAcks)
	c.fastPathH.Free()
	c.fastPathH = c.fastPathH.ReinitMsgSet(c.FQ, accept, free, c.handleFastAndSlowAcks)
}

func (c *Client) handleMsgs() {
	for {
		select {
		case m := <-c.cs.replyChan:
			reply := m.(*MReply)
			c.handleReply(reply)

		/*case m := <-c.cs.readReplyChan:
			reply := m.(*MReadReply)
			c.handleReadReply(reply)*/

		case m := <-c.cs.acceptChan:
			accept := m.(*MAccept)
			c.handleAccept(accept)

		case m := <-c.cs.fastAckChan:
			fastAck := m.(*MFastAck)
			c.handleFastAck(fastAck, false)

		case m := <-c.cs.lightSlowAckChan:
			lightSlowAck := m.(*MLightSlowAck)
			c.handleLightSlowAck(lightSlowAck)

		case m := <-c.cs.acksChan:
			acks := m.(*MAcks)
			for _, f := range acks.FastAcks {
				c.handleFastAck(copyFastAck(&f), false)
			}
			for _, s := range acks.LightSlowAcks {
				ls := s
				c.handleLightSlowAck(&ls)
			}

		case m := <-c.cs.optAcksChan:
			optAcks := m.(*MOptAcks)
			for _, ack := range optAcks.Acks {
				fastAck := newFastAck()
				fastAck.Replica = optAcks.Replica
				fastAck.Ballot = optAcks.Ballot
				fastAck.CmdId = ack.CmdId
				if !IsNilDepOfCmdId(ack.CmdId, ack.Dep) {
					fastAck.Checksum = ack.Checksum
				} else {
					if _, exists := c.alreadySlow[fastAck.CmdId]; !c.Reading && !exists {
						c.slowPaths++
						c.alreadySlow[fastAck.CmdId] = struct{}{}
					}
					fastAck.Checksum = nil
				}
				c.handleFastAck(fastAck, false)
				if _, exists := c.delivered[fastAck.CmdId]; !exists && fastAck.Checksum == nil {
					//c.Println("got slowAck", fastAck.Replica, fastAck.CmdId)
					//c.Println("wanna add", fastAck.Replica, fastAck.CmdId, false)
					c.slowPathH.Add(fastAck.Replica, false, fastAck)
				}
			}
		}
	}
}

func (c *Client) handleFastAck(f *MFastAck, fromLeader bool) bool {
	if c.ballot == -1 {
		c.ballot = f.Ballot
	} else if c.ballot < f.Ballot {
		c.ballot = f.Ballot
		c.reinitFastAndSlowAcks()
	} else if c.ballot > f.Ballot {
		return false
	}

	if _, exists := c.delivered[f.CmdId]; exists {
		return false
	}

	//c.Println("got fastAck", f.Replica, f.CmdId, f.Checksum)
	c.fastPathH.Add(f.Replica, fromLeader, f)
	return true
}

func (c *Client) handleLightSlowAck(ls *MLightSlowAck) {
	if _, exists := c.delivered[ls.CmdId]; exists {
		return
	}

	if _, exists := c.alreadySlow[ls.CmdId]; !c.Reading && !exists {
		c.slowPaths++
		c.alreadySlow[ls.CmdId] = struct{}{}
	}

	f := newFastAck()
	f.Replica = ls.Replica
	f.Ballot = ls.Ballot
	f.CmdId = ls.CmdId
	f.Checksum = nil
	c.handleFastAck(f, false)
	if _, exists := c.delivered[f.CmdId]; !exists {
		//c.Println("wanna add", f.Replica, f.CmdId, false)
		//c.Println("got slowAck", f.Replica, f.CmdId)
		c.slowPathH.Add(f.Replica, false, f)
	}
	// if c.handleFastAck(f, false) {
	// 	c.slowPathH.Add(f.Replica, false, f)
	// }
}

func (c *Client) handleFastAndSlowAcks(leaderMsg interface{}, msgs []interface{}) {
	if leaderMsg == nil {
		return
	}

	cmdId := leaderMsg.(*MFastAck).CmdId
	if _, exists := c.delivered[cmdId]; exists {
		return
	}
	c.delivered[cmdId] = struct{}{}

	c.Println("Slow Paths:", c.slowPaths)
	c.Println("Returning:", c.val.String())
	c.reinitFastAndSlowAcks()
	c.ResChan <- c.val
	c.ready <- struct{}{}
}

func (c *Client) handleReply(r *MReply) {
	if _, exists := c.delivered[r.CmdId]; exists {
		return
	}
	f := newFastAck()
	f.Replica = r.Replica
	f.Ballot = r.Ballot
	f.CmdId = r.CmdId
	f.Checksum = r.Checksum
	c.val = r.Rep
	c.handleFastAck(f, true)
	if _, exists := c.delivered[f.CmdId]; !exists {
		//c.Println("wanna add", f.Replica, f.CmdId, true)
		//c.Println("got Reply", r.Replica, r.CmdId)
		c.slowPathH.Add(f.Replica, true, f)
	}
	// if c.handleFastAck(f, true) {
	// 	c.slowPathH.Add(f.Replica, true, f)
	// }
}

func (c *Client) handleAccept(a *MAccept) {
	if _, exists := c.delivered[a.CmdId]; exists {
		return
	}
	c.delivered[a.CmdId] = struct{}{}

	c.val = a.Rep
	c.Println("Slow Paths:", c.slowPaths)
	c.Println("Returning:", c.val.String())
	c.reinitFastAndSlowAcks()
	c.ResChan <- c.val
	c.ready <- struct{}{}
}

// func (c *Client) handleReadReply(r *MReadReply) {
// 	if _, exists := c.delivered[r.CmdId]; exists {
// 		return
// 	}
// 	f := newFastAck()
// 	f.Replica = r.Replica
// 	f.Ballot = r.Ballot
// 	f.CmdId = r.CmdId
// 	f.Checksum = nil
// 	c.val = r.Rep
// 	c.handleFastAck(f, true) // <-- this `true` is not a bug
// }
