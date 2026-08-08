package main

import (
	"fmt"
	"net/http"

	gowebsocket "github.com/rautNishan/system-design/go-websocket"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request")
	fmt.Printf("%+v\n", r)
	w.Write([]byte("hi\n"))
}

func serverSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request in socket")
	socket, err := gowebsocket.NewServerWebSocket(w, r)
	if err != nil {
		fmt.Errorf("Error while createing socket: %+v", err)
	}
	fmt.Println(socket)
	// go handleWebSocketCommunication(socket)
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/ws", serverSocket)
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		panic(err)
	}
	return
}

// func handleWebSocketCommunication(*gowebsocket.Conn) {

// }
