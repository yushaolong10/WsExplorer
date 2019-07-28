package service

import (
	"lib/convert"
	"lib/logger"
	"protocol/hub"
	"repo/store"
	"repo/ws"
)

func Invoke(action string, request *hub.WsRequest) error {
	data := convert.Bytes2String(request.Data)
	hostVal, _ := store.GetUniqIdGrpcHost(request.FromId)
	logger.Info("[Invoke] receive message. requestId:%s, action:%s, fromId:%d, hostVal:%s, data:%s", request.RequestId, action, request.FromId, hostVal, data)
	if err := ws.Publish2RelevantWS(hostVal, request.RequestId, request.FromId, request.Data); err != nil {
		logger.Info("[Invoke] publish message error. requestId:%s", request.RequestId)
		return err
	}
	return nil
}
