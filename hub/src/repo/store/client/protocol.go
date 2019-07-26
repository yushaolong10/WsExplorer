package client

import (
	"bytes"
	"lib/json"
	"strings"
)

//set
type set struct {
	Key string
	Val string
}

func (s *set) Encode() []byte {
	b := bytes.Buffer{}
	b.WriteString("SET ")
	b.WriteString(strings.TrimSpace(s.Key))
	b.WriteString(" ")
	b.WriteString(strings.TrimSpace(s.Val))
	return b.Bytes()
}

//delete
type del struct {
	Key string
}

func (d *del) Encode() []byte {
	b := bytes.Buffer{}
	b.WriteString("DELETE ")
	b.WriteString(strings.TrimSpace(d.Key))
	return b.Bytes()
}

//get
type get struct {
	Key string
}

func (g *get) Encode() []byte {
	b := bytes.Buffer{}
	b.WriteString("GET ")
	b.WriteString(strings.TrimSpace(g.Key))
	return b.Bytes()
}

//response
type storeProtocol struct {
	Err  int                    `json:"err"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func decodeStoreResp(msg []byte) (*storeProtocol, error) {
	var s storeProtocol
	err := json.Unmarshal(msg, &s)
	return &s, err
}