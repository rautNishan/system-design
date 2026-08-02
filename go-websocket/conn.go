package gowebsocket

import (
	"net"
)

type Conn struct {
	conn     net.Conn
	isServer bool
}

func newConnection(conn net.Conn, isServer bool) *Conn {
	c := &Conn{
		conn:     conn,
		isServer: isServer,
	}
	return c
}
