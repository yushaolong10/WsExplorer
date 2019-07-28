package connection

import (
	"fmt"
	"lib/logger"
	"sync"
)

func initPool(groupCount, maxConnCount int) error {
	if groupCount < 1 {
		groupCount = 1
	}
	if maxConnCount < 1 {
		maxConnCount = 1
	}
	pool = &ConnPool{
		groups:            make(map[int]*ConnGroup, groupCount),
		MaxLimitUserCount: int32(maxConnCount),
		GroupCount:        int32(groupCount),
		calGroupIndex: func(uniqId int) int {
			return uniqId % groupCount
		},
	}
	for i := 0; i < groupCount; i++ {
		pool.groups[i] = &ConnGroup{
			mux:     new(sync.Mutex),
			list:    make(map[int]*WsConnInfo),
			GroupId: i,
		}
	}
	return nil
}

func AddWsConn2Pool(conn *WsConnInfo) error {
	if pool == nil {
		return fmt.Errorf("pool not init")
	}
	index := pool.calGroupIndex(conn.UniqId)
	group, err := pool.GetIndexGroup(index)
	if err != nil {
		return err
	}
	return group.AddWsConn(conn)

}

func GetWsConnFromPool(uniqId int) (*WsConnInfo, error) {
	if pool == nil {
		return nil, fmt.Errorf("pool not init")
	}
	index := pool.calGroupIndex(uniqId)
	group, err := pool.GetIndexGroup(index)
	if err != nil {
		return nil, err
	}
	return group.GetWsConn(uniqId)
}

func DeleteWsConnFromPool(uniqId int) error {
	if pool == nil {
		return fmt.Errorf("pool not init")
	}
	index := pool.calGroupIndex(uniqId)
	group, err := pool.GetIndexGroup(index)
	if err != nil {
		return err
	}
	return group.DeleteWsConn(uniqId)
}

func PingWsConn(uniqId int) error {
	conn, err := GetWsConnFromPool(uniqId)
	if err != nil {
		return err
	}
	return conn.Ping()
}

func PongWsConn(uniqId int) error {
	conn, err := GetWsConnFromPool(uniqId)
	if err != nil {
		return err
	}
	return conn.Pong()
}

var pool *ConnPool

type ConnPool struct {
	groups                map[int]*ConnGroup
	GroupCount            int32
	MaxLimitUserCount     int32
	CurrentTotalUserCount int32
	calGroupIndex         func(int) int
}

func (p *ConnPool) GetIndexGroup(index int) (*ConnGroup, error) {
	if group, ok := p.groups[index]; ok {
		return group, nil
	}
	return nil, fmt.Errorf("index[%d] group not exist", index)
}

type ConnGroup struct {
	mux                   *sync.Mutex
	list                  map[int]*WsConnInfo
	GroupId               int
	CurrentGroupUserCount int32
}

func (g *ConnGroup) AddWsConn(conn *WsConnInfo) error {
	g.mux.Lock()
	old, ok := g.list[conn.UniqId]
	g.list[conn.UniqId] = conn
	g.mux.Unlock()
	if !ok {
		return nil
	}
	if err := old.Close(); err != nil {
		logger.Error("[AddWsConn] old conn close error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	logger.Debug("[AddWsConn] old conn close. uniqId:%d", conn.UniqId)
	//client.NewStoreClient(config.Global.Store.Host)
	return nil
}

func (g *ConnGroup) GetWsConn(uniqId int) (*WsConnInfo, error) {
	g.mux.Lock()
	conn, ok := g.list[uniqId]
	g.mux.Unlock()
	if ok {
		return conn, nil
	}
	return nil, fmt.Errorf("not found conn in group. uniqId:%d,groupId:%d", uniqId, g.GroupId)
}

func (g *ConnGroup) DeleteWsConn(uniqId int) error {
	g.mux.Lock()
	conn, ok := g.list[uniqId]
	if ok {
		delete(g.list, uniqId)
	}
	g.mux.Unlock()
	if !ok {
		return fmt.Errorf("not found conn in group. uniqId:%d,groupId:%d", uniqId, g.GroupId)
	}
	if err := conn.Close(); err != nil {
		logger.Error("[DeleteWsConn] conn close error. uniqId:%d, err:%s", conn.UniqId, err.Error())
	}
	logger.Debug("[DeleteWsConn] delete conn close. uniqId:%d", conn.UniqId)
	return nil
}
