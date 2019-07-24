package server

import (
	"fmt"
	"logger"
	"net"
	"server/connection"
)

func Listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen fail, err: %v\n", err)
	}
	fmt.Println("server listen success")
	logger.Info("server listen success")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("[Listen] listener.Accept error. err:%s", err.Error())
			continue
		}
		//handle each connect
		connection.Handle(conn)
	}
}
