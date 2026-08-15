package gowebsocket

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// magicGUID is defined by RFC 6455 1.3. It's concatenated onto the client's
// Sec-WebSocket-Key before hashing, specifically so that a server which does NOT
// understand WebSocket can't accidentally produce a valid-looking accept value.
const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// NewServerWebSocket is the public entrypoint your handler calls.
func NewServerWebSocket(w http.ResponseWriter, r *http.Request) (Socket, error) {
	return upgrade(w, r, nil)
}

func upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	if r.Method != http.MethodGet {
		return nil, fmt.Errorf("websocket: method must be GET, got %s", r.Method)
	}

	// Connection header must contain the "Upgrade" token
	if !headerContainsToken(r.Header.Get("Connection"), "upgrade") {
		return nil, fmt.Errorf("websocket: 'Connection' header must contain 'Upgrade'")
	}

	// Upgrade header must literally be "websocket"
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return nil, fmt.Errorf("websocket: 'Upgrade' header must be 'websocket'")
	}

	// Sec-WebSocket-Version must be 13
	if r.Header.Get("Sec-Websocket-Version") != "13" {
		return nil, fmt.Errorf("websocket: unsupported version %q, want 13",
			r.Header.Get("Sec-Websocket-Version"))
	}

	// Sec-WebSocket-Key must be present and decode to exactly 16 bytes
	key := r.Header.Get("Sec-Websocket-Key")
	if key == "" {
		return nil, fmt.Errorf("websocket: missing 'Sec-WebSocket-Key' header")
	}
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("websocket: invalid Sec-WebSocket-Key encoding: %w", err)
	}
	if len(decoded) != 16 {
		return nil, fmt.Errorf("websocket: Sec-WebSocket-Key must decode to 16 bytes, got %d", len(decoded))
	}

	// At this point validation is done and we're committing: once Hijack()
	// succeeds, `w` can no longer be used for anything.
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		// Happens e.g. on HTTP/2, or non-hijackable ResponseWriter implementations.
		return nil, fmt.Errorf("websocket: response writer does not support hijacking")
	}
	netConn, rw, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("websocket: hijack failed: %w", err)
	}

	// Compute Sec-WebSocket-Accept and write the 101 response by hand
	// From here on we are fully responsible for every byte on the wire — no more
	// automatic headers, no automatic Content-Length, nothing from net/http.
	accept := computeAcceptKey(key)
	var sb strings.Builder
	sb.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	sb.WriteString("Upgrade: websocket\r\n")
	sb.WriteString("Connection: Upgrade\r\n")
	sb.WriteString("Sec-WebSocket-Accept: " + accept + "\r\n")

	for k, values := range responseHeader {
		for _, v := range values {
			sb.WriteString(k + ": " + v + "\r\n")
		}
	}
	sb.WriteString("\r\n")

	if _, err := rw.WriteString(sb.String()); err != nil {
		netConn.Close()
		return nil, fmt.Errorf("websocket: failed writing handshake response: %w", err)
	}

	if err := rw.Flush(); err != nil {
		netConn.Close()
		return nil, fmt.Errorf("websocket: failed flushing handshake response: %w", err)
	}

	return &Conn{
		conn:     netConn,
		isServer: true,
	}, nil
}

func computeAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte(magicGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func headerContainsToken(header, token string) bool {
	for _, v := range strings.Split(header, ",") {
		if strings.EqualFold(strings.TrimSpace(v), token) {
			return true
		}
	}
	return false
}
