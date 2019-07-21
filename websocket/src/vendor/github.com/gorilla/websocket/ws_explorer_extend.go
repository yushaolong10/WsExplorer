package websocket

import "net"

//wsExporerExtend
//get net.Conn
func (c *Conn) GetNetConn() net.Conn {
	return c.conn
}
