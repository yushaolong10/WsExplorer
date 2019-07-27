package connection

import (
	"context"
	"github.com/mailru/easygo/netpoll"
	"lib/logger"
	"net"
	"server/routine"
	"server/store"
	"time"
)

func Init() error {
	initEpoller()
	return nil
}

func Handle(netConn net.Conn) error {
	conn := &storeConn{netConn: netConn}
	if isNetDegrade() {
		logger.Info("[Handle] net is degrade, use go-routine.")
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		err := routine.Start(ctx, func(t *routine.Task) (err error) {
			degradeProcess(conn)
			return nil
		})
		if err != nil {
			logger.Error("[Handle] routine start error.err:%s", err.Error())
			return err
		}
		return nil
	}
	conn.epollFd, _ = netpoll.HandleReadOnce(conn.netConn)
	handleFunc := func(ctx context.Context) error {
		if deadLine, ok := ctx.Deadline(); ok {
			conn.SetReadDeadline(deadLine)
			conn.SetWriteDeadline(deadLine)
		}
		msg, err := conn.Read()
		if err != nil {
			conn.Close()
			logger.Info("[Handle] rend conn msg error. err:%s", err.Error())
			return err
		}
		result, err := store.Execute(msg)
		if err != nil {
			logger.Error("[Handle] store execute error. err:%s", err.Error())
			return err
		}
		conn.Write(result)
		return nil
	}
	return epollStart(conn, time.Second, handleFunc)
}

func degradeProcess(conn *storeConn) {
	for {
		msg, err := conn.Read()
		if err != nil {
			conn.Close()
			logger.Info("[degradeProcess] rend conn msg error. err:%s", err.Error())
			break
		}
		result, err := store.Execute(msg)
		if err != nil {
			conn.Close()
			logger.Error("[degradeProcess] store execute error. err:%s", err.Error())
			break
		}
		conn.SetWriteDeadline(time.Now().Add(time.Second))
		conn.Write(result)
	}
}
