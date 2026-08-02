package gowebsocket

import (
	"fmt"
	"net/http"
)

// type WebSocketServer struct {
// }

func NewServerWebSocket(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	return upgrade(w, r, nil)
}

func upgrade(w http.ResponseWriter, r *http.Request, reponseHeader http.Header) (*Conn, error) {
	fmt.Println("In upgrade")
	header := r.Header
	fmt.Println("This is header: %+v\n", header)
	return nil, nil
}
