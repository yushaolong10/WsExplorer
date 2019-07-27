package service

import (
	"lib/convert"
	"lib/logger"
	"protocol/hub"
	"repo/store"
)

func ParseActionAndExecute(action string, request *hub.WsRequest) error {
	data := convert.Bytes2String(request.Data)
	hostVal, _ := store.GetUniqIdGrpcHost(int(request.FromId))
	logger.Info("receive message. requstId:%s, action:%s, fromId:%d, hostVal:%s, data:%s", request.RequestId, action, request.FromId, hostVal, data)

	return nil
}
