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
	//sevice degrade
	degrade bool
)

func initEpoller() {
	poller, err := netpoll.New(nil)
	if err != nil {
		degrade = true
		fmt.Printf("netpoll degrade:%s", err.Error())
		return
	}
	epoller, degrade = poller, false
}

func isNetDegrade() bool {
	return degrade
}

func epollStart(conn *WsConnInfo, timeout time.Duration, f func(ctx context.Context) error) error {
	return epoller.Start(conn.EpollFd, func(event netpoll.Event) {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		err := routine.Start(ctx, func(t *routine.Task) (err error) {
			if event&netpoll.EventReadHup != 0 {
				if err := DeleteWsConnFromPool(conn.UniqId); err != nil {
					logger.Error("[epollStart] ws conn close error. uniqId:%d,taskId:%s,err:%s", conn.UniqId, t.GetTaskId(), err.Error())
				}
				return
			}
			err = f(ctx)
			//resume
			if err = epoller.Resume(conn.EpollFd); err != nil {
				logger.Error("[epollStart] epoller resume error. uniqId:%d,taskId:%s,err:%s", conn.UniqId, t.GetTaskId(), err.Error())
			}
			return
		})
		if err != nil {
			logger.Error("[epollStart] routine start error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
	})
}

func epollStop(conn *WsConnInfo) error {
	if conn.EpollFd == nil {
		return nil
	}
	if err := epoller.Stop(conn.EpollFd); err != nil {
		logger.Error("[epollStop] epoller stop fd error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	if err := conn.EpollFd.Close(); err != nil {
		logger.Error("[epollStop] epollfd close error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	return nil
}
