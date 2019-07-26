package connection

import (
	"context"
	"github.com/mailru/easygo/netpoll"
	"lib/logger"
	"runtime/debug"
	"server/routine"
	"time"
)

func Init(groupCount, maxConnCount int) error {
	err := initPool(groupCount, maxConnCount)
	if err != nil {
		return err
	}
	initEpoller()
	return nil
}

func Monitor(conn *WsConnInfo) error {
	if isNetDegrade() {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		//need fix
		//config routine count limit max connections.
		err := routine.Start(ctx, func(t *routine.Task) error {
			degradeProcess(conn)
			return nil
		})
		return err
	}
	fd, _ := netpoll.HandleReadOnce(conn.WsConn.GetNetConn())
	//epoll process
	//set timeout ws
	f := func(ctx context.Context) error {
		if t, ok := ctx.Deadline(); ok {
			conn.WsConn.SetReadDeadline(t)
		}
		message, err := conn.Read()
		if err != nil {
			logger.Error("[Monitor] epoll event driven, but read msg error. uniqId:%d, err:%s", conn.UniqId, err.Error())
			return err
		}
		return HandleRead(ctx, conn, message)
	}
	conn.EpollFd = fd
	return epollStart(conn, wsLogicTimeOut, f)
}

//degrade service ensure the service run normally
func degradeProcess(conn *WsConnInfo) {
	defer func() {
		if pErr := recover(); pErr != nil {
			logger.Error("[degradeProcess] #PANIC# error. err:%v, stack:%s", pErr, string(debug.Stack()))
		}
		if err := DeleteWsConnFromPool(conn.UniqId); err != nil {
			logger.Error("[degradeProcess] ws conn closed [DeleteWsConnFromPool] error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
		logger.Info("[degradeProcess] ws conn closed. uniqId:%d", conn.UniqId)
	}()
	conn.WsConn.SetReadLimit(wsMaxMessageSize)
	for {
		message, err := conn.Read()
		if err != nil {
			logger.Error("[degradeProcess] read error. uniqId:%d,err:%s", conn.UniqId, err.Error())
			break
		}
		ctx, _ := context.WithTimeout(context.Background(), wsLogicTimeOut)
		if err := HandleRead(ctx, conn, message); err != nil {
			logger.Error("[degradeProcess] HandleRead handle error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
	}
}
