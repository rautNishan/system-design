package main

import (
	"errors"
	"fmt"
	"net/http"

	gowebsocket "github.com/rautNishan/system-design/go-websocket"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	w.Write([]byte("hi\n"))
}

func serverSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := gowebsocket.NewServerWebSocket(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if socket == nil {
		return
	}
	go handleWebSocketCommunication(socket)
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/ws", serverSocket)
	fmt.Println("Listing on port 3000")
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		panic(err)
	}
}

func handleWebSocketCommunication(socket gowebsocket.Socket) {
	data, err := socket.Read()
	if err != nil {
		if errors.Is(err, gowebsocket.ErrUnmaskedFrame) {
			fmt.Println("Frame not masked")
		}
		return
	}
	fmt.Printf("Server side Data: %s\n", string(data))
	socket.Write(data)
}
