package clc

import (
	"encoding/binary"
	"log"
)

const (
	// MaxMessageSize is the maximum allowed CLC message size in bytes
	// (for sanity checks)
	MaxMessageSize = 1024
)

// Message is a type for all clc messages
type Message interface {
	Parse([]byte)
	String() string
	Reserved() string
	Dump() string
}

// NewMessage checks buf for a clc message and returns an empty message of
// respective type and its length in bytes. Parse the new message before
// actually using it
func NewMessage(buf []byte) (Message, uint16) {
	// check eyecatcher first
	if !hasEyecatcher(buf) {
		return nil, 0
	}

	// make sure message is not too big
	length := binary.BigEndian.Uint16(buf[5:7])
	if length > MaxMessageSize {
		log.Println("Error parsing CLC header: message too big")
		errDump(buf[:HeaderLen])
		return nil, 0
	}

	// return new (empty) message of correct type
	typ := buf[4]
	path := path(buf[7] & 0b00000011)
	switch typ {
	case typeProposal:
		return &proposal{}, length
	case typeAccept:
		// check path to determine if it's smc-r or smc-d
		switch path {
		case smcTypeR:
			return &AcceptSMCR{}, length
		case smcTypeD:
			return &AcceptSMCD{}, length
		}
	case typeConfirm:
		// check path to determine if it's smc-r or smc-d
		switch path {
		case smcTypeR:
			return &confirmSMCR{}, length
		case smcTypeD:
			return &ConfirmSMCD{}, length
		}
	case typeDecline:
		return &decline{}, length
	}

	return nil, 0
}
