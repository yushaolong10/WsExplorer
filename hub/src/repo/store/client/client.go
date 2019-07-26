package client

import (
	"fmt"
	"net"
	"time"
)

type StoreClient struct {
	conn net.Conn
}

func NewStoreClient(addr string) (*StoreClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("net dial failed, err:%s", err.Error())
	}
	return &StoreClient{conn: conn}, nil
}

func (c *StoreClient) setDefaultDeadline() {
	c.conn.SetDeadline(time.Now().Add(time.Second))
}

func (c *StoreClient) read() ([]byte, error) {
	var buf [70]byte
	n, err := c.conn.Read(buf[:])
	if err != nil {
		return nil, fmt.Errorf("read from connect failed, err: %s", err.Error())
	}
	return buf[:n], nil
}

func (c *StoreClient) Set(key, val string) (bool, error) {
	c.setDefaultDeadline()
	s := set{Key: key, Val: val}
	c.conn.Write(s.Encode())
	read, err := c.read()
	if err != nil {
		return false, err
	}
	store, err := decodeStoreResp(read)
	if err != nil {
		return false, err
	}
	if store.Err != 0 {
		return false, fmt.Errorf("set error:%s", store.Msg)
	}
	ret := store.Data["set"].(bool)
	return ret, nil
}

func (c *StoreClient) Get(key string) (string, bool, error) {
	c.setDefaultDeadline()
	g := get{Key: key}
	c.conn.Write(g.Encode())
	read, err := c.read()
	if err != nil {
		return "", false, err
	}
	store, err := decodeStoreResp(read)
	if err != nil {
		return "", false, err
	}
	if store.Err != 0 {
		return "", false, fmt.Errorf("get error:%s", store.Msg)
	}
	ret := store.Data["get"].(bool)
	val := store.Data["value"].(string)
	return val, ret, nil
}

func (c *StoreClient) Delete(key string) (bool, error) {
	c.setDefaultDeadline()
	d := del{Key: key}
	c.conn.Write(d.Encode())
	read, err := c.read()
	if err != nil {
		return false, err
	}
	store, err := decodeStoreResp(read)
	if err != nil {
		return false, err
	}
	if store.Err != 0 {
		return false, fmt.Errorf("delete error:%s", store.Msg)
	}
	ret := store.Data["delete"].(bool)
	return ret, nil
}

func (c *StoreClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
