# ITU-DISYS2021-Exercise2

Mandatory exercise 2 for the DISYS course at ITU for Autumn 2021

We make use of **Docker**, **docker-compose** and **GNU make** in this project.

## Notes on the project

The solution is prone to hitting a deadlock. This is observed when two nodes in rapid succession wishes to enter the critical section. They manage to communicate to each other and correctly either sends a reply or queues the other, depending on which request has the lower timestamp. The third node however does not replied to either before the deadlock is reached, at which point the third notes has only been observed to be in state RELEASED. This, however, only happens sporadically, and most times the program will handle rapid requests with no problem. We have unfortunately not been able to identify any further distinguishing features of the deadlock, that might help discover the causing problem.

We've tried covering this up with a questionable solution of having a deadlock poller, where we've tried to manually cover up the deadlock, but our assumption of replying to the one message left in the queue in the single node were wrong.

## Prerequisites

* Docker
* docker-compose
* make
* protoc: Protocol Buffers compiler with Golang support

### Usage

#### Building the Protocol Buffers spec's golang files

```bash
make proto
```

#### Running the project natively

```bash
make run
```

#### Building the docker image

```bash
make docker-build
```

#### Running the project with 3 nodes, with docker-compose

```bash
make up
```

#### Stopping the 3 nodes

```bash
make down
```

#### Building the docker image and booting the 3 nodes (best for development)

```bash
make refresh
```

#### Log output from short run

