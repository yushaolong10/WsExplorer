package connection

import (
	"context"
	"fmt"
	"github.com/mailru/easygo/netpoll"
	"logger"
	"server/routine"
	"time"
)

//use event-driven-model
var (
	//netpoll manager
	eplloer netpoll.Poller
	//sevice degrade
	degrade bool
)

func InitEpoller() {
	poller, err := netpoll.New(nil)
	if err != nil {
		degrade = true
		fmt.Printf("netpoll degrade:%s", err.Error())
		return
	}
	eplloer, degrade = poller, false
}

func IsNetPollDegrade() bool {
	return degrade
}

func EpollReadStart(conn *WsConnInfo, timeout time.Duration, f func(ctx context.Context) error) error {
	return eplloer.Start(conn.EpollFd, func(event netpoll.Event) {
		err := routine.Start(context.Background(), func(t *routine.Task) (err error) {
			if event&netpoll.EventReadHup != 0 {
				if err := conn.Close(); err != nil {
					logger.Error("[EpollReadStart] ws conn close error. uniqId:%d,taskId:%s,err:%s", conn.UniqId, t.GetTaskId(), err.Error())
				}
				return
			}
			ctx, _ := context.WithTimeout(context.Background(), timeout)
			err = f(ctx)
			//resume
			if err = eplloer.Resume(conn.EpollFd); err != nil {
				logger.Error("[EpollReadStart] eplloer resume error. uniqId:%d,taskId:%s,err:%s", conn.UniqId, t.GetTaskId(), err.Error())
			}
			return
		})
		if err != nil {
			logger.Error("[EpollReadStart] routine start error. uniqId:%d,err:%s", conn.UniqId, err.Error())
		}
	})
}

func EpollStop(conn *WsConnInfo) error {
	return eplloer.Stop(conn.EpollFd)
}
