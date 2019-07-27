package connection

import (
	"context"
	"lib/logger"
	"lib/uuid"
	"repo/hub"
	"time"
)

//grpc send
func HandleRead(ctx context.Context, conn *WsConnInfo, message []byte) error {
	requestId := uuid.MakeUUId()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := hub.SendWsRawByte(ctx, requestId, conn.UniqId, message); err != nil {
		logger.Error("[HandleRead] send ws byte error. requestId:%s,uniqId:%d", requestId, conn.UniqId)
		return err
	}
	logger.Info("send message success. requestId:%s,uniqId:%d, msg:%s", requestId, conn.UniqId, string(message))
	return nil
}
