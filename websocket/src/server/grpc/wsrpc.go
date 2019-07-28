package grpc

import (
	"context"
	"lib/convert"
	"lib/logger"
	"protocol/ws"
	"server/connection"
)

type WsRpc struct {
}

func (ws *WsRpc) Publish(ctx context.Context, in *ws.ServeRequest) (*ws.ServeReply, error) {
	logger.Info("[Publish] receive message. requestId:%s, uniqId:%d, data:%s", in.RequestId, in.UniqId, convert.Bytes2String(in.Data))
	wsconn, err := connection.GetWsConnFromPool(int(in.UniqId))
	if err != nil {
		logger.Error("[Publish] get uniqId error. requestId:%s, uniqId:%d, err:%s", in.RequestId, in.UniqId, err.Error())
		return setWsReply(in.RequestId, 100101, err.Error()), nil
	}
	if err := wsconn.Write(in.Data); err != nil {
		logger.Error("[Publish] ws conn write error. requestId:%s, uniqId:%d, err:%s", in.RequestId, in.UniqId, err.Error())
		return setWsReply(in.RequestId, 100102, err.Error()), nil
	}
	return setWsReply(in.RequestId, 0, ""), nil
}

func setWsReply(requestId string, errCode int32, message string) *ws.ServeReply {
	return &ws.ServeReply{
		RequestId: requestId,
		ErrCode:   errCode,
		Message:   message,
	}
}
