package connection

import (
	"context"
	"fmt"
	"github.com/mailru/easygo/netpoll"
	"lib/logger"
	"server/routine"
	"time"
)

//use event-driven-model
var (
	//netpoll manager
	epoller netpoll.Poller
	//net service degrade
	degrade bool
)

func initEpoller() error {
	poller, err := netpoll.New(nil)
	if err != nil {
		degrade = true
		fmt.Printf("netpoll degrade:%s", err.Error())
		return nil
	}
	epoller, degrade = poller, false
	return nil
}

func epollStart(conn *storeConn, timeout time.Duration, f func(ctx context.Context) error) error {
	return epoller.Start(conn.epollFd, func(event netpoll.Event) {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		err := routine.Start(ctx, func(t *routine.Task) (err error) {
			if event&netpoll.EventReadHup != 0 {
				if err := conn.Close(); err != nil {
					logger.Error("[epollStart] netConn close error. taskId:%s,err:%s", t.GetTaskId(), err.Error())
				}
				return
			}
			err = f(ctx)
			//resume
			if err = epoller.Resume(conn.epollFd); err != nil {
				logger.Error("[epollStart] eplloer resume error. taskId:%s,err:%s", t.GetTaskId(), err.Error())
			}
			return
		})
		if err != nil {
			logger.Error("[epollStart] routine start error.err:%s", err.Error())
		}
	})
}

func epollStop(conn *storeConn) error {
	if conn.epollFd == nil {
		return nil
	}
	if err := epoller.Stop(conn.epollFd); err != nil {
		logger.Error("[epollStop] eplloer stop fd error. err:%s", err.Error())
	}
	if err := conn.epollFd.Close(); err != nil {
		logger.Error("[epollStop] epollfd close error. err:%s", err.Error())
	}
	return nil
}

func isNetDegrade() bool {
	return degrade
}
