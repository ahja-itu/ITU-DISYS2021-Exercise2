package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andreaswachs/ITU-DISYS2021-Exercise2/src/service"
	"google.golang.org/grpc"
)

type ReplyHandle = chan uint64

func StartServer() {
	port := os.Getenv("PORT")
	address := fmt.Sprintf("0.0.0.0:%v", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error while attempting to listen on port %v: %v", port, err)
	}

	log.Printf("Started server on %s", address)
	grpcServer := grpc.NewServer()

	server = Server{
		nodes: make(map[string]service.ServiceClient),
		addr:  address,
	}
	go server.ConnectNodes()

	service.RegisterServiceServer(grpcServer, &server)

	// It looks as if Serve is a blocking function.
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server failed to serve: %v", err)
	}
}

type Server struct {
	service.UnimplementedServiceServer

	nodes     map[string]service.ServiceClient
	addr      string
	nodesLock sync.Mutex
}

func (s *Server) ConnectNodes() {
	otherPeersAddresses := strings.Split(os.Getenv("OTHERPEERSADDRS"), ",")
	time.Sleep(2 * time.Second)
	for _, peer := range otherPeersAddresses {
		go s.connectNode(peer)
	}
}

func (s *Server) connectNode(nodeAddr string) {
	log.Printf("Connecting to %s..\n", nodeAddr)

	var foundConn *grpc.ClientConn = nil

	for {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		conn, _ := grpc.DialContext(ctx, nodeAddr, grpc.WithInsecure(), grpc.WithBlock())

		if conn != nil {
			foundConn = conn
			break
		}
	}

	client := service.NewServiceClient(foundConn)
	log.Printf("Connected to %s\n", nodeAddr)

	s.nodesLock.Lock()
	s.nodes[nodeAddr] = client
	s.nodesLock.Unlock()
}

func (s *Server) Peers() map[string]service.ServiceClient {
	s.nodesLock.Lock()
	defer s.nodesLock.Unlock()

	peers := make(map[string]service.ServiceClient)
	for peerAddr, peer := range s.nodes {
		if peerAddr == s.addr {
			continue
		}
		peers[peerAddr] = peer
	}
	return peers
}

func (s *Server) Req(context context.Context, message *service.RAMessage) (*service.RAReply, error) {
	handle := make(ReplyHandle)
	go receive(message.Timestamp, message.Pid, handle)
	timestamp := <-handle
	return &service.RAReply{Timestamp: timestamp}, nil
}
