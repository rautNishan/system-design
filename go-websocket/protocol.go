type Opcode uint8
type Frame struct {
	FIN        bool
	RSV1       bool
	RSV2       bool
	RSV3       bool
	Opcode     Opcode
	Mask       bool
	PayloadLen uint64
	MaskingKey [4]byte
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