```text
client1_1  | 2021/11/16 12:24:08 [Lamport: 0] Initialising node..
client1_1  | 2021/11/16 12:24:08 [Lamport: 0] Initialised node, state := RELEASED
client3_1  | 2021/11/16 12:24:08 [Lamport: 0] Initialising node..
client3_1  | 2021/11/16 12:24:08 [Lamport: 0] Initialised node, state := RELEASED
client2_1  | 2021/11/16 12:24:08 [Lamport: 0] Initialising node..
client2_1  | 2021/11/16 12:24:08 [Lamport: 0] Initialised node, state := RELEASED
client1_1  | 2021/11/16 12:24:08 Started server on 0.0.0.0:5000
client2_1  | 2021/11/16 12:24:08 Started server on 0.0.0.0:5001
client3_1  | 2021/11/16 12:24:08 Started server on 0.0.0.0:5002
client1_1  | 2021/11/16 12:24:10 Connecting to client3:5002..
client1_1  | 2021/11/16 12:24:10 Connecting to client2:5001..
client1_1  | 2021/11/16 12:24:10 Connected to client2:5001
client1_1  | 2021/11/16 12:24:10 Connected to client3:5002
client1_1  | 2021/11/16 12:24:10 [Lamport: 1] Attempting to enter critical section, state := WANTED, enterTimestamp := 1..
client1_1  | 2021/11/16 12:24:10 [Lamport: 2] Sending request to client2:5001..
client2_1  | 2021/11/16 12:24:10 [Lamport: 2] Received request to enter critical section from node client1. T = 0, Ti = 1
client2_1  | 2021/11/16 12:24:10 [Lamport: 2] Replying to client1..
client2_1  | 2021/11/16 12:24:10 [Lamport: 3] Exiting receive
client1_1  | 2021/11/16 12:24:10 [Lamport: 3] Sending request to client3:5002..
client3_1  | 2021/11/16 12:24:10 [Lamport: 2] Received request to enter critical section from node client1. T = 0, Ti = 1
client3_1  | 2021/11/16 12:24:10 [Lamport: 2] Replying to client1..
client3_1  | 2021/11/16 12:24:10 [Lamport: 3] Exiting receive
client1_1  | 2021/11/16 12:24:10 [Lamport: 3] Replied to client3:5002
client1_1  | 2021/11/16 12:24:10 [Lamport: 4] Received reply from client3:5002
client1_1  | 2021/11/16 12:24:10 [Lamport: 4] Replied to client2:5001
client1_1  | 2021/11/16 12:24:10 [Lamport: 5] Received reply from client2:5001
client1_1  | 2021/11/16 12:24:10 [Lamport: 5] Entered critical section, state := HELD
client1_1  | 2021/11/16 12:24:10 [Lamport: 5] Node is now in the critical section!
client2_1  | 2021/11/16 12:24:10 Connecting to client3:5002..
client2_1  | 2021/11/16 12:24:10 Connecting to client1:5000..
client2_1  | 2021/11/16 12:24:10 Connected to client1:5000
client2_1  | 2021/11/16 12:24:10 Connected to client3:5002
client3_1  | 2021/11/16 12:24:10 Connecting to client1:5000..
client3_1  | 2021/11/16 12:24:10 Connecting to client2:5001..
client2_1  | 2021/11/16 12:24:10 [Lamport: 4] Attempting to enter critical section, state := WANTED, enterTimestamp := 4..
client2_1  | 2021/11/16 12:24:10 [Lamport: 5] Sending request to client3:5002..
client3_1  | 2021/11/16 12:24:10 Connected to client1:5000
client2_1  | 2021/11/16 12:24:10 [Lamport: 6] Sending request to client1:5000..
client1_1  | 2021/11/16 12:24:10 [Lamport: 6] Received request to enter critical section from node client2. T = 1, Ti = 4
client1_1  | 2021/11/16 12:24:10 [Lamport: 6] Queued reply to client2
client1_1  | 2021/11/16 12:24:10 [Lamport: 6] State: 2, T: 1, P: client1, Ti: 4, Pi: client2
client1_1  | 2021/11/16 12:24:10 [Lamport: 6] Exiting receive
client2_1  | 2021/11/16 12:24:10 [Lamport: 6] Replied to client3:5002
client2_1  | 2021/11/16 12:24:10 [Lamport: 7] Received reply from client3:5002
client3_1  | 2021/11/16 12:24:10 [Lamport: 5] Received request to enter critical section from node client2. T = 0, Ti = 4
client3_1  | 2021/11/16 12:24:10 [Lamport: 5] Replying to client2..
client3_1  | 2021/11/16 12:24:10 [Lamport: 6] Exiting receive
client3_1  | 2021/11/16 12:24:10 Connected to client2:5001
client3_1  | 2021/11/16 12:24:10 [Lamport: 7] Attempting to enter critical section, state := WANTED, enterTimestamp := 7..
client3_1  | 2021/11/16 12:24:10 [Lamport: 8] Sending request to client1:5000..
client3_1  | 2021/11/16 12:24:10 [Lamport: 9] Sending request to client2:5001..
client2_1  | 2021/11/16 12:24:10 [Lamport: 8] Received request to enter critical section from node client3. T = 4, Ti = 7
client2_1  | 2021/11/16 12:24:10 [Lamport: 8] Queued reply to client3
client1_1  | 2021/11/16 12:24:10 [Lamport: 8] Received request to enter critical section from node client3. T = 1, Ti = 7
client1_1  | 2021/11/16 12:24:10 [Lamport: 8] Queued reply to client3
client1_1  | 2021/11/16 12:24:10 [Lamport: 8] State: 2, T: 1, P: client1, Ti: 7, Pi: client3
client1_1  | 2021/11/16 12:24:10 [Lamport: 8] Exiting receive
client2_1  | 2021/11/16 12:24:10 [Lamport: 8] State: 1, T: 4, P: client2, Ti: 7, Pi: client3
client2_1  | 2021/11/16 12:24:10 [Lamport: 8] Exiting receive
client1_1  | 2021/11/16 12:24:11 [Lamport: 8] Attempting to exit critical section, state := RELEASED..
client1_1  | 2021/11/16 12:24:11 [Lamport: 8] Replying to all in queue..
client1_1  | 2021/11/16 12:24:11 [Lamport: 10] Exited critical section
client3_1  | 2021/11/16 12:24:11 [Lamport: 9] Replied to client1:5000
client3_1  | 2021/11/16 12:24:11 [Lamport: 11] Received reply from client1:5000
client2_1  | 2021/11/16 12:24:11 [Lamport: 8] Replied to client1:5000
client2_1  | 2021/11/16 12:24:11 [Lamport: 10] Received reply from client1:5000
client2_1  | 2021/11/16 12:24:11 [Lamport: 10] Entered critical section, state := HELD
client2_1  | 2021/11/16 12:24:11 [Lamport: 10] Node is now in the critical section!
client1_1  | 2021/11/16 12:24:11 [Lamport: 11] Attempting to enter critical section, state := WANTED, enterTimestamp := 11..
client1_1  | 2021/11/16 12:24:11 [Lamport: 13] Sending request to client3:5002..
client1_1  | 2021/11/16 12:24:11 [Lamport: 12] Sending request to client2:5001..
client3_1  | 2021/11/16 12:24:11 [Lamport: 12] Received request to enter critical section from node client1. T = 7, Ti = 11
client3_1  | 2021/11/16 12:24:11 [Lamport: 12] Queued reply to client1
client3_1  | 2021/11/16 12:24:11 [Lamport: 12] State: 1, T: 7, P: client3, Ti: 11, Pi: client1
client3_1  | 2021/11/16 12:24:11 [Lamport: 12] Exiting receive
client2_1  | 2021/11/16 12:24:11 [Lamport: 12] Received request to enter critical section from node client1. T = 4, Ti = 11
client2_1  | 2021/11/16 12:24:11 [Lamport: 12] Queued reply to client1
client2_1  | 2021/11/16 12:24:11 [Lamport: 12] State: 2, T: 4, P: client2, Ti: 11, Pi: client1
client2_1  | 2021/11/16 12:24:11 [Lamport: 12] Exiting receive
client2_1  | 2021/11/16 12:24:12 [Lamport: 12] Attempting to exit critical section, state := RELEASED..
client2_1  | 2021/11/16 12:24:12 [Lamport: 12] Replying to all in queue..
client2_1  | 2021/11/16 12:24:12 [Lamport: 14] Exited critical section
client1_1  | 2021/11/16 12:24:12 [Lamport: 13] Replied to client2:5001
client1_1  | 2021/11/16 12:24:12 [Lamport: 15] Received reply from client2:5001
client3_1  | 2021/11/16 12:24:12 [Lamport: 12] Replied to client2:5001
client3_1  | 2021/11/16 12:24:12 [Lamport: 14] Received reply from client2:5001
client3_1  | 2021/11/16 12:24:12 [Lamport: 14] Entered critical section, state := HELD
client3_1  | 2021/11/16 12:24:12 [Lamport: 14] Node is now in the critical section!
client2_1  | 2021/11/16 12:24:12 [Lamport: 15] Attempting to enter critical section, state := WANTED, enterTimestamp := 15..
client2_1  | 2021/11/16 12:24:12 [Lamport: 16] Sending request to client3:5002..
client2_1  | 2021/11/16 12:24:12 [Lamport: 17] Sending request to client1:5000..
client1_1  | 2021/11/16 12:24:12 [Lamport: 16] Received request to enter critical section from node client2. T = 11, Ti = 15
client1_1  | 2021/11/16 12:24:12 [Lamport: 16] Queued reply to client2
client1_1  | 2021/11/16 12:24:12 [Lamport: 16] State: 1, T: 11, P: client1, Ti: 15, Pi: client2
client1_1  | 2021/11/16 12:24:12 [Lamport: 16] Exiting receive
client3_1  | 2021/11/16 12:24:12 [Lamport: 16] Received request to enter critical section from node client2. T = 7, Ti = 15
client3_1  | 2021/11/16 12:24:12 [Lamport: 16] Queued reply to client2
client3_1  | 2021/11/16 12:24:12 [Lamport: 16] State: 2, T: 7, P: client3, Ti: 15, Pi: client2
client3_1  | 2021/11/16 12:24:12 [Lamport: 16] Exiting receive
client3_1  | 2021/11/16 12:24:13 [Lamport: 16] Attempting to exit critical section, state := RELEASED..
client3_1  | 2021/11/16 12:24:13 [Lamport: 16] Replying to all in queue..
client3_1  | 2021/11/16 12:24:13 [Lamport: 18] Exited critical section
client2_1  | 2021/11/16 12:24:13 [Lamport: 17] Replied to client3:5002
client2_1  | 2021/11/16 12:24:13 [Lamport: 19] Received reply from client3:5002
client1_1  | 2021/11/16 12:24:13 [Lamport: 16] Replied to client3:5002
client1_1  | 2021/11/16 12:24:13 [Lamport: 18] Received reply from client3:5002
client1_1  | 2021/11/16 12:24:13 [Lamport: 18] Entered critical section, state := HELD
client1_1  | 2021/11/16 12:24:13 [Lamport: 18] Node is now in the critical section!
client3_1  | 2021/11/16 12:24:13 [Lamport: 19] Attempting to enter critical section, state := WANTED, enterTimestamp := 19..
client3_1  | 2021/11/16 12:24:13 [Lamport: 20] Sending request to client1:5000..
client3_1  | 2021/11/16 12:24:13 [Lamport: 21] Sending request to client2:5001..
client1_1  | 2021/11/16 12:24:13 [Lamport: 20] Received request to enter critical section from node client3. T = 11, Ti = 19
client1_1  | 2021/11/16 12:24:13 [Lamport: 20] Queued reply to client3
client1_1  | 2021/11/16 12:24:13 [Lamport: 20] State: 2, T: 11, P: client1, Ti: 19, Pi: client3
client1_1  | 2021/11/16 12:24:13 [Lamport: 20] Exiting receive
client2_1  | 2021/11/16 12:24:13 [Lamport: 20] Received request to enter critical section from node client3. T = 15, Ti = 19
client2_1  | 2021/11/16 12:24:13 [Lamport: 20] Queued reply to client3
client2_1  | 2021/11/16 12:24:13 [Lamport: 20] State: 1, T: 15, P: client2, Ti: 19, Pi: client3
client2_1  | 2021/11/16 12:24:13 [Lamport: 20] Exiting receive
^CGracefully stopping... (press Ctrl+C again to force)
Killing itu-disys2021-exercise2_client2_1  ... done
Killing itu-disys2021-exercise2_client3_1  ... done
Killing itu-disys2021-exercise2_client1_1  ... done
```