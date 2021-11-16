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
	state              AgentState
	queue              *Queue = NewQueue()
	clock              Clock  = NewClock()
	server             Server
	stateLock          sync.Mutex
	queueLock          sync.Mutex
	enterTimestampLock sync.Mutex
	enterTimestamp     uint64 = 0
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
	go deadlockPoller()
	StartServer()
}

/**
This function has its purpose to manually wedge the system out of a deadlock,
but it is not successful in its attempt - even when it is a very bad way
to get out of the deadlock.
*/
func deadlockPoller() {
	timestamp := clock.GetCount()

	for {
		time.Sleep(5 * time.Second)

		timestampBuffer := clock.GetCount()
		if timestampBuffer == timestamp {
			log.Printf("[Lamport: %d] Deadlock detected in client: %s. Enter timestamp: %d, length of queue: %d, state: %d, nanoclock: %d", timestamp, server.addr, enterTimestamp, queue.Size(), state, clock.nanoclock.Nanosecond())
			// os.Exit(1)

			// Attempt to wedge the system out of the deadlock
			// Spoiler alert: this is not a good way to get out of the deadlock
			// Spoiler alert: this does not work :(

			if queue.Size() > 0 {
				for i := 0; i < 5; i++ {
					log.Println("Node is counting down to empty the queue")
					time.Sleep(time.Second)
				}

				// empty the queue to hopefully free up the deadlock
				for !queue.IsEmpty() {
					handle, _ := queue.Dequeue()
					handle <- clock.Increment()
				}
			}

		}

		timestamp = timestampBuffer
	}
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
	for {
		enter()
		log.Printf("[Lamport: %d] Node is now in the critical section!", clock.GetCount())
		time.Sleep(1 * time.Second)
		exit()
		time.Sleep(100 * time.Millisecond)
	}
}

func enter() {
	clock.Increment()

	enterTimestampLock.Lock()
	enterTimestamp = clock.GetCount()
	stateLock.Lock()
	state = Wanted
	log.Printf("[Lamport: %d] Attempting to enter critical section, state := WANTED, enterTimestamp := %d..\n", clock.GetCount(), enterTimestamp)
	stateLock.Unlock()
	enterTimestampLock.Unlock()

	// Multicast and wait for replies
	var wg sync.WaitGroup
	for nodeAddr, node := range server.Peers() {
		wg.Add(1)
		go func(addr string, node service.ServiceClient) {
			clock.Increment()
			log.Printf("[Lamport: %d] Sending request to %s..\n", clock.GetCount(), addr)
			msg := &service.RAMessage{Timestamp: enterTimestamp, Pid: server.addr}
			reply, err := node.Req(context.Background(), msg)
			log.Printf("[Lamport: %d] Replied to %s\n", clock.GetCount(), addr)
			if err != nil {
				log.Fatalf("[Lamport: %d] Request to %s failed with %v\n", clock.GetCount(), addr, err)
			}
			clock.Update(reply.Timestamp)
			log.Printf("[Lamport: %d] Received reply from %s\n", clock.GetCount(), addr)
			wg.Done()
		}(nodeAddr, node)
	}
	wg.Wait()

	stateLock.Lock()
	state = Held
	log.Printf("[Lamport: %d] Entered critical section, state := HELD\n", clock.GetCount())
	stateLock.Unlock()
}

func receive(Ti uint64, Pi string, handle ReplyHandle) {
	clock.Update(Ti)
	stateLock.Lock()
	enterTimestampLock.Lock()
	T, P := enterTimestamp, server.addr
	log.Printf("[Lamport: %d] Received request to enter critical section from node %s. T = %d, Ti = %d\n", clock.GetCount(), Pi, T, Ti)
	enterTimestampLock.Unlock()
	if state == Held || (state == Wanted && (T < Ti || (T == Ti && P < Pi))) {
		log.Printf("[Lamport: %d] Queued reply to %s\n", clock.GetCount(), Pi)
		log.Printf("[Lamport: %d] State: %v, T: %d, P: %s, Ti: %d, Pi: %s\n", clock.GetCount(), state, T, P, Ti, Pi)
		// Queue the reply to let the client node wait
		queueLock.Lock()
		queue.Enqueue(handle)
		queueLock.Unlock()
	} else {
		log.Printf("[Lamport: %d] Replying to %s..\n", clock.GetCount(), Pi)
		handle <- clock.Increment()

	}
	stateLock.Unlock()
	log.Printf("[Lamport: %d] Exiting receive\n", clock.GetCount())
}

func exit() {
	stateLock.Lock()
	state = Released
	log.Printf("[Lamport: %d] Attempting to exit critical section, state := RELEASED..\n", clock.GetCount())
	stateLock.Unlock()

	queueLock.Lock()
	defer queueLock.Unlock()
	log.Printf("[Lamport: %d] Replying to all in queue..\n", clock.GetCount())
	for !queue.IsEmpty() {
		handle, err := queue.Dequeue()
		if err != nil {
			log.Fatalf("[Lamport: %d] Queue crashed: %v", clock.GetCount(), err)
		}
		handle <- clock.Increment()
	}

	log.Printf("[Lamport: %d] Exited critical section\n", clock.GetCount())
}
