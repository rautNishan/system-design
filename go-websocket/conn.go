package gowebsocket

import (
	"bufio"
	"io"
	"net"
)

type Conn struct {
	conn     net.Conn
	isServer bool
	reader   io.ReadCloser
	writer   io.WriteCloser
	br       *bufio.Reader
}

func newConnection(conn net.Conn, isServer bool) *Conn {
	c := &Conn{
		conn:     conn,
		isServer: isServer,
		br:       bufio.NewReader(conn),
	}
	return c
}

func (conn *Conn) Reader() {}
