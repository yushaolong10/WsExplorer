package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
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

func (c *StoreClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *StoreClient) Set(key, val string) bool {
	b := bytes.Buffer{}
	b.WriteString("SET ")
	b.WriteString(strings.TrimSpace(key))
	b.WriteString(" ")
	b.WriteString(strings.TrimSpace(val))
	b.Bytes()
	c.conn.Write(b.Bytes())
	ret,err := c.Read()
	if err != nil {
		fmt.Printf("set error. err:%s", err.Error())
	}
	fmt.Printf("set response:%s\n", string(ret))
	return false
}

func (c *StoreClient) Get(key string) (string, bool) {
	b := bytes.Buffer{}
	b.WriteString("GET ")
	b.WriteString(strings.TrimSpace(key))
	b.Bytes()
	c.conn.Write(b.Bytes())
	ret,err := c.Read()
	if err != nil {
		fmt.Printf("get error. err:%s", err.Error())
	}
	fmt.Printf("get response:%s\n", string(ret))
	return "", false
}

func (c *StoreClient) Delete(key string) bool {
	b := bytes.Buffer{}
	b.WriteString("DELETE ")
	b.WriteString(strings.TrimSpace(key))
	b.Bytes()
	c.conn.Write(b.Bytes())
	ret,err := c.Read()
	if err != nil {
		fmt.Printf("delete error. err:%s", err.Error())
	}
	fmt.Printf("delete response:%s\n", string(ret))
	return false
}

func (c *StoreClient) Read() ([]byte, error) {
	var buf [70]byte
	n, err := c.conn.Read(buf[:])
	if err != nil {
		return nil, fmt.Errorf("read from connect failed, err: %s", err.Error())
	}
	return buf[:n], nil
}
