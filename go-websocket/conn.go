package gowebsocket

import (
	"encoding/binary"
	"net"
)

type Conn struct {
	conn     net.Conn
	isServer bool
}

func newConnection(conn net.Conn, isServer bool) *Conn {
	c := &Conn{
		conn:     conn,
		isServer: isServer,
	}
	return c
}

type Opcode uint8
type Frame struct {
	FIN        bool
	RSV1       bool
	RSV2       bool
	RSV3       bool
	Opcode     Opcode
	Mask       bool
	MaskingKey [4]byte
	PayloadLen uint64
	Payload    []byte
}

const (
	OpContinuation Opcode = 0x0
	OpText         Opcode = 0x1
	OpBinary       Opcode = 0x2
	OpClose        Opcode = 0x8
	OpPing         Opcode = 0x9
	OpPong         Opcode = 0xA
)

func (c *Conn) parseIncomingRequest() {

}

// TODO implement continious stream
func (c *Conn) makeClientFrame(data []byte) Frame {
	payload := Frame{
		FIN:        true,
		RSV1:       false,
		RSV2:       false,
		RSV3:       false,
		Opcode:     OpText,
		Mask:       false,
		PayloadLen: uint64(len(data)),
		MaskingKey: [4]byte{},
		Payload:    data,
	}
	return payload
}

func (f Frame) Encode() []byte {

	b0 := byte(f.Opcode) //4 bits exactly so if opcode is 1 (0000 0001)

	if f.FIN {
		b0 |= 0x80
	}
	if f.RSV1 {
		b0 |= 0x40
	}
	if f.RSV2 {
		b0 |= 0x20
	}
	if f.RSV3 {
		b0 |= 0x10
	}
	var header []byte

	switch {
	case f.PayloadLen <= 125:
		header = []byte{b0, byte(f.PayloadLen)}
	case f.PayloadLen > 125 && f.PayloadLen < 65535:
		header = make([]byte, 4)
		header[0] = b0
		header[1] = 126
		binary.BigEndian.PutUint16(header[2:4], uint16(f.PayloadLen))
	default:
		header = make([]byte, 4)
		header[0] = b0
		header[1] = 127
		binary.BigEndian.PutUint64(header[2:4], f.PayloadLen)
	}

	
}

func (c *Conn) Write(data []byte) {
	frame := c.makeClientFrame(data)
	c.conn.Write(frame)
}

func (c *Conn) Read() {}
