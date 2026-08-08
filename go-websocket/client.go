package gowebsocket

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// type WebScoket struct {
// 	Host string
// 	conn *Conn
// }

func NewWebSocket(host string) (*Conn, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("Invalud URL")
	}
	//Need to impement wss and proxy on the client side
	conn, err := net.Dial("tcp", url.Host)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting to the socket server")
	}
	request := http.Request{
		Method: "GET",
		URL:    url,
		Header: make(http.Header),
		Proto:  "HTTP/1.1",
		Host:   url.Host,
	}
	key, err := generateWebSocketKey()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error generating key: %w", err)
	}
	request.Header["Upgrade"] = []string{"websocket"}
	request.Header["Connection"] = []string{"Upgrade"}
	request.Header.Set("Sec-WebSocket-Key", key)
	request.Header.Set("Sec-WebSocket-Version", "13")

	err = request.Write(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error writing handshake request: %w", err)
	}

	c := newConnection(conn, false)
	return c, nil
}

func generateWebSocketKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
