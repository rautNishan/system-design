package gowebsocket

import (
	"fmt"
	"strings"
)

type WsUri struct {
	WS     string
	Host   string
	Port   string
	Path   string
	Params []string
}

// ws://localhost:3000/ws?hehe
func uriParser(path string) (*WsUri, error) {
	chunks := strings.Split(path, ":")
	if len(chunks) == 0 || len(chunks) < 2 {
		return nil, fmt.Errorf("Invalid URI")
	}
	// ["ws","//localhost","3000/ws?heehe"]
	ws := chunks[0]
	host := chunks[1][2:]
	thirdPortion := strings.Split(chunks[2], "/")
	// ["3000","ws?hehe"]
	port := thirdPortion[0]
	fourthProtion := strings.Split(thirdPortion[1], "?")
	path = fourthProtion[0]
	params := fourthProtion[1:]
	return &WsUri{
		WS:     ws,
		Host:   host,
		Port:   port,
		Path:   path,
		Params: params,
	}, nil
}
