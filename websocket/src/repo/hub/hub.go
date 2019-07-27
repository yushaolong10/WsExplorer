package hub

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"lib/logger"
	"lib/pool"
	"time"
)

type hubClient struct {
	conn   *grpc.ClientConn
	client HubCenterClient
}

func (h *hubClient) Close() error {
	if h.conn != nil {
		return h.conn.Close()
	}
	return nil
}

var innerHubPool pool.Pool

func Init(addr []string, minOpen, maxOpen int, maxLifeTime, timeout int) error {
	newFunc := func() (object pool.ObjectItem, e error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		conn, err := grpc.DialContext(ctx, addr[0], grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		client := NewHubCenterClient(conn)
		return &hubClient{conn: conn, client: client}, nil
	}
	maxLifeTimeSec := time.Second * time.Duration(maxLifeTime)
	timeoutMs := time.Millisecond * time.Duration(timeout)
	p, err := pool.NewNormalizePool(minOpen, maxOpen, maxLifeTimeSec, timeoutMs, newFunc)
	if err != nil {
		return err
	}
	innerHubPool = p
	return nil
}

func SendWsRawByte(ctx context.Context, requestId string, uniqId int, data []byte) error {
	object, err := innerHubPool.Acquire()
	if err != nil {
		logger.Error("[SendWsRawByte] acquire pool error. requestId:%s,uniqId:%d,err:%s", requestId, uniqId, err.Error())
		return err
	}
	conn := object.(*hubClient)
	in := WsRequest{
		RequestId: requestId,
		FromId:    int64(uniqId),
		Data:      data,
	}
	reply, err := conn.client.SendWsRawByte(ctx, &in)
	if err := innerHubPool.Release(conn); err != nil {
		logger.Error("[SendWsRawByte] release pool error. requestId:%s,uniqId:%d,err:%s", requestId, uniqId, err.Error())
	}
	if err != nil {
		logger.Error("[SendWsRawByte] send error. requestId:%s,uniqId:%d,err:%s", requestId, uniqId, err.Error())
		return err
	}
	if reply.ErrCode != 0 {
		logger.Error("[SendWsRawByte] reply error. requestId:%s,uniqId:%d,message:%s", requestId, uniqId, reply.Message)
		return fmt.Errorf("reply err:%s", reply.Message)
	}
	return nil
}
