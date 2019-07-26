package connection

import (
	"fmt"
	"github.com/mailru/easygo/netpoll"
	"lib/logger"
	"net"
	"time"
)

type storeConn struct {
	netConn net.Conn
	epollFd *netpoll.Desc
}

func (sc *storeConn) Read() ([]byte, error) {
	var buf [70]byte
	n, err := sc.netConn.Read(buf[:])
	if err != nil {
		return nil, fmt.Errorf("read from connect failed, err: %s", err.Error())
	}
	return buf[:n], nil
}

func (sc *storeConn) Write(buf []byte) error {
	_, err := sc.netConn.Write(buf)
	return err
}

func (sc *storeConn) SetReadDeadline(t time.Time) error {
	return sc.netConn.SetReadDeadline(t)
}

func (sc *storeConn) SetWriteDeadline(t time.Time) error {
	return sc.netConn.SetWriteDeadline(t)
}

func (sc *storeConn) Close() error {
	if err := sc.netConn.Close(); err != nil {
		logger.Error("[Close] net conn close error. err:%s", err.Error())
	}
	if err := epollStop(sc); err != nil {
		logger.Error("[Close] epoll Stop error. err:%s", err.Error())
	}
	return nil
}
