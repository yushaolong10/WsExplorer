package connection

import (
	"context"
	"fmt"
	"logger"
)

//todo grpc send
func HandleRead(ctx context.Context, conn *WsConnInfo, message []byte) error {
	logger.Info("get user msg. uniqId:%d, msg:%s", conn.UniqId, string(message))
	resp := fmt.Sprintf("server resp:%s", string(message))
	conn.Send([]byte(resp))
	return nil
}
