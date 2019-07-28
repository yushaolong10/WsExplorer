package ws

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"lib/logger"
)

type wsClient struct {
	client WsRpcClient
	conn   *grpc.ClientConn
}

func NewWsClient(ctx context.Context, addr string) (*wsClient, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := NewWsRpcClient(conn)
	return &wsClient{conn: conn, client: client}, nil
}

func (ws *wsClient) Publish(ctx context.Context, requestId string, uniqId int64, data []byte) error {
	in := ServeRequest{
		RequestId: requestId,
		UniqId:    uniqId,
		Data:      data,
	}
	reply, err := ws.client.Publish(ctx, &in)
	if err != nil {
		logger.Error("[Publish] ws client publish error. requestId:%s,uniqId:%d,err:%s", requestId, uniqId, err.Error())
		return err
	}
	if reply.ErrCode != 0 {
		logger.Error("[Publish] ws client reply error. requestId:%s,uniqId:%d,message:%s", requestId, uniqId, reply.Message)
		return fmt.Errorf("ws reply err:%s", reply.Message)
	}
	return nil
}

func (ws *wsClient) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}
