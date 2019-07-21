package connection

import (
	"context"
	"github.com/mailru/easygo/netpoll"
	"logger"
	"runtime/debug"
	"server/routine"
	"time"
)

func Init(groupCount, maxConnCount int) error {
	err := InitPool(groupCount, maxConnCount)
	if err != nil {
		return err
	}
	InitEpoller()
	return nil
}

func MonitorWsConn(conn *WsConnInfo) error {
	if IsNetPollDegrade() {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		//need fix
		//config routine count limit max connections.
		err := routine.Start(ctx, func(t *routine.Task) error {
			degradeMonitorRead(conn)
			return nil
		})
		return err
	}
	fd, err := netpoll.HandleReadOnce(conn.WsConn.GetNetConn())
	if err != nil {
		logger.Error("[Monitor] begin degrade because of netpoll hand read error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err = routine.Start(ctx, func(t *routine.Task) error {
			degradeMonitorRead(conn)
			return nil
		})
		return err
	}
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
	return EpollReadStart(conn, wsLogicTimeOut, f)
}

//degrade service ensure the service run normally
func degradeMonitorRead(conn *WsConnInfo) {
	defer func() {
		if pErr := recover(); pErr != nil {
			logger.Error("[degradeMonitorRead] #PANIC# error. err:%v, stack:%s", pErr, string(debug.Stack()))
		}
		if err := DeleteWsConnFromPool(conn.UniqId); err != nil {
			logger.Error("[degradeMonitorRead] ws conn closed [DeleteWsConnFromPool] error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
		logger.Info("[degradeMonitorRead] ws conn closed. uniqId:%d", conn.UniqId)
	}()
	conn.WsConn.SetReadLimit(wsMaxMessageSize)
	for {
		message, err := conn.Read()
		if err != nil {
			logger.Error("[degradeMonitorRead] read error. uniqId:%d,err:%s", conn.UniqId, err.Error())
			break
		}
		ctx, _ := context.WithTimeout(context.Background(), wsLogicTimeOut)
		if err := HandleRead(ctx, conn, message); err != nil {
			logger.Error("[degradeMonitorRead] HandleRead handle error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
	}
}
