package server

import (
	"fmt"
	"google.golang.org/grpc"
	"lib/logger"
	"net"
	"protocol/hub"
)


func Listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen fail, err: %v\n", err)
	}
	//register service
	grpcServer := grpc.NewServer()
	hub.RegisterHubCenterServer(grpcServer, &HubServer{})
	
	fmt.Println("hub grpc listen success")
	logger.Info("hub grpc listen success")
	return grpcServer.Serve(listener)
}
