package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vonaka/shreplic/master/defs"
	"github.com/vonaka/shreplic/server/smr"
)

var (
	portnum  = flag.Int("port", 7087, "Port to listen on")
	numNodes = flag.Int("N", 3, "Number of replicas")
)

type Master struct {
	N          int
	nodeList   []string
	addrList   []string
	portList   []int
	lock       *sync.Mutex
	nodes      []*rpc.Client
	leader     []bool
	alive      []bool
	latencies  []float64
	finishInit bool
	initCond   *sync.Cond
	nextLeader int
}

func main() {
	flag.Parse()

	log.Printf("Master starting on port %d", *portnum)
	log.Printf("...waiting for %d replicas", *numNodes)

	master := &Master{
		N:          *numNodes,
		nodeList:   make([]string, 0, *numNodes),
		addrList:   make([]string, 0, *numNodes),
		portList:   make([]int, 0, *numNodes),
		lock:       new(sync.Mutex),
		nodes:      make([]*rpc.Client, *numNodes),
		leader:     make([]bool, *numNodes),
		alive:      make([]bool, *numNodes),
		latencies:  make([]float64, *numNodes),
		finishInit: false,
		nextLeader: -1,
	}
	master.initCond = sync.NewCond(master.lock)

	rpc.Register(master)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
	if err != nil {
		log.Fatal("Master listen error:", err)
	}

	go master.run()

	http.Serve(l, nil)
}

func (master *Master) run() {
	for {
		master.lock.Lock()
		if len(master.nodeList) == master.N {
			master.lock.Unlock()
			break
		}
		master.lock.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)

	for i := 0; i < master.N; {
		var err error
		addr := fmt.Sprintf("%s:%d", master.addrList[i], master.portList[i]+1000)
		master.nodes[i], err = rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Printf("Error connecting to replica %d (%v), retrying...", i, addr)
			time.Sleep(1 * time.Second)
		} else {
			btlReply := smr.NewBeTheLeaderReply()
			if master.leader[i] {
				err = master.nodes[i].Call("Replica.BeTheLeader",
					new(smr.BeTheLeaderArgs), btlReply)
				if err != nil {
					log.Fatal("Not today Zurg!")
				}
				smr.UpdateBeTheLeaderReply(btlReply)
				if btlReply.Leader != -1 && btlReply.Leader != int32(i) {
					master.leader[i] = false
					master.leader[int(btlReply.Leader)] = true
				}
				master.nextLeader = int(btlReply.NextLeader)
			}
			i++
		}
	}

	var new_leader bool
	pingNode := func(i int, node *rpc.Client) {
		err := node.Call("Replica.Ping", new(smr.PingArgs), new(smr.PingReply))
		if err != nil {
			master.alive[i] = false
			if master.leader[i] {
				new_leader = true
				master.leader[i] = false
			}
		} else {
			master.alive[i] = true
		}
	}
	master.lock.Lock()
	for i, node := range master.nodes {
		pingNode(i, node)
	}
	// initialization is finished
	// (i.e., slice `alive` has been computed)
	master.finishInit = true
	master.initCond.Broadcast()
	master.lock.Unlock()

	beTheLeader := func(i int) error {
		if master.alive[i] {
			btlReply := smr.NewBeTheLeaderReply()
			err := master.nodes[i].Call("Replica.BeTheLeader",
				new(smr.BeTheLeaderArgs), btlReply)
			if err == nil {
				smr.UpdateBeTheLeaderReply(btlReply)
				leaderI := i
				if btlReply.Leader != -1 {
					leaderI = int(btlReply.Leader)
				}
				master.leader[leaderI] = true
				master.nextLeader = int(btlReply.NextLeader)
				log.Printf("Replica %d is the new leader", leaderI)
				return nil
			}
			return err
		}
		return errors.New("dead")
	}

	for {
		time.Sleep(3 * time.Second)
		new_leader = false
		for i, node := range master.nodes {
			pingNode(i, node)
		}

		if !new_leader {
			continue
		}
		if master.nextLeader != -1 {
			if beTheLeader(master.nextLeader) == nil {
				continue
			}
		}
		for i := range master.nodes {
			if beTheLeader(i) == nil {
				break
			}
		}
	}
}

func (master *Master) Register(args *defs.RegisterArgs, reply *defs.RegisterReply) error {
	master.lock.Lock()
	defer master.lock.Unlock()

	nlen := len(master.nodeList)
	index := nlen

	addrPort := fmt.Sprintf("%s:%d", args.Addr, args.Port)

	for i, ap := range master.nodeList {
		if addrPort == ap {
			index = i
			break
		}
	}

	if index == nlen {
		master.nodeList = master.nodeList[0 : nlen+1]
		master.nodeList[nlen] = addrPort
		master.addrList = master.addrList[0 : nlen+1]
		master.addrList[nlen] = args.Addr
		master.portList = master.portList[0 : nlen+1]
		master.portList[nlen] = args.Port
		master.leader[index] = false
		nlen++

		addr := args.Addr
		if addr == "" {
			addr = "127.0.0.1"
		}
		out, err := exec.Command("ping", addr, "-c 2", "-q").Output()
		if err == nil {
			master.latencies[index], _ =
				strconv.ParseFloat(strings.Split(string(out), "/")[4], 64)
			log.Printf("node %v [%v] -> %vms away", index,
				master.nodeList[index], master.latencies[index])
		} else {
			log.Fatal("cannot determine ping latency to " + addr)
		}
	}

	if nlen == master.N {
		reply.Ready = true
		reply.ReplicaId = index
		reply.NodeList = master.nodeList
		reply.IsLeader = false

		minLatency := math.MaxFloat64
		leader := 0

		for i := 0; i < len(master.leader); i++ {
			if master.latencies[i] < minLatency {
				minLatency = master.latencies[i]
				leader = i
			}
		}

		if leader == index {
			log.Printf("Replica %d is the new leader", index)
			master.leader[index] = true
			reply.IsLeader = true
		}

	} else {
		reply.Ready = false
	}

	return nil
}

func (master *Master) GetLeader(args *defs.GetLeaderArgs, reply *defs.GetLeaderReply) error {
	master.lock.Lock()
	defer master.lock.Unlock()

	for i, l := range master.leader {
		if l {
			*reply = defs.GetLeaderReply{
				LeaderId: i,
			}
			break
		}
	}
	return nil
}

func (master *Master) GetReplicaList(args *defs.GetReplicaListArgs, reply *defs.GetReplicaListReply) error {
	master.lock.Lock()

	for !master.finishInit {
		master.initCond.Wait()
	}

	if len(master.nodeList) == master.N {
		reply.Ready = true
	} else {
		reply.Ready = false
	}

	reply.ReplicaList = make([]string, 0)
	reply.AliveList = make([]bool, 0)
	for i, node := range master.nodeList {
		reply.ReplicaList = append(reply.ReplicaList, node)
		reply.AliveList = append(reply.AliveList, master.alive[i])
	}

	log.Printf("nodes list %v", reply.ReplicaList)
	master.lock.Unlock()
	return nil
}
