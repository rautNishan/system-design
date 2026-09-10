package gowebsocket

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

var (
	ErrUnmaskedFrame         = errors.New("gowebsocket: received unmasked frame from client")
	ErrMaskedFrameFromServer = errors.New("gowebsocket: received masked frame from server")
)

type Conn struct {
	conn     net.Conn
	isServer bool
	opcode   Opcode
}

type Socket interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}

func newConnection(conn net.Conn, isServer bool, opcode Opcode) *Conn {
	c := &Conn{
		conn:     conn,
		isServer: isServer,
		opcode:   opcode,
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

// TODO implement continious stream
func (c *Conn) makeClientFrame(data []byte) Frame {
	var key [4]byte
	if _, err := rand.Read(key[:]); err != nil {
		return Frame{}
	}
	payload := Frame{
		FIN:        true,
		RSV1:       false,
		RSV2:       false,
		RSV3:       false,
		Opcode:     c.opcode,
		Mask:       true,
		PayloadLen: uint64(len(data)),
		MaskingKey: [4]byte{},
		Payload:    data,
	}
	return payload
}

func (f Frame) encode() []byte {

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
		header = make([]byte, 10)
		header[0] = b0
		header[1] = 127
		binary.BigEndian.PutUint64(header[2:4], f.PayloadLen)
	}
	var body []byte
	if f.Mask {
		header[1] |= 0x80
		header = append(header, f.MaskingKey[:]...)
		body = make([]byte, len(f.Payload))
		for i, b := range f.Payload {
			body[i] = b ^ f.MaskingKey[i%4]
		}
	} else {
		body = f.Payload
	}
	return append(header, body...)
}

func (c *Conn) Write(data []byte) error {
	frame := c.makeClientFrame(data)
	_, err := c.conn.Write(frame.encode())
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) Read() ([]byte, error) {
	frame, err := c.readFromConn()
	if err != nil {
		return nil, err
	}
	return frame.Payload, nil
}

func (c *Conn) readFromConn() (Frame, error) {
	var f Frame
	// Read the first 2 bytes
	header := make([]byte, 2)
	err := c.readFull(header)
	if err != nil {
		return Frame{}, err
	}

	b0 := header[0]
	b1 := header[1]

	f.FIN = b0&0x80 != 0
	f.RSV1 = b0&0x40 != 0
	f.RSV2 = b0&0x20 != 0
	f.RSV3 = b0&0x10 != 0

	f.Opcode = Opcode(b0 & 0x0F)

	f.Mask = b1&0x80 != 0
	lenField := b1 & 0x7F
	if c.isServer && !f.Mask {
		c.closeWithCode(OpClose, "Frame must be masked")
		return Frame{}, ErrUnmaskedFrame
	}

	switch {
	case lenField <= 125:
		f.PayloadLen = uint64(lenField)
	case lenField == 126:
		ext := make([]byte, 2)
		err := c.readFull(ext)
		if err != nil {
			return Frame{}, err
		}
		f.PayloadLen = uint64(binary.BigEndian.Uint16(ext))
	case lenField == 127:
		ext := make([]byte, 8)
		err := c.readFull(ext)
		if err != nil {
			return Frame{}, err
		}
		f.PayloadLen = uint64(binary.BigEndian.Uint16(ext))
	}
	if f.Mask {
		err = c.readFull(f.MaskingKey[:])
		if err != nil {
			return Frame{}, err
		}
	}

	payload := make([]byte, f.PayloadLen)
	err = c.readFull(payload)
	if err != nil {
		return Frame{}, err
	}

	//Because server does not need to mask
	if f.Mask {
		for i := range payload {
			payload[i] ^= f.MaskingKey[i%4]
		}
	}

	f.Payload = payload
	return f, nil
}

func (c *Conn) readFull(buf []byte) error {
	totalRead := 0
	for totalRead < len(buf) {
		n, err := c.conn.Read(buf[totalRead:])
		if err != nil {
			return err
		}
		totalRead += n
	}
	return nil
}

func (c *Conn) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) closeWithCode(code Opcode, message string) {
	defer c.Close()
	frame := Frame{
		FIN:        true,
		Opcode:     code,
		Mask:       !c.isServer,
		PayloadLen: uint64(len(message)),
		Payload:    []byte(message),
	}
	c.conn.Write(frame.encode())
	time.Sleep(200 * time.Millisecond)
}
