package grpc

import (
	"config"
	"fmt"
	"google.golang.org/grpc"
	"lib/logger"
	"net"
	"protocol/ws"
)

func Start(errChan chan<- error) {
	listener, err := net.Listen("tcp", config.Global.Grpc.Addr)
	if err != nil {
		errChan <- fmt.Errorf("listen fail, err: %v", err)
	}
	//register service
	grpcServer := grpc.NewServer()
	ws.RegisterWsRpcServer(grpcServer, &WsRpc{})
	go func() {
		fmt.Println("grpc server start...")
		logger.Info("grpc server start addr:%s", config.Global.Grpc.Addr)
		errChan <- grpcServer.Serve(listener)
	}()
}
