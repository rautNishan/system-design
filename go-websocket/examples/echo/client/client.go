package main

import (
	"fmt"
	"net/url"

	gowebsocket "github.com/rautNishan/system-design/go-websocket"
)

func main() {
	url := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "ws"}
	socket, err := gowebsocket.NewWebSocket(url.String())
	if err != nil {
		fmt.Printf("Error while connectiong: %+v", err)
	}
	fmt.Println(socket)
	data := make([]byte, 2)
	data = append(data, 'h')
	data = append(data, 'i')
	socket.Write(data)
	data, err = socket.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Client side data: %s\n", string(data))
	defer socket.Close()
}
