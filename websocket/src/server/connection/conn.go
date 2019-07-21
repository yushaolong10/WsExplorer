package connection

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
	"logger"
	"runtime/debug"
	"server/routine"
	"time"
)

const (
	wsWriteTimeOut = time.Second * 3
	wsReadTimeOut  = time.Second * 3

	wsMaxMessageSize = 4096

	//logic process timeout
	wsLogicTimeOut = time.Second * 2
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WsConnInfo struct {
	UniqId    int
	Actor     int
	Timestamp int64
	EpollFd   *netpoll.Desc
	WsConn    *websocket.Conn
}

func (conn *WsConnInfo) Monitor() error {
	if IsNetPollDegrade() {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		//need fix
		//config routine count limit max connections.
		err := routine.Start(ctx, func(t *routine.Task) error {
			conn.DegradeRead()
			return nil
		})
		return err
	}
	fd, err := netpoll.HandleReadOnce(conn.WsConn.GetNetConn())
	if err != nil {
		logger.Error("[Monitor] begin degrade because of netpoll hand read error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err = routine.Start(ctx, func(t *routine.Task) error {
			conn.DegradeRead()
			return nil
		})
		return err
	}
	//epoll process
	//set timeout ws
	f := func(ctx context.Context) error {
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

func (conn *WsConnInfo) Ping() error {
	conn.WsConn.SetWriteDeadline(time.Now().Add(wsWriteTimeOut))
	if err := conn.WsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
		logger.Error("[Ping] ws conn write error. uniqId:%d, err:%s", conn.UniqId, err.Error())
		return err
	}
	return nil
}

func (conn *WsConnInfo) Pong() error {
	conn.WsConn.SetWriteDeadline(time.Now().Add(wsWriteTimeOut))
	if err := conn.WsConn.WriteMessage(websocket.PongMessage, nil); err != nil {
		logger.Error("[Pong] ws conn write error. uniqId:%d, err:%s", conn.UniqId, err.Error())
		return err
	}
	return nil
}

func (conn *WsConnInfo) DegradeRead() {
	defer func() {
		if pErr := recover(); pErr != nil {
			logger.Error("[DegradeRead] #PANIC# error. err:%v, stack:%s", pErr, string(debug.Stack()))
		}
		if err := conn.Close(); err != nil {
			logger.Error("[DegradeRead] ws conn closed [conn.Close] error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
		if err := DeleteWsConnFromPool(conn.UniqId); err != nil {
			logger.Error("[DegradeRead] ws conn closed [DeleteWsConnFromPool] error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
		logger.Info("[DegradeRead] ws conn closed. uniqId:%d", conn.UniqId)
	}()
	conn.WsConn.SetReadLimit(wsMaxMessageSize)
	for {
		message, err := conn.Read()
		if err != nil {
			logger.Error("[DegradeRead] read error. uniqId:%d,err:%s", conn.UniqId, err.Error())
			break
		}
		ctx, _ := context.WithTimeout(context.Background(), wsLogicTimeOut)
		if err := HandleRead(ctx, conn, message); err != nil {
			logger.Error("[DegradeRead] HandleRead handle error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
	}
}

func (conn *WsConnInfo) Read() ([]byte, error) {
	_, message, err := conn.WsConn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return nil, fmt.Errorf("ws conn unexpected close, detail:[%s]", err.Error())
		}
		return nil, err
	}
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	return message, nil
}

func (conn *WsConnInfo) Send(data []byte) error {
	conn.WsConn.SetWriteDeadline(time.Now().Add(wsWriteTimeOut))
	if err := conn.WsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		logger.Error("[Send] ws conn write error. uniqId:%d, err:%s", conn.UniqId, err.Error())
		return err
	}
	return nil
}

func (conn *WsConnInfo) Close() (err error) {
	if err = EpollStop(conn); err != nil {
		logger.Error("[Close] eplloer stop fd error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	if err = conn.EpollFd.Close(); err != nil {
		logger.Error("[Close] epollfd close error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	if err = conn.WsConn.Close(); err != nil {
		logger.Error("[Close] ws conn close error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	//todo
	//clear redis cache cluster
	return
}
