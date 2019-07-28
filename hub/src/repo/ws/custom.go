package ws

import (
	"context"
	"lib/logger"
	"time"
)

func Publish2RelevantWS(addr string, requestId string, uniqId int64, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	client, err := NewWsClient(ctx, addr)
	if err != nil {
		logger.Error("[Publish2RelevantWs] new client error. requestId:%s,err:%s", requestId, err.Error())
		return err
	}
	if err := client.Publish(ctx, requestId, uniqId, data); err != nil {
		logger.Error("[Publish2RelevantWs] client publish error. requestId:%s,err:%s", requestId, err.Error())
		return err
	}
	if err := client.Close(); err != nil {
		logger.Error("[Publish2RelevantWs] client close error. requestId:%s,err:%s", requestId, err.Error())
		return err
	}
	return nil
}
