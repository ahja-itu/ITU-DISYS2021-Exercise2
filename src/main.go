package main

import (
	"log"

	qtip "github.com/phf/go-queue/queue" // we chose a good name, we know ;)
)

type AgentState int64

const (
	Released AgentState = 0
	Wanted   AgentState = 1
	Held     AgentState = 2
)

var (
	state AgentState
	queue *qtip.Queue = qtip.New()
	clock Clock       = NewClock()
)

// On initialisation of the node - i.e. when the program is started
// Go automatically runs this before main.
func init() {
	log.Println("Initialising node..")
	state = Released
	log.Println("Initialised node, state := RELEASED")
}

func main() {

}

func enter() {
	clock.Increment()

	// 	On enter do
	//   state := WANTED;
	//   “multicast ‘req(T,p)’”, where T := time of ‘req’
	//   wait for N-1 replies
	//   state := HELD;
	// End on
	log.Println("Attempting to enter critical section..")
	state = Wanted
	// multicast stuff begin

	// multicast stuff end
	state = Held
	log.Println("Entered critical section")
}

func receive( /* req (T_i, P_i) */ ) {
	// clock.Update(otherClock)

	nodeId := "TODO"
	log.Printf("Received request to enter critical section from node %s\n", nodeId)

	if state == Held || (state == Wanted /* && (T, P) < (T_i, P_i) */) {
		// Queue req
		// req = queue.PushBack()
	} else {
		// Reply to req
	}

}

func exit() {
	clock.Increment()

	log.Println("Attempting to exit critical section..")
	state = Released
	//reply to all in queue

	log.Println("Exited critical section")
}
