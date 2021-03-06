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
	if !HasEyecatcher(buf) {
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
	ver := buf[7] >> 4
	path := Path(buf[7] & 0b00000011)
	switch typ {
	case TypeProposal:
		if ver == SMCv2 {
			return &ProposalV2{}, length
		}
		return &Proposal{}, length
	case TypeAccept:
		// check path to determine if it's smc-r or smc-d
		switch path {
		case SMCTypeR:
			return &AcceptSMCR{}, length
		case SMCTypeD:
			if ver == SMCv2 {
				return &AcceptSMCDv2{}, length
			}
			return &AcceptSMCD{}, length
		}
	case TypeConfirm:
		// check path to determine if it's smc-r or smc-d
		switch path {
		case SMCTypeR:
			return &ConfirmSMCR{}, length
		case SMCTypeD:
			if ver == SMCv2 {
				return &ConfirmSMCDv2{}, length
			}
			return &ConfirmSMCD{}, length
		}
	case TypeDecline:
		if ver == SMCv2 {
			return &DeclineV2{}, length
		}
		return &Decline{}, length
	}

	return nil, 0
}
