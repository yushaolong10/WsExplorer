package server

import (
	"context"
	"lib/convert"
	"lib/json"
	"lib/logger"
	"protocol/hub"
	"repo/store"
	"service"
	"strconv"
)

type HubServer struct {
}

func (h *HubServer) SendWsRawByte(ctx context.Context, request *hub.WsRequest) (*hub.HubReply, error) {
	reqData := convert.Bytes2String(request.Data)
	var reqMap = make(map[string]interface{})
	if err := json.Unmarshal(request.Data, &reqMap); err != nil {
		logger.Error("[SendWsRawByte] json unmarshal error. requestId:%s, fromId:%d, data:%s", request.RequestId, request.FromId, reqData)
		return setHubReply(request.RequestId, 200001, "json unmarshal error"), nil
	}
	dataAction,ok := reqMap["action"]
	if  !ok {
		logger.Error("[SendWsRawByte] action key not exist. requestId:%s, fromId:%d, data:%s", request.RequestId, request.FromId, reqData)
		return setHubReply(request.RequestId, 200002, "[action] key not exist"), nil

	}
	action, ok := dataAction.(string)
	if !ok {
		logger.Error("[SendWsRawByte] action not string. requestId:%s, fromId:%d, data:%s", request.RequestId, request.FromId, reqData)
		return setHubReply(request.RequestId, 200003, "[action] value not string"), nil
	}
	logger.Error("[SendWsRawByte]  request info. requestId:%s, fromId:%d, data:%s", request.RequestId, request.FromId, reqData)
	if err := service.ParseActionAndExecute(action, request); err != nil {
		logger.Error("[SendWsRawByte]  ParseActionAndExecute error. requestId:%s, fromId:%d, data:%s, err:%s", request.RequestId, request.FromId, reqData, err.Error())
		return setHubReply(request.RequestId, 200004, err.Error()), nil

	}
	return setHubReply(request.RequestId, 0, ""), nil
}

func (h *HubServer) SendAppData(ctx context.Context, request *hub.AppRequest) (*hub.HubReply, error) {
	toId := strconv.FormatInt(request.ToId, 10)
	host, has := store.GetStore(toId)
	if !has {
		logger.Error("[SendAppData] store not found toId key:%s", toId)
		return setHubReply(request.RequestId, 200005, "not found toId cache"), nil
	}
	//send host
	logger.Info("get toId:%s, host:%s", toId, host)
	return setHubReply(request.RequestId, 0, ""), nil

}


func setHubReply(requestId string, code int32, message string) *hub.HubReply {
	reply := hub.HubReply{
		RequestId:requestId,
		ErrCode:code,
		Message:message,
	}
	return &reply
}