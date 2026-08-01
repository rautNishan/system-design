package gowebsocket

import (
	"io"
	"net"
)

type Conn struct {
	conn   net.Conn
	reader io.ReadCloser
	writer io.WriteCloser
}
