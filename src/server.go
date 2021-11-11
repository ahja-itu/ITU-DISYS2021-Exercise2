package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/andreaswachs/ITU-DISYS2021-Exercise2/src/service"
	"google.golang.org/grpc"
)

type ReplyHandle = chan uint64

var nodeAddresses = []string{
	"127.0.0.1:5000",
	"127.0.0.1:5001",
	"127.0.0.1:5002",
}

func StartServer() {
	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%v", port)

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
	server.ConnectNodes()

	service.RegisterServiceServer(grpcServer, &server)
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
	otherPeersPorts := strings.Split(os.Getenv("OTHERPEERS"), ",")

	for _, peer := range otherPeersPorts {
		address := fmt.Sprintf("127.0.0.1:%s", peer)
		go s.connectNode(address)
	}
}

func (s *Server) connectNode(nodeAddr string) {
	log.Printf("Connecting to %s..\n", nodeAddr)
	conn, err := grpc.Dial(nodeAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect to %s\n", nodeAddr)
	}

	client := service.NewServiceClient(conn)
	log.Printf("Connected to %s\n", nodeAddr)

	s.nodesLock.Lock()
	defer s.nodesLock.Unlock()
	s.nodes[nodeAddr] = client
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
	receive(message.Timestamp, message.Pid, handle)
	timestamp := <-handle
	return &service.RAReply{Timestamp: timestamp}, nil
}
