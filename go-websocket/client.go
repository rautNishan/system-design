package gowebsocket

import (
	"fmt"
	"net"
)

type WebScoket struct {
	Host string
	conn *Conn
}

func NewWebSocket(host string) (*Conn, error) {
	wsUri, err := uriParser(host)
	if err != nil {
		return nil, fmt.Errorf("Invalid uri")
	}
	conn, err := net.Dial(wsUri.Host, wsUri.Port)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting to the server %+v", err)
	}
	c := newConnection(conn, false)
	return c, nil
}
