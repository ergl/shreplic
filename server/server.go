package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/vonaka/shreplic/curp"
	"github.com/vonaka/shreplic/epaxos"
	"github.com/vonaka/shreplic/master/defs"
	"github.com/vonaka/shreplic/n2paxos"
	"github.com/vonaka/shreplic/paxoi"
	"github.com/vonaka/shreplic/paxos"
	"github.com/vonaka/shreplic/unistore"
)

var (
	portnum     = flag.Int("port", 7070, "Port # to listen on")
	masterAddr  = flag.String("maddr", "", "Master address")
	masterPort  = flag.Int("mport", 7087, "Master port")
	myAddr      = flag.String("addr", "", "Server address (this machine)")
	doEpaxos    = flag.Bool("epaxos", false, "Use EPaxos as the replication protocol")
	doPaxoi     = flag.Bool("paxoi", false, "Use Paxoi as the replication protocol")
	doN2paxos   = flag.Bool("n2paxos", false, "Use n²Paxos as the replication protocol")
	doCurp      = flag.Bool("curp", false, "Use CURP as the replication protocol")
	doOptCurp   = flag.Bool("curpOpt", false, "Use optimized CURP as the replication protocol")
	doUnistore = flag.Bool("unistore", false, "Use unistore as the replication protocol")
	cpuprofile  = flag.String("cpuprofile", "", "Cpu profile")
	thrifty     = flag.Bool("thrifty", false, "Use only as many messages as strictly required")
	exec        = flag.Bool("exec", true, "Execute commands")
	optExec     = flag.Bool("optexec", false, "Execute commands optimistically")
	lread       = flag.Bool("lread", false, "Fast reads")
	dreply      = flag.Bool("dreply", true, "Reply to client only after command has been executed")
	beacon      = flag.Bool("beacon", false, "Send beacons to other replicas to compare their relative speeds")
	maxfailures = flag.Int("maxfailures", -1, "Maximum number of failures")
	durable     = flag.Bool("durable", false, "Log to a stable store")
	batchWait   = flag.Int("batchwait", 0, "Milliseconds to wait before sending a batch")
	tConf       = flag.Bool("tconf", true, "Conflict relation is transitive")
	proxy       = flag.String("proxy", "", "File with the list of clients IPs for this server")
	qfile       = flag.String("qfile", "", "Quorum config file")
	descNum     = flag.Int("desc", 100, "Number of command descriptors (only for Paxoi and n²Paxos)")
	poolLevel   = flag.Int("pool", 1, "Level of pool usage from 0 to 2 (only for Paxoi and n²Paxos)")
	AQreconf    = flag.Bool("AQreconf", true, "Automatically reconfigure Paxoi's slow active quorum")
	args        = flag.String("args", "", "Custom arguments")
)

func main() {
	flag.Parse()

	ps := make(map[string]struct{})
	if *proxy != "" {
		f, err := os.Open(*proxy)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			ps[s.Text()] = struct{}{}
		}
		f.Close()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGUSR1)
		go catchKill(interrupt)
	}

	log.Printf("Server starting on port %d", *portnum)
	fullAddr := fmt.Sprintf("%s:%d", *masterAddr, *masterPort)
	replicaId, nodeList, isLeader, err := registerWithMaster(fullAddr, 10, 100)
	if err != nil {
		log.Fatal("Couldn't connect to master, aborting")
		return
	}

	if *maxfailures == -1 {
		*maxfailures = (len(nodeList) - 1) / 2
	}
	log.Printf("Tolerating %d max. failures", *maxfailures)

	if *doEpaxos {
		log.Println("Starting Egalitarian Paxos replica...")
		rep := epaxos.NewReplica(replicaId, nodeList, *thrifty, *exec, *lread,
			*dreply, *beacon, *durable, *batchWait, *tConf, *maxfailures, ps)
		rpc.Register(rep)
	} else if *doUnistore {
		log.Println("Starting Unistore replica...")
		rep := unistore.NewReplica(replicaId, nodeList, *maxfailures, *exec, *dreply, *args, ps)
		rpc.Register(rep)
	} else if *doPaxoi {
		log.Println("Starting Paxoi replica...")
		paxoi.MaxDescRoutines = *descNum
		rep := paxoi.NewReplica(replicaId, nodeList, *exec, *lread,
			*dreply, *optExec, *AQreconf, *poolLevel, *maxfailures, *qfile, ps)
		rpc.Register(rep)
	} else if *doN2paxos {
		log.Println("Starting n²Paxos replica...")
		n2paxos.MaxDescRoutines = *descNum
		rep := n2paxos.NewReplica(replicaId, nodeList, *exec,
			*dreply, *optExec, *poolLevel, *maxfailures, *qfile, ps)
		rpc.Register(rep)
	} else if *doCurp {
		log.Println("Starting CURP replica...")
		curp.MaxDescRoutines = *descNum
		rep := curp.NewReplica(replicaId, nodeList, *exec,
			*dreply, *poolLevel, *maxfailures, *qfile, false, ps)
		rpc.Register(rep)
	} else if *doOptCurp {
		log.Println("Starting optimized CURP replica...")
		curp.MaxDescRoutines = *descNum
		rep := curp.NewReplica(replicaId, nodeList, *exec,
			*dreply, *poolLevel, *maxfailures, *qfile, true, ps)
		rpc.Register(rep)
	} else {
		log.Println("Starting Paxos replica...")
		rep := paxos.NewReplica(replicaId, nodeList, isLeader, *thrifty, *exec,
			*lread, *dreply, *durable, *batchWait, *maxfailures, ps)
		rpc.Register(rep)
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *portnum+1000))
	if err != nil {
		log.Fatal("listen error:", err)
	}

	http.Serve(l, nil)
}

func registerWithMaster(masterAddr string, retries int, backoff_ms int) (replicaId int, nodeList []string, isLeader bool, exit_err error) {
	var reply defs.RegisterReply
	args := &defs.RegisterArgs{
		Addr: *myAddr,
		Port: *portnum,
	}
	log.Printf("connecting to: %v", masterAddr)

	current_retry := 0
	for {
		mcli, err := rpc.DialHTTP("tcp", masterAddr)
		if err == nil {
			for {
				// TODO: This is an active wait, not cool.
				err = mcli.Call("Master.Register", args, &reply)
				if err == nil && reply.Ready {
					replicaId = reply.ReplicaId
					nodeList = reply.NodeList
					isLeader = reply.IsLeader
					return
				}
				if current_retry == retries {
					exit_err = err
					return
				}
				current_retry++
				log.Printf("Master.Register error: %v, retrying", err)
				time.Sleep(time.Duration(backoff_ms) * time.Millisecond)
			}
		}
		if current_retry == retries {
			exit_err = err
			return
		}
		current_retry++
		log.Printf("rpc.DialHTTP error: %v, retrying", err)
		time.Sleep(time.Duration(backoff_ms) * time.Millisecond)
	}
}

func catchKill(interrupt chan os.Signal) {
	<-interrupt
	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}
	fmt.Println("profing")
}
