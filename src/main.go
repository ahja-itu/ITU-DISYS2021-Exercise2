package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/andreaswachs/ITU-DISYS2021-Exercise2/src/service"
)

type AgentState int64

const (
	Released AgentState = 0
	Wanted   AgentState = 1
	Held     AgentState = 2
)

var (
	state          AgentState
	queue          *Queue = NewQueue()
	clock          Clock  = NewClock()
	server         Server
	stateLock      sync.Mutex
	queueLock      sync.Mutex
	TLock          sync.Mutex
	enterTimestamp uint64 = 0
)

// On initialisation of the node - i.e. when the program is started
// Go automatically runs this before main.
func init() {
	log.Printf("[Lamport: %d] Initialising node..\n", clock.GetCount())
	stateLock.Lock()
	defer stateLock.Unlock()
	state = Released
	log.Printf("[Lamport: %d] Initialised node, state := RELEASED\n", clock.GetCount())
}

func main() {
	go grabCriticalSection()
	StartServer()

	neverExit := make(chan int)
	<-neverExit
}

func waitForNodesToComeOnline() {
	n, err := strconv.Atoi(os.Getenv("NODES"))

	if err != nil {
		log.Fatalf("NODES env variable was not able to be converted to number. It was: %s", os.Getenv("NODES"))
	}

	for len(server.nodes) < n {
		time.Sleep(25 * time.Millisecond)
	}
}

func grabCriticalSection() {
	waitForNodesToComeOnline()
	// get the idea (not too often) that I want to enter the critical section
	for {

		// Now I want to enter the critical section
		// log.Printf("[Lamport: %d] Node is thinking about entering the critical section.", clock.GetCount())
		enter()
		log.Printf("[Lamport: %d] Node is now in the critical section!", clock.GetCount())
		time.Sleep(1 * time.Millisecond)
		// log.Printf("[Lamport: %d] Node is now bored and wants to get out of the critical section.", clock.GetCount())
		exit()
		// log.Printf("[Lamport: %d] Node has now exited the critical section.", clock.GetCount())

		time.Sleep(50 * time.Millisecond)
	}
}

func enter() {
	clock.Increment()
	log.Println("enter()")
	TLock.Lock()
	enterTimestamp = clock.GetCount()
	log.Printf("Set enterTimestap to %d", enterTimestamp)
	TLock.Unlock()

	log.Printf("[Lamport: %d] Attempting to enter critical section, state := WANTED..\n", clock.GetCount())
	stateLock.Lock()
	state = Wanted
	stateLock.Unlock()

	// Multicast and wait for replies
	var wg sync.WaitGroup
	for nodeAddr, node := range server.Peers() {
		wg.Add(1)
		go func(addr string, node service.ServiceClient) {
			log.Printf("[Lamport: %d] Sending request to %s..", clock.GetCount(), addr)
			// We set the PID to the server's address. This will be unique on
			// the local machine, but not necessarily over the network (for
			// example if two nodes both host on localhost:5000).
			clock.Increment()
			msg := &service.RAMessage{Timestamp: enterTimestamp, Pid: server.addr}
			reply, err := node.Req(context.Background(), msg)
			if err != nil {
				log.Fatalf("[Lamport: %d] Request to %s failed with %v", clock.GetCount(), addr, err)
			}
			clock.Update(reply.Timestamp)
			log.Printf("[Lamport: %d] Received reply from %s", clock.GetCount(), addr)
			wg.Done()
		}(nodeAddr, node)
	}

	wg.Wait()

	stateLock.Lock()
	state = Held
	stateLock.Unlock()

	log.Printf("[Lamport: %d] Entered critical section, state := HELD\n", clock.GetCount())
}

func receive(Ti uint64, Pi string, handle ReplyHandle) {
	log.Println("receive()")
	clock.Update(Ti)

	log.Printf("[Lamport: %d] Received request to enter critical section from node %s\n", clock.GetCount(), Pi)
	stateLock.Lock()
	TLock.Lock()
	T, P := enterTimestamp, server.addr
	log.Printf("T in recieve %d", T)
	TLock.Unlock()
	if state == Held || (state == Wanted && (T < Ti || (T == Ti && P < Pi)) /* && (T, P) < (T_i, P_i) */) {
		log.Printf("[Lamport: %d] Queued reply to %s", clock.GetCount(), Pi)
		log.Printf("[Lamport: %d] State: %v, T: %d, P: %s, Ti: %d, Pi: %s", clock.GetCount(), state, T, P, Ti, Pi)
		// Queue the reply to let the client node wait
		queueLock.Lock()
		queue.Enqueue(handle)
		queueLock.Unlock()
	} else {
		log.Printf("[Lamport: %d] Replying to %s..", clock.GetCount(), Pi)
		// log.Printf("[Lamport: %d] T: %d, P: %s, Ti: %d, Pi: %s", T, T, P, Ti, Pi)
		// Reply to req
		handle <- clock.Increment()
		log.Printf("[Lamport: %d] Replied to %s..", clock.GetCount(), Pi)
	}
	log.Printf("[Lamport: %d] Exiting receive", clock.GetCount())
	stateLock.Unlock()
}

func exit() {
	log.Println("exit()")
	clock.Increment()

	log.Printf("[Lamport: %d] Attempting to exit critical section, state := RELEASED..\n", clock.GetCount())
	stateLock.Lock()
	state = Released
	stateLock.Unlock()

	// Reply to all in queue
	queueLock.Lock()
	defer queueLock.Unlock()
	for !queue.IsEmpty() {
		handle, err := queue.Dequeue()
		if err != nil {
			log.Fatalf("[Lamport: %d] Queue crashed: %v", clock.GetCount(), err)
		}
		handle <- clock.Increment()
	}

	log.Printf("[Lamport: %d] Exited critical section\n", clock.GetCount())
}
