package connection

import (
	"context"
	"github.com/mailru/easygo/netpoll"
	"logger"
	"net"
	"server/routine"
	"server/store"
	"time"
)

func HandleConn(netConn net.Conn) {
	conn := &storeConn{netConn: netConn}
	if IsSrvDegrade() {
		err := routine.Start(context.Background(), func(t *routine.Task) (err error) {
			degradeProcess(conn)
			return nil
		})
		if err != nil {
			logger.Error("[handleConn] routine start error.err:%s", err.Error())
		}
		return
	}
	conn.epollFd,_ = netpoll.HandleReadOnce(conn.netConn)
	epollStart(conn, time.Second, func(ctx context.Context) error {
		msg, err := conn.Read()
		if err != nil {
			conn.Close()
			logger.Info("[handleConn] rend conn msg error. err:%s", err.Error())
			return err
		}
		result, err := store.Execute(msg)
		if err != nil {
			logger.Error("[handleConn] store execute error. err:%s", err.Error())
			return err
		}
		conn.Write(result)
		return nil
	})
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
		conn.Write(result)
	}
}
