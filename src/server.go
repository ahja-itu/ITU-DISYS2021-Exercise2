package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/andreaswachs/ITU-DISYS2021-Exercise2/src/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port        = 5000
	hostAddress = fmt.Sprintf("localhost:%d", port)
	srv         = ServiceServer{}
)

func StartServer() {
	listener, err := net.Listen("tcp", hostAddress)
	if err != nil {
		log.Fatalf("Error while attempting to listen on port 3333: %v", err)
	}

	log.Println("Started server")
	server := grpc.NewServer()

	service.RegisterServiceServer(server, &srv)
	server.Serve(listener)
}

type ServiceServer struct {
	service.UnimplementedServiceServer
}

func (s *ServiceServer) Req(context context.Context, message *service.RAMessage) (*service.RAReply, error) {
	// This agent is getting this request on themselves, from another agent
	// message.Pid
	// message.Timestamp

	receive()
	return nil, status.Errorf(codes.Unimplemented, "method Req not implemented")
}
