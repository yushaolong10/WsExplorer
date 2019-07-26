package connection

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
	"lib/logger"
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

func (conn *WsConnInfo) Write(data []byte) error {
	conn.WsConn.SetWriteDeadline(time.Now().Add(wsWriteTimeOut))
	if err := conn.WsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		logger.Error("[Write] ws conn write error. uniqId:%d, err:%s", conn.UniqId, err.Error())
		return err
	}
	return nil
}

func (conn *WsConnInfo) Close() (err error) {
	if err = epollStop(conn); err != nil {
		logger.Error("[Close] eplloer stop error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	if err = conn.WsConn.Close(); err != nil {
		logger.Error("[Close] ws conn close error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	//todo
	//clear redis cache cluster
	return
}
