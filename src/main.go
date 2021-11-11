package main

import (
	"context"
	"log"
	"math/rand"
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
	state     AgentState
	queue     *Queue = NewQueue()
	clock     Clock  = NewClock()
	server    Server
	stateLock sync.Mutex
	queueLock sync.Mutex
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
	StartServer()

	go grabCriticalSection()
	neverExit := make(chan int)
	<-neverExit
}

func grabCriticalSection() {
	time.Sleep(5 * time.Second)
	// get the idea (not too often) that I want to enter the critical section
	for {
		number := rand.Intn(10)
		if number == 3 {
			// Now I want to enter the critical section
			log.Printf("[Lamport: %d] Node is thinking about entering the critical section.", clock.GetCount())
			enter()
			log.Printf("[Lamport: %d] Node is now in the critical section!", clock.GetCount())
			time.Sleep(5 * time.Millisecond)
			log.Printf("[Lamport: %d] Node is now bored and wants to get out of the critical section.", clock.GetCount())
			exit()
			log.Printf("[Lamport: %d] Node has now exited the critical section.", clock.GetCount())
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func enter() {
	clock.Increment()

	log.Printf("[Lamport: %d] Attempting to enter critical section..\n", clock.GetCount())
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
			msg := &service.RAMessage{Timestamp: clock.GetCount(), Pid: server.addr}
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

	log.Printf("[Lamport: %d] Entered critical section\n", clock.GetCount())
}

func receive(Ti uint64, Pi string, handle ReplyHandle) {
	clock.Update(Ti)
	T, P := clock.GetCount(), server.addr

	log.Printf("[Lamport: %d] Received request to enter critical section from node %s\n", clock.GetCount(), Pi)
	stateLock.Lock()

	if state == Held || (state == Wanted && (T < Ti || (T == Ti && P < Pi)) /* && (T, P) < (T_i, P_i) */) {
		log.Printf("[Lamport: %d] Queued reply to %s", clock.GetCount(), Pi)

		// Queue the reply to let the client node wait
		queueLock.Lock()
		queue.Enqueue(handle)
		queueLock.Unlock()
	} else {
		log.Printf("[Lamport: %d] Replying to %s..", clock.GetCount(), Pi)
		// Reply to req
		handle <- clock.Increment()
		log.Printf("[Lamport: %d] Replied to %s..", clock.GetCount(), Pi)
	}
	stateLock.Unlock()
}

func exit() {
	clock.Increment()

	log.Printf("[Lamport: %d] Attempting to exit critical section..\n", clock.GetCount())
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
